package main

import (
	"github.com/gin-gonic/gin"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	r := gin.Default()

	r.POST("/message", generateEmbeddings)
	r.POST("/tfconfig", generateTfconfig)
	r.Run(":8082")
}

// func main() {
// 	r := gin.Default()

// 	r.POST("/message", func(c *gin.Context) {
// 		var m Message
// 		fmt.Println("post called ...")
// 		if err := c.ShouldBindJSON(&m); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		response, err := fetchConfigFromVecStore(m.Text)
// 		if err != nil {
// 			fmt.Println("Error getting response from VectorStore Embeddings", err)
// 			c.JSON(http.StatusBadGateway, gin.H{"error": err})
// 		}

// 		c.JSON(http.StatusOK, gin.H{"response": response})
// 	})

// 	r.POST("/turboIntegrated", func(c *gin.Context) {
// var req json.RawMessage
// if err := c.ShouldBindJSON(&req); err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	return
// }
// // userQuery, err := generateUserQuery(req)
// // response, err := getConfigFromVecStore(userQuery)

// // raw := json.RawMessage([]byte(response))
// // refinementResponse, err := refineEmbeddings(raw)

// refinementResponse, err := integratedHandler(req)

// if err != nil {
// 	fmt.Println("Error getting response from VectorStore Embeddings", err)
// 	c.JSON(http.StatusBadGateway, gin.H{"error": err})
// 	return
// }

// c.JSON(http.StatusOK, gin.H{"response": refinementResponse})
// 	})

// 	r.Run(":8082") // listen and serve on 0.0.0.0:8080

// }
