package objects

import (
	"Distributed_Object_Storage/dataServer/locate"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func getFile(name string) string {
	// 根据给定的文件名在指定目录下查找文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	// 如果找到的文件数量不等于1，则返回空字符串
	if len(files) != 1 {
		return ""
	}
	// 获取找到的第一个文件的路径
	file := files[0]
	// 创建一个新的sha256哈希对象
	h := sha256.New()
	// 调用sendFile函数计算文件的哈希值
	sendFile(h, file)
	// 将计算得到的哈希值进行base64编码，并转义为URL安全的格式
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	// 从文件路径中提取哈希值
	hash := strings.Split(file, ".")[2]
	// 如果计算得到的哈希值与文件中存储的哈希值不匹配
	if d != hash {
		// 打印日志信息，表示哈希值不匹配，需要删除该文件
		log.Println("object hash mismatch, remove", file)
		// 从定位信息中删除该文件的哈希值
		locate.Del(hash)
		// 删除文件
		os.Remove(file)
		// 返回空字符串
		return ""
	}
	// 返回文件的路径
	return file
}
