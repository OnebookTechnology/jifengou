package _interface

import (
	"image"
	"time"
)

type CaptchaInfo struct {
	Text       string
	CreateTime time.Time
	ShownTimes int
}

type Captcha interface {
	GetImage(code string) (image.Image, error)
	GetKey(length int) (string, string, error)
}
