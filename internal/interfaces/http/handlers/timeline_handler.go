package handlers

import (
	"ChallengeUALA/internal/application/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type TimelineHandler struct {
	timelineService *services.TimelineService
}

func NewTimelineHandler(timelineService *services.TimelineService) *TimelineHandler {
	return &TimelineHandler{
		timelineService: timelineService,
	}
}

func (h *TimelineHandler) GetTimeline(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	// Obtener el timeline del usuario usando el servicio
	timeline, err := h.timelineService.GetTimeline(c.Context(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("error getting timeline: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(timeline)
}
