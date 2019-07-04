package admin

import (
    "github.com/LongMarch7/higo-web/db/object/models"
    "github.com/LongMarch7/higo-web/models/utils"
    "github.com/LongMarch7/higo/auth"
    "github.com/LongMarch7/higo/db"
    "google.golang.org/grpc/grpclog"
    "strconv"
    "strings"
)


func GetMenuList(role_name string) []models.MicroAdminMenu{
    engine := db.NewDb(db.DefaultNAME).Engine()
    var menuList = make([]models.MicroAdminMenu, 0)
    engine.Table(models.MicroAdminMenu{}.TableName()).Find(&menuList)
    j := len(menuList)
    casbin := auth.NewCasbin()
    j = 0
    for _,value := range menuList{
        if casbin.Enforcer().Enforce(role_name, value.Url + ":" + value.Func, value.Method){
            menuList[j] = value
            j++
        }
    }
    return menuList[:j]
}

func GetMenuViewList(role_name string) []models.MicroAdminMenu{
    list :=GetMenuList(role_name)
    j :=0
    for _, value:= range list{
        if value.Type == 0 || value.Type == 1  {
            list[j] = value
            j++
        }
    }
    return list[:j]
}



func GetMenuPower(check []string, no_check string) (list []utils.MenuPower, success bool){
    engine := db.NewDb(db.DefaultNAME).Engine()
    success = false

    for _,value := range check{
        var menu = models.MicroAdminMenu{}
        id,err :=strconv.Atoi(value)
        if err != nil{
            grpclog.Error("Get menu power err:" , err.Error())
            return
        }
        _,err =engine.Table(models.MicroAdminMenu{}.TableName()).ID(value).Get(&menu)
        if err != nil{
            grpclog.Error("Get menu power err:" , err.Error())
            return
        }else if menu.Id != id {
            grpclog.Error("Get menu power err1:",menu)
            return
        }
        list = append(list, utils.MenuPower{Pattern:menu.Url + ":" + menu.Func,Method:menu.Method,Status:true})
    }
    noCheck := strings.Split(no_check,",")
    for _,value := range noCheck{
        var menu = models.MicroAdminMenu{}
        id,err :=strconv.Atoi(value)
        if err != nil{
            grpclog.Error("Get menu power err:" , err.Error())
            return
        }
        _,err =engine.Table(models.MicroAdminMenu{}.TableName()).ID(value).Get(&menu)
        if err != nil {
            grpclog.Error("get menu power err2:",err.Error())
            return
        }else if menu.Id != id {
            grpclog.Error("Get menu power err1:",menu)
            return
        }
        list = append(list, utils.MenuPower{Pattern:menu.Url + ":" + menu.Func,Method:menu.Method,Status:false})
    }
    success = true
    return
}
