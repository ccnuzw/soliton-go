package handlers

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

// RunGoModTidyRequest is the request body for running go mod tidy.
type RunGoModTidyRequest struct {
	ProjectPath string `json:"project_path" binding:"required"`
}

// RunGoModTidy handles POST /api/projects/tidy
// Runs go mod tidy in the specified project directory
func RunGoModTidy(c *gin.Context) {
	var req RunGoModTidyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate project path
	if _, err := os.Stat(req.ProjectPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project path does not exist"})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Run go mod tidy
	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Dir = req.ProjectPath
	cmd.Env = append(os.Environ(), "GOWORK=off")

	// Capture output
	output := ""
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create stdout pipe"})
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create stderr pipe"})
		return
	}

	if err := cmd.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("failed to start go mod tidy: %v", err),
		})
		return
	}

	// Read stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			output += scanner.Text() + "\n"
		}
	}()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			output += scanner.Text() + "\n"
		}
	}()

	// Wait for command to finish
	err = cmd.Wait()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   fmt.Sprintf("go mod tidy failed: %v", err),
			"output":  output,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dependencies downloaded successfully",
		"output":  output,
	})
}
