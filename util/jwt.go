package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
)

var (
	DefaultHeader = JwtHeader{
		Algo: "HS256",
		Type: "JWT",
	}
)

type JwtHeader struct {
	Algo string `json:"alg"` //加密算法一般为sha256
	Type string `json:"typ"` //令牌类型统一写JWT
}

type JwtPayload struct {
	ID          string         `json:"jti"` // JWT ID用于标识该JWT
	Issue       string         `json:"iss"` // 发行人。比如微信
	Audience    string         `json:"aud"` // 受众人。比如王者荣耀
	Subject     string         `json:"sub"` // 主题
	IssueAt     int64          `json:"iat"` // 发布时间,精确到秒
	NotBefore   int64          `json:"nbf"` // 在此之前不可用,精确到秒
	Expiration  int64          `json:"exp"` // 到期时间,精确到秒
	UserDefined map[string]any `json:"ud"`  // 用户自定义的其他字段
}

const (
	JWT_SECRET = "myblog_secret"
)

func GenJwt(header *JwtHeader, payload *JwtPayload, secret string) (string, error) {

	marshal, err := sonic.Marshal(header)
	if err != nil {
		return "", err
	}
	part1 := base64.RawURLEncoding.EncodeToString(marshal) //base64.RawURLEncoding不会在编码结果里添加任何填充字符(=)，也不会把+和/替换成-和_，所以生成的JWT里就不会有这些字符了
	marshal, err = sonic.Marshal(payload)
	if err != nil {
		return "", err
	}
	part2 := base64.RawURLEncoding.EncodeToString(marshal)

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(part1 + "." + part2))
	sign := hash.Sum(nil)
	part3 := base64.RawURLEncoding.EncodeToString(sign)
	return part1 + "." + part2 + "." + part3, nil
}

func VerifyJwt(token, secret string) (*JwtHeader, *JwtPayload, error) {
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return nil, nil, fmt.Errorf("invalid token format")
	}
	//hash 验证
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(split[0] + "." + split[1]))
	sign := hash.Sum(nil)
	part3 := base64.RawURLEncoding.EncodeToString(sign)
	if part3 != split[2] {
		return nil, nil, fmt.Errorf("verify token failed")
	}
	//解析 header 和 payload
	decodeString1, err := base64.RawURLEncoding.DecodeString(split[0])
	if err != nil {
		return nil, nil, fmt.Errorf("header decode failed")
	}
	header := &JwtHeader{}
	if err := sonic.Unmarshal(decodeString1, header); err != nil {
		return nil, nil, fmt.Errorf("header unmarshal failed")
	}
	decodeString2, err := base64.RawURLEncoding.DecodeString(split[1])
	if err != nil {
		return nil, nil, fmt.Errorf("payload decode failed")
	}
	payload := &JwtPayload{}
	if err := sonic.Unmarshal(decodeString2, payload); err != nil {
		return nil, nil, fmt.Errorf("payload unmarshal failed")
	}
	return header, payload, nil
}
