package global

import (
	"path"
	"runtime"
)

var (
	ProjectRootPath string = getoncurrentPath() + "/../"
)

func getoncurrentPath() string {
	_, file, _, _ := runtime.Caller(0) //返回第0级(偏移量为0)调用堆栈信息
	return path.Dir(file)              //返回路径的目录部分
}
