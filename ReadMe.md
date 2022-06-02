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
docker build . -t httpserver:0.0.1
docker run -d -p 8080:80 httpserver:0.0.1
```