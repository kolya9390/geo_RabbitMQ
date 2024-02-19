package user_client

import (
	"context"
	"fmt"
	"log"

	"github.com/kolya9390/gRPC_GeoProvider/client_Proxy/config"
	userApp "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/protoc/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcUserClient struct {
	client *grpc.ClientConn
}

func NewUserClient() UserGetter {

	config := config.NewAppConf("client_app/.env").UserRPC

	host_port := fmt.Sprintf("%s:%s",config.Host,config.Port)
	client, err := grpc.Dial(host_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %s", err)
	}
	return &grpcUserClient{client: client}

}

func (u *grpcUserClient) GetUserIDs(ctx context.Context, user_id int64) (User, error) {
	clientUser := userApp.NewUserServiceClient(u.client)

	resp, err := clientUser.GetUserProfileIDs(ctx, &userApp.RequestUserID{UserId: user_id})

	if err != nil {
		return User{}, err

	}

	return User{
		ID:       resp.UserId,
		Name:     resp.Name,
		Email:    resp.Email,
		Password: []byte(resp.Password)}, nil
}

func (u *grpcUserClient) GetListUsers(ctx context.Context) ([]User, error) {

	clientUser := userApp.NewUserServiceClient(u.client)

	var users []User

	resp, err := clientUser.GetListUser(ctx, &userApp.Empty{})

	if err != nil {
		return nil, err

	}

	for _, user := range resp.Users {
		users = append(users, User{
			ID:       user.UserId,
			Name:     user.Name,
			Email:    user.Email,
			Password: []byte(user.Password),
		})
	}

	return users, nil

}
