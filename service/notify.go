package service

import (
	"fmt"
	"log"
	"time"
)

type StaffMessage struct {
	// 丰巢快递编码
	Code string `json:"code"`
	// 客户手机号
	ClientPhone string `json:"clientMobile"`
	// 快递员手机号
	StaffPhone string `json:"staffMobile"`
	// 快递柜地址
	GoodsAddress string `json:"address"`
	// 快递公司名称
	CompanyName string `json:"companyName"`
	// 状态
	BoxStatus int32 `json:"boxStatus"`
}

type StaffNotify struct {
}

func NotifyStaff() {
	itemList := GetFcBoxInfo()
	if itemList != nil {
		for index := range itemList {
			item := itemList[index]
			isSend := IsSend(item.Code)
			if !isSend {
				log.Println("已经发送短信，本次不发送")
				// 不发生短信
				continue
			}
			tencent := new(TencentMessageSend)
			tencent.StaffMobile = item.StaffPhone
			tencent.CompanyName = item.CompanyName
			tencent.Code = item.Code
			tencent.Address = item.GoodsAddress
			tencent.ClientMobile = item.ClientPhone
			log.Printf("%s\n", tencent)
			fmt.Println(tencent)
			tencent.send()
			Bark("快递通知["+item.Code+"]", "已通知快递员送货上门")
			Set(item.Code)
			// 腾讯云有频率限制 这里停顿60s
			time.Sleep(60 * time.Second)
		}
	}
}
