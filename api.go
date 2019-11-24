package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"

	//"google.golang.org/appengine/log"
	"log"
	"net/http"
)

func addAPIHandlers(r *gin.RouterGroup, db *sql.DB) {
	h := apiHandlers{db}
	r.GET("/user/:id/orders", h.getUserOrders)
	r.GET("/order/:id", h.getOrderDetails)
	r.POST("/order", h.postOrder)
}

type apiHandlers struct {
	db *sql.DB
}

func (h apiHandlers) getUserOrders(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	orders, err := getUserOrders(c.Request.Context(), h.db, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h apiHandlers) getOrderDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	order, err := getOrder(c.Request.Context(), h.db, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if order == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h apiHandlers) postOrder(c *gin.Context) {
	var order struct {
		UserID    int `json:"user_id" binding:"required"`
		ProductID int `json:"product_id" binding:"required"`
	}
	if err := c.BindJSON(&order); err != nil {
		log.Println(err)
		return
	}

	orderID, err := createOrder(c.Request.Context(), h.db, &Order{UserID: order.UserID, ProductID: order.ProductID})
	if err != nil {
		err := errors.Wrap(err, "failed to create order")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": orderID})
}

