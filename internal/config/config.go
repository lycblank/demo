package config

import (
    "fmt"
    yaml "gopkg.in/yaml.v3"
    "os"
)

type Config struct {
    Push PushConfig `yaml:"push"`
}

type PushConfig struct {
    DingDing DingDingConfig `yaml:"dingding"`
}

type DingDingConfig struct {
    Webhook string `yaml:"webhook"`
}

func GetConfig() *Config {
    confPath := os.Getenv("WEATHER_CONFIG_PATH")
    if confPath == "" {
        confPath = "configs/service.yaml"
    }
    confDatas, err := os.ReadFile(confPath)
    if err != nil {
        fmt.Printf("read conf path failed. err: %+v, path: %s\n", err, confPath)
        return nil
    }

    conf := &Config{}
    if err := yaml.Unmarshal(confDatas, conf); err != nil {
        fmt.Printf("unmarshal conf failed. err: %+v, data: %s\n", err, string(confDatas))
        return nil
    }
    return conf
}
