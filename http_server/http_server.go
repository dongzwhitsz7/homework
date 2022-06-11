package http_server

import (
	"context"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
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
		log.Info("set loglevel to： %s", logLevel)
		log.SetLevel(level)
	}
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
			log.Info("  [%s]:    http status: %d", request.RemoteAddr, status)
		})
	}
	log.Info("httpserver listen at 8888 por t")
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
