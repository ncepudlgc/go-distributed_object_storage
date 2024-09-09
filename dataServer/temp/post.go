package temp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

func post(w http.ResponseWriter, r *http.Request) {
	// 执行 uuidgen 命令生成 UUID，并获取输出
	// uuidgen 是一个在 Unix/Linux 系统中可用的命令行工具，用于生成 UUID
	output, _ := exec.Command("uuidgen").Output()
	// 去除 UUID 字符串末尾的换行符
	uuid := strings.TrimSuffix(string(output), "\n")
	// 从请求的 URL 中解析出文件名
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 从请求头中获取文件大小，并转换为 int64 类型
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 创建一个 tempInfo 结构体实例，并设置其属性
	t := tempInfo{uuid, name, size}
	// 将 tempInfo 实例写入文件
	e = t.writeToFile()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 创建一个空文件，文件名为 UUID + ".dat"
	os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid + ".dat")
	// 将 UUID 作为响应返回给客户端
	w.Write([]byte(uuid))
}

func (t *tempInfo) writeToFile() error {
	// 创建一个文件，文件路径为环境变量"STORAGE_ROOT"指定的路径加上"/temp/"和tempInfo的Uuid
	f, e := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid)
	if e != nil {
		// 如果创建文件失败，则返回错误
		return e
	}
	// 在函数返回前关闭文件
	defer f.Close()
	// 将tempInfo对象序列化为JSON格式的字节数组
	b, _ := json.Marshal(t)
	// 将字节数组写入文件
	f.Write(b)
	// 返回nil表示成功
	return nil
}
