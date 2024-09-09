package temp

import (
	"net/http"
	"os"
	"strings"
)

func del(w http.ResponseWriter, r *http.Request) {
	// 从请求的URL中获取uuid
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 拼接临时信息文件的路径
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	// 拼接数据文件的路径
	datFile := infoFile + ".dat"
	// 删除临时信息文件
	os.Remove(infoFile)
	// 删除数据文件
	os.Remove(datFile)
}
