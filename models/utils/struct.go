package utils

import (
    "github.com/LongMarch7/higo-web/db/object/models"
    "github.com/LongMarch7/higo/db"
    "github.com/LongMarch7/higo/util/define"
    "github.com/bradfitz/gomemcache/memcache"
    "google.golang.org/grpc/grpclog"
)

type MenuPower struct {
    Pattern  string
    Method   string
    Status   bool
}


func UpdateRoleStatus(){
    var roleStatus = make(map[int]bool)
    engine :=db.NewDb(db.DefaultNAME).Engine()
    list := make([]models.MicroRole, 0)
    err :=engine.Table(models.MicroRole{}.TableName()).Find(&list)
    if err != nil {
        grpclog.Error("Load role list err: ",err.Error())
        return
    }
    for _,value := range list{
        if value.RoleStatus == 1{
            roleStatus[value.Id] = true
        }else{
            roleStatus[value.Id] = false
        }
    }

    for _,value := range list{
        if value.ParentId >0 {
            if stat,ok :=roleStatus[value.ParentId]; ok && stat == false{
                roleStatus[value.Id] = false
            }
        }
        buf :=define.ToByte(roleStatus[value.Id])
        db.MemcacheClient.Set(&memcache.Item{Key: define.RolePrefix + value.RoleName, Value: buf})
    }
}