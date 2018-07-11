package captcha_sdk

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cxt90730/gocaptcha"
	"github.com/robfig/config"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
	"time"
)

const (
	DefaultPicWidth  = 120
	DefaultPicHeight = 40
)

type CaptchaService struct {
	w         int
	h         int
	fontDir   string
	imgFormat int
}

func NewCaptchaService(cfgFile string) (*CaptchaService, error) {
	cs := new(CaptchaService)
	c, err := config.ReadDefault(cfgFile)
	if err != nil {
		return nil, err
	}
	cs.w, err = c.Int("OneBookCaptcha", "width")
	if err != nil {
		cs.w = DefaultPicWidth
	}
	cs.h, err = c.Int("OneBookCaptcha", "height")
	if err != nil {
		cs.h = DefaultPicHeight
	}

	cs.fontDir, err = c.String("OneBookCaptcha", "font_dir")
	if err != nil {
		return nil, err
	}

	imgFormat, err := c.String("OneBookCaptcha", "img_format")
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(imgFormat) {

	case "png":
		cs.imgFormat = gocaptcha.ImageFormatPng
	case "gif":
		cs.imgFormat = gocaptcha.ImageFormatGif
	case "jpg":
		fallthrough
	case "jpeg":
		fallthrough
	default:
		cs.imgFormat = gocaptcha.ImageFormatJpeg
	}

	err = gocaptcha.ReadFonts("./fonts", "ttf")
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (cs *CaptchaService) GetImage(code string) (image.Image, error) {
	i, err := gocaptcha.NewCaptchaImage(cs.w, cs.h, gocaptcha.RandLightColor())
	if err != nil {
		return nil, err
	}
	i.Drawline(3)
	i.DrawSineLine()
	i.DrawHollowLine()
	i.DrawBorder(gocaptcha.ColorToRGB(0x7A7A7A))
	i.DrawNoise(gocaptcha.CaptchaComplexHigh)
	i.DrawText(code)
	return i.NRGBA, nil
}

//Get origin code and random text
//Next you need to save to storage
func (cs *CaptchaService) GetKey(length int) (string, string, error) {
	if length < 1 {
		return "", "", errors.New("code length is less than 1")
	}
	code := gocaptcha.RandNumber(length)
	return code, encode(code), nil
}

func encode(code string) string {
	key := fmt.Sprintf("%s-%s-%x", code, gocaptcha.RandNumber(20), time.Now().UnixNano())
	key = hex.EncodeToString(md5.New().Sum([]byte(key)))
	key = key[:32]
	return key
}

//保存图片对象
func (cs *CaptchaService) SaveImage(image image.Image, w io.Writer) error {

	if cs.imgFormat == gocaptcha.ImageFormatPng {
		return png.Encode(w, image)
	}
	if cs.imgFormat == gocaptcha.ImageFormatJpeg {
		return jpeg.Encode(w, image, &jpeg.Options{100})
	}
	if cs.imgFormat == gocaptcha.ImageFormatGif {
		return gif.Encode(w, image, &gif.Options{NumColors: 256})
	}

	return errors.New("not supported image format")
}
