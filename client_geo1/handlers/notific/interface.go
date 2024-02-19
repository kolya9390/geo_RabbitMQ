package handler_notific

import "net/http"

type NotificServiceController interface {
	GetSMS(w http.ResponseWriter, r *http.Request)
	GetEMail(w http.ResponseWriter, r *http.Request)
}