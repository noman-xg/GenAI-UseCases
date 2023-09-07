package drawtf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Message struct represents a JSON message with 'text' and 'path' fields.
type Message struct {
	Text string `json:"text"`
	Path string `json:"path"`
}

// HTTP handler:  that generates embeddings based on the input message.
func GenerateEmbeddings(c *gin.Context) {
	var m Message
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Fetch configuration from VecStore.
	response, err := fetchConfigFromVecStore(m.Text, m.Path, false)
	if err != nil {
		fmt.Println("error is", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
}

// HTTP handler:  generates a Terraform configuration based on query parameters and JSON request.
func GenerateTfconfig(c *gin.Context) {

	// Retrieve 'path' query parameter.
	docsPath := c.Query("path")

	// Retrieve 'rag' query parameter and parse it as a boolean.

	isRag, err := strconv.ParseBool(c.Query("rag"))
	if err != nil {
		fmt.Println("error is", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Param. Rag should be boolean"})
		return
	}

	// Parse the incoming JSON request.
	var req json.RawMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call an integrated handler function with the parsed data.
	tfConfig, err := configGenerator(req, docsPath, isRag)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": tfConfig})
}

// StartGradio is an HTTP handler that launches Gradio.
func StartGradio(c *gin.Context) {
	err := launchGradio()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err})
	}
}
