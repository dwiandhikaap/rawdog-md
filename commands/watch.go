package commands

import (
	"fmt"
	"path/filepath"
	"rawdog-md/global"
)

func Watch(relativePath string) error {
	rootAbs, err := filepath.Abs(relativePath)
	if err != nil {
		return err
	}

	config := global.ConfigType{
		RootRelativePath: relativePath,
		RootAbsolutePath: rootAbs,
	}
	global.SetGlobalConfig(config)

	fmt.Println("Watching '" + relativePath + "' for changes...")

	return nil
}
