package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"os/exec"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/exec-be", func(c *gin.Context) {
		path := c.DefaultQuery("path", "/opt/playground/backend/unittest")
		cmd := exec.Command("go", "test", "-v")
		cmd.Dir = path
		stdoutStderr, _ := cmd.CombinedOutput()

		c.String(http.StatusOK, string(stdoutStderr))
	})

	r.GET("/exec-fe", func(c *gin.Context) {
		path := c.DefaultQuery("path", "/opt/playground/frontend/unittest")
		cmd := exec.Command("npm", "run", "test", "--prefix", "/opt/playground", path, "--", "--watchAll=false")
		stdoutStderr, _ := cmd.CombinedOutput()

		c.String(http.StatusOK, string(stdoutStderr))
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
