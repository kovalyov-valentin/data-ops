package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kovalyov-valentin/data-ops/internal/config"
	"github.com/kovalyov-valentin/data-ops/internal/http/v1/goods"
	"github.com/kovalyov-valentin/data-ops/internal/services"
	"log/slog"
)

type Handler struct {
	config   *config.Config
	services *services.Service
	log      *slog.Logger
}

func NewHandler(config *config.Config, services *services.Service) *Handler {
	return &Handler{
		config:   config,
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	r := router.Group("/v1")

	goodsHandle := goods.New(h.log, h.services.Goods)
	good := r.Group("/good")
	{
		good.POST("/create", goodsHandle.Create)
		good.PATCH("/update", goodsHandle.Update)
		good.GET("/", goodsHandle.GetItem)
		good.DELETE("/remove", goodsHandle.Delete)
	}
	r.GET("/goods/list", goodsHandle.GetItems)

	return router
}
