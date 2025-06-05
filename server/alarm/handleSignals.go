package alarm

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// 处理信号
func signals() chan struct{} {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("[Zabbix] Received system signal: %v", sig)
		close(shutdownSignal)
	}()

	return shutdownSignal
}

func InitSignal() chan struct{} {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		once.Do(func() {
			close(shutdownSignal)
		})
	}()

	return shutdownSignal
}
