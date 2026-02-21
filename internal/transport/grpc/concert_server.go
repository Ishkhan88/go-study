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

// --- helpers ---

func concertToProto(c model.Concert) *gostudyv1.Concert {
	return &gostudyv1.Concert{
		Id:             int64(c.ID),
		Title:          c.Title,
		Date:           timestamppb.New(c.Date),
		Location:       c.Location,
		TicketPrice:    c.TicketPrice,
		TicketsTotal:   int64(c.TicketsTotal),
		TicketsLeft:    int64(c.TicketsLeft),
		OrganizerEmail: c.OrganizerEmail,
		CreatedAt:      timestamppb.New(c.CreatedAt),
		UpdatedAt:      timestamppb.New(c.UpdatedAt),
	}
}

func concertFromProto(p *gostudyv1.Concert) model.Concert {
	if p == nil {
		return model.Concert{}
	}
	c := model.Concert{
		ID:             int(p.Id),
		Title:          p.Title,
		Location:       p.Location,
		TicketPrice:    p.TicketPrice,
		TicketsTotal:   int(p.TicketsTotal),
		TicketsLeft:    int(p.TicketsLeft),
		OrganizerEmail: p.OrganizerEmail,
	}
	if p.Date != nil {
		c.Date = p.Date.AsTime()
	}
	if p.CreatedAt != nil {
		c.CreatedAt = p.CreatedAt.AsTime()
	}
	if p.UpdatedAt != nil {
		c.UpdatedAt = p.UpdatedAt.AsTime()
	}
	return c
}

func mapServiceErrConcert(err error) error {
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

// --- ConcertsService methods ---

func (s *ConcertsServer) Create(ctx context.Context, req *gostudyv1.CreateConcertRequest) (*gostudyv1.ConcertResponse, error) {
	c := concertFromProto(req.GetConcert())

	created, err := service.CreateConcert(c)
	if err != nil {
		return nil, mapServiceErrConcert(err)
	}

	return &gostudyv1.ConcertResponse{Concert: concertToProto(created)}, nil
}

func (s *ConcertsServer) Get(ctx context.Context, req *gostudyv1.IdRequest) (*gostudyv1.ConcertResponse, error) {
	got, err := service.GetConcert(int(req.GetId()))
	if err != nil {
		return nil, mapServiceErrConcert(err)
	}
	return &gostudyv1.ConcertResponse{Concert: concertToProto(got)}, nil
}

func (s *ConcertsServer) Update(ctx context.Context, req *gostudyv1.UpdateConcertRequest) (*gostudyv1.ConcertResponse, error) {
	upd := concertFromProto(req.GetConcert())

	updated, err := service.UpdateConcert(int(req.GetId()), upd)
	if err != nil {
		return nil, mapServiceErrConcert(err)
	}

	return &gostudyv1.ConcertResponse{Concert: concertToProto(updated)}, nil
}

func (s *ConcertsServer) Delete(ctx context.Context, req *gostudyv1.IdRequest) (*gostudyv1.DeleteResponse, error) {
	err := service.DeleteConcert(int(req.GetId()))
	if err != nil {
		return nil, mapServiceErrConcert(err)
	}
	return &gostudyv1.DeleteResponse{Deleted: true, Id: req.GetId()}, nil
}

func (s *ConcertsServer) List(ctx context.Context, req *gostudyv1.ListRequest) (*gostudyv1.ListConcertsResponse, error) {
	concerts := service.ListConcerts()
	out := make([]*gostudyv1.Concert, 0, len(concerts))
	for _, c := range concerts {
		out = append(out, concertToProto(c))
	}
	return &gostudyv1.ListConcertsResponse{Concerts: out}, nil
}
