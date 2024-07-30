package global

type BuildMode string

const (
	Development BuildMode = "development"
	Release     BuildMode = "release"
)

// create global variable to store the global config
type ConfigType struct {
	RootRelativePath string
	RootAbsolutePath string
	BuildMode        BuildMode
}

var Config ConfigType

func SetGlobalConfig(config ConfigType) {
	Config = config
}
