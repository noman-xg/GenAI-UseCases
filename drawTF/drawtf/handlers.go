package drawtf

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Text string `json:"text"`
	Path string `json:"path"`
}

func GenerateEmbeddings(c *gin.Context) {
	var m Message
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	response, err := fetchConfigFromVecStore(m.Text, m.Path)
	if err != nil {
		fmt.Println("error is", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
}

func GenerateTfconfig(c *gin.Context) {
	var req json.RawMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	refinementResponse, err := integratedHandler(req)

	if err != nil {
		fmt.Println("Error getting response from VectorStore Embeddings", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": refinementResponse})
}

func StartGradio(c *gin.Context) {
	err := launchGradio()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err})
	}
}
