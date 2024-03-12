package grpcservice

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"google.golang.org/grpc"
	"log"
	"my-app/common"
	query2 "my-app/module/category/query"
	"my-app/proto/category"
	"net"
)

type service struct {
	port int
	sctx sctx.ServiceContext
}

func NewCatGRPCService(port int, sctx sctx.ServiceContext) *service {
	return &service{port: port, sctx: sctx}
}

type categoryServer struct {
	category.UnimplementedCategoryServer
	sctx sctx.ServiceContext
}

func NewCategoryServer(sctx sctx.ServiceContext) *categoryServer {
	return &categoryServer{sctx: sctx}
}

func (cs *categoryServer) GetCategoriesByIds(ctx context.Context, request *category.GetCatIdsRequest) (*category.CatIdsResp, error) {
	var cats []query2.CategoryDTO

	dbContext := cs.sctx.MustGet(common.KeyGorm).(common.DbContext)

	ids := make([]uuid.UUID, len(request.Ids))

	for i := range ids {
		ids[i] = uuid.MustParse(request.Ids[i])
	}

	if err := dbContext.GetDB().Table(query2.CategoryDTO{}.TableName()).
		Where("id in (?)", ids).
		Find(&cats).Error; err != nil {
		return nil, core.ErrBadRequest.WithError("cannot list categories").WithDebug(err.Error())
	}

	results := make([]*category.CategoryDTO, len(cats))

	for i := range results {
		results[i] = &category.CategoryDTO{
			Id:    cats[i].Id.String(),
			Title: cats[i].Title,
		}
	}

	return &category.CatIdsResp{Data: results}, nil
}

func (s *service) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))

	if err != nil {
		return err
	}

	// Create a gRPC server object
	grpcServ := grpc.NewServer()
	// Attach the Greeter service to the server

	category.RegisterCategoryServer(grpcServ, NewCategoryServer(s.sctx))

	log.Println(fmt.Sprintf("Serving gRPC on 0.0.0.0:%d", s.port))
	log.Fatal(grpcServ.Serve(lis))

	return nil
}
