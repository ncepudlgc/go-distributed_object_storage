package locate

import (
	"Distributed_Object_Storage/src/lib/rabbitmq"
	"Distributed_Object_Storage/src/lib/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// map的 key存放 对象的散列值，value 存放分片id
var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(hash string) int {
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

func StartLocate() {
	// 创建 RabbitMQ 连接
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()

	// 绑定队列 "dataServers"
	q.Bind("dataServers")

	// 创建消费通道
	c := q.Consume()

	// 循环接收消息
	for msg := range c {
		// 去除消息体中的引号，并转换为字符串
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}

		// 调用 Locate 函数查找 id
		id := Locate(hash)

		// 如果 id 不为 -1，则发送消息回复
		if id != -1 {
			// 发送消息到指定的 ReplyTo 队列
			q.Send(msg.ReplyTo, types.LocateMessage{Addr: os.Getenv("LISTEN_ADDRESS"), Id: id})
		}
	}
}

// CollectObjects 函数用于收集存储根目录下所有对象文件的哈希值和ID，并将其存入objects映射中
func CollectObjects() {
	// 获取存储根目录下的所有对象文件路径
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	// 遍历文件路径列表
	for i := range files {
		// 获取文件名（不带路径）
		file := strings.Split(filepath.Base(files[i]), ".")
		// 判断文件名格式是否正确（是否包含三部分）
		if len(file) != 3 {
			panic(files[i])
		}
		// 提取哈希值
		hash := file[0]
		// 将文件名中的第二部分转换为整数类型的ID
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		// 将哈希值和ID存入objects映射中
		objects[hash] = id
	}
}
