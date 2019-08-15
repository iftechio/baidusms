package baidusms

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// BaiduSMS is config of sms service. AccessKey, SecretKey, Region should be provided
// Region is one of "bj" and "gz"
type BaiduSMS struct {
	AccessKey string
	SecretKey string
	Region    string
}

// SuccessResponse is success body of baidu response
type SuccessResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}

var (
	// Version of baidusms
	Version = "1.0.2"
)

func (bd BaiduSMS) sendRequest(method string, path string, body string) (s SuccessResponse, err error) {
	now := time.Now()
	auth := auth{bd.AccessKey, bd.SecretKey}
	var host string
	if strings.ToLower(bd.Region) == "gz" {
		host = "sms.gz.baidubce.com"
	} else {
		host = "sms.bj.baidubce.com"
	}
	targetURL := fmt.Sprintf("https://%s%s", host, path)
	req, err := http.NewRequest(method, targetURL, strings.NewReader(body))
	req.Header.Add("User-Agent", fmt.Sprintf("bce-sdk-go/%s", Version))
	req.Header.Add("Host", host)
	req.Header.Add("Connection", "close")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	req.Header.Add("x-bce-date", getCanonicalTime(now))
	sum := sha256.Sum256([]byte(body))
	req.Header.Add("x-bce-content-sha256", hex.EncodeToString(sum[:]))
	headers := req.Header
	req.Header.Add("Authorization", auth.generateAuthorization(method, path, headers, url.Values{}, now))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			err = readErr
			return
		}
		err = json.Unmarshal(bodyBytes, &s)
		if err != nil {
			return
		}
		return
	}
	err = fmt.Errorf("Request SMS error, code: %d", resp.StatusCode)
	return
}

type requestBody struct {
	InvokeID     string            `json:"invokeId"`
	PhoneNumber  string            `json:"phoneNumber"`
	TemplateCode string            `json:"templateCode"`
	ContentVar   map[string]string `json:"contentVar"`
}

// SendSMSCode will call HTTP request to Baidu API to send a sms
func (bd BaiduSMS) SendSMSCode(invokeID string, mobilePhoneNumber string, templateCode string, contentVar map[string]string) (s SuccessResponse, err error) {
	path := "/bce/v2/message"
	body := requestBody{invokeID, mobilePhoneNumber, templateCode, contentVar}
	bodyStr, err := json.Marshal(body)
	return bd.sendRequest("POST", path, string(bodyStr))
}
