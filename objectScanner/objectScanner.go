package main

import (
	"Distributed_Object_Storage/apiServer/objects"
	"Distributed_Object_Storage/src/lib/es"
	"Distributed_Object_Storage/src/lib/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")

	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		verify(hash)
	}
}

// verify 函数用于验证传入的哈希值是否正确
// 参数：
// hash：待验证的哈希值
// 返回值：无
func verify(hash string) {
	// 打印验证哈希值的日志
	log.Println("verify", hash)

	// 调用es.SearchHashSize函数查询哈希值对应的大小
	size, e := es.SearchHashSize(hash)
	if e != nil {
		// 如果查询出错，则打印错误日志并返回
		log.Println(e)
		return
	}

	// 调用objects.GetStream函数获取指定哈希值和大小的数据流
	// 这个过程实际上会自动进行数据修复，如果修复后哈希值不匹配，则说明修复失败
	stream, e := objects.GetStream(hash, size)
	if e != nil {
		// 如果获取数据流出错，则打印错误日志并返回
		log.Println(e)
		return
	}

	// 调用utils.CalculateHash函数计算数据流的哈希值
	d := utils.CalculateHash(stream)
	if d != hash {
		// 如果计算得到的哈希值与请求的哈希值不匹配，则打印不匹配日志
		log.Printf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}

	// 关闭数据流
	stream.Close()
}
