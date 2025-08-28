package controllers

import (
	"net/http"

	"go-ticketing/services"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketService services.TicketService
}

func NewTicketController(ticketService services.TicketService) *TicketController {
	return &TicketController{ticketService}
}

type BuyTicketRequest struct {
	// UserID   uint `json:"user_id" binding:"required"`
	EventID  uint `json:"event_id" binding:"required"`
	Quantity uint `json:"quantity" binding:"required"`
}

type CancelTicketRequest struct {
	UserID   uint `json:"user_id" binding:"required"`
	TicketID uint `json:"ticket_id" binding:"required"`
}

// POST /tickets/buy
func (c *TicketController) BuyTicket(ctx *gin.Context) {
	var req BuyTicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil user_id dari JWT claim
	userIDVal, _ := ctx.Get("user_id") // middleware JWT set "user_id"
	userID := userIDVal.(uint)

	//before ambil dari request
	// ticket, err := c.ticketService.BuyTicket(req.UserID, req.EventID, req.Quantity)
	ticket, err := c.ticketService.BuyTicket(userID, req.EventID, req.Quantity)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ticket)
}

// POST /tickets/cancel
func (c *TicketController) CancelTicket(ctx *gin.Context) {
	var req CancelTicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ticketService.CancelTicket(req.TicketID, req.UserID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ticket cancelled"})
}
