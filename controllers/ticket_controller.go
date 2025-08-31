package controllers

import (
	"net/http"
	"strconv"

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

// GET /tickets
func (c *TicketController) GetMyTickets(ctx *gin.Context) {
	userIDVal, _ := ctx.Get("user_id")
	userID := userIDVal.(uint)

	tickets, err := c.ticketService.GetTicketsByUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tickets)
}

// POST /tickets/cancel
// POST /tickets/cancel/:id
func (c *TicketController) CancelTicket(ctx *gin.Context) {
	userIDVal, _ := ctx.Get("user_id") // dari JWT middleware
	userID := userIDVal.(uint)
	role := ctx.GetString("role") // tambahkan ambil role dari JWT

	idParam := ctx.Param("id")
	ticketID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	if err := c.ticketService.CancelTicket(uint(ticketID64), userID, role); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ticket cancelled successfully"})
}
