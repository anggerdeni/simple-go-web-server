package main

import (
	"net/http"
	"strings"

	"os/exec"

	"github.com/gin-gonic/gin"
)

type ExecRequest struct {
	Command string `json:"command" binding:"required"`
	Path    string `json:"path" binding:"required"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/exec", func(c *gin.Context) {
		execRequest := ExecRequest{}
		if err := c.ShouldBindJSON(&execRequest); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		commandSplitted := strings.Split(execRequest.Command, " ")
		cmd := exec.Command(commandSplitted[0], commandSplitted[1:]...)
		cmd.Dir = execRequest.Path
		stdoutStderr, _ := cmd.CombinedOutput()

		c.String(http.StatusOK, string(stdoutStderr))
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
