package image

import (
	"context"
	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"my-app/common"
	"net/http"
)

type httpService struct {
	serviceCtx sctx.ServiceContext
}

func NewHTTPService(serviceCtx sctx.ServiceContext) httpService {
	return httpService{serviceCtx: serviceCtx}
}

func (s httpService) handleUploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		f, err := c.FormFile("file")

		if err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		file, err := f.Open()

		if err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		defer file.Close()

		fileData := make([]byte, f.Size)

		if _, err := file.Read(fileData); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		dto := UploadDTO{
			Name:     c.PostForm("name"),
			FileName: f.Filename,
			FileType: http.DetectContentType(fileData), //  f.Header.Get("Content-Type")
			FileSize: int(f.Size),
			FileData: fileData,
		}

		uploader := s.serviceCtx.MustGet(common.KeyAWSS3).(Uploader)
		dbContext := s.serviceCtx.MustGet(common.KeyGorm).(common.DbContext)

		uc := NewUseCase(uploader, NewRepo(dbContext.GetDB()))

		media, err := uc.UploadImage(c.Request.Context(), dto)

		if err != nil {
			common.WriteErrorResponse(c, err)
			return
		}

		media.SetCDNDomain(uploader.GetDomain())

		c.JSON(http.StatusOK, core.ResponseData(media))
	}
}

type mockImageRepo struct{}

func (mockImageRepo) Create(ctx context.Context, entity *Image) error {
	return nil
}

func (s httpService) Routes(group *gin.RouterGroup) {
	group.POST("/upload-image", s.handleUploadImage())
}
