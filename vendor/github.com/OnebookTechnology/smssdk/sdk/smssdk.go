package sdk

import (
	"errors"
	"fmt"
	"github.com/robfig/config"
	"github.com/ying32/alidayu"
	"regexp"
)

//var sms *SMS

type SMSService struct {
	signName        string
	accessKeyId     string
	accessKeySecret string
	templateMap     map[string]string
}

func NewSMSService(confPath string) (*SMSService, error) {
	c, err := config.ReadDefault(confPath)
	if err != nil {
		return nil, err
	}
	smsService := new(SMSService)
	smsService.templateMap = make(map[string]string)
	smsService.signName, err = c.String("OneBookSMS", "signName")
	if err != nil {
		return nil, err
	}
	smsService.templateMap["verify"], err = c.String("OneBookSMS", "templateCodeVerify")
	if err != nil {
		return nil, err
	}
	smsService.templateMap["borrow"], err = c.String("OneBookSMS", "templateCodeBorrow")
	if err != nil {
		return nil, err
	}
	smsService.templateMap["return"], err = c.String("OneBookSMS", "templateCodeReturn")
	if err != nil {
		return nil, err
	}
	smsService.templateMap["purchase"], err = c.String("OneBookSMS", "templateCodePurchase")
	if err != nil {
		return nil, err
	}
	//smsService.templateMap["purchase_borrowed"], err = c.String("OneBookSMS", "templateCodePurchaseBorrowed")
	//if err != nil {
	//	return nil, err
	//}
	smsService.templateMap["heartbeat_error"], err = c.String("OneBookSMS", "templateCodeHeartBeatErr")
	if err != nil {
		return nil, err
	}
	smsService.templateMap["overdue"], err = c.String("OneBookSMS", "templateCodeOverdue")
	if err != nil {
		return nil, err
	}
	smsService.accessKeyId, err = c.String("OneBookSMS", "accessKeyId")
	if err != nil {
		return nil, err
	}
	smsService.accessKeySecret, err = c.String("OneBookSMS", "accessKeySecret")
	if err != nil {
		return nil, err
	}
	return smsService, nil
}

// send verification code
func (s *SMSService) SendVerificationCode(phoneNumber, code string) (bool, string, error) {
	if !phoneNumberValidation(phoneNumber) {
		return false, "phone:" + phoneNumber + " is invalid", errors.New("invalid phoneNumber: " + phoneNumber)
	}
	var paramString = "{\"code\":\"" + code + "\"}"
	ok, msg, err := alidayu.SendSMS(phoneNumber, s.signName, s.templateMap["verify"], paramString, s.accessKeyId, s.accessKeySecret)
	return ok, msg, err
}

// send return book message
func (s *SMSService) SendReturnMessage(phoneNumber, bookName string) (bool, string, error) {
	if !phoneNumberValidation(phoneNumber) {
		return false, "phone:" + phoneNumber + " is invalid", errors.New("invalid phoneNumber: " + phoneNumber)
	}
	var paramString = "{\"bookname\":\"" + bookName + "\"}"
	ok, msg, err := alidayu.SendSMS(phoneNumber, s.signName, s.templateMap["return"], paramString, s.accessKeyId, s.accessKeySecret)
	return ok, msg, err
}

// send purchase book message
func (s *SMSService) SendPurchaseMessage(phoneNumber, bookName, price string) (bool, string, error) {
	if !phoneNumberValidation(phoneNumber) {
		return false, "phone:" + phoneNumber + " is invalid", errors.New("invalid phoneNumber: " + phoneNumber)
	}
	var paramString = "{\"bookname\":\"" + bookName + "\",\"price\":\"" + price + "\"}"
	ok, msg, err := alidayu.SendSMS(phoneNumber, s.signName, s.templateMap["purchase"], paramString, s.accessKeyId, s.accessKeySecret)
	return ok, msg, err
}

// send purchase borrowed book message
func (s *SMSService) SendPurchaseBorrowedMessage(phoneNumber, bookName, price string) (bool, string, error) {
	if !phoneNumberValidation(phoneNumber) {
		return false, "phone:" + phoneNumber + " is invalid", errors.New("invalid phoneNumber: " + phoneNumber)
	}
	var paramString = "{\"bookname\":\"" + bookName + "\",\"price\":\"" + price + "\"}"
	ok, msg, err := alidayu.SendSMS(phoneNumber, s.signName, s.templateMap["purchase_borrowed"], paramString, s.accessKeyId, s.accessKeySecret)
	return ok, msg, err
}

// send borrow book message
func (s *SMSService) SendBorrowMessage(phoneNumber, bookName, endTime string) (bool, string, error) {
	if !phoneNumberValidation(phoneNumber) {
		return false, "phone:" + phoneNumber + " is invalid", errors.New("invalid phoneNumber: " + phoneNumber)
	}
	var paramString = "{\"bookname\":\"" + bookName + "\",\"freeperiod\":\"" + endTime + "\"}"
	ok, msg, err := alidayu.SendSMS(phoneNumber, s.signName, s.templateMap["borrow"], paramString, s.accessKeyId, s.accessKeySecret)
	return ok, msg, err
}

// send borrow book message
func (s *SMSService) SendHeatBeatErrMessage(phoneNumber, storeName, storeType, dateTime string) (bool, string, error) {
	if !phoneNumberValidation(phoneNumber) {
		return false, "phone:" + phoneNumber + " is invalid", errors.New("invalid phoneNumber: " + phoneNumber)
	}
	var paramString = "{\"storename\":\"" + storeName + "\",\"type\":\"" + storeType + "\",\"datetime\":\"" + dateTime + "\"}"
	ok, msg, err := alidayu.SendSMS(phoneNumber, s.signName, s.templateMap["heartbeat_error"], paramString, s.accessKeyId, s.accessKeySecret)
	return ok, msg, err
}

// send borrow book message
func (s *SMSService) SendOverdueMessage(phoneNumber, bookName, dateTime string) (bool, string, error) {
	if !phoneNumberValidation(phoneNumber) {
		return false, "phone:" + phoneNumber + " is invalid", errors.New("invalid phoneNumber: " + phoneNumber)
	}
	var paramString = "{\"bookname\":\"" + bookName + "\",\"datetime\":\"" + dateTime + "\"}"
	ok, msg, err := alidayu.SendSMS(phoneNumber, s.signName, s.templateMap["overdue"], paramString, s.accessKeyId, s.accessKeySecret)
	return ok, msg, err
}

// check phone number validation
func phoneNumberValidation(phoneNumber string) bool {
	var pattern = "^1[0-9]{10}$"
	result, err := regexp.MatchString(pattern, phoneNumber)
	if err != nil {
		fmt.Println(err)
	}
	return result
}
