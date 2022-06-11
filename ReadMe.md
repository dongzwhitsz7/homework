# 目录结构
* 程序模块在modulex,实现的逻辑
* 运行模块在test/modulex下面，为测试代码或者启动代码
* 当前测试使用的go版本：go 1.17

# 模块一
略

# 模块二
通过运行`test/module2/exercise_test.go`里面的`TestExcersizeHttpServer`方法启动一个HttpServer
* 通过反射将所有的函数注册到路径上
* 通过为ItemService添加方法来扩展逻辑，
* HttpServer的日志打印和路径映射被httpServer封装
* 每一个ItemService的方法都会被映射成一个全小写的路径

已实现的功能
* 接收客户端 request，并将 request 中带的 header 写入 response header
* 读取当前系统的环境变量中的 VERSION 配置，并写入 response header，如果没有返回默认信息
* Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
* 当访问 localhost/healthz 时，应返回 200

测试client脚本：
```shell
curl -v --location --request GET "http://127.0.0.1:80/healthz" --header "Content-Type: application/json"
```

# 模块三
```shell
docker build . -t httpserver
docker run -d -p 8888:8888 httpserver
curl "http://127.0.0.1:8888/healthz"


# 上传
docker build . -t httpserver
docker tag httpserver dongzw/httpserver:v4.0
docker push dongzw/httpserver:v4.0

docker tag c5507fd0cdbf dongzw/httpserver:v2.0 # 为容器镜像添加标签
docker push dongzw/httpserver:v2.0
```

# 模块八

1、优雅启动

      readinessProbe:
        httpGet:
          path: /healthz
          port: 8888
        initialDelaySeconds: 30
        periodSeconds: 5
        successThreshold: 2   

2、优雅终止
    
在代码中新增一个处理信号的协程
```shell
go func() {
    // 优雅终止
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    select {
    case <-sigs:
        log.Info("notify sigs\n")
        if err := httpServer.Shutdown(context.Background()); err != nil {
            log.Error(err)
        }
        log.Info("http shutdown gracefully\n")
    }
}()
```
3、资源需求和QoS保证
    
制定cpu和内存的requests和limits，不一致时为Burstable类型的QoS

    resources:
        limits:
          cpu: 500m
          memory: 512Mi
        requests:
          cpu: 250m
          memory: 256Mi

4、探活

      # 探活
      livenessProbe:
        httpGet:
          path: /healthz
          port: 8888
        initialDelaySeconds: 30
        periodSeconds: 5

5、日常运维需求，日志等级


6、配置和代码分离
