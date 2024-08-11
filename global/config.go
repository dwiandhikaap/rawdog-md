package global

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/dwiandhikaap/rawdog-md/helper"
	"gopkg.in/yaml.v3"
)

//go:embed default.yaml
var defaultUserConfigSource []byte

type BuildMode string

const (
	Development BuildMode = "development"
	Release     BuildMode = "release"
)

// Create global variable to store the global config
var Config ConfigType

func SetGlobalConfig(config ConfigType) {
	Config = config
}

type ConfigType struct {
	RootRelativePath string
	RootAbsolutePath string
	BuildMode        BuildMode

	UserConfig UserConfigType
}

// Type masturbation incoming 🗣🗣🔥🔥
type UserConfigType struct {
	Version         int
	MarkdownPlugins MarkdownPluginFieldType
}

type MarkdownPluginFieldType struct {
	Highlighting HighlightingType
	Enclave      EnclaveType
	Anchor       AnchorType
}

type HighlightingType struct {
	Enabled        bool
	Style          *string // Actually optional
	UseLineNumbers bool
}

type EnclaveType struct {
	Enabled bool
}

type AnchorType struct {
	Enabled  bool
	Position string
	Text     string
	Class    string
}

// Pointer-based version of the structs above
type ParseUserConfigType struct {
	Version         *int                          `yaml:"version"`
	MarkdownPlugins *ParseMarkdownPluginFieldType `yaml:"markdownPlugins"`
}

type ParseMarkdownPluginFieldType struct {
	Highlighting *ParseHighlightingType `yaml:"highlighting"`
	Enclave      *ParseEnclaveType      `yaml:"enclave"`
	Anchor       *ParseAnchorType       `yaml:"anchor"`
}

type ParseHighlightingType struct {
	Enabled        bool    `yaml:"enabled"`
	Style          *string `yaml:"style,omitempty"` // Optional
	UseLineNumbers bool    `yaml:"useLineNumbers"`
}

type ParseEnclaveType struct {
	Enabled bool `yaml:"enabled"`
}

type ParseAnchorType struct {
	Enabled  bool    `yaml:"enabled"`
	Position *string `yaml:"position,omitempty"` // Optional
	Text     *string `yaml:"text,omitempty"`     // Optional
	Class    *string `yaml:"class,omitempty"`    // Optional
}

func LoadUserConfig() error {
	// Load default config
	defaultUserConfig, err := loadUserConfig(defaultUserConfigSource)
	if err != nil {
		return fmt.Errorf("error loading default config: %w", err)
	}

	// Check if user config file exists, if not use default config
	if _, err := os.Stat(Config.RootAbsolutePath + "/rawdog.yaml"); os.IsNotExist(err) {
		Config.UserConfig = convertParsedConfig(defaultUserConfig)
		return nil
	}

	// Load user-defined config
	userConfig, err := loadUserConfigFromFile(Config.RootAbsolutePath + "/rawdog.yaml")
	if err != nil {
		return fmt.Errorf("error loading user config: %w", err)
	}

	// Merge user config with default config
	mergedConfig := mergeUserConfig(defaultUserConfig, userConfig)

	Config.UserConfig = convertParsedConfig(mergedConfig)

	return nil
}

// Convert struct with pointer fields to struct with non-pointer fields
// this is gonna blow up if one of the field in the default config is absent
func convertParsedConfig(parsed ParseUserConfigType) UserConfigType {
	return UserConfigType{
		Version: *parsed.Version,
		MarkdownPlugins: MarkdownPluginFieldType{
			Highlighting: HighlightingType{
				Enabled:        parsed.MarkdownPlugins.Highlighting.Enabled,
				Style:          parsed.MarkdownPlugins.Highlighting.Style,
				UseLineNumbers: parsed.MarkdownPlugins.Highlighting.UseLineNumbers,
			},
			Enclave: EnclaveType{
				Enabled: parsed.MarkdownPlugins.Enclave.Enabled,
			},
			Anchor: AnchorType{
				Enabled:  parsed.MarkdownPlugins.Anchor.Enabled,
				Position: *parsed.MarkdownPlugins.Anchor.Position,
				Text:     *parsed.MarkdownPlugins.Anchor.Text,
				Class:    *parsed.MarkdownPlugins.Anchor.Class,
			},
		},
	}
}

func loadUserConfigFromFile(filePath string) (ParseUserConfigType, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ParseUserConfigType{}, err
	}
	return loadUserConfig(data)
}

func loadUserConfig(data []byte) (ParseUserConfigType, error) {
	var userConfig ParseUserConfigType

	err := yaml.Unmarshal(data, &userConfig)
	if err != nil {
		return ParseUserConfigType{}, err
	}

	err = validateConfig(userConfig)
	if err != nil {
		return ParseUserConfigType{}, err
	}

	return userConfig, nil
}

func mergeUserConfig(defaultConfig, userConfig ParseUserConfigType) ParseUserConfigType {
	mergedConfig := ParseUserConfigType{
		Version:         userConfig.Version, // User version takes precedence
		MarkdownPlugins: mergeMarkdownPlugins(defaultConfig.MarkdownPlugins, userConfig.MarkdownPlugins),
	}
	return mergedConfig
}

func mergeMarkdownPlugins(defaultMarkdownPlugins, userMarkdownPlugins *ParseMarkdownPluginFieldType) *ParseMarkdownPluginFieldType {
	if userMarkdownPlugins == nil {
		return defaultMarkdownPlugins
	}

	merged := &ParseMarkdownPluginFieldType{
		Highlighting: mergeHighlighting(defaultMarkdownPlugins.Highlighting, userMarkdownPlugins.Highlighting),
		Enclave:      mergeEnclave(defaultMarkdownPlugins.Enclave, userMarkdownPlugins.Enclave),
		Anchor:       mergeAnchor(defaultMarkdownPlugins.Anchor, userMarkdownPlugins.Anchor),
	}

	return merged
}

func mergeHighlighting(defaultHighlighting, userHighlighting *ParseHighlightingType) *ParseHighlightingType {
	if userHighlighting == nil {
		return defaultHighlighting
	}

	merged := &ParseHighlightingType{
		Enabled:        userHighlighting.Enabled,
		UseLineNumbers: userHighlighting.UseLineNumbers,
	}

	if userHighlighting.Style != nil {
		merged.Style = userHighlighting.Style
	} else {
		merged.Style = defaultHighlighting.Style
	}

	return merged
}

func mergeEnclave(defaultEnclave, userEnclave *ParseEnclaveType) *ParseEnclaveType {
	if userEnclave == nil {
		return defaultEnclave
	}

	return &ParseEnclaveType{
		Enabled: userEnclave.Enabled,
	}
}

func mergeAnchor(defaultAnchor, userAnchor *ParseAnchorType) *ParseAnchorType {
	if userAnchor == nil {
		return defaultAnchor
	}

	merged := &ParseAnchorType{
		Enabled: userAnchor.Enabled,
	}

	if userAnchor.Position != nil {
		merged.Position = userAnchor.Position
	} else {
		merged.Position = defaultAnchor.Position
	}

	if userAnchor.Text != nil {
		merged.Text = userAnchor.Text
	} else {
		merged.Text = defaultAnchor.Text
	}

	if userAnchor.Class != nil {
		merged.Class = userAnchor.Class
	} else {
		merged.Class = defaultAnchor.Class
	}

	return merged
}

func validateConfig(userConfig ParseUserConfigType) error {
	if userConfig.Version == nil {
		return fmt.Errorf("config version not provided")
	}

	var supportedVersions = []int{1}

	if !helper.SliceContainsInt(supportedVersions, *userConfig.Version) {
		return fmt.Errorf("config version %d not supported", *userConfig.Version)
	}

	if userConfig.MarkdownPlugins == nil {
		return nil
	}

	if userConfig.MarkdownPlugins.Highlighting != nil {
		err := validateHighlightingConfig(userConfig.MarkdownPlugins.Highlighting)
		if err != nil {
			return fmt.Errorf("highlighting config: %w", err)
		}
	}

	if userConfig.MarkdownPlugins.Anchor != nil {
		err := validateAnchorConfig(userConfig.MarkdownPlugins.Anchor)
		if err != nil {
			return fmt.Errorf("enclave config: %w", err)
		}
	}

	return nil
}

func validateHighlightingConfig(config *ParseHighlightingType) error {
	var styleList = []string{
		"abap", "algol", "algol_nu", "arduino", "autumn", "average", "base16-snazzy",
		"borland", "bw", "catppuccin-frappe", "catppuccin-latte", "catppuccin-macchiato",
		"catppuccin-mocha", "colorful", "doom-one", "doom-one2", "dracula", "emacs",
		"evergarden", "friendly", "fruity", "github", "github-dark", "gruvbox",
		"gruvbox-light", "hr_high_contrast", "hrdark", "igor", "lovelace", "manni",
		"modus-operandi", "modus-vivendi", "monokai", "monokailight", "murphy",
		"native", "nord", "nordi", "onedark", "onesenterprise", "paraiso-dark",
		"paraiso-light", "pastie", "perldoc", "pygments", "rainbow_dash", "rose-pine",
		"rose-pine-dawn", "rose-pine-moon", "rrt", "solarized-dark", "solarized-dark256",
		"solarized-light", "swapoff", "tango", "tokyonight-day", "tokyonight-moon",
		"tokyonight-night", "tokyonight-storm", "trac", "vim", "vs", "vulcan",
		"witchhazel", "xcode", "xcode-dark",
	}

	// Validate style field only if it's provided
	if config.Style != nil && !helper.SliceContainsString(styleList, *config.Style) {
		return fmt.Errorf("style '%s' not supported. visit 'https://xyproto.github.io/splash/docs/' for available styles", *config.Style)
	}

	return nil
}

func validateAnchorConfig(config *ParseAnchorType) error {
	var positionList = []string{
		"left", "right",
	}

	// Validate position field only if it's provided
	if config.Position != nil && !helper.SliceContainsString(positionList, *config.Position) {
		return fmt.Errorf("position '%s' not supported. available positions: %s", *config.Position, positionList)
	}

	return nil
}
