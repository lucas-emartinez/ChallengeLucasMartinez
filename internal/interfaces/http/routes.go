package http

import (
	"ChallengeUALA/internal/interfaces/http/handlers"

	"ChallengeUALA/internal/application/services"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configura las rutas de nuestra app
func SetupRoutes(
	app *fiber.App,
	tweetService *services.TweetService,
	followService *services.FollowService,
	timelineService *services.TimelineService,
) {

	tweetHandler := handlers.NewTweetHandler(tweetService)
	followHandler := handlers.NewFollowHandler(followService)
	timelineHandler := handlers.NewTimelineHandler(timelineService)

	router := app.Group("/api")
	router.Post("/tweets", tweetHandler.PostTweet)
	router.Post("/follow", followHandler.Follow)
	router.Get("/timeline/:userID", timelineHandler.GetTimeline)
}
