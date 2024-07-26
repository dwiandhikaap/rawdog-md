package global

// create global variable to store the global config
type ConfigType struct {
	RootRelativePath string
	RootAbsolutePath string
}

var Config ConfigType

func SetGlobalConfig(config ConfigType) {
	Config = config
}
