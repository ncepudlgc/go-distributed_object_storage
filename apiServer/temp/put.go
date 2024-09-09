package temp

import (
	"Distributed_Object_Storage/apiServer/locate"
	"Distributed_Object_Storage/src/lib/es"
	"Distributed_Object_Storage/src/lib/rs"
	"Distributed_Object_Storage/src/lib/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	// 从请求的URL中获取token
	token := strings.Split(r.URL.EscapedPath(), "/")[2]

	// 根据token创建可恢复上传流
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 获取当前流的大小
	current := stream.CurrentSize()
	if current == -1 {
		// 如果当前流大小为-1，则表示流不存在
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 从请求头中获取偏移量
	offset := utils.GetOffsetFromHeader(r.Header)
	if current != offset {
		// 如果当前流大小与偏移量不匹配，则返回请求范围不可满足的状态码
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// 创建一个字节切片用于读取数据
	bytes := make([]byte, rs.BLOCK_SIZE)

	for {
		// 从请求体读取数据到字节切片中
		n, e := io.ReadFull(r.Body, bytes)

		if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 更新当前流大小
		current += int64(n)

		// 如果当前流大小超过了目标大小，则回滚操作并返回禁止访问的状态码
		if current > stream.Size {
			stream.Commit(false)
			log.Println("resumable put exceed size")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// 如果读取的字节数不等于块大小且当前流大小不等于目标大小，则直接返回
		if n != rs.BLOCK_SIZE && current != stream.Size {
			return
		}

		// 将读取的数据写入流中
		stream.Write(bytes[:n])

		// 如果当前流大小等于目标大小，则执行后续操作
		if current == stream.Size {
			// 刷新流中的数据到服务器
			stream.Flush()

			// 根据流的信息创建一个可恢复下载流
			getStream, e := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size)

			// 计算可恢复下载流的哈希值
			hash := url.PathEscape(utils.CalculateHash(getStream))

			// 如果哈希值不匹配，则回滚操作并返回禁止访问的状态码
			if hash != stream.Hash {
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// 检查哈希值是否存在于某个列表中
			if locate.Exist(url.PathEscape(hash)) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}

			// 将流的信息添加到版本控制系统中
			e = es.AddVersion(stream.Name, stream.Hash, stream.Size)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}
	}
}
