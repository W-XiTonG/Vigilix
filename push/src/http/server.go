package Http

import (
	"Push/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	customMessage = 1
	agentAlarm    = 2
	zabbixAlarm   = 3
	agentStatus   = 4
)

func Server(push func(Content string), LineBreaks, Port, DeleteString, AuthenticationKeyword string,
	asyncConfig AsyncConfig, LineBreaksStatus, DeleteStringStatus, AuthenticationStatus bool) {
	// 初始化异步处理器
	asyncProcessor := NewAsyncProcessor(push, asyncConfig)
	asyncProcessor.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Client connected: %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		// 检查是否是POST
		if r.Method != http.MethodPost {
			http.Error(w, "Error No POST!", http.StatusMethodNotAllowed)
			log.Printf("[ -> Push] Error No POST! %s", r.Method)
			return
		}
		// 读取请求
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Read request body failed: "+err.Error(), http.StatusBadRequest)
			log.Println("[ -> Push] Read request body failed: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer func() {
			if err = r.Body.Close(); err != nil {
				log.Printf("close body failed :%v", err)
			}
		}()
		log.Printf("[ -> Push] %s", string(body))

		var receiveContent customData
		if err := json.Unmarshal(body, &receiveContent); err != nil {
			http.Error(w, "[ -> Push] Error: data in wrong format: "+err.Error(), http.StatusBadRequest)
			log.Println("[ -> Push] Error: data in wrong format: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 检查Content字段是否为空（处理客户端传递空值的情况）
		if string(receiveContent.Content) == "" {
			http.Error(w, "[ -> Push] Error: Content字段为空", http.StatusBadRequest)
			log.Println("[ -> Push] Error: Content字段为空")
			return
		}
		jsonStr, err := processMessageByType(receiveContent)
		if err != nil {
			http.Error(w, "[ -> Push] Error: data in wrong format: "+err.Error(), http.StatusBadRequest)
		}
		//jsonStr := receiveContent.Content
		// 将请求体内容转换为字符串
		//jsonStr := string(body)
		//Contents := util.AddLineBreaks(jsonStr, LineBreaks)
		var Contents string
		if AuthenticationStatus {
			ok := strings.Contains(jsonStr, AuthenticationKeyword)
			if !ok {
				log.Printf("Authentication keyword error: %s", jsonStr)
				http.Error(w, "认证失败", http.StatusUnauthorized)
				return
			}
		}
		// 替换换行符
		Contents = Conversion(LineBreaksStatus, jsonStr, LineBreaks)
		// 删除指定字符串
		if DeleteStringStatus {
			Contents = util.RemoveChars(Contents, DeleteString)
		}

		// 返回响应
		select {
		case asyncProcessor.taskQueue <- Contents: // 成功入队
			w.WriteHeader(http.StatusAccepted)                                       // 202 状态码
			_, err = fmt.Fprint(w, "The server receives the request successfully\n") // 修改响应内容
			if err != nil {
				log.Printf("write response failed :%v", err)
			}
		default: // 队列已满
			http.Error(w, "系统繁忙，请稍后重试", http.StatusServiceUnavailable)
		}
	})

	// 监听端口
	//port := ":" + YamlConfig.Listening.Port
	//log.Printf("Listening on port %s", port)
	// http.HandleFunc 注册之后
	srv := &http.Server{
		Addr:         ":" + Port,        // 注意变量名一致性
		ReadTimeout:  5 * time.Second,   // 防止慢客户端攻击
		WriteTimeout: 10 * time.Second,  // 响应超时控制
		IdleTimeout:  120 * time.Second, // 长连接超时
		Handler:      nil,               // 使用默认mux
	}
	// 优雅关闭逻辑
	done := make(chan bool)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("\n接收到关闭信号，关闭程序...")

		// 1. 关闭HTTP服务器（停止接收新请求）
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP服务关闭失败: %v", err)
		}

		// 2. 关闭异步处理器（带超时控制）
		asyncStopCtx, asyncCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer asyncCancel()
		asyncDone := make(chan struct{})
		go func() {
			defer close(asyncDone)
			asyncProcessor.Stop() // 停止处理队列任务
		}()

		select {
		case <-asyncDone:
			log.Println("异步处理器已关闭")
		case <-asyncStopCtx.Done():
			log.Println("异步处理器关闭超时，强制退出")
		}

		// 3. 通知主线程退出
		close(done)
	}()

	log.Printf("Listening on port %s", Port)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server start error: %v", err)
	}
}
