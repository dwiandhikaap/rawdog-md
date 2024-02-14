package global

// create global variable to store the global config
type ConfigType struct {
	RootAbsolutePath string
	RootRelativePath string
}

var Config ConfigType

func SetGlobalConfig(config ConfigType) {
	Config = config
}
