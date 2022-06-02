package module2

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type httpServer struct {
	handler *ItemService
}

func NewHttpServer() *httpServer {
	return &httpServer{handler: &ItemService{}}
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
			log.Printf("  [%s]:    http status: %d\n", request.RemoteAddr, status)
		})
	}
	log.Printf("httpserver listen at 80 port")
	if err := http.ListenAndServe("localhost:80", nil); err != nil {
		log.Fatal(err)
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
