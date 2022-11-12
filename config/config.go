package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Notify `yaml:"notify"`
	Redis  `yaml:"redis"`
	FcBox  `yaml:"fcbox"`
}

type FcBox struct {
	Authorization string `yaml:"authorization"`
	Api           string `yaml:"api"`
	ContentType   string `yaml:"contentType"`
	Size          string `yaml:"size"`
}

type Redis struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type Notify struct {
	Frequency int64 `yaml:"frequency"`
	Tencent   `yaml:"tencent"`
	Bark      `yaml:"bark"`
}

type Bark struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

type Tencent struct {
	SecretId   string `yaml:"secretId"`
	SecretKey  string `yaml:"secretKey"`
	SdkAppId   string `yaml:"sdkAppId"`
	SignName   string `yaml:"signName"`
	TemplateId string `yaml:"templateId"`
}

// ReadYaml 读取yaml配置
func ReadYaml() *Config {
	content, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		panic(err)
	}
	var config = new(Config)
	err = yaml.Unmarshal(content, config)
	if err != nil {
		panic(err)
	}
	return config
}
