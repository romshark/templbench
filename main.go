package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	fEnvVars    FlagMulti
	fVegetaArgs FlagMulti
	fSrvSrcPath = flag.String("run", "./cmd/std", "Server source path. "+
		"When empty the server is not built and started, "+
		"instead the server is expected to already be running on '-host'")
	fMethod    = flag.String("method", "GET", "HTTP request method")
	fScheme    = flag.String("scheme", "http", "Templ server host address scheme")
	fHost      = flag.String("host", "127.0.0.1:9090", "Templ server host address")
	fPath      = flag.String("path", "/helloworld", "URL path")
	fPingDelay = flag.Duration("ping-delay", 100*time.Millisecond, "Server ping delay")
)

func main() {
	flag.Var(&fVegetaArgs, "veg", "Vegeta CLI arguments, "+
		"see https://github.com/tsenart/vegeta?tab=readme-ov-file#usage-manual "+
		"(multiple values accepted, e.g. '-veg rate=0 -veg duration=5s')")
	flag.Var(&fEnvVars, "env", "Environment variables "+
		"(multiple values accepted, e.g. '-env GOMAXPROCS=1 -env FOO=BAR')")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatal("ERR: ", err.Error())
	}
}

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	serverURL := fmt.Sprintf("%s://%s%s", *fScheme, *fHost, *fPath)

	var serverCmd *exec.Cmd
	if *fSrvSrcPath != "" {
		log.Print("Building server")
		execPath, err := buildServer(ctx)
		if err != nil {
			return fmt.Errorf("building server executable: %w", err)
		}
		defer func() {
			// Clean up executable.
			_ = os.Remove(execPath)
		}()

		serverCmd = exec.CommandContext(ctx, execPath, "-host", *fHost)
		serverCmd.Env = make([]string, len(fEnvVars))
		copy(serverCmd.Env, fEnvVars)
		log.Printf("Running server with env: %v", serverCmd.Env)
		serverCmd.Stdout, serverCmd.Stderr = os.Stdout, os.Stderr

		if err := serverCmd.Start(); err != nil {
			return fmt.Errorf("starting server: %w", err)
		}
		log.Printf("Server PID: %d", serverCmd.Process.Pid)
	} else {
		log.Printf("No server specified, assuming the server is running on %s", serverURL)
	}

	// Wait for server to initialize.
	if err := waitForServer(
		ctx, http.Client{Timeout: 1 * time.Second}, *fMethod, serverURL,
	); err != nil {
		return fmt.Errorf("pinging server: %w", err)
	}
	log.Print("Server OK")

	if err := runVegetaBenchmark(ctx, *fMethod, serverURL); err != nil {
		return fmt.Errorf("running benchmark: %w", err)
	}

	cancel()
	if serverCmd == nil {
		return nil // Server process wasn't started by this process.
	}
	if err := serverCmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
				// Process was killed by a signal (e.g., SIGKILL or SIGTERM).
				// This is fine, ignore this error.
				if s := status.Signal(); s == syscall.SIGKILL || s == syscall.SIGTERM {
					return nil
				}
			}
		}
		return fmt.Errorf("server exited with error: %w", err)
	}
	return nil
}

func buildServer(ctx context.Context) (execPath string, err error) {
	tempDir := os.TempDir()
	now := time.Now()
	const timeFormat = "2006_01_02_15_04_05_999999999Z07_00"
	serverExecFileName := fmt.Sprintf("./server_%s", now.Format(timeFormat))
	p := filepath.Join(tempDir, serverExecFileName)
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", p, *fSrvSrcPath)
	buildCmd.Stdout, buildCmd.Stderr = os.Stdout, os.Stderr
	if err := buildCmd.Run(); err != nil {
		return "", fmt.Errorf("executing go build: %w", err)
	}
	return p, nil
}

const pingRetries = 10

// waitForServer attempts to ping the server until
// it returns a 200 OK or retries are exhausted.
func waitForServer(ctx context.Context, c http.Client, method, url string) error {
	for i := 0; i < pingRetries; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		req, err := http.NewRequest(method, url, http.NoBody)
		if err != nil {
			return fmt.Errorf("initializing ping request: %w", err)
		}
		req = req.WithContext(ctx)

		resp, err := c.Do(req)
		if err != nil {
			time.Sleep(*fPingDelay)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		time.Sleep(*fPingDelay)
	}
	return fmt.Errorf("server didn't return 200 OK after %d retries", pingRetries)
}

// runVegetaBenchmark runs Vegeta benchmark against the provided URL.
func runVegetaBenchmark(ctx context.Context, method, url string) error {
	// Prepare Vegeta attack command
	args := []string{"run", "github.com/tsenart/vegeta", "attack"}
	for _, arg := range fVegetaArgs {
		args = append(args, "-"+arg)
	}

	attackCmd := exec.CommandContext(ctx, "go", args...)
	attackCmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%s %s", method, url))

	resultBuffer := new(bytes.Buffer)
	attackCmd.Stdout, attackCmd.Stderr = resultBuffer, os.Stderr

	log.Printf("Running benchmark with args: %v\n", args[2:])
	if err := attackCmd.Run(); err != nil {
		return fmt.Errorf("vegeta attack failed: %w", err)
	}

	// Run Vegeta report
	reportCmd := exec.Command("go", "run", "github.com/tsenart/vegeta", "report")
	reportCmd.Stdin = resultBuffer
	reportCmd.Stdout, reportCmd.Stderr = os.Stdout, os.Stderr

	log.Print("Generating report...")
	if err := reportCmd.Run(); err != nil {
		return fmt.Errorf("vegeta report failed: %w", err)
	}

	return nil
}

type FlagMulti []string

func (i *FlagMulti) String() string { return fmt.Sprintf("%v", *i) }
func (i *FlagMulti) Set(value string) error {
	*i = append(*i, value)
	return nil
}
