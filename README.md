# baidusms

`baidusms` package provides Go support for Baidu SMS API.

The document of Baidu SMS API is [here](https://cloud.baidu.com/doc/SMS/s/3jwvxrwjx)

The document of the generation of `Authorization` header is [here](https://cloud.baidu.com/doc/Reference/s/njwvz1yfu)

# Example

```go
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
	// example template "Your SMS code is ${code}, expires in ${hour} hours"
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
```
