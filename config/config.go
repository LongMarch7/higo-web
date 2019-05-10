package config

import (
    "encoding/json"
    "google.golang.org/grpc/grpclog"
    "io/ioutil"
)

const SvrMode = 1
const CliMode = 2

type Middleware struct{
    //ZipkinName          string    `json:"zipkin_name"`
    ZipkinUrl           string    `json:"zipkin_url"`
    ZipkinhostPort      string    `json:"zipkinhost_port"`
    ZipkinDebug         bool      `json:"zipkin_debug"`

    RatelimitInterval   int64     `json:"ratelimit_interval"`
    RatelimitBurst      int       `json:"ratelimit_burst"`

    //PrometheusSubsystem string    `json:"prometheus_subsystem"`
    //PrometheusName      string    `json:"prometheus_name"`

    //HystrixName         string    `json:"hystrix_name"`
    HystrixTimeout      int       `json:"hystrix_timeout"`
    HystrixMCR          int       `json:"hystrix_mcr"`    //maxConcurrentRequests
    HystrixRVT          int       `json:"hystrix_rvt"`    //requestVolumeThreshold
    HystrixSW           int       `json:"hystrix_sw"`     //sleepWindow
    HystrixEPT          int       `json:"hystrix_ept"`    //errorPercentThreshold

    //LogName             string    `json:"log_name"`
    //LogMethodName       string    `json:"log_method_name"`
}
type ServiceList struct {
    Name         string  `json:"name"`
    Addr         string  `json:"addr"`
    Port         string  `json:"port"`
    Count        string  `json:"count"`   //connect max num
    AdAddr       string  `json:"ad_addr"` //consul advertise address
    AdPort       string  `json:"ad_port"` //consul advertise port
}
type Config struct{
    Domain       string        `json:"domain"`
    SslKey       string        `json:"ssl_key"`
    SslCrt       string        `json:"ssl_crt"`
    Port         string        `json:"port"`
    ConsulServer string        `json:"consul_server"`
    ServiceList  []ServiceList `json:"service_list"`
    RetryCount   int           `json:"retry_count"`
    RetryTime    int64         `json:"retry_time"`
    Middleware   Middleware    `json:"middleware"`
    StaticPath   string        `json:"static_path"`
    UploadPath   string        `json:"upload_path"`
    TemplatePath   string      `json:"template_path"`
}

type Configer struct {
    Mode       int
    Name       string
    Path       string
    Config     Config
}
func ConfigInit(cli bool, svr bool, name string, path string) *Configer{
    config :=new(Configer)
    if cli && svr {
        grpclog.Error("Can't set client and server modes at the same time")
        return nil
    }else if !cli && !svr{
        grpclog.Error("Not set mode")
        return nil
    }else if cli {
        config.Mode = CliMode
    }else{
        config.Mode = SvrMode
    }
    if len(name)==0 {
        grpclog.Error("Not set server name")
        return nil
    }
    config.Name = name
    if len(path) == 0 {
        path = "config.json"
    }
    config.Path = path
    config.Config = Config{}
    err := Load(path, &config.Config)
    if err != nil {
        grpclog.Error("Read config failed")
        return nil
    }
    return config
}

func Load(filename string, v interface{}) error{
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        grpclog.Error(err.Error())
        return err
    }
    err = json.Unmarshal(data, v)
    if err != nil {
        grpclog.Error(err.Error())
        return err
    }
    return nil
}
