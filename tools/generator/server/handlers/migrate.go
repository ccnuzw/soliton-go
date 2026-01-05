package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RunMigrationRequest is the request body for running migrations.
type RunMigrationRequest struct {
	ProjectPath    string `json:"project_path" binding:"required"`
	AutoTidy       bool   `json:"auto_tidy"`
	TimeoutSeconds int    `json:"timeout_seconds"`
}

// MigrationLogEntry is a structured log line for migration output.
type MigrationLogEntry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Step    string `json:"step"`
	Message string `json:"message"`
}

// MigrationResult is the response body for migration runs.
type MigrationResult struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message,omitempty"`
	Logs       []MigrationLogEntry`json:"logs"`
	DurationMs int64              `json:"duration_ms"`
	ExitCode   int                `json:"exit_code"`
	Command    string             `json:"command"`
	StartedAt  string             `json:"started_at"`
	FinishedAt string             `json:"finished_at"`
}

// RunMigration handles POST /api/projects/migrate
func RunMigration(c *gin.Context) {
	var req RunMigrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := os.Stat(req.ProjectPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project path does not exist"})
		return
	}

	timeout := req.TimeoutSeconds
	if timeout <= 0 {
		timeout = 300
	}
	if timeout > 1800 {
		timeout = 1800
	}

	startedAt := time.Now()
	logs := make([]MigrationLogEntry, 0, 128)
	var mu sync.Mutex

	appendLog := func(level, step, message string) {
		mu.Lock()
		logs = append(logs, MigrationLogEntry{
			Time:    time.Now().Format(time.RFC3339),
			Level:   level,
			Step:    step,
			Message: message,
		})
		mu.Unlock()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	appendLog("info", "system", "migration started")
	appendLog("info", "system", fmt.Sprintf("project: %s", req.ProjectPath))

	if req.AutoTidy {
		appendLog("info", "tidy", "running go mod tidy")
		_, err := runCommandWithLogs(ctx, req.ProjectPath, "tidy", []string{"go", "mod", "tidy"}, &logs, &mu)
		if err != nil {
			appendLog("error", "tidy", fmt.Sprintf("go mod tidy failed: %v", err))
		}
	}

	args, display, err := detectMigrationCommand(req.ProjectPath)
	if err != nil {
		appendLog("error", "system", err.Error())
		finishedAt := time.Now()
		c.JSON(http.StatusOK, MigrationResult{
			Success:    false,
			Message:    err.Error(),
			Logs:       logs,
			DurationMs: finishedAt.Sub(startedAt).Milliseconds(),
			ExitCode:   -1,
			Command:    "",
			StartedAt:  startedAt.Format(time.RFC3339),
			FinishedAt: finishedAt.Format(time.RFC3339),
		})
		return
	}

	appendLog("info", "system", fmt.Sprintf("command: %s", display))
	exitCode, cmdErr := runCommandWithLogs(ctx, req.ProjectPath, "migrate", args, &logs, &mu)

	finishedAt := time.Now()
	success := cmdErr == nil
	message := "migration completed"
	if !success {
		message = fmt.Sprintf("migration failed: %v", cmdErr)
		appendLog("error", "system", message)
	} else {
		appendLog("info", "system", message)
	}

	c.JSON(http.StatusOK, MigrationResult{
		Success:    success,
		Message:    message,
		Logs:       logs,
		DurationMs: finishedAt.Sub(startedAt).Milliseconds(),
		ExitCode:   exitCode,
		Command:    display,
		StartedAt:  startedAt.Format(time.RFC3339),
		FinishedAt: finishedAt.Format(time.RFC3339),
	})
}

func detectMigrationCommand(projectPath string) ([]string, string, error) {
	migrateDir := filepath.Join(projectPath, "cmd", "migrate")
	if info, err := os.Stat(migrateDir); err == nil && info.IsDir() {
		return []string{"go", "run", "./cmd/migrate"}, "go run ./cmd/migrate", nil
	}

	legacyFile := filepath.Join(projectPath, "cmd", "migrate.go")
	if info, err := os.Stat(legacyFile); err == nil && !info.IsDir() {
		return []string{"go", "run", "./cmd/migrate.go"}, "go run ./cmd/migrate.go", nil
	}

	return nil, "", fmt.Errorf("migration entry not found: cmd/migrate/main.go or cmd/migrate.go")
}

func runCommandWithLogs(ctx context.Context, dir, step string, args []string, logs *[]MigrationLogEntry, mu *sync.Mutex) (int, error) {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return -1, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return -1, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return -1, fmt.Errorf("failed to start command: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go scanOutput(&wg, stdout, "info", step, logs, mu)
	go scanOutput(&wg, stderr, "error", step, logs, mu)

	err = cmd.Wait()
	wg.Wait()

	exitCode := -1
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return exitCode, fmt.Errorf("command timeout")
		}
		return exitCode, err
	}

	return exitCode, nil
}

func scanOutput(wg *sync.WaitGroup, reader io.ReadCloser, level, step string, logs *[]MigrationLogEntry, mu *sync.Mutex) {
	defer wg.Done()
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		mu.Lock()
		*logs = append(*logs, MigrationLogEntry{
			Time:    time.Now().Format(time.RFC3339),
			Level:   level,
			Step:    step,
			Message: line,
		})
		mu.Unlock()
	}
	if err := scanner.Err(); err != nil {
		mu.Lock()
		*logs = append(*logs, MigrationLogEntry{
			Time:    time.Now().Format(time.RFC3339),
			Level:   "error",
			Step:    step,
			Message: fmt.Sprintf("log scan error: %v", err),
		})
		mu.Unlock()
	}
}
