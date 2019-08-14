package baidusms

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

// TestHash tests HMAC SHA256 hash function
func TestHash(t *testing.T) {
	data := "bce-auth-v1/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/2015-04-27T08:23:49Z/1800"
	key := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	result := "1d5ce5f464064cbee060330d973218821825ac6952368a482a592e6615aef479"
	testResult := hash(data, key)
	if result == testResult {
		t.Log("hash test success")
	} else {
		t.Error("hash test fail")
	}
}

// TestStringNormalize tests string normalize
// rule: https://cloud.baidu.com/doc/Reference/s/njwvz1yfu
func TestStringNormalize(t *testing.T) {
	data1 := "NFzcPqhviddjRNnSOGo4rw=="
	data2 := "text/plain"
	data3 := "Mon, 27 Apr 2015 16:23:49 +0800"

	result1 := "NFzcPqhviddjRNnSOGo4rw%3D%3D"
	result2 := "text%2Fplain"
	result3 := "Mon%2C%2027%20Apr%202015%2016%3A23%3A49%20%2B0800"

	testRes1 := stringNormalize(data1, true)
	testRes2 := stringNormalize(data2, true)
	testRes2WithoutSlash := stringNormalize(data2, false)
	testRes3 := stringNormalize(data3, true)

	errorMsg := "string normalize test fail %s"
	if testRes1 != result1 {
		t.Errorf(errorMsg, testRes1)
	}
	if testRes2 != result2 {
		t.Errorf(errorMsg, testRes2)
	}
	if testRes2WithoutSlash != data2 {
		t.Errorf(errorMsg, testRes2WithoutSlash)
	}
	if testRes3 != result3 {
		t.Errorf(errorMsg, testRes3)
	}

	t.Log("string normalize test success")
}

func TestQueryStringCanonicalization(t *testing.T) {
	qs := url.Values{
		"text":   []string{},
		"text1":  []string{"测试"},
		"text10": []string{"test"},
	}
	testRes := queryStringCanonicalization(qs)
	if testRes == "text10=test&text1=%E6%B5%8B%E8%AF%95&text=" {
		t.Log("qs Canonicalization test success")
	} else {
		t.Errorf("qs Canonicalization test fail %s", testRes)
	}
}

func TestHeadersCanonicalization(t *testing.T) {
	headers := http.Header{
		"Host":           []string{"bj.bcebos.com"},
		"Date":           []string{"Mon, 27 Apr 2015 16:23:49 +0800"},
		"Content-Type":   []string{"text/plain"},
		"Content-Length": []string{"8"},
		"Content-Md5":    []string{"NFzcPqhviddjRNnSOGo4rw=="},
		"x-bce-date":     []string{"2015-04-27T08:23:49Z"},
	}
	res, signedHeaders := headersCanonicalization(headers)
	if res == "content-length:8\ncontent-md5:NFzcPqhviddjRNnSOGo4rw%3D%3D\ncontent-type:text%2Fplain\nhost:bj.bcebos.com\nx-bce-date:2015-04-27T08%3A23%3A49Z" && signedHeaders == "content-length;content-md5;content-type;host;x-bce-date" {
		t.Log("header Canonicalization test success")
	} else {
		t.Errorf("header Canonicalization test fail %s %s", res, signedHeaders)
	}
}

func TestGenerateAuthorization(t *testing.T) {
	method := "PUT"
	headers := http.Header{
		"Host":           []string{"bj.bcebos.com"},
		"Date":           []string{"Mon, 27 Apr 2015 16:23:49 +0800"},
		"Content-Type":   []string{"text/plain"},
		"Content-Length": []string{"8"},
		"Content-Md5":    []string{"NFzcPqhviddjRNnSOGo4rw=="},
		"x-bce-date":     []string{"2015-04-27T08:23:49Z"},
	}
	// body := "Example"
	qs := url.Values{
		"partNumber": []string{"9"},
		"uploadId":   []string{"a44cc9bab11cbd156984767aad637851"},
	}
	path := "/v1/test/myfolder/readme.txt"
	auth := auth{
		AccessKey: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		SecretKey: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	}
	currentTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41Z")
	testRes := auth.generateAuthorization(method, path, headers, qs, currentTime)
	if testRes == "bce-auth-v1/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/2012-11-01T22:08:41Z/1800/content-length;content-md5;content-type;host;x-bce-date/b47ae03e099e2a1553ef56fbba0c23fd5b632c99bc118df7f652b0fe1b4c5988" {
		t.Log("auth test success")
	} else {
		t.Errorf("auth test fail %s", testRes)
	}
}
