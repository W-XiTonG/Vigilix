package util

import (
	"log"
	"strings"
)

// AddLineBreaks 转换换行符
func AddLineBreaks(data, lineBreak string) string {
	// 打印接收到的数据
	log.Printf("JSON data received: %s", data)
	// 将数据按换行符分割成切片
	lines := strings.Split(data, lineBreak)
	result := strings.Join(lines, "\n")
	return result
}

// RemoveChars 删除指定字符
func RemoveChars(str, chars string) string {
	filter := func(r rune) rune {
		if strings.ContainsRune(chars, r) {
			return -1
		}
		return r
	}
	return strings.Map(filter, str)
}
