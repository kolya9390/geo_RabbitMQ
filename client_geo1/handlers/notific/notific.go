package handler_notific

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	not_service "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/notifications"
)

type NotificHandler struct {
	client not_service.Notific
}

func NewNotificHandler(client not_service.Notific) *NotificHandler {
	return &NotificHandler{client: client}
}

func (ns *NotificHandler) GetSMS(w http.ResponseWriter, r *http.Request) {

	resp,err := ns.client.SendSMS(context.Background())

	if err != nil {
		log.Printf("ns.client.SendSMS\nError:%v",err)
		return
	}

	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("Error writing JSON response:", err)
		return
	}
}

func (ns *NotificHandler) GetEMail(w http.ResponseWriter, r *http.Request){

	resp,err := ns.client.SendEmail(context.Background())

	if err != nil {
		log.Printf("ns.client.SendEmail\nError:%v",err)
		return
	}

	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("Error writing JSON response:", err)
		return
	}
}