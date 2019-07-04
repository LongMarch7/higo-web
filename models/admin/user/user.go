package user

import (
    "github.com/LongMarch7/higo-web/db/object/models"
    "github.com/LongMarch7/higo-web/models/utils"
    "github.com/LongMarch7/higo/auth"
    "github.com/LongMarch7/higo/db"
    "google.golang.org/grpc/grpclog"
    "strconv"
    "github.com/LongMarch7/higo/util/define"
    "time"
)

type UserRole struct {
    RoleName     string   `json:"role_name"`
    RoleId       int      `json:"role_id"`
    UserId       int64    `json:"user_id"`
}

func GetSysUserListWhereSql(where map[string]string) (string, []interface{}) {
    var sql = "1=1"
    var value []interface{}
    if v, ok := where["user_name"]; ok && v != "" {
        keywords := define.TrimString(v)
        sql += " AND user_login like ?"
        value = append(value,"%"+keywords+"%")
    }

    if v, ok := where["start_time"]; ok && v != "" {
        startTime := define.GetTimestamp(v)
        sql += " AND create_time >= ?"
        value = append(value,strconv.Itoa(int(startTime)))
    }

    if v, ok := where["end_time"]; ok && v != "" {
        endTime := define.GetTimestamp(v)
        sql += " AND create_time <= ?"
        value = append(value,strconv.Itoa(int(endTime)))
    }
    return sql, value
}

type userList struct{
    models.MicroUser
    Power        int `json:"power"`
}

func GetUserList(where map[string]string, page int, rows int, setter_role_id int)  (user_list []userList, count int){
    engine := db.NewDb(db.DefaultNAME).Engine()
    list := make([]models.MicroUser, 0)
    count = 0
    sql , value:= GetSysUserListWhereSql(where)

    err :=engine.Table(models.MicroUser{}.TableName()).Limit(rows,(page-1)*rows).Where(sql, value...).Find(&list)
    if err != nil {
        grpclog.Error("get user list err: ",err.Error())
        return
    }
    ret,err :=engine.Table(models.MicroUser{}.TableName()).Where(sql, value...).Count(&models.MicroUser{})
    if err != nil {
        grpclog.Error("get user count err: ",err.Error())
        return
    }
    for _,value := range list{
        userInfo :=userList{MicroUser:value}
        if ok,_ :=UserPowerCheck(setter_role_id, strconv.FormatInt(value.Id,10)); ok{
            userInfo.Power = 1
        }else{
            userInfo.Power = 0
        }
        user_list = append(user_list, userInfo)
    }
    count = int(ret)
    return
}

type AuthorizedRole struct{
    ID       int
    Name     string
    Status   int
}
func HasPower(user_id string, roles map[int]AuthorizedRole) (ret map[int]AuthorizedRole){
    engine := db.NewDb(db.DefaultNAME).Engine()
    roleUser := models.MicroRoleUser{}
    ret = roles
    has,err :=engine.Table(models.MicroRoleUser{}.TableName()).Where("user_id = ?", user_id).Get(&roleUser)
    if err != nil{
        return
    }else if !has{
        return
    }
    if value,ok := ret[roleUser.RoleId]; ok{
        ret[roleUser.RoleId] = AuthorizedRole{ID:value.ID,Name:value.Name,Status:1}
    }
    return
}
func GetAuthorizedRoleList(role_id int, role_name string)(roles map[int]AuthorizedRole){
    engine := db.NewDb(db.DefaultNAME).Engine()
    listRole := make([]models.MicroRole, 0)
    err :=engine.Table(models.MicroRole{}.TableName()).Find(&listRole)
    if err != nil {
        grpclog.Error("Load user list err: ",err.Error())
        return
    }
    roles = make(map[int]AuthorizedRole)
    roles[role_id] = AuthorizedRole{ID:role_id,Name:role_name}
    for _,value :=range listRole{
        if _,ok := roles[value.ParentId]; ok{
            roles[value.Id] = AuthorizedRole{ID:value.Id,Name:value.RoleName}
        }
    }
    return
}
func GetUserInfo(user_id string) models.MicroUser{
    engine := db.NewDb(db.DefaultNAME).Engine()
    userInfo := models.MicroUser{}

    has, err :=engine.Table(models.MicroUser{}.TableName()).ID(user_id).Get(&userInfo)
    if err != nil{
        grpclog.Error("Found user failed")
    }else if !has {
        grpclog.Error("Not found user")
    }
    return userInfo
}
func DelUsers(users_id []string, users_name []string, setter_role_id int) (success bool, msg string){
    engine := db.NewDb(db.DefaultNAME).Engine()
    session := engine.NewSession()
    defer session.Close()
    msg ="执行出错"
    success = false
    for index,value := range users_id{
        if ok,message :=UserPowerCheck(setter_role_id, value); !ok{
            success = ok
            msg = message
            return
        }
        err := session.Begin()
        user :=models.MicroUser{}
        if _,err =session.Table(models.MicroUser{}.TableName()).ID(value).Delete(&user); err != nil{
            session.Rollback()
            grpclog.Error("delete failed !",err.Error())
            msg = users_name[index] + "删除用户失败"
            return
        }
        roleUser :=models.MicroRoleUser{}
        if _,err =session.Table(models.MicroRoleUser{}.TableName()).Where("user_id = ?", value).Delete(&roleUser); err != nil{
            session.Rollback()
            grpclog.Error("delete failed !",err.Error())
            msg = users_name[index] + "删除用户失败"
            return
        }
        if err = session.Commit();err != nil{
            session.Rollback()
            grpclog.Error("delete  failed by session !",err.Error())
            return
        }
        db.MemcacheClient.Delete(define.UserPrefix + users_name[index])
    }
    msg = "操作成功"
    success = true
    return
}

func UsersStatusChange(users_id []string, users_status int, setter_role_id int) (success bool, msg string){
    engine := db.NewDb(db.DefaultNAME).Engine()
    user :=models.MicroUser{}
    success = false
    msg = "操作出错"
    for _,value := range users_id{
        if ok,message :=UserPowerCheck(setter_role_id, value); !ok{
            success = ok
            msg = message
            return
        }
        user.UserStatus = users_status
        if _,err :=engine.Table(models.MicroUser{}.TableName()).ID(value).Cols("user_status").Update(&user); err != nil{
            grpclog.Error("update user failed ", err.Error())
            return
        }
    }
    msg = "操作成功"
    success = true
    return
}

func UserPowerCheck(setter_role_id int, user_id string)  (success bool, msg string){
    engine := db.NewDb(db.DefaultNAME).Engine()
    success = false
    msg = "操作出错"
    id, err := strconv.ParseInt(user_id, 10, 64)
    if err != nil{
        grpclog.Error("string to int64 error")
        return
    }else if id == 1{
        msg = "super user"
        return
    }
    var user = UserRole{}
    ok, err := engine.Table(models.MicroUser{}.TableName()).Alias("u").Where("u.id = ?", id).
        Join("INNER",[]string{models.MicroRoleUser{}.TableName(),"ur"},"u.id = ur.user_id").
        Join("INNER",[]string{models.MicroRole{}.TableName(),"r"},"ur.role_id = r.id").
        Get(&user)
    if err != nil{
        msg = "查询失败"
        grpclog.Error("UserPowerCheck Exec sql failed",err.Error())
        return
    }else if !ok {
        msg = "用户不存在"
        grpclog.Error("UserPowerCheck Not found user")
        return
    }else if setter_role_id == user.RoleId{
        success = true
        return
    }
    var role_id = user.RoleId
    for{
        roleInfo := models.MicroRole{}
        has, err :=engine.Table(models.MicroRole{}.TableName()).ID(role_id).Get(&roleInfo)
        if err != nil{
            grpclog.Error("Found role failed")
            return
        }else if !has {
            grpclog.Error("Not found role")
            msg = "查询出错"
            return
        }

        if roleInfo.ParentId == setter_role_id{
            success = true
            return
        }
        if roleInfo.ParentId == 0{
            msg = "权限不足"
            return
        }
        role_id = roleInfo.ParentId
    }
}


func UserEdit(list []utils.MenuPower, user_status int, user_remark string, user_id string,user_name string, setter_user_name string, setter_role_id int) (success bool,msg string){
    casbin := auth.NewCasbin().Enforcer()
    engine := db.NewDb(db.DefaultNAME).Engine()
    session := engine.NewSession()
    defer session.Close()
    err := session.Begin()
    msg = "执行错误"
    success = false

    nowTime := int(time.Now().Unix())
    user :=models.MicroUser{UserStatus:user_status,UpdateTime:int(time.Now().Unix()),Signature:user_remark,UserNickname:user_name}
    if len(user_id) >0{
        if ok,message :=UserPowerCheck(setter_role_id, user_id); !ok{
            success = ok
            msg = message
            return
        }
        if  _,err = session.Table(models.MicroUser{}.TableName()).ID(user_id).Cols("user_status","update_time","remark","user_name").Update(&user); err != nil{
            session.Rollback()
            grpclog.Error("User update:",err.Error())
            msg = "更新出错"
            return
        }
    }else{
        has,err :=session.Table(models.MicroUser{}.TableName()).Where("user_name = ?", user_name).Exist(&models.MicroUser{})
        if err == nil && has{
            grpclog.Error("User name exits:")
            msg = "角色名存在"
            return
        }
        user.CreateTime = nowTime
        if  _,err = session.Table(models.MicroUser{}.TableName()).Insert(&user); err != nil{
            session.Rollback()
            grpclog.Error("User insert:",err.Error())
            return
        }
    }
    for _,value := range list{
        policy :=makePolicy("p",user_name,value.Pattern, value.Method)
        sql,where :=policyWhere(policy)
        if value.Status && casbin.Enforce(setter_user_name,value.Pattern, value.Method) {
            has,err :=session.Table(models.MicroCasbinRule{}.TableName()).Where(sql,where...).Exist(&models.MicroCasbinRule{})
            if err == nil && has{
                continue
            }
            if _,err = session.Table(models.MicroCasbinRule{}.TableName()).Insert(&policy); err != nil{
                session.Rollback()
                grpclog.Error("User insert policy err:",err.Error())
                return
            }
        }else if !value.Status {
            if _,err = session.Table(models.MicroCasbinRule{}.TableName()).Where(sql,where...).Delete(&policy); err != nil{
                session.Rollback()
                grpclog.Error("User delete policy err:",err.Error())
                return
            }
        }
    }

    if err = session.Commit();err != nil{
        session.Rollback()
        return
    }
    success = true
    return
}

func UpdateUserStatus(){
    //var userStatus = make(map[int]bool)
    engine :=db.NewDb(db.DefaultNAME).Engine()
    list := make([]models.MicroUser, 0)
    err :=engine.Table(models.MicroUser{}.TableName()).Find(&list)
    if err != nil {
        grpclog.Error("Load user list err: ",err.Error())
        return
    }
    //for _,value := range list{
    //    if value.UserStatus == 1{
    //        userStatus[value.Id] = true
    //    }else{
    //        userStatus[value.Id] = false
    //    }
    //}
    //
    //for _,value := range list{
    //    if value.ParentId >0 {
    //        if stat,ok :=userStatus[value.ParentId]; ok && stat == false{
    //            userStatus[value.Id] = false
    //        }
    //    }
    //    buf :=define.ToByte(userStatus[value.Id])
    //    db.MemcacheClient.Set(&memcache.Item{Key: value.UserName, Value: buf})
    //}
}