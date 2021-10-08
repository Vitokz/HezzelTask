package proto

import (
	"HezzelTask/handler"
	"HezzelTask/models"
	"context"
	"github.com/pkg/errors"
	"log"

	//	"google.golang.org/grpc"
)

type GRPCServer struct {
	Handler *handler.Handler
}

func (g *GRPCServer) AddUser(ctx context.Context, user *AddUserRequesst) (*AddUserResponse, error) {
	switch {
	case user.User.Name == "":
		return nil, errors.New("User name is empty")
	case user.User.Email == "":
		return nil, errors.New("User email is empty")
	case user.User.Phone == "":
		return nil, errors.New("User phone is empty")
	}
    log.Println("start creating user!")
	err := g.Handler.AddUser(ctx, user.User.Name,user.User.Email,user.User.Phone)
	if err != nil {
		log.Fatal(err)
		return &AddUserResponse{}, err
	}
	log.Println("user created!")
	return &AddUserResponse{User: user.User}, nil
}

func (g *GRPCServer) DeleteUSer(ctx context.Context, email *DeleteUserRequest) (*DeleteUserResponse, error) {
	if email.Email == "" {
		return nil,errors.New("Email is empty")
	}

	err :=g.Handler.DeleteUser(ctx,email.Email)
	if err != nil {
		return nil, err
	}

	return &DeleteUserResponse{}, nil
}

func (g *GRPCServer) UserList(ctx context.Context, in *UserListRequest) (*UserListResponse, error) {
	users , err := g.Handler.UserList(ctx)
	if err != nil {
		return nil, err
	}
	return &UserListResponse{Users: converting(users)}, nil
}

func (g *GRPCServer) mustEmbedUnimplementedHezzelUsersServer() {
}


func converting(users *[]models.User) []*User {
	result := make([]*User,0)
	for _,v := range *users {
		result = append(result, &User{
			Name: v.Name,
			Email: v.Email,
			Phone: v.Phone,
		})
	}
	return result
}