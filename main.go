package main

import (
    "flag"
    "github.com/LongMarch7/higo-web/config"
    "github.com/LongMarch7/higo/util/log"
    "google.golang.org/grpc/grpclog"
)

func main(){
    cliMode    := flag.Bool("c",false,"is client mode")
    svrMode    := flag.Bool("s",false,"is service mode")
    svrName    := flag.String("n","","service name")
    configPath := flag.String("p","./config.json","config path")
    flag.Parse()

    grpclog.SetLoggerV2(zap.NewDefaultLoggerConfig().NewLogger())
    conf := config.ConfigInit(*cliMode, *svrMode,*svrName,*configPath)
    if conf == nil{
        grpclog.Error("star error by init config")
        return
    }
    if conf.Mode == config.CliMode {
        GateWay(conf)
    }else{
        Server(conf)
    }
}