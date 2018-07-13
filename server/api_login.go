package server

import (
	"database/sql"
	"errors"
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/cxt90730/xxtea-go/xxtea"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"image/png"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const XXTEA_KEY = "JiFenGou"
const SessionPrefix = "/session/"
const UserSessionPrefix = "/usession/"
const VerifyCodePrefix = "/vcode/"
const CaptchaPrefix = "/captcha/"
const MaxSessionTimeout = 7200 //2 hours

func Login(ctx *gin.Context) {
	crossDomain(ctx)
	user := ctx.PostForm("username")
	pwd := ctx.PostForm("password")

	u, err := server.DB.FindUser(user, pwd)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.String(http.StatusUnauthorized, "user or password is incorrect.")
			return
		}
		ctx.String(http.StatusInternalServerError, "err: %s", err.Error())
		return
	}

	session := xxtea.EncryptStdToURLString(user+"/"+nowTimestampString(), XXTEA_KEY)

	ts := strconv.FormatInt(time.Now().Add(MaxSessionTimeout*time.Second).Unix(), 10)
	err = server.Consist.Put(SessionPrefix+session, ts, MaxSessionTimeout*time.Second)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "err: %s", err)
	}
	ctx.Header("SESSION", session)
	ctx.String(http.StatusOK, "%s", u.UserName)
}

func CheckUserSession(c *gin.Context) error {
	session := c.GetHeader("SESSION")
	if len(session) == 0 {
		return errors.New("empty session.")
	}

	//TODO: Consist session
	// Check if etcd contain this sessionKey
	outTime, err := server.Consist.Get(UserSessionPrefix + session)
	if len(outTime) == 0 {
		return errors.New("invalid session key:" + session)
	}
	if err != nil {
		return errors.New("Consist.Get err:" + session)
	}
	// check if login object's timestamp is over time
	now := time.Now().Unix()
	sessionTime, _ := strconv.ParseInt(outTime, 10, 64)
	diff := now - sessionTime
	if diff > MaxSessionTimeout {
		// over time means current user without any operation over 30 minutes
		return errors.New("Usession time out:" + session)
	}
	return nil
}

func CheckUserSessionWithPhone(c *gin.Context) (int, error) {
	session := c.GetHeader("SESSION")
	if len(session) == 0 {
		return 0, errors.New("empty session.")
	}

	//TODO: Consist session
	// Check if etcd contain this sessionKey
	outTime, err := server.Consist.Get(UserSessionPrefix + session)
	if len(outTime) == 0 {
		return 0, errors.New("invalid session key:" + session)
	}
	if err != nil {
		return 0, errors.New("Consist.Get err:" + session)
	}

	// check if login object's timestamp is over time
	now := time.Now().Unix()
	sessionTime, _ := strconv.ParseInt(outTime, 10, 64)
	diff := now - sessionTime
	if diff > MaxSessionTimeout {
		// over time means current user without any operation over 30 minutes
		return 0, errors.New("Usession time out:" + session)
	}

	//strconv.FormatUint(vReq.PhoneNumber, 10)+":"+nowTimestampString(), XXTEA_KEY
	s, err := xxtea.DecryptURLToStdString(session, XXTEA_KEY)
	if err != nil {
		return 0, errors.New("invalid session:" + session)
	}

	phoneStr := strings.Split(s, ":")[0]
	phoneNumber, err := strconv.Atoi(phoneStr)
	if err != nil {
		return 0, errors.New("invalid session:" + session)
	}
	return phoneNumber, nil
}

func CheckSession(c *gin.Context) error {
	session := c.GetHeader("SESSION")
	if len(session) == 0 {
		return errors.New("empty session:" + session)
	}

	//TODO: Consist session
	// Check if etcd contain this sessionKey
	outTime, err := server.Consist.Get(SessionPrefix + session)
	if len(outTime) == 0 {
		return errors.New("invalid session key:" + session)
	}
	if err != nil {
		return errors.New("Consist.Get err:" + session)
	}
	// check if login object's timestamp is over time
	now := time.Now().Unix()
	sessionTime, _ := strconv.ParseInt(outTime, 10, 64)
	diff := now - sessionTime
	if diff > MaxSessionTimeout {
		// over time means current user without any operation over 30 minutes
		return errors.New("Csession time out:" + session)
	}
	return nil
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		crossDomain(c)
		// Check if header contain SESSION
		session := c.GetHeader("SESSION")
		if len(session) == 0 {
			sendFailedResponse(c, SessionErr, "invalid session.")
			c.Abort()
			return
		}

		//TODO: Consist session
		// Check if etcd contain this sessionKey
		outTime, err := server.Consist.Get(SessionPrefix + session)
		if len(outTime) == 0 {
			sendFailedResponse(c, SessionErr, "invalid session key:", session)
			c.Abort()
			return
		}
		if err != nil {
			sendFailedResponse(c, SessionErr, "Consist.Get err:", err)
			c.Abort()
			return
		}
		// check if login object's timestamp is over time
		now := time.Now().Unix()
		sessionTime, _ := strconv.ParseInt(outTime, 10, 64)
		diff := now - sessionTime
		if diff > MaxSessionTimeout {
			// over time means current user without any operation over 30 minutes
			sendFailedResponse(c, SessionErr, "session time out")
			c.Abort()
			return
		}
		c.Next()
	}
}

type CaptchaInfo struct {
	Text       string
	CreateTime time.Time
	ShownTimes int
}

func GetKey(ctx *gin.Context) {
	crossDomain(ctx)
	origin, key, err := server.Captcha.GetKey(4)
	if err != nil {
		sendFailedResponse(ctx, Err, err.Error())
		return
	}
	info := &CaptchaInfo{
		Text:       origin,
		CreateTime: time.Now(),
		ShownTimes: 0,
	}
	infoJson, _ := jsoniter.MarshalToString(info)
	err = server.Consist.Put(CaptchaPrefix+key, infoJson, 3*time.Minute)
	if err != nil {
		sendFailedResponse(ctx, Err, err.Error())
		return
	}
	sendSuccessResponseWithMessage(ctx, key, nil)
	return
}

// 返回图片验证码
func ShowImage(ctx *gin.Context) {
	crossDomain(ctx)
	key := ctx.Param("key")
	codeJson, err := server.Consist.Get(CaptchaPrefix + key)
	if err != nil {
		sendFailedResponse(ctx, Err, err.Error())
		return
	}
	if len(codeJson) == 0 {
		logger.Error("empty captcha code. key:", key)
		sendFailedResponse(ctx, Err, "empty captcha code. key: %s", key)
		return
	}

	//Text       string
	//CreateTime time.Time
	//ShownTimes int
	i := new(CaptchaInfo)
	err = jsoniter.UnmarshalFromString(codeJson, i)
	if err != nil {
		sendFailedResponse(ctx, Err, "unmarshal ShowImage data err", err.Error(), "data:", codeJson)
		return
	}

	//If duplicate request
	if i.ShownTimes > 0 {
		origin, _, _ := server.Captcha.GetKey(len(i.Text))
		i.Text = origin
	}

	i.ShownTimes++
	newInfo, _ := jsoniter.MarshalToString(i)
	err = server.Consist.Put(CaptchaPrefix+key, newInfo, 3*time.Minute)
	if err != nil {
		sendFailedResponse(ctx, Err, err.Error())
		return
	}

	img, err := server.Captcha.GetImage(i.Text)
	if err != nil {
		sendFailedResponse(ctx, Err, err.Error())
		return
	}

	ctx.Header("Content-Type", "image/png")
	png.Encode(ctx.Writer, img)
}

type VerifyRequest struct {
	Code        string `json:"code"`
	Key         string `json:"key"`
	PhoneNumber uint64 `json:"phone_number"`
}

type VerifyCodeInfo struct {
	Code        string        `json:"code"`
	PhoneNumber uint64        `json:"phone_number"`
	CreateTime  time.Duration `json:"create_time"`
	UserInfo    string        `json:"user_info"`
}

// VerifyCaptcha
func Verify(ctx *gin.Context) {
	crossDomain(ctx)
	var vReq VerifyRequest
	if err := ctx.ShouldBindJSON(&vReq); err == nil {
		infoJson, err := server.Consist.Get(CaptchaPrefix + vReq.Key)
		if err != nil {
			sendFailedResponse(ctx, Err, err.Error())
			return
		}
		if len(infoJson) == 0 {
			sendFailedResponse(ctx, Err, "empty captcha. key:", vReq.Key)
			return
		}

		info := new(CaptchaInfo)
		err = jsoniter.Unmarshal([]byte(infoJson), info)
		if err != nil {
			sendFailedResponse(ctx, Err, "unmarshal infoJson data err", err.Error(), " data:", infoJson)
			return
		}

		if vReq.Code != info.Text {
			logger.Warning("captcha info is not match!", "req.code:", vReq.Code, "code", info.Text)
			sendFailedResponse(ctx, Err, "captcha info is not match!")
			return
		}

		//Generate SMS verify code
		phoneNumber := vReq.PhoneNumber
		vcode, key, err := server.Captcha.GetKey(6)
		if err != nil {
			sendFailedResponse(ctx, Err, err.Error())
			return
		}
		logger.Info("get verify code. phone:", phoneNumber, "key:", key, "code:", vcode)

		vi := &VerifyCodeInfo{
			Code:        vcode,
			PhoneNumber: phoneNumber,
		}
		viJson, _ := jsoniter.MarshalToString(vi)

		err = server.Consist.Put(VerifyCodePrefix+key, viJson, 3*time.Minute)
		if err != nil {
			sendFailedResponse(ctx, Err, err.Error())
			return
		}

		sendSuccessResponseWithMessage(ctx, key, nil)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}

}

//发送短信验证码
func SendVerifyCode(ctx *gin.Context) {
	crossDomain(ctx)
	key := ctx.Param("key")
	viJson, err := server.Consist.Get(VerifyCodePrefix + key)
	if err != nil {
		sendFailedResponse(ctx, Err, err.Error())
		return
	}
	if len(viJson) == 0 {
		sendFailedResponse(ctx, Err, "empty verify code. key:", key)
		return
	}
	var vReq VerifyCodeInfo
	err = jsoniter.UnmarshalFromString(viJson, &vReq)
	if err != nil {
		sendFailedResponse(ctx, Err, "UnmarshalFromString err:", err)
		return
	}
	go func(phone, code string) {
		success, msg, err := server.SMS.SendVerificationCode(phone, code)
		if err != nil {
			logger.Error("SendIdentifyingCode err:", err, "phone:", phone, "code:", code)
		}
		if !success {
			logger.Warning("SendIdentifyingCode Failed!", "phone:", phone, "msg:", msg, "code:", code)
		}
	}(strconv.Itoa(int(vReq.PhoneNumber)), vReq.Code)

	sendSuccessResponse(ctx, nil)
}

//验证短信验证码
func VerifyVCode(ctx *gin.Context) {
	crossDomain(ctx)
	var vReq VerifyRequest
	if err := ctx.ShouldBindJSON(&vReq); err == nil {
		viJson, err := server.Consist.Get(VerifyCodePrefix + vReq.Key)
		if err != nil {
			sendFailedResponse(ctx, Err, err.Error())
			return
		}
		if len(viJson) == 0 {
			sendFailedResponse(ctx, Err, "empty verify code. key:", vReq.Key)
			return
		}

		vi := new(VerifyCodeInfo)
		err = jsoniter.Unmarshal([]byte(viJson), vi)
		if err != nil {
			sendFailedResponse(ctx, Err, "unmarshal viJson data err:", err.Error(), "data:", viJson)
			return
		}

		if vReq.Code != vi.Code {
			logger.Warning("SMS info is not match!", "req.code:", vReq.Code, "code", vi.Code)
			sendFailedResponse(ctx, Err, "captcha info is not match!")
			return
		}
		_, err = server.DB.FindMobileUser(vi.PhoneNumber)
		if err != nil {
			//if NO user found, register user
			if err == sql.ErrNoRows {
				_, err = registerUser(vi)
				if err != nil {
					sendFailedResponse(ctx, Err, "db error when RegisterMobileUser err:", err.Error(), "Phone:",
						vi.PhoneNumber)
					return
				}
				goto SUCCESS
			} else {
				sendFailedResponse(ctx, Err, "db error when RegisterMobileUser Phone:", vi.PhoneNumber, "err:",
					err.Error())
				return
			}
		}

		//if user can be found, update user login time
		goto SUCCESS
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}

SUCCESS:
	userSession := xxtea.EncryptStdToURLString(strconv.FormatUint(vReq.PhoneNumber, 10)+":"+nowTimestampString(), XXTEA_KEY)
	//Save session
	err := server.Consist.Put(UserSessionPrefix+userSession, nowTimestampString(), time.Duration(86400*time.Second))
	if err != nil {
		sendFailedResponse(ctx, Err, err.Error())
		return
	}
	ctx.Header("SESSION", userSession)
	sendSuccessResponse(ctx, nil)
	return
}

func registerUser(vi *VerifyCodeInfo) (*models.MobileUser, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	logger.Info("New User! phone:", vi.PhoneNumber, "time:", now)

	//phone_number,login_time,register_time,balance,deposit,credit
	newUser := &models.MobileUser{
		PhoneNumber:  vi.PhoneNumber,
		RegisterTime: now,
	}
	err := server.DB.RegisterMobileUser(newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
