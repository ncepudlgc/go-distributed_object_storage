package main

import (
	"Distributed_Object_Storage/src/lib/es"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// main 函数遍历存储目录下的所有对象文件
func main() {
	// 获取所有存储对象的文件路径
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")

	// 遍历每个文件路径
	for i := range files {
		// 提取文件名中的哈希值
		hash := strings.Split(filepath.Base(files[i]), ".")[0]

		// 检查哈希值是否存在于元数据中
		hashInMetadata, e := es.HasHash(hash)
		if e != nil {
			log.Println(e)
			return
		}

		// 如果哈希值不存在于元数据中，则删除对应的文件
		if !hashInMetadata {
			del(hash)
		}
	}
}

func del(hash string) {
	log.Println("delete", hash)
	url := "http://" + os.Getenv("LISTEN_ADDRESS") + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	client.Do(request)
}
