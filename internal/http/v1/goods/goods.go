package goods

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kovalyov-valentin/data-ops/internal/http/v1/goods/response"
	"github.com/kovalyov-valentin/data-ops/internal/services"
	"log/slog"
	"net/http"
	"strconv"
)

type Handlers struct {
	log     *slog.Logger
	service services.Goods
}

func NewHandlers(log *slog.Logger, repo services.Goods) *Handlers {
	return &Handlers{
		log:     log,
		service: repo,
	}
}

func (h *Handlers) Create(c *gin.Context) {

	projectsId, err := strconv.Atoi(c.Query("projectsId"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var req response.RequestCreate
	err = c.ShouldBindJSON(&req)

	if err != nil || req.Name == "" {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	goods, err := h.service.CreateGoods(c.Request.Context(), req.Name, projectsId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goods)
}

func (h *Handlers) Update(c *gin.Context) {
	goodsId, projectsId, err := h.parseParam(c)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var req response.RequestUpdate
	err = c.ShouldBindJSON(&req)
	if err != nil || req.Name == "" {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	item, err := h.service.UpdateGoods(c.Request.Context(), req.Name, req.Description, goodsId, projectsId)
	if errors.Is(err, errors.New("not found goods")) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "errors.item.notFound", "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handlers) Delete(c *gin.Context) {
	goodsId, projectsId, err := h.parseParam(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"goodsId": goodsId, "projectsId": projectsId, "err": err})
		return
	}

	it, err := h.service.DeleteGoods(c.Request.Context(), goodsId, projectsId)
	if errors.Is(err, errors.New("not found goods")) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "errors.item.notFound", "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ResponseRemove{Removed: it.Removed, Id: goodsId, ProjectsId: projectsId})
}

func (h *Handlers) GetGood(c *gin.Context) {
	goodsId, projectsId, err := h.parseParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	allItems, err := h.service.GetGood(c.Request.Context(), goodsId, projectsId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, allItems)
}

func (h *Handlers) GetGoods(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	var (
		limitInt  = 10
		offsetInt = 1
		err       error
	)

	if limitStr != "" {
		limitInt, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
			return
		}
	}

	if offsetStr != "" {
		offsetInt, err = strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "offset must be an integer"})
			return
		}
	}

	goodsResponse, err := h.service.GetGoods(c.Request.Context(), limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goodsResponse)
}

func (h *Handlers) parseParam(c *gin.Context) (goodsId, projectsId int, err error) {
	projectsId, err = strconv.Atoi(c.Query("projectsId"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	goodsId, err = strconv.Atoi(c.Query("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	return
}
