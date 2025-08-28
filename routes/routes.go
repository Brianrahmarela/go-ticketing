package routes

import (
	"go-ticketing/controllers"
	"go-ticketing/middleware"
	"go-ticketing/repositories"
	"go-ticketing/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes daftar semua endpoint
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api/v1")

	// auth routes (registrasi, login, dsb.)
	SetupAuthRoutes(api, db)

	// --------- Protected routes (butuh JWT) ---------
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// === Ticket setup ===
	ticketRepo := repositories.NewTicketRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	ticketService := services.NewTicketService(db, ticketRepo, eventRepo)
	ticketCtrl := controllers.NewTicketController(ticketService)

	protected.POST("/tickets", ticketCtrl.BuyTicket)
	protected.DELETE("/tickets/:id", ticketCtrl.CancelTicket)
	// protected.GET("/tickets", ticketCtrl.GetMyTickets)

	// === Event setup ===
	eventRepo2 := repositories.NewEventRepository(db)
	eventService := services.NewEventService(eventRepo2)

	eventCtrl := controllers.NewEventController(eventService)

	// hanya admin yang boleh CRUD event
	admin := protected.Group("/events")
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.POST("", eventCtrl.Create)
		admin.PUT("/:id", eventCtrl.Update)
		admin.DELETE("/:id", eventCtrl.Delete)
	}

	// user biasa bisa lihat daftar event
	protected.GET("/events", eventCtrl.GetAll)
	protected.GET("/events/:id", eventCtrl.GetByID)
}
