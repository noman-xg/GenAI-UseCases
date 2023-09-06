package main

// Import necessary packages for the application.
import (
	"github.com/gin-gonic/gin"
	"github.com/noman-xg/GenAI-UseCases/drawtf"
)

func main() {
	// Create a new Gin router with default middleware.
	r := gin.Default()

	// Define HTTP routes and associate them with specific handler functions.
	// These routes are used to handle incoming HTTP POST requests.
	// - "/message" route is associated with the "GenerateEmbeddings" function from the "drawtf" package.
	// - "/tfconfig" route is associated with the "GenerateTfconfig" function from the "drawtf" package.
	// - "/start-gradio" route is associated with the "StartGradio" function from the "drawtf" package.
	r.POST("/message", drawtf.GenerateEmbeddings)
	r.POST("/tfconfig", drawtf.GenerateTfconfig)
	r.POST("/start-gradio", drawtf.StartGradio)

	// Start the Gin HTTP server on port 8082.
	r.Run(":8082")
}
