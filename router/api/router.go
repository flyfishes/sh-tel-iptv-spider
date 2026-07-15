package api

import (
	"iptv-spider-sh/global"
	"iptv-spider-sh/modules/auth"
	"iptv-spider-sh/utils"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
)

func InitApiRouters(rg iris.Party) {
	rg.Get("/schedule", schedule)
	rg.Get("/run", func(ctx iris.Context) {
		taskName := ctx.FormValue("task")

		go func() {
			global.ConcurrencyControl.Do("", func() (interface{}, error) {
				switch taskName {
				case "clean-ch":
					auth.CleanChannelData()
				case "clean-chi":
					auth.CleanChannelInfoData()
				case "clean-epg":
					auth.CleanEPGDetailsData()
				case "clean":
					auth.CleanChannelData()
					auth.CleanChannelInfoData()
					auth.CleanEPGDetailsData()
				case "update-chi":
					auth.GetGlobalClient().FetchChannelList()
				case "update-epg":
					auth.GetGlobalClient().FetchChannelProg(false)
				case "upload-m3u":
					auth.GenerateAndUploadM3u()
				case "upload-xmltv":
					auth.GenerateAndUploadXmlTv()
				case "upload-xmltv7":
					auth.GenerateAndUploadXmlTvDays7()
				case "upload-epgjson":
					auth.GenerateAndUploadEpgJson()
				case "upload-epgjson7":
					auth.GenerateAndUploadEpgJsonDays7()
				}
				return nil, nil
			})
		}()
		ctx.WriteString("OK")
	})

	rg.Get("/m3u8", generateM3u8)

	rg.Get("/tsM3u8", generateTsM3u8)
	rg.Get("/diyp", generateDiypTxt)

	rg.Get("/epg", generateXmlTv)
	rg.Get("/epgjson", generateEpgJson)

}

func generateEpgJson(ctx iris.Context) {
	days := ctx.FormValue("days")
	daysAgo, err := strconv.Atoi(days)
	if err != nil {
		daysAgo = 1 // default to 1 day
	}

	ref := ctx.FormValue("ref")
	reqMD5Key := utils.CalcMD5KeyForRequest("generateEpgJson", days)

	// 缓存机制
	if ref != "true" && global.CACHE.IsExist(reqMD5Key) {
		ctx.ContentType("application/json")
		ctx.Binary(global.CACHE.Get(reqMD5Key).([]byte))
		return
	}

	// 并发时合并请求
	resp, _, _ := global.ConcurrencyControl.Do(reqMD5Key, func() (interface{}, error) {
		epgBytes, err := auth.GenerateEpgJson(daysAgo)
		if err != nil {
			global.LOG.Error("生成EPG JSON失败: " + err.Error())
			return nil, err
		}
		timeOut := time.Duration(global.CONFIG.Cache.DefTimeOut)
		global.CACHE.Put(reqMD5Key, epgBytes, time.Minute*timeOut)
		SaveToLogDir(epgBytes, "epg.json")
		return epgBytes, nil
	})

	if resp == nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	ctx.ContentType("application/json")
	ctx.Binary(resp.([]byte))
}

func schedule(ctx iris.Context) {
	type s struct {
		ID       int
		PreTime  time.Time
		NextTime time.Time
	}
	var schedule []s
	for _, entry := range global.CRON.Entries() {
		schedule = append(schedule, s{
			ID:       int(entry.ID),
			PreTime:  entry.Prev,
			NextTime: entry.Next,
		})
	}
	ctx.JSON(schedule)
}

// 生成m3u8文件 节目去重
func generateM3u8(ctx iris.Context) {
	// 获取query参数
	udpxy := ctx.FormValue("udpxy")
	scheme := ctx.FormValue("scheme")
	xteve := ctx.FormValue("xteve")
	all := ctx.FormValue("all")
	ref := ctx.FormValue("ref")

	var bufStr string
	if xteve == "true" {
		bufStr = "xteve"
	} else if udpxy != "" {
		bufStr = udpxy
	} else if scheme != "" {
		bufStr = scheme
	}
	if all == "true" {
		bufStr += all
	}
	if ctx == nil {
		respBytes := auth.GenerateM3u8(udpxy, scheme, xteve, all)
		SaveToLogDir(respBytes, "iptv.m3u")
		return
	}
	reqMD5Key := utils.CalcMD5KeyForRequest("generateM3u8", bufStr)
	// 缓存机制
	if ref != "true" && global.CACHE.IsExist(reqMD5Key) {
		ctx.Header("Content-Disposition", "attachment; filename=iptv.m3u")
		ctx.Binary(global.CACHE.Get(reqMD5Key).([]byte))
		return
	}
	// 并发时合并请求
	resp, _, _ := global.ConcurrencyControl.Do(reqMD5Key, func() (interface{}, error) {
		respBytes := auth.GenerateM3u8(udpxy, scheme, xteve, all)
		timeOut := time.Duration(global.CONFIG.Cache.DefTimeOut)
		global.CACHE.Put(reqMD5Key, respBytes, time.Minute*timeOut)
		return respBytes, nil
	})
	ctx.Header("Content-Disposition", "attachment; filename=iptv.m3u")
	ctx.Binary(resp.([]byte))
}

func generateTsM3u8(ctx iris.Context) {
	ref := ctx.FormValue("ref")
	udpxy := ctx.FormValue("udpxy")
	scheme := ctx.FormValue("scheme")
	xteve := ctx.FormValue("xteve")
	all := ctx.FormValue("all")

	var bufStr string
	if xteve == "true" {
		bufStr = "xteve"
	} else if udpxy != "" {
		bufStr = udpxy
	} else if scheme != "" {
		bufStr = scheme
	}
	if all == "true" {
		bufStr += all
	}
	reqMD5Key := utils.CalcMD5KeyForRequest("generateTsM3u8", bufStr)
	// 缓存机制
	if ref != "true" && global.CACHE.IsExist(reqMD5Key) {
		ctx.Header("Content-Disposition", "attachment; filename=iptv-ts.m3u")
		ctx.Binary(global.CACHE.Get(reqMD5Key).([]byte))
		return
	}
	// 并发时合并请求
	resp, _, _ := global.ConcurrencyControl.Do(reqMD5Key, func() (interface{}, error) {
		respBytes := auth.GenerateTimeShiftM3u8(udpxy, scheme, xteve, all)
		timeOut := time.Duration(global.CONFIG.Cache.DefTimeOut)
		global.CACHE.Put(reqMD5Key, respBytes, time.Minute*timeOut)
		SaveToLogDir(respBytes, "iptv-ts.m3u")
		return respBytes, nil
	})
	ctx.Header("Content-Disposition", "attachment; filename=iptv-ts.m3u")
	ctx.Binary(resp.([]byte))
}

func generateDiypTxt(ctx iris.Context) {
	ref := ctx.FormValue("ref")
	udpxy := ctx.FormValue("udpxy")
	scheme := ctx.FormValue("scheme")
	xteve := ctx.FormValue("xteve")
	all := ctx.FormValue("all")

	var bufStr string
	if xteve == "true" {
		bufStr = "xteve"
	} else if udpxy != "" {
		bufStr = udpxy
	} else if scheme != "" {
		bufStr = scheme
	}
	if all == "true" {
		bufStr += all
	}

	if ctx == nil {
		respBytes := auth.GenerateDiyp(udpxy, scheme, xteve, all)
		SaveToLogDir(respBytes, "iptv-diyp.txt")
		return
	}
	reqMD5Key := utils.CalcMD5KeyForRequest("generateTsM3u8", bufStr)
	// 缓存机制
	if ref != "true" && global.CACHE.IsExist(reqMD5Key) {
		ctx.Header("Content-Disposition", "attachment; filename=iptv-diyp.txt")
		ctx.Binary(global.CACHE.Get(reqMD5Key).([]byte))
		return
	}
	// 并发时合并请求
	resp, _, _ := global.ConcurrencyControl.Do(reqMD5Key, func() (interface{}, error) {
		respBytes := auth.GenerateDiyp(udpxy, scheme, xteve, all)
		timeOut := time.Duration(global.CONFIG.Cache.DefTimeOut)
		global.CACHE.Put(reqMD5Key, respBytes, time.Minute*timeOut)
		return respBytes, nil
	})
	ctx.Header("Content-Disposition", "attachment; filename=iptv-diyp.txt")
	ctx.Binary(resp.([]byte))
}

func SaveToLogDir(data []byte, filename string) {
	logDir := global.CONFIG.Zap.Director
	if ok, _ := utils.PathExists(logDir); !ok {
		os.MkdirAll(logDir, os.ModePerm)
	}
	filePath := path.Join(logDir, filename)
	os.WriteFile(filePath, data, 0644)
}
