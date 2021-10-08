package main

import (
	"HezzelTask/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

const(
	port = "8080"
)

func main() {

	conn,err := grpc.Dial(":"+port,grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	c := proto.NewHezzelUsersClient(conn)

	_,err =c.AddUser(context.Background(),&proto.AddUserRequesst{User: &proto.User{
		Name: "John1",
		Email: "JohnDoe@mail.ru",
		Phone: "+39487322343",
	}})
	if err != nil {
		log.Fatal(err)
	}
	_,err =c.AddUser(context.Background(),&proto.AddUserRequesst{User: &proto.User{
		Name: "John2",
		Email: "JohnBoe@mail.ru",
		Phone: "+39487322343",
	}})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Users create work!")

	res,err := c.UserList(context.Background(),&proto.UserListRequest{})
	if err !=nil {
		log.Fatal(err)
	}
	for _,v :=range res.Users {
		fmt.Println(v)
	}
	log.Println("Users list work!")

	_,err = c.DeleteUSer(context.Background(),&proto.DeleteUserRequest{Email: "JohnDoe@mail.ru"})
	if err !=nil {
		log.Fatal(err)
	}
	_,err = c.DeleteUSer(context.Background(),&proto.DeleteUserRequest{Email: "JohnBoe@mail.ru"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Users delete work!")
}