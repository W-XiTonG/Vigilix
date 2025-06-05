package dingding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// Calculate 计算签名函数
func Calculate(timestamp, secret string) (string, error) {
	// 将时间戳和密钥拼接成一个字符串
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, secret)
	// 创建一个 HMAC-SHA256 哈希对象
	h := hmac.New(sha256.New, []byte(secret))
	// 向哈希对象写入要签名的字符串
	_, err := h.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}
	// 计算哈希值
	signatureBytes := h.Sum(nil)
	// 对哈希值进行 Base64 编码
	signature := base64.StdEncoding.EncodeToString(signatureBytes)
	return signature, nil
}
