package cat

import (
	"context"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

type Service struct {
}

func (s *Service) GetMyCat(ctx context.Context, msg *GetMyCatMessage) (*MyCatResponse, error) {
	switch msg.TargetCat {
	case "tama":
		return &MyCatResponse{
			Name: msg.TargetCat,
			Kind: "Mainecoon",
		}, nil
	case "mike":
		return &MyCatResponse{
			Name: msg.TargetCat,
			Kind: "Norwegian Forest Cat",
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "not found your cat")
}

func NewCatService() CatServer {
	return &Service{}
}
