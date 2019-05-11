package main

import (
    "github.com/LongMarch7/higo-web/config"
    "github.com/LongMarch7/higo/app"
    "github.com/LongMarch7/higo/middleware"
    "github.com/LongMarch7/higo/middleware/hystrix"
    "github.com/LongMarch7/higo/middleware/ratelimit"
    "github.com/LongMarch7/higo/middleware/zipkin"
    "strconv"
    "time"
    "github.com/LongMarch7/higo/router"
    "context"
)

func middlewareResolving(config *config.Configer) []middleware.MOption{
    var mOpts []middleware.MOption

    var hOptions []hystrix.HOption
    if config.Config.Middleware.HystrixEPT != 0 {
        hOptions = append(hOptions,hystrix.ErrorPercentThreshold(config.Config.Middleware.HystrixEPT))
    }
    if config.Config.Middleware.HystrixMCR != 0 {
        hOptions = append(hOptions,hystrix.MaxConcurrentRequests(config.Config.Middleware.HystrixMCR))
    }
    if config.Config.Middleware.HystrixRVT != 0 {
        hOptions = append(hOptions,hystrix.RequestVolumeThreshold(config.Config.Middleware.HystrixRVT))
    }
    if config.Config.Middleware.HystrixSW != 0 {
        hOptions = append(hOptions,hystrix.SleepWindow(config.Config.Middleware.HystrixSW))
    }
    if config.Config.Middleware.HystrixTimeout != 0 {
        hOptions = append(hOptions,hystrix.Timeout(config.Config.Middleware.HystrixTimeout))
    }
    if len(hOptions) >0 {
        mOpts = append(mOpts,middleware.HOptions(hOptions))
    }

    var rOptions   []ratelimit.ROption
    if config.Config.Middleware.RatelimitBurst > 0 {
        rOptions = append(rOptions,ratelimit.Burst(config.Config.Middleware.RatelimitBurst))
    }
    if config.Config.Middleware.RatelimitInterval > 0 {
        rOptions = append(rOptions,ratelimit.Interval(time.Duration(config.Config.Middleware.RatelimitInterval)))
    }
    if len(rOptions) >0 {
        mOpts = append(mOpts,middleware.ROptions(rOptions))
    }

    var zOptions   []zipkin.ZOption
    zOptions = append(zOptions,zipkin.Debug(config.Config.Middleware.ZipkinDebug))
    if len(config.Config.Middleware.ZipkinhostPort) > 0 {
        zOptions = append(zOptions,zipkin.HostPort(config.Config.Middleware.ZipkinhostPort))
    }
    if len(config.Config.Middleware.ZipkinUrl) > 0 {
        zOptions = append(zOptions,zipkin.Url(config.Config.Middleware.ZipkinUrl))
    }
    if len(zOptions) >0 {
        mOpts = append(mOpts,middleware.ZOptions(zOptions))
    }
    return mOpts
}
func cliResolving(config *config.Configer) *cliConfig{
    cliConf := new(cliConfig)
    cliConf.router = router.NewRouter()

    cliConf.mOpts = append(cliConf.mOpts,middleware.Prefix(config.Name),middleware.MethodName("request"))
    resolves := middlewareResolving(config)
    if len(resolves) > 0 {
        cliConf.mOpts = append(cliConf.mOpts,resolves...)
    }
    if config.Config.RetryTime > 0 {
        cliConf.cOpts = append(cliConf.cOpts,app.CRetryTime(time.Duration(config.Config.RetryTime)))
    }
    if config.Config.RetryCount > 0 {
        cliConf.cOpts = append(cliConf.cOpts,app.CRetryCount(config.Config.RetryCount))
    }
    if len(config.Config.ConsulServer) > 0 {
        cliConf.cOpts = append(cliConf.cOpts,app.CConsulAddr(config.Config.ConsulServer))
    }
    if len(config.Config.StaticPath) > 0 {
        cliConf.staticPath = config.Config.StaticPath
    }else{
        cliConf.staticPath = "public/static"
    }
    if len(config.Config.UploadPath) > 0 {
        cliConf.uploadPath = config.Config.UploadPath
    }else{
        cliConf.uploadPath = "public/upload"
    }
    if len(config.Config.Port) > 0 {
        cliConf.port = config.Config.Port
    }else{
        cliConf.port = "8080"
    }
    cliConf.serviceList = config.Config.ServiceList
    return cliConf
}

func svrResolving(config *config.Configer) *svrConfig{
    svrConf :=new(svrConfig)

    svrConf.mOpts = append(svrConf.mOpts,middleware.Prefix(config.Name))
    resolves := middlewareResolving(config)
    if len(resolves) > 0 {
        svrConf.mOpts = append(svrConf.mOpts,resolves...)
    }

    if len(config.Name) > 0 {
        svrConf.sOpts = append(svrConf.sOpts,app.SPrefix(config.Name))
    }
    if len(config.Config.ConsulServer) > 0 {
        svrConf.sOpts = append(svrConf.sOpts,app.SConsulAddr(config.Config.ConsulServer))
    }

    svrConf.sOpts = append(svrConf.sOpts,app.SCtx(context.Background()))
    var zOptions   []zipkin.ZOption
    zOptions = append(zOptions,zipkin.Name(config.Name),zipkin.Debug(config.Config.Middleware.ZipkinDebug))
    if len(config.Config.Middleware.ZipkinhostPort) > 0 {
        zOptions = append(zOptions,zipkin.HostPort(config.Config.Middleware.ZipkinhostPort))
    }
    if len(config.Config.Middleware.ZipkinUrl) > 0 {
        zOptions = append(zOptions,zipkin.Url(config.Config.Middleware.ZipkinUrl))
    }
    svrConf.sOpts = append(svrConf.sOpts,app.SzOptions(zOptions))

    if len(config.Config.TemplatePath) > 0 {
        svrConf.TemplatePath = config.Config.TemplatePath
    }else{
        svrConf.TemplatePath = "template"
    }
    var notfouned = true
    for _,value := range config.Config.ServiceList{
        if value.Name == config.Name{
            if len(value.Addr) > 0 {
                svrConf.sOpts = append(svrConf.sOpts,app.SServerAddr(value.Addr))
            }
            if len(value.Port) > 0 {
                port,err := strconv.Atoi(value.Port)
                if err == nil{
                    svrConf.sOpts = append(svrConf.sOpts,app.SServerPort(port))
                }
            }
            if len(value.AdAddr) > 0 {
                svrConf.sOpts = append(svrConf.sOpts,app.SAdvertiseAddress(value.AdAddr))
            }
            if len(value.AdPort) > 0 {
                svrConf.sOpts = append(svrConf.sOpts,app.SAdvertisePort(value.AdPort))
            }
            if len(value.Count) > 0 {
                svrConf.sOpts = append(svrConf.sOpts,app.SMaxThreadCount(value.Count))
            }else{
                svrConf.sOpts = append(svrConf.sOpts,app.SMaxThreadCount("1024"))
            }
            notfouned = false
            break
        }
    }
    if notfouned {
        return nil
    }
    return svrConf
}
