package objects

import (
	"Distributed_Object_Storage/dataServer/locate"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func del(w http.ResponseWriter, r *http.Request) {
	// 从请求的URL中获取哈希值
	hash := strings.Split(r.URL.EscapedPath(), "/")[2]

	// 根据哈希值搜索对应的文件，并获取文件列表
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + hash + ".*")

	// 如果文件列表长度不等于1，则直接返回
	if len(files) != 1 {
		return
	}

	// 从定位信息中删除该哈希值
	locate.Del(hash)

	// 将文件重命名为垃圾文件夹下的文件名，并将文件移动到垃圾文件夹中
	os.Rename(files[0], os.Getenv("STORAGE_ROOT")+"/garbage/"+filepath.Base(files[0]))
}
