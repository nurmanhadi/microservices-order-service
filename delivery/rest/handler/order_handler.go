package handler

import (
	"order-service/internal/dto"
	"order-service/internal/service"
	"order-service/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	AddOrder(ctx *gin.Context)
	GetAllOrder(ctx *gin.Context)
	GetOrderByID(ctx *gin.Context)
	GetAllOrderByUserID(ctx *gin.Context)
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

// GetAllOrder godoc
// @Summary Fetch all orders
// @Description Fetch all orders
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200
// @Failure 404
// @Router /orders [get]
func (h *orderHandler) GetAllOrder(ctx *gin.Context) {
	result, err := h.orderService.GetAllOrder()
	if err != nil {
		ctx.Error(err)
		return
	}
	response.Success(ctx, 200, result)
}

// GetOrderByID godoc
// @Summary Get order by id
// @Description Get order by spesific order id and include items
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "order id"
// @Success 200
// @Failure 404
// @Router /orders/{id} [get]
func (h *orderHandler) GetOrderByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := h.orderService.GetOrderByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}
	response.Success(ctx, 200, result)
}

// GetAllOrderByUserID godoc
// @Summary get all order by user id
// @Description Get all order by spesific user id
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200
// @Failure 404
// @Router /orders/users/{id} [get]
func (h *orderHandler) GetAllOrderByUserID(ctx *gin.Context) {
	userID := ctx.Param("id")
	result, err := h.orderService.GetAllOrderByUserID(userID)
	if err != nil {
		ctx.Error(err)
		return
	}
	response.Success(ctx, 200, result)
}
