package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(text string) string { //md5加密
	hash := md5.New()                 //生成加密变量
	hash.Write([]byte(text))          //写入加密变量
	digest := hash.Sum(nil)           //加密生成加密结果
	ans := hex.EncodeToString(digest) //字节数组转成字符串
	return ans
}
