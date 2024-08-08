package commands

import (
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dwiandhikaap/rawdog-md/global"
	"github.com/dwiandhikaap/rawdog-md/helper"
	"github.com/dwiandhikaap/rawdog-md/internal"

	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"
)

var watcherServer = internal.NewWatcherServer()

var pathStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00b0ff"))

type WatcherCallbacks struct {
	Write  func(string, *internal.Project) error
	Create func(string, *internal.Project) error
	Remove func(string, *internal.Project) error
	Rename func(string, *internal.Project) error
}

func Watch(relativePath string, port int) error {
	rootAbs, err := filepath.Abs(relativePath)
	if err != nil {
		return err
	}

	rootAbs = strings.ReplaceAll(rootAbs, "\\", "/")

	config := global.ConfigType{
		RootRelativePath: relativePath,
		RootAbsolutePath: rootAbs,
		BuildMode:        global.Development,
	}
	global.SetGlobalConfig(config)

	// Watch for changes in the './pages' and './templates' directories
	pageWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	pageWatcherCallback := WatcherCallbacks{
		Write:  pageWriteCallback,
		Create: pageCreateCallback,
		Remove: pageRemoveCallback,
		Rename: pageRenameCallback,
	}
	defer pageWatcher.Close()

	// Watch for changes in the './assets' directory
	assetWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	assetWatcherCallback := WatcherCallbacks{
		Write:  assetCallback,
		Create: assetCallback,
		Remove: assetCallback,
		Rename: assetCallback,
	}
	defer assetWatcher.Close()

	project, err := internal.NewProject()
	if err != nil {
		log.Fatal(err)
	}

	err = project.ForceRebuild()
	if err != nil {
		log.Println(err)
	}

	go runWatcher(pageWatcher, pageWatcherCallback, project)

	err = registerWatcherWalk(pageWatcher, filepath.Join(relativePath, "pages"))
	if err != nil {
		log.Fatal(err)
	}

	err = registerWatcherWalk(pageWatcher, filepath.Join(relativePath, "templates"))
	if err != nil {
		log.Fatal(err)
	}

	go runWatcher(assetWatcher, assetWatcherCallback, project)

	err = registerWatcherWalk(assetWatcher, filepath.Join(relativePath, "static"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press Ctrl+C to stop.")

	go internal.Serve(watcherServer, filepath.Join(relativePath, "build"), port)

	<-make(chan int)

	fmt.Println("Shutting down...")

	return nil
}

// TODO: delay rebuild for file rename, create, remove events to prevent multiple unexpected rebuilds
func runWatcher(w *fsnotify.Watcher, cb WatcherCallbacks, project *internal.Project) {
	mu := sync.Mutex{}
	timers := make(map[string]*time.Timer)

	for {

		select {
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			fmt.Println("Error:", err)

		case event, ok := <-w.Events:
			if !ok {
				return
			}

			mu.Lock()
			t, ok := timers[event.Name]
			mu.Unlock()

			if !ok {
				t = time.AfterFunc(math.MaxInt64, func() {
					mu.Lock()
					delete(timers, event.Name)
					mu.Unlock()

					time.Sleep(10 * time.Millisecond)

					var err error = nil
					if event.Op&fsnotify.Write == fsnotify.Write {
						// If there is a new directory, add it to the watcher
						if helper.IsPathDir(event.Name) {
							err := registerWatcherWalk(w, event.Name)
							if err != nil {
								log.Fatal(err)
							}
						}

						err = cb.Write(event.Name, project)
					}

					if event.Op&fsnotify.Create == fsnotify.Create {
						err = cb.Create(event.Name, project)
					}

					if event.Op&fsnotify.Remove == fsnotify.Remove {
						err = cb.Remove(event.Name, project)
					}

					if event.Op&fsnotify.Rename == fsnotify.Rename {
						err = cb.Rename(event.Name, project)
					}

					if err != nil {
						fmt.Println(err)
					}
				})
				t.Stop()

				mu.Lock()
				timers[event.Name] = t
				mu.Unlock()
			}

			t.Reset(100 * time.Millisecond)
		}
	}
}

func registerWatcherWalk(watcher *fsnotify.Watcher, path string) error {
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			err = watcher.Add(path)
			fmt.Println("Watching", pathStyle.Render(path))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func eventRelativeRoot(eventPath string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	abs := filepath.Join(cwd, eventPath)
	relPath, err := filepath.Rel(global.Config.RootAbsolutePath, abs)
	if err != nil {
		log.Fatal(err)
	}

	return relPath
}

func pageWriteCallback(eventPath string, project *internal.Project) error {
	startTime := time.Now()
	err := project.ForceRebuild()
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		fmt.Println("Modified:", pathStyle.Render(eventRelativeRoot(eventPath)))
		return fmt.Errorf("build error: %v", err)
	}

	fmt.Println("Modified:", pathStyle.Render(eventRelativeRoot(eventPath)), "(rebuild took", durationMs, "ms)")
	watcherServer.Broadcast("reload")

	return nil
}
func pageCreateCallback(eventPath string, project *internal.Project) error {
	startTime := time.Now()
	err := project.ForceRebuild()
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		fmt.Println("Created:", pathStyle.Render(eventRelativeRoot(eventPath)))
		return fmt.Errorf("build error: %v", err)
	}

	fmt.Println("Created:", pathStyle.Render(eventRelativeRoot(eventPath)), "(rebuild took", durationMs, "ms)")
	watcherServer.Broadcast("reload")
	return nil
}
func pageRemoveCallback(eventPath string, project *internal.Project) error {
	startTime := time.Now()
	err := project.ForceRebuild()
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		fmt.Println("Removed:", pathStyle.Render(eventRelativeRoot(eventPath)))
		return fmt.Errorf("build error: %v", err)
	}

	fmt.Println("Removed:", pathStyle.Render(eventRelativeRoot(eventPath)), "(rebuild took", durationMs, "ms)")
	watcherServer.Broadcast("reload")
	return nil
}

func pageRenameCallback(eventPath string, project *internal.Project) error {
	startTime := time.Now()
	err := project.ForceRebuild()
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		fmt.Println("Renamed:", pathStyle.Render(eventRelativeRoot(eventPath)))
		return fmt.Errorf("build error: %v", err)
	}

	fmt.Println("Renamed:", pathStyle.Render(eventRelativeRoot(eventPath)), "(rebuild took", durationMs, "ms)")
	watcherServer.Broadcast("reload")
	return nil
}

func assetCallback(eventPath string, project *internal.Project) error {
	startTime := time.Now()
	err := project.ForceRebuild()
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		fmt.Println("Static file changes", pathStyle.Render(eventRelativeRoot(eventPath)))
		return fmt.Errorf("build error: %v", err)
	}

	fmt.Println("Static file changes", pathStyle.Render(eventRelativeRoot(eventPath)), "(rebuild took", durationMs, "ms)")
	if strings.HasSuffix(eventPath, ".css") {
		log.Println("Reloading CSS")
		watcherServer.Broadcast("refreshcss")
	} else {
		watcherServer.Broadcast("reload")
	}

	return nil
}
