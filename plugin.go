package testtraefikgolog

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Config 插件配置
type Config struct {
	LogRequest  bool `json:"logRequest,omitempty"`
	LogResponse bool `json:"logResponse,omitempty"`
	LogHeaders  bool `json:"logHeaders,omitempty"`
}

// CreateConfig 创建默认配置
func CreateConfig() *Config {
	return &Config{
		LogRequest:  true,
		LogResponse: true,
		LogHeaders:  false,
	}
}

// LogPlugin 日志插件结构体
type LogPlugin struct {
	next   http.Handler
	name   string
	config *Config
}

// New 创建新的插件实例
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &LogPlugin{
		next:   next,
		name:   name,
		config: config,
	}, nil
}

// ServeHTTP 处理HTTP请求
func (p *LogPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()

	// 记录请求信息
	if p.config.LogRequest {
		p.logRequest(req)
	}

	// 创建自定义ResponseWriter来捕获状态码
	customRW := &responseWriter{
		ResponseWriter: rw,
		statusCode:     http.StatusOK, // 默认状态码
	}

	// 调用下一个处理器
	p.next.ServeHTTP(customRW, req)

	// 记录响应信息
	if p.config.LogResponse {
		duration := time.Since(start)
		p.logResponse(req, customRW.statusCode, duration)
	}
}

// logRequest 记录请求日志
func (p *LogPlugin) logRequest(req *http.Request) {
	log.Printf("[REQUEST] %s %s %s", req.Method, req.URL.Path, req.Proto)
	log.Printf("[REQUEST] Host: %s", req.Host)
	log.Printf("[REQUEST] RemoteAddr: %s", req.RemoteAddr)
	log.Printf("[REQUEST] User-Agent: %s", req.UserAgent())

	if p.config.LogHeaders {
		log.Printf("[REQUEST] Headers:")
		for name, values := range req.Header {
			for _, value := range values {
				log.Printf("  %s: %s", name, value)
			}
		}
	}
}

// logResponse 记录响应日志
func (p *LogPlugin) logResponse(req *http.Request, statusCode int, duration time.Duration) {
	log.Printf("[RESPONSE] %s %s - Status: %d - Duration: %v",
		req.Method, req.URL.Path, statusCode, duration)
}

// responseWriter 自定义ResponseWriter用于捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 重写WriteHeader方法以捕获状态码
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
