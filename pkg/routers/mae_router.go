package routers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/MAE/pkg/handlers"
)

type MaeWebHookServer struct {
	Server *gin.Engine
	Addr   string
	Port   int
	Cert   string
	Key    string
}

func NewMaeWebHookServer(addr string, port int, cert string, key string) *MaeWebHookServer {
	return &MaeWebHookServer{
		Server: gin.Default(),
		Addr:   addr,
		Port:   port,
		Cert:   cert,
		Key:    key,
	}
}

func (m *MaeWebHookServer) Start() error {
	m.AddPostRoute("/mutate", handlers.MaeHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", m.Addr, m.Port),
		Handler: m.Server,
	}

	go func() {
		if err := srv.ListenAndServeTLS(m.Cert, m.Key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)

	// 接收到中断信号时关闭服务器
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	// 5秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")

	return nil
}

func (m *MaeWebHookServer) AddPostRoute(route string, handler gin.HandlerFunc) {
	m.Server.POST(route, handler)
}

func (m *MaeWebHookServer) Validate() error {
	if m.Addr == "" {
		m.Addr = "0.0.0.0"
	}

	if m.Port == 0 {
		m.Port = 8443
	}

	if m.Cert == "" {
		return errors.New("cert file is required")
	}

	if m.Key == "" {
		return errors.New("key file is required")
	}

	return nil
}
