package temp

import (
	"Distributed_Object_Storage/src/lib/rs"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func head(w http.ResponseWriter, r *http.Request) {
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

	// 设置响应头的Content-Length为当前流的大小
	w.Header().Set("content-length", fmt.Sprintf("%d", current))
}
