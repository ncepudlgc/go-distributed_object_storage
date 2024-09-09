package temp

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	// 从请求的URL中获取uuid
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]

	// 从文件中读取临时信息
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 拼接临时信息文件的路径
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	// 拼接数据文件的路径
	datFile := infoFile + ".dat"

	// 打开数据文件
	f, e := os.Open(datFile)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// 获取文件信息
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 获取文件实际大小
	actual := info.Size()

	// 删除临时信息文件
	os.Remove(infoFile)

	// 如果实际大小与临时信息中的大小不匹配
	if actual != tempinfo.Size {
		// 删除数据文件
		os.Remove(datFile)
		log.Println("actual size mismatch, expect", tempinfo.Size, "actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 提交临时对象
	commitTempObject(datFile, tempinfo)
}
