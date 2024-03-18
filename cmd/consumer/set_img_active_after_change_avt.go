package consumer

import (
	"context"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/component/gormc"
	"log"
	"my-app/common"
	"my-app/common/pubsub"
	"my-app/component"
	"my-app/module/image"
)

var SetImgActAfterChangeAvtCmd = &cobra.Command{
	Use:   "SetImgActiveAfterChangeAvt",
	Short: "Start consumer: SetImgActiveAfterChangeAvt",
	Run: func(cmd *cobra.Command, args []string) {
		service := sctx.NewServiceContext(
			sctx.WithName("SetImgActiveAfterChangeAvt"),
			sctx.WithComponent(gormc.NewGormDB(common.KeyGorm, "")),
			sctx.WithComponent(component.NewNATSComponent(common.KeyNATS)),
		)

		if err := service.Load(); err != nil {
			log.Fatalln(err)
		}

		ps := service.MustGet(common.KeyNATS).(pubsub.PubSub)

		//"UserChangedAvatar"
		ctx := context.Background()
		ch, _ := ps.Subscribe(ctx, common.TopicUserChangedAvt)

		dbCtx := service.MustGet(common.KeyGorm).(common.DbContext)

		for msg := range ch {
			mapData := msg.Data()
			imgId := uuid.MustParse(mapData["img_id"].(string))

			repo := image.NewRepo(dbCtx.GetDB())

			if err := repo.SetImageStatusActivated(ctx, imgId); err != nil {
				log.Println(err)
			}
		}
	},
}
