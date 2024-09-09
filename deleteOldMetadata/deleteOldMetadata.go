package main

import (
	"Distributed_Object_Storage/src/lib/es"
	"log"
)

const MIN_VERSION_COUNT = 5

func main() {
	// 调用 es.SearchVersionStatus 函数获取版本状态桶，并返回错误
	buckets, e := es.SearchVersionStatus(MIN_VERSION_COUNT + 1)
	if e != nil {
		log.Println(e)
		return
	}

	// 遍历版本状态桶
	for i := range buckets {
		bucket := buckets[i]

		// 遍历每个桶中的文档数，从 MIN_VERSION_COUNT+1 开始到文档数结束
		for v := 0; v < bucket.Doc_count-MIN_VERSION_COUNT; v++ {
			// 调用 es.DelMetadata 函数删除指定版本的元数据
			es.DelMetadata(bucket.Key, v+int(bucket.Min_version.Value))
		}
	}
}
