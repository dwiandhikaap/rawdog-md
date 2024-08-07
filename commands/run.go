package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dwiandhikaap/rawdog-md/global"
	"github.com/dwiandhikaap/rawdog-md/internal"

	"github.com/charmbracelet/lipgloss"
)

func Run(relativePath string) error {
	rootAbs, err := os.Getwd()
	if err != nil {
		return err
	}

	rootAbs = filepath.Join(rootAbs, relativePath)
	rootAbs = strings.ReplaceAll(rootAbs, "\\", "/")

	config := global.ConfigType{
		RootRelativePath: relativePath,
		RootAbsolutePath: rootAbs,
		BuildMode:        global.Release,
	}
	global.SetGlobalConfig(config)

	style1 := lipgloss.NewStyle()
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00b0ff"))

	projectDirStr := relativePath
	if projectDirStr == "." {
		projectDirStr = "current directory"
	}

	fmt.Println(style1.Render("🚀 building at ") +
		style2.Render(projectDirStr) +
		style1.Render("..."))

	startTime := time.Now()

	project, err := internal.NewProject()
	if err != nil {
		return err
	}

	err = project.ForceRebuild()
	if err != nil {
		return err
	}

	duration := time.Since(startTime)

	durationText := fmt.Sprintf("%.2fs", duration.Seconds())

	// send message to user that the project has been built successfully at the rootAbs
	fmt.Println(style1.Render("✨ done in ") +
		style2.Render(durationText) +
		style1.Render("."))

	return nil
}
