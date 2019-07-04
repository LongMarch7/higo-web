package user

import (
    "github.com/LongMarch7/higo-web/db/object/models"
    "github.com/LongMarch7/higo-web/models/utils"
    "github.com/LongMarch7/higo/auth"
    "github.com/LongMarch7/higo/db"
    "github.com/bradfitz/gomemcache/memcache"
    "google.golang.org/grpc/grpclog"
    "strconv"
    "github.com/LongMarch7/higo/util/define"
    "time"
)

func GetSysRoleListWhereSql(where map[string]string) (string, []interface{}) {
    var sql = "1=1"
    var value []interface{}
    if v, ok := where["role_name"]; ok && v != "" {
        keywords := define.TrimString(v)
        sql += " AND role_name like ?"
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

type roleList struct{
    models.MicroRole
    Power   int `json:"power"`
}

func GetRoleList(where map[string]string, page int, rows int,setter_role_id int)  (role_list []roleList, count int){
    engine := db.NewDb(db.DefaultNAME).Engine()
    list := make([]models.MicroRole, 0)
    count = 0
    sql , value:= GetSysRoleListWhereSql(where)

    err :=engine.Table(models.MicroRole{}.TableName()).Limit(rows,(page-1)*rows).Where(sql, value...).Find(&list)
    if err != nil {
        grpclog.Error("get role list err: ",err.Error())
        return
    }
    ret,err :=engine.Table(models.MicroRole{}.TableName()).Where(sql, value...).Count(&models.MicroRole{})
    if err != nil {
        grpclog.Error("get role count err: ",err.Error())
        return
    }
    for _,value := range list{
        if ok,_ :=RolePowerCheck(setter_role_id, strconv.Itoa(value.Id)); ok{
            role_list = append(role_list,roleList{MicroRole:value,Power: 1})
        }else{
            role_list = append(role_list,roleList{MicroRole:value,Power: 0})
        }
    }
    count = int(ret)
    return
}

func GetRoleInfo(role_id string) models.MicroRole{
    engine := db.NewDb(db.DefaultNAME).Engine()
    roleInfo := models.MicroRole{}

    has, err :=engine.Table(models.MicroRole{}.TableName()).ID(role_id).Get(&roleInfo)
    if err != nil{
        grpclog.Error("Found role failed")
    }else if !has {
        grpclog.Error("Not found user")
    }
    return roleInfo
}
func DelRoles(roles_id []string, roles_name []string, setter_role_id int) (success bool, msg string){
    engine := db.NewDb(db.DefaultNAME).Engine()
    session := engine.NewSession()
    defer session.Close()
    role :=models.MicroRole{}
    msg ="执行出错"
    success = false
    for index,value := range roles_id{
        if ok,message :=RolePowerCheck(setter_role_id, value); !ok{
            success = ok
            msg = message
            return
        }
        err := session.Begin()
        has,err :=session.Table(models.MicroRole{}.TableName()).Where("parent_id = ?",value).Exist(&models.MicroRole{})
        if err == nil && has{
            grpclog.Error("delete failed ,because role has child !")
            msg = roles_name[index] + "角色删除失败，存在子角色依赖"
            return
        }
        if _,err =session.Table(models.MicroRole{}.TableName()).ID(value).Delete(&role); err != nil{
            session.Rollback()
            grpclog.Error("delete failed !",err.Error())
            msg = roles_name[index] + "删除角色失败"
            return
        }
        if _,err = session.Table(models.MicroCasbinRule{}.TableName()).Where("v0 = ?",roles_name[index]).Delete(&models.MicroCasbinRule{}); err != nil{
            session.Rollback()
            grpclog.Error("Role delete policy err:",err.Error())
            msg = roles_name[index] + "角色规则删除失败"
            return
        }
        if err = session.Commit();err != nil{
            session.Rollback()
            grpclog.Error("delete  failed by session !",err.Error())
            return
        }
        db.MemcacheClient.Delete(define.RolePrefix + roles_name[index])
    }
    msg = "操作成功"
    success = true
    return
}

func RolesStatusChange(roles_id []string, roles_status int, setter_role_id int) (success bool, msg string){
    engine := db.NewDb(db.DefaultNAME).Engine()
    role :=models.MicroRole{}
    success = false
    msg = "操作出错"
    for _,value := range roles_id{
        if ok,message :=RolePowerCheck(setter_role_id, value); !ok{
            success = ok
            msg = message
            return
        }
        role.RoleStatus = roles_status
        if _,err :=engine.Table(models.MicroRole{}.TableName()).ID(value).Cols("role_status").Update(&role); err != nil{
            return
        }
    }
    msg = "操作成功"
    success = true
    return
}

func RolePowerCheck(setter_role_id int, role_id string)  (success bool, msg string){
    engine := db.NewDb(db.DefaultNAME).Engine()
    success = false
    msg = "操作出错"
    id,err := strconv.Atoi(role_id)
    if err != nil {
        return
    }
    if setter_role_id == id{
        msg = "无法对自身角色操作"
        return
    }
    for{
        roleInfo := models.MicroRole{}
        has, err :=engine.Table(models.MicroRole{}.TableName()).ID(id).Get(&roleInfo)
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
        id = roleInfo.ParentId
    }
}

func MenuPowerCheck(list []models.MicroAdminMenu, setter_role_name string, role_name string) []models.MicroAdminMenu{
    casbin := auth.NewCasbin().Enforcer()
    j := 0
    for _,value :=range list{
        if casbin.Enforce(setter_role_name,value.Url + ":" + value.Func, value.Method) {
            value.Status = 0
            if len(role_name)>0 && casbin.Enforce(role_name,value.Url + ":" + value.Func, value.Method){
                value.Status = 1
            }
            list[j] = value
            j++
        }
    }
    return list[:j]
}

func makePolicy(list ...string) models.MicroCasbinRule{
    line := models.MicroCasbinRule{}

    if len(list) > 0 {
        line.PType = list[0]
    }
    if len(list) > 1 {
        line.V0 = list[1]
    }
    if len(list) > 2 {
        line.V1 = list[2]
    }
    if len(list) > 3 {
        line.V2 = list[3]
    }
    if len(list) > 4 {
        line.V3 = list[4]
    }
    if len(list) > 5 {
        line.V4 = list[5]
    }
    if len(list) > 6 {
        line.V5 = list[6]
    }
    return line
}
func policyWhere(line models.MicroCasbinRule) (string, []interface{}) {
    queryArgs := []interface{}{line.PType}
    queryStr := "p_type = ?"
    if line.V0 != "" {
        queryStr += " and v0 = ?"
        queryArgs = append(queryArgs, line.V0)
    }
    if line.V1 != "" {
        queryStr += " and v1 = ?"
        queryArgs = append(queryArgs, line.V1)
    }
    if line.V2 != "" {
        queryStr += " and v2 = ?"
        queryArgs = append(queryArgs, line.V2)
    }
    if line.V3 != "" {
        queryStr += " and v3 = ?"
        queryArgs = append(queryArgs, line.V3)
    }
    if line.V4 != "" {
        queryStr += " and v4 = ?"
        queryArgs = append(queryArgs, line.V4)
    }
    if line.V5 != "" {
        queryStr += " and v5 = ?"
        queryArgs = append(queryArgs, line.V5)
    }
    return queryStr,queryArgs
}
func RoleEdit(list []utils.MenuPower, role_status int, role_remark string, role_id string,role_name string, setter_role_name string, setter_role_id int) (success bool,msg string){
    casbin := auth.NewCasbin().Enforcer()
    engine := db.NewDb(db.DefaultNAME).Engine()
    session := engine.NewSession()
    defer session.Close()
    err := session.Begin()
    msg = "执行错误"
    success = false

    nowTime := int(time.Now().Unix())
    role :=models.MicroRole{RoleStatus:role_status,UpdateTime:int(time.Now().Unix()),Remark:role_remark,RoleName:role_name}
    if len(role_id) >0{
        if ok,message :=RolePowerCheck(setter_role_id, role_id); !ok{
            success = ok
            msg = message
            return
        }
        if  _,err = session.Table(models.MicroRole{}.TableName()).ID(role_id).Cols("role_status","update_time","remark","role_name").Update(&role); err != nil{
            session.Rollback()
            grpclog.Error("Role update:",err.Error())
            msg = "更新出错"
            return
        }
    }else{
        has,err :=session.Table(models.MicroRole{}.TableName()).Where("role_name = ?", role_name).Exist(&models.MicroRole{})
        if err == nil && has{
            grpclog.Error("Role name exits:")
            msg = "角色名存在"
            return
        }
        role.CreateTime = nowTime
        role.ParentId = setter_role_id
        if  _,err = session.Table(models.MicroRole{}.TableName()).Insert(&role); err != nil{
            session.Rollback()
            grpclog.Error("Role insert:",err.Error())
            return
        }
    }
    for _,value := range list{
        policy :=makePolicy("p",role_name,value.Pattern, value.Method)
        sql,where :=policyWhere(policy)
        if value.Status && casbin.Enforce(setter_role_name,value.Pattern, value.Method) {
            has,err :=session.Table(models.MicroCasbinRule{}.TableName()).Where(sql,where...).Exist(&models.MicroCasbinRule{})
            if err == nil && has{
                continue
            }
            if _,err = session.Table(models.MicroCasbinRule{}.TableName()).Insert(&policy); err != nil{
                session.Rollback()
                grpclog.Error("Role insert policy err:",err.Error())
                return
            }
        }else if !value.Status {
            if _,err = session.Table(models.MicroCasbinRule{}.TableName()).Where(sql,where...).Delete(&policy); err != nil{
                session.Rollback()
                grpclog.Error("Role delete policy err:",err.Error())
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