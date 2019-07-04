package main

import (
    "github.com/LongMarch7/higo/config"
    "github.com/LongMarch7/higo/app"
    "github.com/LongMarch7/higo/middleware"
    "github.com/LongMarch7/higo/service/base"
    "github.com/LongMarch7/higo/router"
    "github.com/LongMarch7/higo/service/web"
    "google.golang.org/grpc/grpclog"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "time"
)

type cliConfig struct{
    cOpts        []app.COption
    mOpts        []middleware.MOption
    staticPath   string
    uploadPath   string
    serviceList  []config.ServiceList
    router       *router.Router
    port         string
    domain       string
}

var c chan os.Signal
var wg sync.WaitGroup

func Producer(){
Loop:
    for{
        select {
        case s := <-c:
            grpclog.Info("Producer get", s)
            break Loop
        default:
        }
        time.Sleep(500 * time.Millisecond)
    }
    wg.Done()
}

func GateWay(config *config.Configer) {

    grpclog.SetLoggerV2(LogConfig(config.Config.Logger,config.Name + ".log").NewLogger())

    if len(config.Config.ServiceList) == 0 {
        grpclog.Error("not have service")
        return
    }
    cliConf := cliResolving(config)
    //mw := middleware.NewMiddleware(middleware.Prefix("gateway"),middleware.MethodName("request"))
    mw := middleware.NewMiddleware(cliConf.mOpts...)
    client := app.NewClient(cliConf.cOpts...)

    for _,service := range cliConf.serviceList{
        htmlHandler := func(pattern string)  func(http.ResponseWriter, *http.Request){
            return base.MakeReqDataMiddleware(
                web.MakeHtmlCallHandler(client.GetClientEndpoint(service.Name),pattern))
        }
        apiHandler := func(pattern string)  func(http.ResponseWriter, *http.Request){
            return base.MakeReqDataMiddleware(
                web.MakeApiCallHandler(client.GetClientEndpoint(service.Name),pattern))
        }
        switch service.Name {
        case "WebServer":
            client.AddEndpoint(app.CMiddleware(mw),app.CServiceName(service.Name))
            cliConf.router.Add([]router.Routs{
                //{"post|get","/admin",base.MakeReqDataMiddleware(
                //    web.MakeHtmlCallHandler(client.GetClientEndpoint(service.Name),"admin:Index"))},
                //{"post|get","/login/{name}",base.MakeReqDataMiddleware(
                //    web.MakeApiCallHandler(client.GetClientEndpoint(service.Name),"admin:Login"))},
                //admin
                {"get","/admin",htmlHandler("admin:Index")},
                {"get","/admin/info",htmlHandler("admin:Info")},
                {"get","/admin/login",htmlHandler("admin:Login")},
                {"get","/admin/user/role_index",htmlHandler("admin/user:RoleIndex")},
                {"get","/admin/user/role_edit",htmlHandler("admin/user:RoleEdit")},
                {"get","/admin/user/user_add",htmlHandler("admin/user:UserAdd")},
                {"get","/admin/user/user_index",htmlHandler("admin/user:UserIndex")},
                {"get","/admin/user/user_edit",htmlHandler("admin/user:UserEdit")},

                {"post","/admin/login_post",apiHandler("admin:LoginPost")},
                {"post","/admin/logout_post",apiHandler("admin:LogoutPost")},
                {"get","/admin/user/role_list",apiHandler("admin/user:RoleList")},
                {"post","/admin/user/role_delete",apiHandler("admin/user:RoleDelete")},
                {"post","/admin/user/role_status_change",apiHandler("admin/user:RoleStatusChange")},
                {"post","/admin/user/role_edit_post",apiHandler("admin/user:RoleEditPost")},
                {"get","/admin/user/update",apiHandler("admin/user:Update")},
                {"get","/admin/user/user_list",apiHandler("admin/user:UserList")},
                {"post","/admin/user/user_delete",apiHandler("admin/user:UserDelete")},
                {"post","/admin/user/user_status_change",apiHandler("admin/user:UserStatusChange")},
                {"post","/admin/user/user_edit_post",apiHandler("admin/user:UserEditPost")},

                //user
                {"get","/",htmlHandler("user:Index")},
                {"get","/index",htmlHandler("user:Index")},
            })
        }
    }
    cliConf.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",http.FileServer(http.Dir(cliConf.staticPath))))
    cliConf.router.PathPrefix("/upload/").Handler(http.StripPrefix("/static/",http.FileServer(http.Dir(cliConf.uploadPath))))

    c = make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)
    wg.Add(1)
    go http.ListenAndServe(":"+cliConf.port, cliConf.router)
    go Producer()
    wg.Wait()
}

