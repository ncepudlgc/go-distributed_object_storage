package objects

import (
	"Distributed_Object_Storage/apiServer/locate"
	"Distributed_Object_Storage/src/lib/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	// 检查哈希值是否已存在于定位信息中
	if locate.Exist(url.PathEscape(hash)) {
		// 如果存在，则返回 HTTP 状态码 200 OK
		return http.StatusOK, nil
	}

	// 创建写入流
	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		// 如果创建写入流失败，则返回 HTTP 状态码 500 Internal Server Error 和错误信息
		return http.StatusInternalServerError, e
	}

	// 创建一个同时读取输入流和写入流的 TeeReader
	reader := io.TeeReader(r, stream)

	// 计算输入流的哈希值
	d := utils.CalculateHash(reader)
	if d != hash {
		// 如果计算得到的哈希值与请求的哈希值不匹配，则回滚写入流，并返回 HTTP 状态码 400 Bad Request 和错误信息
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}

	// 提交写入流
	stream.Commit(true)

	// 返回 HTTP 状态码 200 OK
	return http.StatusOK, nil
}
