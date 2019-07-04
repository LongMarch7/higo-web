package admin

import (
    "encoding/json"
    "github.com/LongMarch7/higo-web/db/object/models"
    "github.com/LongMarch7/higo/db"
    "github.com/LongMarch7/higo/util/define"
    "github.com/LongMarch7/higo/util/token"
    "github.com/bradfitz/gomemcache/memcache"
    "google.golang.org/grpc/grpclog"
    "strings"
)
type UserRole struct {
    RoleName     string   `json:"role_name"`
    RoleId       int      `json:"role_id"`
    UserLogin    string   `json:"user_login"`
    UserPass     string   `json:"user_pass"`
    UserNickname string   `json:"user_nickname"`
    UserStatus   int      `json:"user_status"`
    UserId       int64    `json:"user_id"`
}

//func GetUserInfo(user_token string, user_name string) UserRole{
//   key := user_token + user_name
//   var userInfo = UserRole{}
//   if it, err := db.MemcacheClient.Get(user_token); err == nil && it.Key == user_token{
//       name := string(it.Value)
//       if strings.Compare(user_name,name) == 0 {
//           if it, err := db.MemcacheClient.Get(name); err == nil && it.Key == name{
//               json.Unmarshal(it.Value,&userInfo)
//           }
//           if len(userInfo.RoleName) >0 && global.RolePowerCheck(userInfo.RoleName){
//               return true
//           }
//       }
//   }
//   return userInfo
//}
//func GetUserInfo(user_name string) UserRole{
//    engine := db.NewDb(db.DefaultNAME).Engine()
//    var user = UserRole{}
//    ok, err := engine.Table(models.MicroUser{}.TableName()).Alias("u").Where("user_login = ?",user_name).
//        Join("INNER",[]string{models.MicroRoleUser{}.TableName(),"ur"},"u.id = ur.user_id").
//        Join("INNER",[]string{models.MicroRole{}.TableName(),"r"},"ur.role_id = r.id").
//        Get(&user)
//    if err != nil{
//        grpclog.Error("Exec sql failed")
//    }else if !ok {
//        grpclog.Error("Not found user")
//    }
//    return user
//}
func CheckRoleStatus(role_name string) bool{
    var stat = false
    key  :=define.RolePrefix + role_name
    if it, err :=  db.MemcacheClient.Get(key); err == nil && it.Key == key{
        define.FromByte(it.Value,&stat)
    }
    return stat
}
func IsLogin(user_token string, user_name string) (is_login bool,user_info UserRole){
    user_info = UserRole{}
    is_login = false
    if it, err := db.MemcacheClient.Get(user_token); err == nil && it.Key == user_token{
        name := string(it.Value)
        
        if strings.Compare(user_name,name) == 0 {
            key :=define.UserPrefix + user_name
            if it, err := db.MemcacheClient.Get(key); err == nil && it.Key == key{
                json.Unmarshal(it.Value,&user_info)
            }else{
                engine := db.NewDb(db.DefaultNAME).Engine()
                ok, err := engine.Table(models.MicroUser{}.TableName()).Alias("u").Where("user_login = ?",user_name).
                    Join("INNER",[]string{models.MicroRoleUser{}.TableName(),"ur"},"u.id = ur.user_id").
                    Join("INNER",[]string{models.MicroRole{}.TableName(),"r"},"ur.role_id = r.id").
                    Get(&user_info)
                if err != nil{
                    grpclog.Error("AdminLogin Exec sql failed")
                    return
                }else if !ok {
                    grpclog.Error("AdminLogin Not found user")
                    return
                }
            }
            if user_info.UserStatus == 1 && len(user_info.RoleName) >0 && CheckRoleStatus(user_info.RoleName){
                is_login = true
                return
            }
        }
    }
    return
}

func AdminLogin(user_token string, user_name string, pw string) (succecss bool, t string){
    engine := db.NewDb(db.DefaultNAME).Engine()
    succecss = false
    t = user_token
    if ok,_ := IsLogin(user_token, user_name ); ok {
        succecss = true
        return
    }

    t = token.NewTokenWithTime(user_name)
    var user = UserRole{}
    ok, err := engine.Table(models.MicroUser{}.TableName()).Alias("u").Where("user_login = ? AND user_pass = ?",user_name,token.NewTokenWithSalt2(pw)).
        Join("INNER",[]string{models.MicroRoleUser{}.TableName(),"ur"},"u.id = ur.user_id").
        Join("INNER",[]string{models.MicroRole{}.TableName(),"r"},"ur.role_id = r.id").
        Get(&user)
    if err != nil{
        grpclog.Error("AdminLogin Exec sql failed")
        return
    }else if !ok {
        grpclog.Error("AdminLogin Not found user")
        return
    }

    if user.UserStatus != 1 {
        grpclog.Error("User is Banned")
        return
    }else if len(user.RoleName)==0 || !CheckRoleStatus(user.RoleName){
        grpclog.Error("Role is Banned")
        return
    }
    db.MemcacheClient.Set(&memcache.Item{Key: t, Value: []byte(user_name)})
    if value,err :=json.Marshal(user); err ==nil{
        db.MemcacheClient.Set(&memcache.Item{Key: define.UserPrefix + user_name, Value: value})
    }else{
        grpclog.Info( err.Error() )
    }
    succecss = true
    return
}