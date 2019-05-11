package main
import (
    "github.com/LongMarch7/higo-web/config"
    "github.com/LongMarch7/higo/app"
    "github.com/LongMarch7/higo/middleware"
    "github.com/LongMarch7/higo/service/web"
    _ "github.com/LongMarch7/higo-web/controller/admin"
    "github.com/LongMarch7/higo/view"
)

type svrConfig struct{
    sOpts        []app.SOption
    mOpts        []middleware.MOption
    TemplatePath string
}
func Server(config *config.Configer) {
    svrConf :=svrResolving(config)

    view.NewView(view.Dir(svrConf.TemplatePath))
    server := app.NewServer(svrConf.sOpts...)
    webServer := &web.GrpcServer{}
    webService := web.NewService()
    manager := middleware.NewMiddleware()
    HtmlOpts := append(svrConf.mOpts, middleware.MethodName("HTML"),middleware.Endpoint(web.MakeHtmlCallServerEndpoint(webService)))
    webServer.HtmlCallHandler = manager.AddMiddleware(HtmlOpts...).NewServer()
    ApiOpts := append(svrConf.mOpts, middleware.MethodName("API"), middleware.Endpoint(web.MakeApiCallServerEndpoint(webService)))
    webServer.ApiCallHandler = manager.AddMiddleware(ApiOpts...).NewServer()
    server.RegisterServiceServer(web.MakeRegisteFunc(webServer))
    server.Run()
}

