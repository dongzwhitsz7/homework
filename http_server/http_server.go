package http_server

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
)

type httpServer struct {
	handler *ItemService
	viper   *viper.Viper
}

func NewHttpServer() *httpServer {
	// 动态配置
	viperInstance := viper.New()
	viperInstance.SetConfigFile("/etc/httpserver/httpserver.properties")
	if err := viperInstance.ReadInConfig(); err != nil {
		log.Error("configuration file read error")
	}
	viperInstance.WatchConfig()
	viperInstance.OnConfigChange(func(in fsnotify.Event) {
		log.Info("log level changed")
		logLevel := viperInstance.GetString("log_level")
		if level, err := log.ParseLevel(logLevel); err != nil {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(level)
		}
	})
	logLevel := viperInstance.GetString("log_level")
	if level, err := log.ParseLevel(logLevel); err != nil {
		log.Info("current log level：default as debug")
		log.SetLevel(log.DebugLevel)
	} else {
		log.Infof("set loglevel to： %s\n", logLevel)
		log.SetLevel(level)
	}

	// 注册metrics
	RegisterProm()

	return &httpServer{handler: &ItemService{}, viper: viperInstance}
}

func (s *httpServer) Serve() {
	typeOf := reflect.TypeOf(s.handler)
	valueOf := reflect.ValueOf(s.handler)
	// 通过反射将所有的函数注册到路径上
	for i := 0; i < typeOf.NumMethod(); i++ {
		method := valueOf.Method(i)
		http.HandleFunc("/"+strings.ToLower(typeOf.Method(i).Name), func(writer http.ResponseWriter, request *http.Request) {
			method.Call([]reflect.Value{reflect.ValueOf(writer), reflect.ValueOf(request)})
			status := reflect.ValueOf(writer).Elem().FieldByName("status").Int()
			log.Infof("  [%s]:    http status: %d\n", request.RemoteAddr, status)
		})
	}
	log.Info("httpserver listen at 8888 port")
	httpServer := http.Server{Addr: ":8888", Handler: nil}
	go func() {
		// 优雅终止
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		select {
		case <-sigs:
			log.Info("received signal")
			if err := httpServer.Shutdown(context.Background()); err != nil {
				log.Error(err)
			}
			log.Info("http shutdown gracefully")
		}
	}()
	if err := httpServer.ListenAndServe(); err != nil {
		log.Error(err)
	}
}

// 1、通过为ItemService添加方法来扩展逻辑，
// 2、HttpServer的日志打印和路径映射被httpServer封装
// 3、每一个ItemService的方法都会被映射成一个全小写的路径
type ItemService struct {
}

func (s *ItemService) Healthz(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		w.Header().Set(k, v[0])
	}
	version := "default version, not found in environment"
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "VERSION=") {
			version = strings.Split(v, "=")[1]
			break
		}
	}
	w.Header().Set("VERSION", version)
	w.WriteHeader(http.StatusOK)
}

func (s *ItemService) Ping(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		w.Header().Set(k, v[0])
	}
	hostname, _ := os.Hostname()
	_, _ = w.Write([]byte(fmt.Sprintf("hostname: %s, host address: %s\n", hostname, r.Host)))
	//w.Write([]byte(fmt.Sprintf("remote address: %s", r.RemoteAddr)))
	w.WriteHeader(http.StatusOK)
}


func (s *ItemService) Metrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func (s *ItemService) Hello(w http.ResponseWriter, r *http.Request) {
	log.Info("entering root handler")
	timer := NewTimer()
	defer timer.ObserveTotal()
	user := r.URL.Query().Get("user")
	delay := randInt(10,2000)
	time.Sleep(time.Millisecond*time.Duration(delay))
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	log.Info("Respond in %d ms", delay)
}