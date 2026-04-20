package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	mu        sync.Mutex
	cmd       *exec.Cmd
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the development server with auto-reload",
	Run: func(cmd *cobra.Command, args []string) {
		watch, _ := cmd.Flags().GetBool("watch")

		if !watch {
			startServer()
			return
		}

		color.Magenta("🚀 Kodia Dev Server starting with --watch...")
		
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		done := make(chan bool)
		
		go func() {
			var timer *time.Timer
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					// Ignore non-source files
					if !isSourceFile(event.Name) {
						continue
					}

					if event.Op&fsnotify.Write == fsnotify.Write {
						// Debounce reloads
						if timer != nil {
							timer.Stop()
						}
						timer = time.AfterFunc(800*time.Millisecond, func() {
							color.Cyan("🔄 Change detected in %s. Reloading...", filepath.Base(event.Name))
							restartServer()
						})
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()

		// Watch backend directories
		watchDirs := []string{
			filepath.Join("backend", "cmd"),
			filepath.Join("backend", "internal"),
			filepath.Join("backend", "pkg"),
		}

		for _, dir := range watchDirs {
			addDirRecursive(watcher, dir)
		}

		startServer()
		<-done
	},
}

func isSourceFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".go" || ext == ".html" || ext == ".tmpl" || ext == ".yaml"
}

func addDirRecursive(watcher *fsnotify.Watcher, path string) {
	filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.Contains(walkPath, "node_modules") || strings.Contains(walkPath, ".git") {
				return filepath.SkipDir
			}
			return watcher.Add(walkPath)
		}
		return nil
	})
}

func startServer() {
	mu.Lock()
	defer mu.Unlock()

	serverPath := filepath.Join("backend", "cmd", "server", "main.go")
	
	// Check OS for appropriate shell
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/c"
	} else {
		shell = "sh"
		flag = "-c"
	}

	cmd = exec.Command(shell, flag, "go run "+serverPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Start(); err != nil {
		color.Red("Failed to start server: %v", err)
		return
	}
}

func restartServer() {
	mu.Lock()
	if cmd != nil && cmd.Process != nil {
		// Kill the process group to ensure children are also killed
		if runtime.GOOS == "windows" {
			exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(cmd.Process.Pid)).Run()
		} else {
			cmd.Process.Kill()
		}
	}
	mu.Unlock()
	
	startServer()
}

func init() {
	serveCmd.Flags().BoolP("watch", "w", true, "Auto-reload on changes")
	rootCmd.AddCommand(serveCmd)
}
