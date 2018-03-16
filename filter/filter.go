package filter

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	. "wksw/notify/models/error"
	. "wksw/notify/models/keystone"
)

func init() {
	var FilterAppId = func(ctx *context.Context) {
		appid := ctx.Input.Header("AppId")
		if appid == "" {
			beego.Error("missing appid in header")
			ctx.Output.JSON(MissingAppIdErr, true, false)
			return
		}
	}

	var FilterToken = func(ctx *context.Context) {
		var ersp ErrResponse
		auth_token := ctx.Input.Header("X-Auth-Token")
		if auth_token == "" {
			beego.Error("missing X-Auth-Token in header")
			ersp.ErrCode = 1001
			ersp.ErrMsg = "missing X-Auth-Token in header"
			ctx.Output.JSON(ersp, true, false)

		}
		err := ValidToken(auth_token)
		if err != nil {
			beego.Error(err)
			ersp.ErrCode = 1002
			ersp.ErrMsg = "forbidden"
			ctx.Output.JSON(ersp, true, false)
			return
		}
	}

	beego.InsertFilter("/*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "AppId", "Secret", "X-Auth-Token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.InsertFilter("/miniprogram/*", beego.BeforeRouter, FilterAppId)
	beego.InsertFilter("/publicnum/*", beego.BeforeRouter, FilterAppId)

	beego.InsertFilter("/miniprogram/*", beego.BeforeRouter, FilterToken)
	beego.InsertFilter("/publicnum/*", beego.BeforeRouter, FilterToken)
	beego.InsertFilter("/custometoken/*", beego.BeforeRouter, FilterToken)
	beego.InsertFilter("/project/*", beego.BeforeRouter, FilterToken)
	beego.InsertFilter("/template/*", beego.BeforeRouter, FilterToken)
	beego.InsertFilter("/clusterinfo/*", beego.BeforeRouter, FilterToken)
}
