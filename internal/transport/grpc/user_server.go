package grpc

import (
	"context"
	"errors"

	gostudyv1 "github.com/Ishkhan88/go-study/gen/gostudy/v1"
	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// --- helpers (маппинг proto <-> model) ---

func userToProto(u model.User) *gostudyv1.User {
	return &gostudyv1.User{
		Id:        int64(u.ID),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		AvatarUrl: u.AvatarURL,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func userFromProto(p *gostudyv1.User) model.User {
	if p == nil {
		return model.User{}
	}
	u := model.User{
		ID:        int(p.Id),
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
		Phone:     p.Phone,
		AvatarURL: p.AvatarUrl,
	}
	if p.CreatedAt != nil {
		u.CreatedAt = p.CreatedAt.AsTime()
	}
	if p.UpdatedAt != nil {
		u.UpdatedAt = p.UpdatedAt.AsTime()
	}
	return u
}

func mapServiceErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, service.ErrBadInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// --- UsersService methods ---

func (s *UsersServer) Create(ctx context.Context, req *gostudyv1.CreateUserRequest) (*gostudyv1.UserResponse, error) {
	u := userFromProto(req.GetUser())

	created, err := service.CreateUser(u)
	if err != nil {
		return nil, mapServiceErr(err)
	}

	return &gostudyv1.UserResponse{User: userToProto(created)}, nil
}

func (s *UsersServer) Get(ctx context.Context, req *gostudyv1.IdRequest) (*gostudyv1.UserResponse, error) {
	got, err := service.GetUser(int(req.GetId()))
	if err != nil {
		return nil, mapServiceErr(err)
	}
	return &gostudyv1.UserResponse{User: userToProto(got)}, nil
}

func (s *UsersServer) Update(ctx context.Context, req *gostudyv1.UpdateUserRequest) (*gostudyv1.UserResponse, error) {
	upd := userFromProto(req.GetUser())

	updated, err := service.UpdateUser(int(req.GetId()), upd)
	if err != nil {
		return nil, mapServiceErr(err)
	}

	return &gostudyv1.UserResponse{User: userToProto(updated)}, nil
}

func (s *UsersServer) Delete(ctx context.Context, req *gostudyv1.IdRequest) (*gostudyv1.DeleteResponse, error) {
	err := service.DeleteUser(int(req.GetId()))
	if err != nil {
		return nil, mapServiceErr(err)
	}
	return &gostudyv1.DeleteResponse{Deleted: true, Id: req.GetId()}, nil
}

func (s *UsersServer) List(ctx context.Context, req *gostudyv1.ListRequest) (*gostudyv1.ListUsersResponse, error) {
	users := service.ListUsers()
	out := make([]*gostudyv1.User, 0, len(users))
	for _, u := range users {
		out = append(out, userToProto(u))
	}
	return &gostudyv1.ListUsersResponse{Users: out}, nil
}
