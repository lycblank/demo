package weather

import (
    "bytes"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "github.com/lycblank/demo/internal/logger"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
)

const (
    spider_7d_url string = `http://www.weather.com.cn/weather/`
    spider_1d_url        = `http://www.weather.com.cn/weather1d/`
)

type WeatherContent struct {
    Title       string
    Wea         string
    Sky         string
    Temperature string
    WinType     string
    WinLevel    string
    Sun         string
}

type WeatherData struct {
    Now   time.Time
    City  string
    Datas []WeatherContent
}

type weather struct {
    city     string
    cityCode string
}

func NewWeather(city, cityCode string) *weather {
    return &weather{
        city:     city,
        cityCode: cityCode,
    }
}

func (w *weather) Spider() *WeatherData {
    return w.spider(w.cityCode)
}

func (w *weather) Format(data *WeatherData) string {
    markdown := strings.Builder{}
    markdown.WriteString("## ")
    markdown.WriteString(fmt.Sprintf("天气预报（%s）", data.Now.Format("2006-01-02")))
    markdown.WriteString("\n")

    markdown.WriteString("城市：")
    markdown.WriteString(data.City)
    markdown.WriteString("\n")

    for i := 0; i < len(data.Datas); i++ {
        markdown.WriteString("\n")
        markdown.WriteString("### ")
        markdown.WriteString(data.Datas[i].Title)
        markdown.WriteString("\n")
        markdown.WriteString("\n")

        markdown.WriteString("天气：")
        markdown.WriteString(data.Datas[i].Wea)
        markdown.WriteString("\n")
        markdown.WriteString("\n")

        if data.Datas[i].Sky != "" {
            markdown.WriteString("天空：")
            markdown.WriteString(data.Datas[i].Sky)
            markdown.WriteString("\n")
            markdown.WriteString("\n")
        }

        markdown.WriteString("温度：")
        markdown.WriteString(data.Datas[i].Temperature)
        markdown.WriteString("℃\n")
        markdown.WriteString("\n")

        markdown.WriteString(data.Datas[i].WinType)
        markdown.WriteString("：")
        markdown.WriteString(data.Datas[i].WinLevel)
        markdown.WriteString("\n")
        markdown.WriteString("\n")

        markdown.WriteString(data.Datas[i].Sun)
        markdown.WriteString("\n")
    }
    return markdown.String()
}

func (w *weather) spider(city string) *WeatherData {
    cityHtml, err := w.getCityHtml(city)
    if err != nil {
        logger.Error("get city html failed. city: %s, err: %+v", city, err)
        return nil
    }
    data, err := w.parse(cityHtml)
    if err != nil {
        logger.Error("parse city html failed. err: %+v, city: %s", err, city)
        return nil
    }
    return data
}

func (w *weather) parse(cityHtml []byte) (*WeatherData, error) {
    doc, err := goquery.NewDocumentFromReader(bytes.NewReader(cityHtml))
    if err != nil {
        logger.Error("new document failed. err: %+v", err)
        return nil, err
    }

    data := &WeatherData{
        Now:  time.Now(),
        City: w.city,
    }

    doc.Find(".clearfix li").Each(func(idx int, s *goquery.Selection) {
        title := s.Find("h1").Text()
        wea := s.Find(".wea").Text()
        sky := s.Find(".sky .txt").Text()
        tem := s.Find(".tem span").Text()
        winType := s.Find(".win span").AttrOr("title", "")
        winLevel := s.Find(".win span").Text()
        sun := s.Find(".sun span").Text()
        if wea == "" {
            return
        }
        data.Datas = append(data.Datas, WeatherContent{
            Title:       title,
            Wea:         wea,
            Sky:         sky,
            Temperature: tem,
            WinType:     winType,
            WinLevel:    winLevel,
            Sun:         sun,
        })
    })
    return data, nil
}

func (w *weather) getCityHtml(city string) ([]byte, error) {
    cityUrl := w.getCityURL(city)
    resp, err := http.Get(cityUrl)
    if err != nil {
        logger.Error("get weather html failed. city: %s, city_url: %s, err: %+v", city, cityUrl, err)
        return nil, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Error("get weather html body failed. city: %s, city_url: %s, err: %+v", city, cityUrl, err)
        return nil, err
    }
    return body, nil
}

func (w *weather) getCityURL(city string) string {
    return fmt.Sprintf("%s%s.shtml", spider_1d_url, city)
}
