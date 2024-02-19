package service_notifications

import (
	"context"
	"fmt"
	"log"

	"github.com/kolya9390/gRPC_GeoProvider/client_Proxy/config"
	notific_service "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/protoc/gen/notifications"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcNotificClient struct {
	client *grpc.ClientConn
}

func NewNotificClient() Notific {

	config := config.NewAppConf("client_app/.env").NoticRPC

	log.Println(config.Host,config.Port)

	host_port := fmt.Sprintf("%s:%s", config.Host, config.Port)
	client, err := grpc.Dial(host_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %s", err)
	}
	return &grpcNotificClient{client: client}

}

func (ns *grpcNotificClient) SendSMS(ctx context.Context) ([]SMS,error){

 clientNot := notific_service.NewNotificServiceGRPCClient(ns.client)

 var smsListResp []SMS

 smsList,err := clientNot.SmsSend(ctx,&notific_service.RequestSMS{})

 if err != nil {
	return nil,err
 }

 for _,sms := range smsList.SMSlist{
	smsListResp = append(smsListResp, SMS{
		Phones: sms.Phones,
		Msg: sms.Msg,
	})
 }

 return smsListResp,nil

}
func (ns *grpcNotificClient) SendEmail(ctx context.Context)([]Email,error){

	clientNot := notific_service.NewNotificServiceGRPCClient(ns.client)

	var emailListResp []Email

	emailList,err := clientNot.EmailSend(ctx,&notific_service.RequestEmail{})
   
	if err != nil {
	   return nil,err
	}
   
	for _,email := range emailList.ListEmail{
	   emailListResp = append(emailListResp, Email{
		   Email: email.Email,
		   Msg: email.Msg,
	   })
	}

	return emailListResp,nil
}