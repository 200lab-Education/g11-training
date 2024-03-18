package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/component/gormc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"my-app/builder"
	"my-app/cmd/consumer"
	"my-app/common"
	"my-app/common/pubsub"
	"my-app/component"
	"my-app/middleware"
	"my-app/module/category/infras/grpcservice"
	catHTTP "my-app/module/category/infras/httpservice"
	"my-app/module/image"
	"my-app/module/product/controller"
	productHTTP "my-app/module/product/infras/httpservice"
	productmysql "my-app/module/product/repository/mysql"
	productusecase "my-app/module/product/usecase"
	"my-app/module/user/infras/httpservice"
	"my-app/module/user/infras/repository"
	"my-app/module/user/usecase"
	sub "my-app/sub"
	"net/http"
	"os"
)

func newService() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName("G11"),
		sctx.WithComponent(gormc.NewGormDB(common.KeyGorm, "")),
		sctx.WithComponent(component.NewNATSComponent(common.KeyNATS)),
		sctx.WithComponent(pubsub.NewLocalPubSub(common.KeyLocalPS)),
		sctx.WithComponent(component.NewJWT(common.KeyJWT)),
		sctx.WithComponent(component.NewAWSS3Provider(common.KeyAWSS3)),
		sctx.WithComponent(component.NewConfigComponent(common.KeyConfig)),
	)
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start main service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()

		if err := service.Load(); err != nil {
			log.Fatalln(err)
		}

		db := service.MustGet(common.KeyGorm).(common.DbContext).GetDB()

		r := gin.Default()

		r.Use(middleware.Recovery())

		tokenProvider := service.MustGet(common.KeyJWT).(component.TokenProvider)

		authClient := usecase.NewIntrospectUC(repository.NewUserRepo(db), repository.NewSessionMySQLRepo(db), tokenProvider)

		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		r.DELETE("/v1/revoke-token", middleware.RequireAuth(authClient), func(c *gin.Context) {
			requester := c.MustGet(common.KeyRequester).(common.Requester)

			repo := repository.NewSessionMySQLRepo(db)

			if err := repo.Delete(c.Request.Context(), requester.TokenId()); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"data": true})
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

		//userUC := usecase.NewUseCase(repository.NewUserRepo(db), &common.Hasher{}, tokenProvider, repository.NewSessionMySQLRepo(db))

		userUseCase := usecase.UseCaseWithBuilder(builder.NewSimpleBuilder(db, tokenProvider))

		httpservice.NewUserService(userUseCase, service).SetAuthClient(authClient).Routes(v1)
		image.NewHTTPService(service).Routes(v1)
		catHTTP.NewCategoryHttpService(service).Routes(v1)

		go func() {
			config := service.MustGet(common.KeyConfig).(interface{ GetGRPCServPort() int })
			_ = grpcservice.NewCatGRPCService(config.GetGRPCServPort(), service).Start()
		}()

		opts := grpc.WithTransportCredentials(insecure.NewCredentials())

		config := service.MustGet(common.KeyConfig).(interface{ GetGRPCServer() string })

		cc, err := grpc.Dial(config.GetGRPCServer(), opts)

		if err != nil {
			log.Fatal(err)
		}

		prodService := productHTTP.NewHttpService(service)
		prodService.SetGRPCClientConn(cc)
		prodService.Routes(v1)

		go sub.NewTopicUserChangeAvt(service, service.MustGet(common.KeyLocalPS).(pubsub.PubSub)).Start()

		r.Run(":3000")
	},
}

func Execute() {
	// TransAddPoint outenv as a sub command
	rootCmd.AddCommand(outEnvCmd)

	consumerCmd := &cobra.Command{Use: "consumer", Short: "Start consumer"}
	consumerCmd.AddCommand(consumer.SetImgActAfterChangeAvtCmd)

	rootCmd.AddCommand(consumerCmd)

	//rootCmd.AddCommand(cronjob)
	//
	//rootCmd.AddCommand(startSubUserLikedRestaurantCmd)
	//rootCmd.AddCommand(startSubUserDislikedRestaurantCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
