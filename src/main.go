package main

import (
	"assignment-runner-base/helpers"
	"fmt"
	"net/http"
	"strings"
	"time"

	"os/exec"

	"github.com/gin-gonic/gin"
)

type ExecWithAssignmentBaseRequest struct {
	SubmissionPath string `json:"submissionPath" binding:"required"`
	AssignmentPath string `json:"assignmentPath" binding:"required"`
}

type ExecStandaloneDirectoryRequest struct {
	Path string `json:"path" binding:"required"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/exec-with-assignment-base", func(c *gin.Context) {
		execRequest := ExecWithAssignmentBaseRequest{}
		if err := c.ShouldBindJSON(&execRequest); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		var err error
		timestamp := time.Now().Unix()
		workdir1 := fmt.Sprintf("/tmp/workdir1-%d", timestamp)
		workdir2 := fmt.Sprintf("/tmp/workdir2-%d", timestamp)

		// clean $WORKDIR1 and $WORKDIR2
		err = exec.Command("mkdir", "-p", workdir1).Run()
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		err = exec.Command("mkdir", "-p", workdir2).Run()
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		err = cleanUpPath(fmt.Sprintf("%s/*", workdir1))
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		err = cleanUpPath(workdir2)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		// unzip assignment from assignmentPath to $WORKDIR1
		err = exec.Command("unzip", "-d", workdir1, execRequest.AssignmentPath).Run()
		if err != nil {
			cleanUpPath(workdir1)
			cleanUpPath(workdir2)
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		// copy submissionPath to $WORKDIR2
		err = exec.Command("cp", "-r", execRequest.SubmissionPath, workdir2).Run()
		if err != nil {
			cleanUpPath(workdir1)
			cleanUpPath(workdir2)
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		// remove all unwanted file
		err = exec.Command("find", workdir2, "-type", "f", "-regex", ".*/.*_test\\.go", "-exec", "rm", "{}", ";").Run()
		if err != nil {
			cleanUpPath(workdir1)
			cleanUpPath(workdir2)
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		err = exec.Command("find", workdir2, "-type", "f", "-regex", ".*/.*\\.test\\.js", "-exec", "rm", "{}", ";").Run()
		if err != nil {
			cleanUpPath(workdir1)
			cleanUpPath(workdir2)
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		err = exec.Command("find", workdir2, "-type", "f", "-regex", ".*/assignment-config.json", "-exec", "rm", "{}", ";").Run()
		if err != nil {
			cleanUpPath(workdir1)
			cleanUpPath(workdir2)
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		// cp -r $WORKDIR2/* $WORKDIR1
		err = exec.Command("/bin/sh", "-c", fmt.Sprintf("cp -r %s/* %s", workdir2, workdir1)).Run()
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		assignmentConf, err := helpers.ReadAssignmentConfigFromDirectory(workdir1)
		if err != nil {
			cleanUpPath(workdir1)
			cleanUpPath(workdir2)
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		commandSplitted := strings.Split(assignmentConf.Command, " ")
		cmd := exec.Command(commandSplitted[0], commandSplitted[1:]...)
		cmd.Dir = workdir1
		stdoutStderr, _ := cmd.CombinedOutput()

		cleanUpPath(workdir1)
		cleanUpPath(workdir2)

		c.String(http.StatusOK, string(stdoutStderr))
	})

	r.POST("/exec-standalone-directory", func(c *gin.Context) {
		execRequest := ExecStandaloneDirectoryRequest{}
		if err := c.ShouldBindJSON(&execRequest); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		var err error
		timestamp := time.Now().Unix()
		workdir := fmt.Sprintf("/tmp/workdir-%d", timestamp)

		err = exec.Command("mkdir", "-p", workdir).Run()
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		err = cleanUpPath(workdir)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		err = exec.Command("cp", "-r", execRequest.Path, workdir).Run()
		if err != nil {
			cleanUpPath(workdir)
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		assignmentConf, err := helpers.ReadAssignmentConfigFromDirectory(workdir)
		if err != nil {
			cleanUpPath(workdir)
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		commandSplitted := strings.Split(assignmentConf.Command, " ")
		cmd := exec.Command(commandSplitted[0], commandSplitted[1:]...)
		cmd.Dir = workdir
		stdoutStderr, _ := cmd.CombinedOutput()
		cleanUpPath(workdir)

		c.String(http.StatusOK, string(stdoutStderr))
	})

	return r
}

func cleanUpPath(path string) error {
	return exec.Command("rm", "-rf", path).Run()
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
