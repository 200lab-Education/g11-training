package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"my-app/builder"
	"my-app/component"
	"my-app/module/product/controller"
	productusecase "my-app/module/product/domain/usecase"
	productmysql "my-app/module/product/repository/mysql"
	"my-app/module/user/infras/httpservice"
	"my-app/module/user/usecase"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db = db.Debug()

	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Setup dependencies
	repo := productmysql.NewMysqlRepository(db)
	useCase := productusecase.NewCreateProductUseCase(repo)
	api := controller.NewAPIController(useCase)

	v1 := r.Group("/v1")
	{
		products := v1.Group("/products")
		{
			products.POST("", api.CreateProductAPI(db))
		}

	}

	jwtSecret := os.Getenv("JWT_SECRET")

	tokenProvider := component.NewJWTProvider(jwtSecret,
		60*60*24*7, 60*60*24*14)

	//userUC := usecase.NewUseCase(repository.NewUserRepo(db), &common.Hasher{}, tokenProvider, repository.NewSessionMySQLRepo(db))

	userUseCase := usecase.UseCaseWithBuilder(builder.NewComplexBuilder(builder.NewSimpleBuilder(db, tokenProvider)))

	httpservice.NewUserService(userUseCase).Routes(v1)

	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//type mockSessionRepo struct {
//}
//
//func (m mockSessionRepo) Create(ctx context.Context, data *domain.Session) error {
//	return nil
//}
