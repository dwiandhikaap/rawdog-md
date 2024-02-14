package commands

import (
	"path/filepath"
	"rawdog-md/global"
	"rawdog-md/internal"
)

func Run(relativePath string) error {
	rootAbs, err := filepath.Abs(relativePath)
	if err != nil {
		return err
	}

	config := global.ConfigType{
		RootRelativePath: relativePath,
		RootAbsolutePath: rootAbs,
	}
	global.SetGlobalConfig(config)

	pages, err := internal.LoadProject()
	if err != nil {
		return err
	}

	err = internal.WritePages(pages)
	if err != nil {
		return err
	}

	err = internal.CopyStaticFiles()
	if err != nil {
		return err
	}

	return nil
}
