package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	AppId     = "wx9d18575704ec8e27"
	AppSecret = "f99a2b3d2e2c463e006fa1ffaa92c231"
)

type WxRequest struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	Errcode     int    `json:"errcode,omitempty"`
	Ticket      string `json:"ticket,omitempty"`
	Errmsg      string `json:"errmsg,omitempty"`
}

type WxResponse struct {
	Url         string `json:"url"`
	JsapiTicket string `json:"jsapi_ticket"`
	NonceStr    string `json:"nonceStr"`
	Timestamp   string `json:"timestamp"`
	Signature   string `json:"signature"`
	Appid       string `json:"appid"`
}

func GetWxConfig(ctx *gin.Context) {
	crossDomain(ctx)
	shareUrl := strings.Split(ctx.Query("url"), "#")[0]
	t, err := server.DB.FindToken()
	if err != nil {
		sendFailedResponse(ctx, Err, "db error when FindToken. error:", err)
		return
	}
	et, _ := nowParse(t.ExpireTime)
	//超过时间，返回新的token，并更新
	now := time.Now()
	if now.After(et) {
		url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", AppId, AppSecret)
		res, err := http.Get(url)
		if err != nil {
			sendFailedResponse(ctx, Err, "Get wx token error:", err)
			return
		}
		if res.StatusCode != http.StatusOK {
			sendFailedResponse(ctx, Err, "Get wx token status:", res.Status)
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		wt := new(WxRequest)
		err = jsoniter.Unmarshal(body, wt)
		if err != nil {
			sendFailedResponse(ctx, Err, "Unmarshal wx token err:", err)
			return
		}
		if wt.Errcode != 0 {
			sendFailedResponse(ctx, Err, "errcode:", wt.Errcode, "errmsg:", wt.Errmsg)
			return
		}

		t.Token = wt.AccessToken
		t.ExpireTime = now.Add(time.Duration(wt.ExpiresIn) * time.Second).Format("2006-01-02 15:04:05")
		err = server.DB.UpdateToken(t.Token, t.ExpireTime)
		if err != nil {
			sendFailedResponse(ctx, Err, "db error when UpdateToken. error:", err)
			return
		}
	}

	// 获取ticket
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", t.Token)
	res, err := http.Get(url)
	if err != nil {
		sendFailedResponse(ctx, Err, "Get wx ticket error:", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		sendFailedResponse(ctx, Err, "Get wx ticket status:", res.Status)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	wt := new(WxRequest)
	err = jsoniter.Unmarshal(body, wt)
	if err != nil {
		sendFailedResponse(ctx, Err, "Unmarshal wx token err:", err)
		return
	}
	if wt.Errcode != 0 {
		sendFailedResponse(ctx, Err, "errcode:", wt.Errcode, "errmsg:", wt.Errmsg)
		return
	}

	nonceStr := genNonceStr()
	timestamp := genTimeStamp()
	r := &WxResponse{
		Url:         shareUrl,
		JsapiTicket: wt.Ticket,
		NonceStr:    nonceStr,
		Timestamp:   timestamp,
		Appid:       AppId,
	}
	sign := genSign(r)
	r.Signature = sign
	resData := &ResData{
		WxResponse: r,
	}
	sendSuccessResponse(ctx, resData)
	return
}

func genSign(r *WxResponse) string {
	s := "jsapi_ticket=" + r.JsapiTicket +
		"&noncestr=" + r.NonceStr +
		"&timestamp=" + r.Timestamp +
		"&url=" + strings.Split(r.Url, "#")[0]
	return string(doSHA1([]byte(s)))
}

func genNonceStr() string {
	return time.Now().Format("20060102150405")
}

func genTimeStamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
