package httpservice

import (
	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"google.golang.org/grpc"
	"my-app/common"
	"my-app/module/product/query"
	"my-app/module/product/repository/grpcclient"
	"my-app/proto/category"
	"net/http"
)

type httpService struct {
	sctx           sctx.ServiceContext
	grpcClientConn *grpc.ClientConn
}

func NewHttpService(sctx sctx.ServiceContext) *httpService {
	return &httpService{sctx: sctx}
}

func (s *httpService) handleListProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param query.ListProductParam

		if err := c.Bind(&param); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		//configComp := s.sctx.MustGet(common.KeyConfig).(interface{ GetURLRPCCategory() string })
		//urlRPC := fmt.Sprintf("%s/query-categories-ids", configComp.GetURLRPCCategory())
		//catRepo := rpchttp.NewRpcGetCategoriesByIds(urlRPC)

		catRepo := grpcclient.NewCatGRPCClient(category.NewCategoryClient(s.grpcClientConn))

		result, err := query.NewListProductQuery(s.sctx, catRepo).Execute(c.Request.Context(), &param)

		if err != nil {
			common.WriteErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, core.SuccessResponse(result, param.Paging, param.ListProductFilter))
	}
}

func (s *httpService) Routes(g *gin.RouterGroup) {
	products := g.Group("products")
	{
		products.GET("", s.handleListProduct())
	}

}

func (s *httpService) SetGRPCClientConn(grpcClientConn *grpc.ClientConn) {
	s.grpcClientConn = grpcClientConn
}
