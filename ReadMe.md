# 模块一、模块二
https://gitee.com/dvge/dongzwhom

# 模块三
```shell
docker build . -t httpserver
docker run -d -p 8888:8888 httpserver
curl "http://127.0.0.1:8888/healthz"


# 上传
docker build . -t httpserver
docker tag httpserver dongzw/httpserver:v6.0
docker push dongzw/httpserver:v6.0

docker tag c5507fd0cdbf dongzw/httpserver:v2.0 # 为容器镜像添加标签
docker push dongzw/httpserver:v2.0
```

# 模块八
## 第一部分

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
    
通过 `log "github.com/sirupsen/logrus"`库设置日志等级

6、配置和代码分离

将configmap挂载在系统文件中，通过`"github.com/spf13/viper"`实现动态读取配置文件

## 第二部分
7、service和ingress

Ingress安装
```shell
kubectl create -f nginx-ingress-deployment.yaml
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=cncamp.com/O=cncamp" -addext "subjectAltName = DNS:cncamp.com"
kubectl create secret tls cncamp-tls --cert=./tls.crt --key=./tls.key
kubectl create -f ingress.yaml
```


# 模块10
prometheus安装和使用,代码在module10中
```shell
helm repo add grafana https://grafana.github.io/helm-charts
helm upgrade --install loki grafana/loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
# 修改为NodePort模式
kubectl edit svc loki-grafana -oyaml -n default
kubectl edit svc loki-prometheus-server -oyaml -n default
# 获取password
kubectl get secret loki-grafana -oyaml -n default
echo 'xxx' | base64 -d # password需要使用base64解码
# 用户名也可以获得，但是解码是admin 
```
执行完 http_server:8888/hello方法之后，在prometheus的ui界面，搜索指标：httpserver_execution_latency_seconds_bucket