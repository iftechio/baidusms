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

// ErrSendFail is fail body of baidu response
type ErrSendFail struct {
	HTTPCode  int
	APICode   string
	Message   string
	RequestID string
}

func (e *ErrSendFail) Error() string {
	return fmt.Sprintf("Baidu SMS API error, httpcode: %d, code: %s, message: %s, requestID: %s",
		e.HTTPCode, e.APICode, e.Message, e.RequestID)
}

var (
	// Version of baidusms
	Version = "2.0.3"
)

func (bd BaiduSMS) sendRequest(method string, path string, body string) (*SuccessResponse, error) {
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
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return nil, readErr
		}
		var s SuccessResponse
		err = json.Unmarshal(bodyBytes, &s)
		if err != nil {
			return nil, err
		}
		if s.Code != "1000" {
			// only 1000 is correct
			return nil, &ErrSendFail{
				HTTPCode:  resp.StatusCode,
				APICode:   s.Code,
				Message:   s.Message,
				RequestID: s.RequestID,
			}
		}
		return &s, nil
	}
	return nil, &ErrSendFail{
		HTTPCode: resp.StatusCode,
	}
}

type requestBody struct {
	InvokeID     string            `json:"invokeId"`
	PhoneNumber  string            `json:"phoneNumber"`
	TemplateCode string            `json:"templateCode"`
	ContentVar   map[string]string `json:"contentVar"`
}

// SendSMSCode will call HTTP request to Baidu API to send a sms
func (bd BaiduSMS) SendSMSCode(invokeID string, mobilePhoneNumber string, templateCode string, contentVar map[string]string) (*SuccessResponse, error) {
	path := "/bce/v2/message"
	body := requestBody{invokeID, mobilePhoneNumber, templateCode, contentVar}
	bodyStr, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bd.sendRequest("POST", path, string(bodyStr))
}
