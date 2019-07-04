package user

import (
    "bytes"
    "context"
    "github.com/LongMarch7/higo-web/models/admin"
    "github.com/LongMarch7/higo-web/models/admin/user"
    "github.com/LongMarch7/higo/controller/base"
    "github.com/LongMarch7/higo/util/validator"
    "github.com/LongMarch7/higo/view"
    "golang.org/x/tools/go/ssa/interp/testdata/src/errors"
    "google.golang.org/grpc/grpclog"
    "strconv"
)

func (u* adminUserController)RoleIndex(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,_ := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            out := &bytes.Buffer{}
            view.NewView().Render(out, name + "/role_index", nil)
            return out.String(), nil
        }
    }
    out := &bytes.Buffer{}
    data := make(map[string]interface{})
    data["content"] = "网络繁忙"
    view.NewView().Render(out,"error", data)
    return out.String(), nil
}

func (u* adminUserController)RoleList(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,info := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            where := make(map[string]string)
            where["role_name"] 	= param.GetParams.Get("role_name")
            where["start_time"] = param.GetParams.Get("start_time")
            where["end_time"] 	= param.GetParams.Get("end_time")
            grpclog.Info(where)
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
            setterRoleId := info.RoleId
            if setterRoleId ==0 {
                grpclog.Error("get role id err")
                return base.NewLayuiRet(-3, "参数不合法",0,nil), nil
            }
            roleList, count:= user.GetRoleList(where, page, row, setterRoleId)
            return base.NewLayuiRet(0, "获取成功",count, roleList), nil
        }
    }
    return base.NewLayuiRet(-2, "参数错误",0,nil), nil
}
func (u* adminUserController)RoleDelete(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,info :=admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            id   := param.PostFormParams["id"]
            name := param.PostFormParams["name"]
            grpclog.Info(id)
            grpclog.Info(name)
            if err :=validator.Validate.Var(&name, validator.ArrayAlphanum); err != nil {
                grpclog.Error("role name 不合法",err.Error())
                return base.NewLayuiRet(-5, "参数错误",0,nil), nil
            }
            if err :=validator.Validate.Var(&id, validator.ArrayNumeric); err != nil {
                grpclog.Error("role id 不合法",err.Error())
                return base.NewLayuiRet(-4, "参数不合法-3",0,nil), nil
            }
            setterRoleId := info.RoleId
            if setterRoleId ==0 {
                grpclog.Error("get role id err")
                return base.NewLayuiRet(-3, "参数不合法",0,nil), nil
            }
            if succsess,msg :=user.DelRoles(id, name, setterRoleId); succsess{
                return base.NewLayuiRet(0, "删除成功", 0, nil), nil
            }else{
                return base.NewLayuiRet(-2, msg, 0, nil), nil
            }
        }
    }
    grpclog.Error("参数错误-2")
    return base.NewLayuiRet(-2, "参数错误",0,nil), nil
}

func (u* adminUserController)RoleStatusChange(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,info := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            id := param.PostFormParams["id"]
            grpclog.Info(id)
            if err :=validator.Validate.Var(&id, validator.ArrayNumeric); err != nil {
                grpclog.Error("RoleStatusChange id validate err：", err.Error())
                return base.NewLayuiRet(-6, "参数不合法",0,nil), nil
            }
            status := param.PostFormParams.Get("status")
            sta,err :=strconv.Atoi(status)
            if err != nil{
                grpclog.Error("RoleStatusChange status err：", err.Error())
                return base.NewLayuiRet(-5, "参数不合法",0,nil), nil
            }else if (sta <0) || (sta >1){
                grpclog.Error("RoleStatusChange status2 err")
                return base.NewLayuiRet(-4, "参数不合法",0,nil), nil
            }
            setterRoleId := info.RoleId
            if setterRoleId ==0 {
                grpclog.Error("get role id err")
                return base.NewLayuiRet(-3, "参数不合法",0,nil), nil
            }
            if success,msg :=user.RolesStatusChange(id, sta, setterRoleId); success{
                var message = "角色:"
                for _,value := range id{
                    message += "[" + value + "]"
                }
                if sta == 1 {
                    message += "已启用"
                }else if sta == 0{
                    message += "已禁用"
                }
                grpclog.Info(message)
                return base.NewLayuiRet(0, message, 0, nil), nil
            }else{
                return base.NewLayuiRet(-2, msg,0,nil), nil
            }

        }
    }
    grpclog.Error("RoleStatusChange parameter err")
    return base.NewLayuiRet(-2, "参数错误",0,nil), nil
}

func (u* adminUserController)RoleEdit(ctx context.Context) (rs string , err error){
    out := &bytes.Buffer{}
    data := make(map[string]interface{})
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,info := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            roleStatus := 0
            roleRemark := ""
            setterRoleName := info.RoleName
            setterRoleID := info.RoleId
            if len(setterRoleName) == 0  || setterRoleID == 0{
                grpclog.Error("get user info err")
                return "", errors.New("get user info err")
            }

            roleName := param.GetParams.Get("role_name")
            roleId := param.GetParams.Get("role_id")
            if len(roleId) > 0{
                if success,msg :=user.RolePowerCheck(setterRoleID, roleId); !success{
                    data["content"] = msg
                    view.NewView().Render(out,"error", data)
                    return out.String(), nil
                }
                userInfo := user.GetRoleInfo(roleId)
                roleStatus = userInfo.RoleStatus
                roleRemark = userInfo.Remark
            }
            data["CurrUserMenu"] = user.MenuPowerCheck(admin.GetMenuList(setterRoleName), setterRoleName, roleName)
            data["RoleStatus"] = roleStatus
            data["RoleName"] = roleName
            data["Intro"] = roleRemark
            data["Id"] = roleId
            view.NewView().Render(out, name+"/role_edit", data)
            return out.String(), nil
        }
    }
    data["content"] = "操作异常"
    view.NewView().Render(out,"error", data)
    return out.String(), nil
}

func (u* adminUserController)RoleEditPost(ctx context.Context) (rs string , err error){
    if param := base.GetParamByCtx(ctx); param != nil {
        if ok,info := admin.IsLogin(param.Cookie.T, param.Cookie.U); ok {
            grpclog.Info(param.PostFormParams)
            list,success := admin.GetMenuPower(param.PostFormParams["Power"], param.PostFormParams.Get("NoCheck"))
            if !success {
                grpclog.Error("参数错误-5")
                return base.NewLayuiRet(-5, "参数错误",0,nil), nil
            }
            roleStatus := param.PostFormParams.Get("RoleStatus")
            stat :=0
            if len(roleStatus) > 0{
                stat = 1
            }
            roleName := param.PostFormParams.Get("RoleName")
            if err :=validator.Validate.Var(roleName, validator.Alphanum); err != nil {
                grpclog.Error("参数错误-4",err.Error())
                return base.NewLayuiRet(-4, "参数错误",0,nil), nil
            }
            remark :=param.PostFormParams.Get("Intro")
            id  :=param.PostFormParams.Get("Id")
            setterRoleName := info.RoleName
            setterRoleID := info.RoleId
            if len(setterRoleName) ==0 || setterRoleID == 0{
                grpclog.Error("参数错误-3")
                return base.NewLayuiRet(-3, "参数错误",0,nil), nil
            }
            if success, msg :=user.RoleEdit(list, stat, remark,id , roleName, setterRoleName, setterRoleID); success{
                return base.NewLayuiRet(0, "提交成功",0,nil), nil
            }else{
                return base.NewLayuiRet(-2, msg,0,nil), nil
            }
        }
    }
    grpclog.Error("参数错误-2")
    return base.NewLayuiRet(-2, "参数错误",0,nil), nil
}