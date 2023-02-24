package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kargotech/go-testapp/internal/entity"
	"github.com/kargotech/gokargo/serror"
)

func (h Handler) GetOrderByID(ctx *gin.Context) error {
	id := ctx.Param("id")

	order, err := h.uc.GetOrderByID(ctx.Request.Context(), id)
	if err != nil {
		superErr := serror.ConvertErrToSuperErrItf(err)

		ctx.JSON(superErr.HttpStatus(), gin.H{
			"error": superErr.ErrorMessage(),
			"code":  superErr.ErrorCode(),
		})
		return err
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order": gin.H{
			"ksuid":      order.Ksuid,
			"client_id":  order.ClientID,
			"product_id": order.ProductID,
		},
		"success": "true",
	})
	return nil
}

func (h Handler) CreateOrder(ctx *gin.Context) error {

	var order entity.CreateOrder

	err := ctx.BindJSON(&order)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
			"code":  err.Error(),
		})
	}

	createdOrder, err := h.uc.CreateOrder(ctx.Request.Context(), order)
	if err != nil {
		superErr := serror.ConvertErrToSuperErrItf(err)

		ctx.JSON(superErr.HttpStatus(), gin.H{
			"error": superErr.ErrorMessage(),
			"code":  superErr.ErrorCode(),
		})
		return err
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order": gin.H{
			"ksuid": createdOrder.Ksuid,
		},
		"success": "true",
	})
	return nil
}

func (h Handler) UpdateOrder(ctx *gin.Context) error {
	var order entity.UpdateOrder

	err := ctx.BindJSON(&order)
	order.Ksuid = ctx.Param("id")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
			"code":  err.Error(),
		})
	}

	updatedOrder, err := h.uc.UpdateOrder(ctx.Request.Context(), order)
	if err != nil {
		superErr := serror.ConvertErrToSuperErrItf(err)

		ctx.JSON(superErr.HttpStatus(), gin.H{
			"error": superErr.ErrorMessage(),
			"code":  superErr.ErrorCode(),
		})
		return err
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order": gin.H{
			"ksuid": updatedOrder.Ksuid,
		},
		"success": "true",
	})
	return nil
}
