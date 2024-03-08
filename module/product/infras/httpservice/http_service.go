package httpservice

import (
	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"my-app/common"
	"my-app/module/product/query"
	"net/http"
)

type httpService struct {
	sctx sctx.ServiceContext
}

func NewHttpService(sctx sctx.ServiceContext) httpService {
	return httpService{sctx: sctx}
}

func (s httpService) handleListProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param query.ListProductParam

		if err := c.Bind(&param); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		result, err := query.NewListProductQuery(s.sctx).Execute(c.Request.Context(), &param)

		if err != nil {
			common.WriteErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, core.SuccessResponse(result, param.Paging, param.ListProductFilter))
	}
}

func (s httpService) Routes(g *gin.RouterGroup) {
	products := g.Group("products")
	{
		products.GET("", s.handleListProduct())
	}

}
