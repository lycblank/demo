package main

import (
    "github.com/jasonlvhit/gocron"
    "github.com/lycblank/demo/internal/config"
    "github.com/lycblank/demo/internal/logger"
    "github.com/lycblank/demo/internal/push"
    "github.com/lycblank/demo/internal/weather"
    "os"
)

func main() {
    // 天气预报
    spiderWeather := func(){
        logger.Info("start spider weather...")
        wh := weather.NewWeather("成都", "101270101")
        data := wh.Spider()
        if data == nil {
            logger.Error("spider weather failed. city: %s", "成都")
            return
        }

        conf := config.GetConfig()
        p := push.NewDingDing(conf.Push.DingDing.Webhook)
        p.Push(push.ContentTypeMarkdown, wh.Format(data))
        logger.Info("end spider weather...")
    }

    gocron.Every(1).Day().At("00:00:01").Do(spiderWeather)
    gocron.Every(1).Day().At("08:00:01").Do(spiderWeather)
    gocron.Every(1).Day().At("12:00:01").Do(spiderWeather)
    gocron.Every(1).Day().At("19:00:01").Do(spiderWeather)

    logger.Info("process start %d ...", os.Getpid())
    <-gocron.Start()
    logger.Info("process end %d ...", os.Getpid())
}
