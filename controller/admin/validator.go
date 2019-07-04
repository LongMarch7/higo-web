package admin

import (
    "github.com/LongMarch7/higo/util/validator"
    "google.golang.org/grpc/grpclog"
    "net/url"
)


type UserLogin struct {
    UserName        string `form:"username" validate:"required"`
    Password        string `form:"password" validate:"required,ascii"`
}

func LoginParameterCheck(values url.Values) *UserLogin{
    var user UserLogin
    err := validator.Decoder.Decode(&user, values)
    if err != nil{
        grpclog.Error("LoginParameterCheck form decode error: ",err.Error())
        return nil
    }
    if err := validator.Validate.Struct(user); err != nil {
        grpclog.Error("LoginParameterCheck validator check error: ",err.Error())
        return nil
    }
    return &user
}