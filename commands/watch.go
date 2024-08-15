package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dwiandhikaap/rawdog-md/global"
	"github.com/dwiandhikaap/rawdog-md/internal"
	"github.com/radovskyb/watcher"

	"github.com/charmbracelet/lipgloss"
)

type WatcherCallbacks struct {
	Write  func(string, *internal.Project) error
	Create func(string, *internal.Project) error
	Remove func(string, *internal.Project) error
	Rename func(string, *internal.Project) error
}

var watcherServer = internal.NewWatcherServer()

var pathStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00b0ff"))

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

	// Load user config
	err = global.LoadUserConfig()
	if err != nil {
		return err
	}

	project, err := internal.NewProject()
	if err != nil {
		log.Fatal(err)
	}

	// Watch for changes in the './pages' and './templates' directories
	pageWatcherCallback := WatcherCallbacks{
		Write:  pageWriteCallback,
		Create: pageCreateCallback,
		Remove: pageRemoveCallback,
		Rename: pageRenameCallback,
	}
	go runWatcher(filepath.Join(relativePath, "pages"), pageWatcherCallback, project, false)
	go runWatcher(filepath.Join(relativePath, "templates"), pageWatcherCallback, project, false)

	// Watch for changes in the './static' directory
	assetWatcherCallback := WatcherCallbacks{
		Write:  assetCallback,
		Create: assetCallback,
		Remove: assetCallback,
		Rename: assetCallback,
	}
	go runWatcher(filepath.Join(relativePath, "static"), assetWatcherCallback, project, false)

	// Watch for changes in the './rawdog.yaml' file
	rawdogWatcherCallback := WatcherCallbacks{
		Write:  configWriteCallback,
		Create: configCreateCallback,
		Remove: configRemoveCallback,
		Rename: configRemoveCallback,
	}
	go runWatcher(filepath.Join(relativePath, "rawdog.yaml"), rawdogWatcherCallback, project, true)

	err = project.ForceRebuild()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("✨ Watching for changes in", pathStyle.Render(relativePath))
	fmt.Println("Press Ctrl+C to stop.")

	go internal.Serve(watcherServer, filepath.Join(relativePath, "build"), port)

	<-make(chan int)

	fmt.Println("Shutting down...")

	return nil
}

// if persist is true, the watcher will keep retrying to watch the file/folder if it doesn't exist until it does
func runWatcher(relativePath string, callbacks WatcherCallbacks, project *internal.Project, persist bool) {
	w := watcher.New()

	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write, watcher.Create, watcher.Remove, watcher.Rename)

	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.Op == watcher.Write {
					err := callbacks.Write(event.Path, project)
					if err != nil {
						log.Println(err)
					}
				} else if event.Op == watcher.Create {
					err := callbacks.Create(event.Path, project)
					if err != nil {
						log.Println(err)
					}
				} else if event.Op == watcher.Remove {
					err := callbacks.Remove(event.Path, project)
					if err != nil {
						log.Println(err)
					}
				} else if event.Op == watcher.Rename {
					err := callbacks.Rename(event.Path, project)
					if err != nil {
						log.Println(err)
					}
				}
			case err := <-w.Error:
				if persist {
					callbacks.Remove(relativePath, project) // could be caused by file rename but whatever

					for {
						if _, err := os.Stat(relativePath); err == nil {
							err := callbacks.Create(relativePath, project)
							if err != nil {
								log.Println(err)
							}
							w.AddRecursive(relativePath)
							break
						}

						time.Sleep(500 * time.Millisecond)
					}
				} else {
					log.Fatalln(err)
				}
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(relativePath); err != nil {
		if persist {
			// fmt.Println("watched file or folder doesnt exist, waiting for it to be created")

			for {
				if _, err := os.Stat(relativePath); err == nil {
					err := callbacks.Create(relativePath, project)
					if err != nil {
						log.Println(err)
					}
					w.AddRecursive(relativePath)
					break
				}

				time.Sleep(500 * time.Millisecond)
			}

		} else {
			log.Fatalln(err)
		}
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func eventRelativeRoot(eventPath string) string {
	cwd := global.Config.RootAbsolutePath

	relPath, err := filepath.Rel(cwd, eventPath)
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

func configCreateCallback(eventPath string, project *internal.Project) error {
	err := global.LoadUserConfig()
	if err != nil {
		fmt.Println("🔧 Config File Created")
		return fmt.Errorf("config error: %v", err)
	}

	startTime := time.Now()
	err = project.ForceRebuild()
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		fmt.Println("🔧 Config File Created")
		return fmt.Errorf("build error: %v", err)
	}

	fmt.Println("🔧 Config File Created", "(rebuild took", durationMs, "ms)")
	watcherServer.Broadcast("reload")
	return nil
}

func configWriteCallback(eventPath string, project *internal.Project) error {
	err := global.LoadUserConfig()
	if err != nil {
		fmt.Println("🔧 Config File Created")
		return fmt.Errorf("config error: %v", err)
	}

	startTime := time.Now()
	err = project.ForceRebuild()
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		fmt.Println("🔧 Config File Modified")
		return fmt.Errorf("build error: %v", err)
	}

	fmt.Println("🔧 Config File Modified", "(rebuild took", durationMs, "ms)")
	watcherServer.Broadcast("reload")
	return nil
}

func configRemoveCallback(eventPath string, project *internal.Project) error {
	err := global.LoadUserConfig()
	if err != nil {
		fmt.Println("🔧 Config File Created")
		return fmt.Errorf("config error: %v", err)
	}

	startTime := time.Now()
	err = project.ForceRebuild()
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		fmt.Println("🔧 Config File Removed")
		return fmt.Errorf("build error: %v", err)
	}

	fmt.Println("🔧 Config File Removed, falling back to default config.", "(rebuild took", durationMs, "ms)")
	watcherServer.Broadcast("reload")
	return nil
}
