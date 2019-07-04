package user

import (
    "bytes"
    "github.com/LongMarch7/higo-web/db/object/models"
    "github.com/LongMarch7/higo-web/models/admin"
    "github.com/LongMarch7/higo-web/models/admin/user"
    "github.com/LongMarch7/higo/controller/base"
    "github.com/LongMarch7/higo/db"
    "github.com/LongMarch7/higo/util/define"
    "github.com/LongMarch7/higo/util/validator"
    "github.com/LongMarch7/higo/view"
    "google.golang.org/grpc/grpclog"
    "strconv"
    "context"
    "errors"
)

type adminUserController struct {
}
var Controller = &adminUserController{}
var name = "admin/user"
func init(){
    base.AddController(name, Controller)
}

func (u* adminUserController)UserIndex(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,_ := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            out := &bytes.Buffer{}
            view.NewView().Render(out, name + "/user_index", nil)
            return out.String(), nil
        }
    }
    out := &bytes.Buffer{}
    data := make(map[string]interface{})
    data["content"] = "网络繁忙"
    view.NewView().Render(out,"error", data)
    return out.String(), nil
}

func (u* adminUserController)UserList(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,info := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            if info.UserStatus != 1 {
                grpclog.Error("get user id err")
                return base.NewLayuiRet(-4, "参数不合法",0,nil), nil
            }
            where := make(map[string]string)
            where["user_name"] 	= param.GetParams.Get("user_name")
            where["start_time"] = param.GetParams.Get("start_time")
            where["end_time"] 	= param.GetParams.Get("end_time")
            page_num := param.GetParams.Get("page")
            page,err :=strconv.Atoi(page_num)
            if err != nil || page <=0{
                page = 1
            }
            row_num := param.GetParams.Get("limit")
            row,err :=strconv.Atoi(row_num)
            if err != nil || row <=0{
                row = 10
            }
            setterRoleID:= info.RoleId
            if setterRoleID == 0 {
                grpclog.Error("get role id err")
                return base.NewLayuiRet(-3, "参数不合法",0,nil), nil
            }
            userList, count:= user.GetUserList(where, page, row, setterRoleID)
            return base.NewLayuiRet(0, "获取成功",count, userList), nil
        }
    }
    return base.NewLayuiRet(-2, "参数错误",0,nil), nil
}
func (u* adminUserController)UserDelete(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,info :=admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            id   := param.PostFormParams["id"]
            name := param.PostFormParams["name"]
            grpclog.Info(id)
            grpclog.Info(name)
            if err :=validator.Validate.Var(&name, validator.ArrayAlphanum); err != nil {
                grpclog.Error("user name 不合法",err.Error())
                return base.NewLayuiRet(-5, "参数错误",0,nil), nil
            }
            if err :=validator.Validate.Var(&id, validator.ArrayNumeric); err != nil {
                grpclog.Error("user id 不合法",err.Error())
                return base.NewLayuiRet(-4, "参数不合法-3",0,nil), nil
            }
            setterRoleID := info.RoleId
            if setterRoleID == 0{
                grpclog.Error("get user id err")
                return base.NewLayuiRet(-3, "参数不合法",0,nil), nil
            }
            if succsess,msg :=user.DelUsers(id, name, setterRoleID); succsess{
                return base.NewLayuiRet(0, "删除成功", 0, nil), nil
            }else{
                return base.NewLayuiRet(-2, msg, 0, nil), nil
            }
        }
    }
    grpclog.Error("参数错误-2")
    return base.NewLayuiRet(-2, "参数错误",0,nil), nil
}

func (u* adminUserController)UserStatusChange(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,info := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            id := param.PostFormParams["id"]
            if err :=validator.Validate.Var(&id, validator.ArrayNumeric); err != nil {
                grpclog.Error("UserStatusChange id validate err：", err.Error())
                return base.NewLayuiRet(-6, "参数不合法",0,nil), nil
            }
            status := param.PostFormParams.Get("status")
            sta,err :=strconv.Atoi(status)
            if err != nil{
                grpclog.Error("UserStatusChange status err：", err.Error())
                return base.NewLayuiRet(-5, "参数不合法",0,nil), nil
            }else if (sta <0) || (sta >1){
                grpclog.Error("UserStatusChange status2 err")
                return base.NewLayuiRet(-4, "参数不合法",0,nil), nil
            }
            setterRoleID := info.RoleId
            if setterRoleID==0 {
                grpclog.Error("get role id err")
                return base.NewLayuiRet(-3, "参数不合法",0,nil), nil
            }
            if success,msg :=user.UsersStatusChange(id, sta, setterRoleID); success{
                var message = "用户:"
                for _,value := range id{
                    message += "[" + value + "]"
                }
                if sta == 1 {
                    message += "已启用"
                }else if sta == 0{
                    message += "已禁用"
                }
                db.MemcacheClient.Delete(define.UserPrefix + param.Cookie.U)
                grpclog.Info(message)
                return base.NewLayuiRet(0, message, 0, nil), nil
            }else{
                return base.NewLayuiRet(-2, msg,0,nil), nil
            }

        }
    }
    grpclog.Error("UserStatusChange parameter err")
    return base.NewLayuiRet(-2, "参数错误",0,nil), nil
}

func (u* adminUserController)UserEdit(ctx context.Context) (rs string , err error){
    out := &bytes.Buffer{}
    data := make(map[string]interface{})
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,info := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            setterRoleID := info.RoleId
            setterRoleName := info.RoleName
            if setterRoleID == 0 || len(setterRoleName) == 0{
                grpclog.Error("get user info err")
                return "", errors.New("get user info err")
            }
            userId := param.GetParams.Get("user_id")
            var userInfo models.MicroUser
            rolesMap := user.GetAuthorizedRoleList(setterRoleID, setterRoleName)
            if len(userId) > 0{
                if success,msg :=user.UserPowerCheck(setterRoleID, userId); !success{
                    data["content"] = msg
                    view.NewView().Render(out,"error", data)
                    return out.String(), nil
                }
                userInfo = user.GetUserInfo(userId)
                rolesMap = user.HasPower(userId, rolesMap)
            }
            //data["CurrUserMenu"] = user.MenuPowerCheck(admin.GetMenuList(setteruserName), setteruserName, userName)
            data["UserInfo"] = userInfo
            data["RoleMap"] = rolesMap
            view.NewView().Render(out, name+"/user_edit", data)
            return out.String(), nil
        }
    }
    data["content"] = "操作异常"
    view.NewView().Render(out,"error", data)
    return out.String(), nil
}

func (u* adminUserController)UserEditPost(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,_ := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            //grpclog.Info(param.PostFormParams)
            //userStatus := param.PostFormParams.Get("UserStatus")
            //stat :=0
            //if len(userStatus) > 0{
            //    stat = 1
            //}
            //userName := param.PostFormParams.Get("UserName")
            //if err :=validator.Validate.Var(userName, validator.Alphanum); err != nil {
            //    grpclog.Error("参数错误-4",err.Error())
            //    return base.NewLayuiRet(-4, "参数错误",0,nil), nil
            //}
            //remark :=param.PostFormParams.Get("Intro")
            //id  :=param.PostFormParams.Get("Id")
            //setterUserName := info.UserNickname
            //setterRoleID := info.RoleId
            //if len(setterUserName) ==0 || setterRoleID == 0{
            //    grpclog.Error("参数错误-3")
            //    return base.NewLayuiRet(-3, "参数错误",0,nil), nil
            //}
            //if success, msg :=user.UserEdit(list, stat, remark,id , userName, setterUserName, setterRoleID); success{
            //    return base.NewLayuiRet(0, "提交成功",0,nil), nil
            //}else{
            //    return base.NewLayuiRet(-2, msg,0,nil), nil
            //}
        }
    }
    grpclog.Error("参数错误-2")
    return base.NewLayuiRet(-2, "参数错误",0,nil), nil
}