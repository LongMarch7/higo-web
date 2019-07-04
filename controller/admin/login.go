package admin

import (
    "bytes"
    "github.com/LongMarch7/higo-web/models/admin"
    "github.com/LongMarch7/higo/controller/base"
    "github.com/LongMarch7/higo/util/define"
    "github.com/LongMarch7/higo/util/token"
    "github.com/LongMarch7/higo/view"
    "context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "strings"
)

func (a* adminController)Login(ctx context.Context) (rs string , err error){
    out := &bytes.Buffer{}
    data := make(map[string]interface{})
    if param := base.GetParamByCtx(ctx); param != nil{
       if ok,_ := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok{
           return base.JumpToUrl("已登录", "/admin")
       }
    }
    data["AppName"] = "higo-web"
    data["Version"] = "1.0.1"
    view.NewView().Render(out, name + "/login",data)
    return out.String(), nil
}

func (a* adminController)LoginPost(ctx context.Context) (rs string , err error){
    code := 0
    if param := base.GetParamByCtx(ctx); param != nil{
        if p :=LoginParameterCheck(param.PostFormParams); p !=nil {
            userName := p.UserName
            password := p.Password
            if param.Cookie.F == define.CookieFlagSAVE && strings.Compare(password,param.Cookie.S) == 0 {
            }else{
                password = token.NewTokenWithSalt1(password)
            }
            success,token := admin.AdminLogin(param.Cookie.T,userName,password)
            cookie := &param.Cookie
            if success {
                cookie.UpdateCookie(token, userName, password, 1)
            }else{
                cookie.UpdateCookie("", userName, "", 1)
            }
            header := metadata.Pairs(define.ResCookieName, cookie.Marshal())
            grpc.SetHeader(ctx, header)
            p = nil
            if success {
                return base.NewHtmlRet(code, "登录成功","/admin"), nil
            }
            code = -4
        }else{
            code = -3
        }

    }else{
        code = -2
    }
    return base.NewHtmlRet(code, "用户或者密码错误","/admin/login"), nil
}