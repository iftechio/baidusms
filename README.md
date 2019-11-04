# baidusms

`baidusms` package provides Go support for Baidu SMS API.

This Project is non-official.

The document of Baidu SMS API is [here](https://cloud.baidu.com/doc/SMS/s/3jwvxrwjx)

The document of the generation of `Authorization` header is [here](https://cloud.baidu.com/doc/Reference/s/njwvz1yfu)

# Example

```go
package main

import (
	"errors"
	"log"

	"github.com/iftechio/baidusms/v2"
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
    var sendError *baidusms.ErrSendFail
    // go 1.13
		if errors.As(err, &sendError) {
			switch sendError.APICode {
			case "4621":
				// 4621 手机号配额异常
				log.Printf("baidusms single number frequency exceed")

			case "4503":
				// 4503, 手机号码格式不正确
				log.Printf("baidusms single number format error")

			}
		}
		log.Fatalf("error %s ", err)

	}
	log.Printf("%v", resp)
}
```
