package main

import (
	"github.com/gin-gonic/gin"
	"github.com/noman-xg/GenAI-UseCases/tree/main/drawTF/drawtf"
)

func main() {
	r := gin.Default()

	r.POST("/message", drawtf.GenerateEmbeddings)
	r.POST("/tfconfig", drawtf.GenerateTfconfig)
	r.POST("/start-gradio", drawtf.StartGradio)
	r.Run(":8082")
}
