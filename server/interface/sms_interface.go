package _interface

type SMS interface {
	SendVerificationCode(phoneNumber, code string) (bool, string, error)
}
