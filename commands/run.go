package commands

import (
	"fmt"
	"os"
	"rawdog-md/global"
	"rawdog-md/internal"
	"strings"
)

func Run() error {
	rootAbs, err := os.Getwd()
	if err != nil {
		return err
	}

	rootAbs = strings.ReplaceAll(rootAbs, "\\", "/")

	config := global.ConfigType{
		RootRelativePath: ".",
		RootAbsolutePath: rootAbs,
		BuildMode:        global.Release,
	}
	global.SetGlobalConfig(config)

	fmt.Println(config)

	project, err := internal.NewProject()
	if err != nil {
		return err
	}

	err = project.ForceRebuild()
	if err != nil {
		return err
	}

	return nil
}
