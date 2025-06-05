package Http

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// CORS 中间件函数
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置允许的来源（* 表示允许所有来源，生产环境建议指定具体域名）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许的 HTTP 方法
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		// 允许的请求头（根据实际需求添加，如 Content-Type、Authorization）
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// 允许携带 Cookie（如需使用 Cookie，需设置为 true 且 Allow-Origin 不能为 *）
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求（OPTIONS 方法）
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 辅助函数：格式化响应
func respondSuccess(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, msg)
	log.Printf("响应成功: %s", msg)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprint(w, msg)
	log.Printf("响应错误 %d: %s", code, msg)
}

// 辅助函数：格式化ID列表
func formatIDs(ids []int) string {
	if len(ids) == 0 {
		return ""
	}
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
}
