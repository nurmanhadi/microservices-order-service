package handler

import (
	"order-service/internal/dto"
	"order-service/internal/service"
	"order-service/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	AddOrder(ctx *gin.Context)
}
type orderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) OrderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with items
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body dto.OrderAddRequest true "Order Request"
// @Success 201
// @Failure 400
// @Router /orders [post]
func (h *orderHandler) AddOrder(ctx *gin.Context) {
	request := new(dto.OrderAddRequest)
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.Error(response.Except(400, "failed to parse json"))
		return
	}
	err := h.orderService.AddOrder(*request)
	if err != nil {
		ctx.Error(err)
		return
	}
	response.Success(ctx, 201, "OK")
}
