package main

import (
    "github.com/LongMarch7/higo-web/config"
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

    if len(config.Config.ServiceList) == 0 {
        grpclog.Error("not have service")
        return
    }
    cliConf := cliResolving(config)
    //mw := middleware.NewMiddleware(middleware.Prefix("gateway"),middleware.MethodName("request"))
    mw := middleware.NewMiddleware(cliConf.mOpts...)
    client := app.NewClient(cliConf.cOpts...)

    for _,service := range cliConf.serviceList{
        switch service.Name {
        case "WebServer":
            client.AddEndpoint(app.CMiddleware(mw),app.CServiceName(service.Name))
            cliConf.router.Add([]router.Routs{
                {"post|get","/admin/{name}",base.MakeReqDataMiddleware(
                    web.MakeHtmlCallHandler(client.GetClientEndpoint(service.Name),"admin:Index"))},
                {"post|get","/cookie/{name}",base.MakeReqDataMiddleware(
                    web.MakeApiCallHandler(client.GetClientEndpoint(service.Name),"admin:Login"))},
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

