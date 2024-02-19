package service_notifications

import "context"

type Notific interface {
	SendSMS(ctx context.Context) ([]SMS,error)
	SendEmail(ctx context.Context)([]Email,error)
}

type SMS struct{
	Phones string `json:"phone"`
	Msg    string `json:"msg"`	
}

type RespSMSList struct {
	SMSList []SMS `json:"ListSMS"`
}

type Email struct {
	Email string `json:"email"`
	Msg	  string `json:"msg"`
}