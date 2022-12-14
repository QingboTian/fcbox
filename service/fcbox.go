package service

import (
	"encoding/json"
	"fcbox/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type FcBoxResponse struct {
	Success bool               `json:"success"`
	Code    string             `json:"code"`
	Msg     string             `json:"msg"`
	Data    *FcBoxResponseData `json:"data"`
}

type FcBoxResponseData struct {
	ToPickTotal int32           `json:"toPickTotal"`
	Total       int32           `json:"total"`
	Data        []*StaffMessage `json:"data"`
}

const (
	// DaiQuJian 待取件
	DaiQuJian = 1
	// ZhiLiuZhong 滞留中
	ZhiLiuZhong = 2
)

// GetFcBoxInfo 获取丰巢接口信息
func GetFcBoxInfo() []*StaffMessage {
	yaml := config.ReadYaml()
	params := strings.NewReader(url.Values{"pageNo": {"1"}, "pageSize": {yaml.FcBox.Size}}.Encode())
	request, err := http.NewRequest("POST", yaml.FcBox.Api, params)
	request.Header.Set("Authorization", yaml.FcBox.Authorization)
	request.Header.Set("content-type", yaml.FcBox.ContentType)
	if err != nil {
		ErrorBark("构建丰巢请求参数发生异常")
		panic(err)
	}
	response, _ := http.DefaultClient.Do(request)
	if response == nil {
		ErrorBark("请求丰巢响应结果为空")
		panic("请求丰巢响应结果为空")
	}
	if response.StatusCode == 401 {
		ErrorBark("丰巢认证失败")
		panic("丰巢认证失败")
	}
	body, _ := ioutil.ReadAll(response.Body)
	fcBoxResponse := new(FcBoxResponse)
	err = json.Unmarshal(body, fcBoxResponse)
	if err != nil {
		ErrorBark("解析丰巢响应体失败")
		panic(err)
	}
	var result []*StaffMessage
	if fcBoxResponse.Success {
		fcBoxResponseData := fcBoxResponse.Data
		if fcBoxResponseData != nil && fcBoxResponseData.Data != nil {
			for index := range fcBoxResponseData.Data {
				// 快件未取
				// 1 待取件 2 滞留中
				if fcBoxResponseData.Data[index].BoxStatus == DaiQuJian || fcBoxResponseData.Data[index].BoxStatus == ZhiLiuZhong {
					result = append(result, fcBoxResponseData.Data[index])
				}
			}
		}
	}
	return result
}
