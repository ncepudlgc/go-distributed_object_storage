package temp

import (
	"Distributed_Object_Storage/dataServer/locate"
	"Distributed_Object_Storage/src/lib/utils"
	"compress/gzip"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}

func commitTempObject(datFile string, tempinfo *tempInfo) {
	// 打开数据文件
	f, _ := os.Open(datFile)
	defer f.Close()

	// 计算数据的哈希值
	d := url.PathEscape(utils.CalculateHash(f))

	// 将文件指针重定位到文件开头
	f.Seek(0, io.SeekStart)

	// 创建存储对象的文件
	w, _ := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + tempinfo.Name + "." + d)

	// 创建gzip压缩写入器
	w2 := gzip.NewWriter(w)

	// 将数据文件的内容复制到压缩写入器中
	io.Copy(w2, f)

	// 关闭压缩写入器
	w2.Close()

	// 删除临时数据文件
	os.Remove(datFile)

	// 将对象的哈希值和ID添加到定位信息中
	locate.Add(tempinfo.hash(), tempinfo.id())
}
