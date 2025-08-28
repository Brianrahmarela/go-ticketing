package controllers

import (
	"net/http"
	"strconv"

	"go-ticketing/models"
	"go-ticketing/services"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	service services.EventService
}

func NewEventController(service services.EventService) *EventController {
	return &EventController{service}
}

func (ec *EventController) Create(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ec.service.Create(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// hitung remaining
	event.Remaining = int(event.Capacity) - int(event.Sold)

	c.JSON(http.StatusOK, event)
}

func (ec *EventController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event.ID = uint(id)

	if err := ec.service.Update(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	event.Remaining = int(event.Capacity) - int(event.Sold)

	c.JSON(http.StatusOK, event)
}

func (ec *EventController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := ec.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "event deleted"})
}

func (ec *EventController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	event, err := ec.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	event.Remaining = int(event.Capacity) - int(event.Sold)

	c.JSON(http.StatusOK, event)
}

func (ec *EventController) GetAll(c *gin.Context) {
	events, err := ec.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// hitung remaining tiap event
	for i := range events {
		events[i].Remaining = int(events[i].Capacity) - int(events[i].Sold)
	}

	c.JSON(http.StatusOK, events)
}
