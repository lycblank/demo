package push

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httputil"
)

type dingding struct {
    dingdingWebhook string
}

func NewDingDing(dingdingWebhook string) *dingding {
    return &dingding{
        dingdingWebhook: dingdingWebhook,
    }
}

func (dd *dingding) Push(contentType string, content string) error {
    if contentType == ContentTypeMarkdown {
        return dd.pushMarkdown(content)
    }

    return fmt.Errorf("not support content type(%s)", contentType)
}

func (dd *dingding) pushMarkdown(content string) error {
    body := map[string]interface{}{
        "msgtype": "markdown",
        "markdown": map[string]string{
            "title": "天气预报",
            "text":  content,
        },
    }

    jsonBody, err := json.Marshal(body)
    if err != nil {
        fmt.Printf("json marshal failed. err: %+v\n", err)
        return err
    }

    resp, err := http.Post(dd.dingdingWebhook, "application/json", bytes.NewReader(jsonBody))
    if err != nil {
        fmt.Printf("post webhook failed. err: %+v\n", err)
        return err
    }
    defer resp.Body.Close()
    respContent, err := httputil.DumpResponse(resp, true)
    if err != nil {
        fmt.Printf("dump response failed. err: %+v\n", err)
        return err
    }

    fmt.Printf("post webhook success. content: %s\n", string(respContent))
    return nil
}
