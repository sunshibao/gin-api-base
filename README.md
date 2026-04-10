# gin-api-base-framework

gin接口基础框架

## 测试环境
   ```
   服务器：192.168.1.1
   项目地址：/data/api
   ```

## 生产环境
   ```
   生产环境域名：https://xx.aaa.com
   生产环境备用域名：https://xx.bbb.com:1443
   服务器：192.168.1.1、192.168.1.2
   配置文件：config/config.json
   ```

## 框架使用
1. 路由定义  =>  router目录
2. 业务控制器（主要业务逻辑） =>  controller目录
3. 数据库操作  =>  model目录（遵循最小完整功能原则）
4. 数据库驱动 =>  server目录
5. 配置文件  =>  config目录
   ```
   config区分环境配置，根据环境配置文件名区分环境，如：config.json、configPro.json
   flowLimit.json 为限流器配置，不存在则不启用,文件配置如下:
   enable为是否启用流控，name为限流器名称，logPath为日志文件路径，collectIntervalMs为 0 表示关闭系统指标收集，flushIntervalSec为 0 表示关闭监控日志文件异步输出，
   rule为限流规则，数组中每个元素表示一条限流规则，method为请求方法，route为请求路径，qps为每秒请求数（默认值 1000），
   limitCode为限流返回状态码（默认值 429），limitResBody为限流返回内容（默认值 {"code":500,"msg":"服务器繁忙!"}），如：
   ```
   ```json
    {
      "enable": true,
      "name": "flowLimit",
      "logPath": "./logs",
      "collectIntervalMs": 0,
      "flushIntervalSec": 0,
      "rule": [
        {
          "method": "GET",
          "route": "/api/health"
        },
        {
          "method": "DELETE",
          "route": "/api/login",
          "qps": 1000,
          "limitCode": 429,
          "limitResBody": {
            "code": 500,
            "msg": "服务器繁忙!"
          }
        }
      ]
    }
   ```

## 构建
```go
GOOS=linux GOARCH=amd64 go build -o ginApi main.go
```
