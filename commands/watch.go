package commands

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"rawdog-md/global"
	"rawdog-md/helper"
	"rawdog-md/internal"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type WatcherCallbacks struct {
	Write  func(string, *internal.Project) error
	Create func(string, *internal.Project) error
	Remove func(string, *internal.Project) error
	Rename func(string, *internal.Project) error
}

func Watch(relativePath string) error {
	rootAbs, err := filepath.Abs(relativePath)
	if err != nil {
		return err
	}

	rootAbs = strings.ReplaceAll(rootAbs, "\\", "/")

	config := global.ConfigType{
		RootRelativePath: relativePath,
		RootAbsolutePath: rootAbs,
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

	err = project.WritePages()
	if err != nil {
		log.Fatal(err)
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

	fmt.Println("Watching '" + relativePath + "' for changes...")
	fmt.Println("Press Ctrl+C to stop.")

	<-make(chan int)

	return nil
}

func runWatcher(w *fsnotify.Watcher, cb WatcherCallbacks, project *internal.Project) {
	lastEvent := time.Now()
	for {

		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}

			if time.Since(lastEvent) < 200*time.Millisecond {
				continue
			}

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

			lastEvent = time.Now()
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			fmt.Println("Error:", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func registerWatcherWalk(watcher *fsnotify.Watcher, path string) error {
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			err = watcher.Add(path)
			fmt.Println("Watching", path)
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
	fmt.Println("Modified:", eventRelativeRoot(eventPath))
	err := project.ForceRebuild()
	if err != nil {
		return fmt.Errorf("build error: %v", err)
	}

	return nil
}
func pageCreateCallback(eventPath string, project *internal.Project) error {
	fmt.Println("Created:", eventRelativeRoot(eventPath))
	err := project.ForceRebuild()
	if err != nil {
		return fmt.Errorf("build error: %v", err)
	}

	return nil
}
func pageRemoveCallback(eventPath string, project *internal.Project) error {
	fmt.Println("Removed:", eventRelativeRoot(eventPath))
	err := project.ForceRebuild()
	if err != nil {
		return fmt.Errorf("build error: %v", err)
	}

	return nil
}

func pageRenameCallback(eventPath string, project *internal.Project) error {
	fmt.Println("Renamed:", eventRelativeRoot(eventPath))
	err := project.ForceRebuild()
	if err != nil {
		return fmt.Errorf("build error: %v", err)
	}

	return nil
}

func assetCallback(eventPath string, project *internal.Project) error {
	fmt.Println("Static file changes", eventRelativeRoot(eventPath))
	err := project.ForceRebuild()
	if err != nil {
		return fmt.Errorf("build error: %v", err)
	}

	return nil
}
