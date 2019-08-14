package main

import (
	"log"

	"github.com/iftechio/baidusms"
)

func main() {
	sms := baidusms.BaiduSMS{
		AccessKey: "YOUR_ACCESS_KEY",
		SecretKey: "YOUR_SECRET_KEY",
		Region:    "bj",
	}
	// content var is related your sms template
	// we use example template "Your SMS code is ${code}, expires in ${hour} hours"
	// so contentVAr should contain code and hour
	contentVar := map[string]string{
		"code": "1234",
		"hour": "2",
	}
	resp, err := sms.SendSMSCode("YOUR_INVOKE_ID", "17612233344", "YOUR_TEMPLATE_CODE", contentVar)
	if err != nil {
		log.Fatalf("error %s ", err)
	}
	log.Printf("%v", resp)
}
