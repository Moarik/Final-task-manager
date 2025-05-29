package user

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	userSvc "taskManager/proto/gen/user"
	"taskManager/user/internal/adapter/grpc/service/user/dto"
	"taskManager/user/pkg/security"
)

type Client struct {
	userSvc.UnimplementedUserServiceServer

	uc UserUsecase
}

func New(uc UserUsecase) *Client {
	return &Client{
		uc: uc,
	}
}

func (c *Client) UserCreate(ctx context.Context, req *userSvc.UserCreateRequest) (*userSvc.UserCreateResponse, error) {
	client, err := dto.ToClientFromRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	newClient, err := c.uc.Create(ctx, client)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &userSvc.UserCreateResponse{
		Id:       newClient.ID,
		Name:     newClient.Name,
		Email:    newClient.Email,
		Password: newClient.Password,
	}, nil
}

func (c *Client) UserLogin(ctx context.Context, req *userSvc.UserLoginRequest) (*userSvc.UserLoginResponse, error) {
	client, err := dto.ToClientFromLogin(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if client.Email == "" || client.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email or password is empty")
	}

	id, err := c.uc.Login(ctx, client)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := security.CreateToken(id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	fmt.Println("ID from grpc: ", id)

	return &userSvc.UserLoginResponse{
		Id:    id,
		Token: token,
	}, nil
}

func (c *Client) UserGet(ctx context.Context, req *userSvc.UserIDRequest) (*userSvc.UserCreateResponse, error) {
	client, err := c.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &userSvc.UserCreateResponse{
		Id:       client.ID,
		Name:     client.Name,
		Email:    client.Email,
		Password: client.Password,
	}, nil
}

func (c *Client) UserDelete(ctx context.Context, req *userSvc.UserIDRequest) (*userSvc.Empty, error) {
	err := c.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &userSvc.Empty{}, nil
}
