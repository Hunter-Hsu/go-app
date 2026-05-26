# Go App

這是一個具備可擴展插件架構的領域驅動微框架

## 概述

提供搭建可擴展應用服務的基本功能和結構，隨附常用套件如：gRPC 伺服器，模組化設計可按需替換插件，高度可延展、易用的插件架構且支援依賴注入。

## 開始使用

### 創建方式

參考 [範例程式](#範例程式) 先行準備 `User 模組` 後，以下示範如何創建帶有 User 模組的應用服務。

```golang
package main

import (
  "my_api/internal/user"

  go_app "github.com/shoplineapp/go-app"
  "github.com/shoplineapp/go-app/plugins/grpc/presets"
)

func main() {
  // 建立新應用
  app := go_app.NewApplication()

  // 註冊您的業務邏輯模組
  app.AddModule(&user.UserModule{})

  // 初始化函數：提供所有依賴模組的初始化
  app.Run(func(
    // 在此逐一注入所有需要的依賴模組
    userModule *user.UserModule,
    grpc *presets.DefaultGrpcServerWithNewrelic,
  ) {
  })
}
```

啟動應用

```sh
$ go run cmd/api.go
INFO[0000] PROVIDE plugin *env.Env                      
INFO[0000] PROVIDE plugin *logger.Logger                
INFO[0000] PROVIDE plugin *grpc.GrpcServer              
INFO[0000] PROVIDE plugin *healthcheck.HealthCheckServer 
INFO[0000] PROVIDE plugin *interceptors.RecoveryInterceptor 
INFO[0000] PROVIDE plugin *interceptors.RequestLogInterceptor 
INFO[0000] PROVIDE plugin *interceptors.DeadlineInterceptor 
INFO[0000] PROVIDE plugin *newrelic.NewrelicAgent       
INFO[0000] PROVIDE plugin *stats_handlers.NewrelicStatsHandler 
INFO[0000] PROVIDE plugin *presets.DefaultGrpcServerWithNewrelic 
INFO[0000] PROVIDE plugin *controllers.UsersController  
INFO[0000] PROVIDE plugin *user.UserModule              
INFO[0000] PROVIDE plugin fx.Lifecycle                  
INFO[0000] PROVIDE plugin fx.Shutdowner                 
INFO[0000] PROVIDE plugin fx.DotGraph                   
INFO[0000] = User module init                           
INFO[0000] Application RUNNING                          
INFO[0000] GRPC server is up and running on 0.0.0.0:3000 
^CWARN[0001] Received INTERRUPT                           
INFO[0001] GRPC server gracefully shutting down...      
INFO[0001] Bye.
```

### 範例程式

User 模組

```golang
package user

import (
"my_api/internal/user/controllers"
"my_api/protos"
"github.com/shoplineapp/go-app/plugins/grpc/presets"
"github.com/shoplineapp/go-app/plugins/logger"
go_app "github.com/shoplineapp/go-app"
)

type UserModule struct {
go_app.AppModuleInterface
controller *controllers.UsersController
}

func (m *UserModule) Controllers() []interface{} {
return []interface{}{
    // Register module controller constructors
    controllers.NewUsersController,
}
}

func (m *UserModule) Provide() []interface{} {
return []interface{}{
    // Register module with dependencies
    // Requires all the constructor of structs that you need dependency injection
    func(
    controller *controllers.UsersController,
    grpc *presets.DefaultGrpcServerWithNewrelic,
    logger *logger.Logger,
    ) *UserModule {
    // Register gRPC server with controller
    protos.RegisterUsersServer(grpc.Server(), controller)
    return m
    },
}
}
```

控制器

```golang
package controllers

import (
"context"
"my_api/protos"
"github.com/shoplineapp/go-app/plugins/logger"
)

type UsersController struct {
test.UnimplementedUsersServer
logger *logger.Logger
}

// Constructor of controller with dependencies needed
func NewUsersController(logger *logger.Logger) *UsersController {
c := &UsersController{
    logger: logger,
}
return c
}

// gRPC handler
func (c UsersController) Userinfo(context.Context, *protos.UserinfoRequest) (*protos.UserinfoEmptyResponse, error) {
c.logger.Info("Hello there")
return &protos.UserinfoEmptyResponse{}, nil
}
```

詳細使用資訊請參閱[範例](https://github.com/shoplineapp/go-app/tree/master/examples)。

---

## Plugins

可用的插件放置於 [plugins](https://github.com/shoplineapp/go-app/tree/master/plugins) 資料夾下，例如：

- gRPC Server (帶 `grpc` build tag)
- 常見的 gRPC 攔截器，如請求日誌、伺服器端超時、恢復機制
- Logrus logger
- 環境變數 .env 設定檔和預設值
- Newrelic 整合(帶有 `newrelic` 建置標籤)

插件自動載入，且可透過建構標籤控制選擇是否載入。參考以下指令：

```shell
# Build or run application with build tags
go run -tags sentry cmd/api.go
go build -tags grpc,sentry -o build/api cmd/api.go
```

關於插件更多說明請見[說明文件](https://github.com/shoplineapp/go-app/tree/master/plugins#plugins)。
