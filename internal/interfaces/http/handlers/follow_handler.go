package handlers

import (
	"fmt"
	"net/http"

	"ChallengeUALA/internal/application/services"

	"github.com/gofiber/fiber/v2"
)

type FollowHandler struct {
	followService *services.FollowService
}

func NewFollowHandler(followService *services.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

func (h *FollowHandler) Follow(c *fiber.Ctx) error {
	var request struct {
		FollowerID string `json:"follower_id"`
		FolloweeID string `json:"followee_id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.followService.Follow(c.Context(), request.FollowerID, request.FolloweeID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Errorf("error following user: %w", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User followed successfully",
	})
}
