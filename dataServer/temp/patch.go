package temp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func patch(w http.ResponseWriter, r *http.Request) {
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
	// 以只写和追加模式打开数据文件
	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	// 将请求体内容拷贝到数据文件中
	_, e = io.Copy(f, r.Body)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 获取数据文件的信息
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 获取数据文件实际大小
	actual := info.Size()
	// 如果实际大小超过临时信息中的大小
	if actual > tempinfo.Size {
		// 删除数据文件
		os.Remove(datFile)
		// 删除临时信息文件
		os.Remove(infoFile)
		log.Println("actual size", actual, "exceeds", tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func readFromFile(uuid string) (*tempInfo, error) {
	// 打开文件
	f, e := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
	if e != nil {
		// 如果打开文件失败，则返回错误
		return nil, e
	}
	// 延迟关闭文件
	defer f.Close()
	// 读取文件内容
	b, _ := ioutil.ReadAll(f)
	// 定义临时变量info，用于存储解析后的信息
	var info tempInfo
	// 解析JSON数据到info变量中
	json.Unmarshal(b, &info)
	// 返回解析后的信息
	return &info, nil
}
