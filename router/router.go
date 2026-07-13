package router

import (
	"iptv-spider-sh/middleware"
	"iptv-spider-sh/router/api"

	"github.com/kataras/iris/v12"
)

func InitRouters(app *iris.Application) {
	registerMacros(app)

	app.Use(middleware.AccessLog)

	app.Get("/", index)

	apiRouterGroup := app.Party("/api")
	{
		api.InitApiRouters(apiRouterGroup)
	}
}

func index(ctx iris.Context) {
	ctx.WriteString(`IPTV Spider API

频道列表接口:
  GET /api/m3u8     - 生成M3U8频道列表文件
    参数: udpxy, scheme, xteve, all, ref
  GET /api/tsM3u8   - 生成时移M3U8文件
    参数: ref

节目单接口:
  GET /api/epg      - 生成XMLTV节目单
    参数: daysAgo(默认1), ref

任务管理接口:
  GET /api/schedule - 获取定时任务调度列表
  GET /api/run      - 手动触发任务执行
    参数: task(clean-ch/clean-chi/clean-epg/clean/update-chi/update-epg/upload-m3u/upload-xmltv/upload-xmltv7)
`)
}
