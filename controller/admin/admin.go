package admin

import (
    "bytes"
    "github.com/LongMarch7/higo-web/models/admin"
    "github.com/LongMarch7/higo/controller/base"
    "github.com/LongMarch7/higo/util/define"
    "github.com/LongMarch7/higo/view"
    "golang.org/x/net/context"
    "google.golang.org/grpc/grpclog"
    "runtime"
    "time"
)

type adminController struct {
}
var Controller = &adminController{}
var name = "admin"
func init(){
    base.AddController(name, Controller)
}

func (a* adminController)Index(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,info :=admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            roleName := info.RoleName
            if len(roleName) == 0 {
                roleName = define.DefaultRoleName
            }
            out := &bytes.Buffer{}
            data := make(map[string]interface{})
            data["CurrUserMenu"] = admin.GetMenuViewList(roleName)
            grpclog.Info(data["CurrUserMenu"])
            data["AppName"] = "higo-web"
            data["Version"] = "1.0.1"
            data["UserName"] = param.Cookie.U
            view.NewView().Render(out, name+"/index", data)
            return out.String(), nil
        }
    }
    return base.JumpToUrl("未登录", "/admin/login")
}

func (a* adminController)Info(ctx context.Context) (rs string , err error){
    out := &bytes.Buffer{}
    data := make(map[string]interface{})
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,_ :=admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            data["AppName"]  = "higo-web"
            data["Version"]  = "1.0.1"
            data["UserName"] = param.Cookie.U
            data["CurrTime"] =  time.Now().Format("2006/1/2 15:04:05")
            data["OS"] = runtime.GOOS
            data["GOVersion"] = runtime.Version()
            data["Author"] = "Huangyp"
        }
    }
    view.NewView().Render(out, name + "/info",data)
    return out.String(), nil
}

