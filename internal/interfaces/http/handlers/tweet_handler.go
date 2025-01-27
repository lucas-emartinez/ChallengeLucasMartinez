package handlers

import (
	"fmt"
	"net/http"

	"ChallengeUALA/internal/application/services"

	"github.com/gofiber/fiber/v2"
)

type TweetHandler struct {
	tweetService *services.TweetService
}

func NewTweetHandler(tweetService *services.TweetService) *TweetHandler {
	return &TweetHandler{
		tweetService: tweetService,
	}
}

func (h *TweetHandler) PostTweet(c *fiber.Ctx) error {
	var request struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}

	// Algunos controles
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(request.Content) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content is required",
		})
	}

	if len(request.Content) > 140 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content is too long",
		})
	}

	if err := h.tweetService.PostTweet(c.Context(), request.UserID, request.Content); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Errorf("error posting tweet: %w", err),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Tweet posted successfully",
	})
}
