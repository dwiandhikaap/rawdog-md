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
	Options         UserConfigOptionsType
	MarkdownPlugins MarkdownPluginFieldType
}

type UserConfigOptionsType struct {
	Html   HtmlOptionsType
	Minify MinifyOptionsType
}

type HtmlOptionsType struct {
	Unsafe bool
}

type MinifyOptionsType struct {
	HTML bool
	CSS  bool
	JS   bool
	JSON bool
	XML  bool
	SVG  bool
}

type MarkdownPluginFieldType struct {
	Highlighting   HighlightingType
	Enclave        EnclaveType
	Anchor         AnchorType
	GFM            GFMType
	CJK            CJKType
	DefinitionList DefinitionListType
	Footnote       FootnoteType
	Table          TableType
	Strikethrough  StrikethroughType
	Typographer    TypographerType
	TaskList       TaskListType
	Linkify        LinkifyType
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

type GFMType struct {
	Enabled bool
}

type CJKType struct {
	Enabled bool
}

type DefinitionListType struct {
	Enabled bool
}

type FootnoteType struct {
	Enabled bool
}

type TableType struct {
	Enabled bool
}

type StrikethroughType struct {
	Enabled bool
}

type TypographerType struct {
	Enabled bool
}

type TaskListType struct {
	Enabled bool
}

type LinkifyType struct {
	Enabled bool
}

// Pointer-based version of the structs above
type ParseUserConfigType struct {
	Version         *int                          `yaml:"version"`
	Options         *ParseUserConfigOptionsType   `yaml:"options"`
	MarkdownPlugins *ParseMarkdownPluginFieldType `yaml:"markdownPlugins"`
}

type ParseUserConfigOptionsType struct {
	Html   *ParseHtmlOptionsType   `yaml:"html"`
	Minify *ParseMinifyOptionsType `yaml:"minify"`
}

type ParseHtmlOptionsType struct {
	Unsafe *bool `yaml:"unsafe"`
}

type ParseMinifyOptionsType struct {
	HTML *bool `yaml:"html"`
	CSS  *bool `yaml:"css"`
	JS   *bool `yaml:"js"`
	JSON *bool `yaml:"json"`
	XML  *bool `yaml:"xml"`
	SVG  *bool `yaml:"svg"`
}

type ParseMarkdownPluginFieldType struct {
	Highlighting   *ParseHighlightingType   `yaml:"highlighting"`
	Enclave        *ParseEnclaveType        `yaml:"enclave"`
	Anchor         *ParseAnchorType         `yaml:"anchor"`
	GFM            *ParseGFMType            `yaml:"gfm"`
	CJK            *ParseCJKType            `yaml:"cjk"`
	DefinitionList *ParseDefinitionListType `yaml:"definitionlist"`
	Footnote       *ParseFootnoteType       `yaml:"footnote"`
	Table          *ParseTableType          `yaml:"table"`
	Strikethrough  *ParseStrikethroughType  `yaml:"strikethrough"`
	Typographer    *ParseTypographerType    `yaml:"typographer"`
	TaskList       *ParseTaskListType       `yaml:"tasklist"`
	Linkify        *ParseLinkifyType        `yaml:"linkify"`
}

type ParseHighlightingType struct {
	Enabled        bool    `yaml:"enabled"`
	Style          *string `yaml:"style,omitempty"` // Optional
	UseLineNumbers *bool   `yaml:"useLineNumbers"`
}

type ParseEnclaveType struct {
	Enabled bool `yaml:"enabled"`
}

type ParseGFMType struct {
	Enabled bool `yaml:"enabled"`
}

type ParseCJKType struct {
	Enabled bool `yaml:"enabled"`
}

type ParseDefinitionListType struct {
	Enabled bool `yaml:"enabled"`
}

type ParseFootnoteType struct {
	Enabled bool `yaml:"enabled"`
}

type ParseTableType struct {
	Enabled bool `yaml:"enabled"`
}

type ParseStrikethroughType struct {
	Enabled bool `yaml:"enabled"`
}

type ParseTypographerType struct {
	Enabled bool `yaml:"enabled"`
}

type ParseTaskListType struct {
	Enabled bool `yaml:"enabled"`
}

type ParseLinkifyType struct {
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Failed to parse default config. This should never happen if the default config is correct.")
			panic(r)
		}
	}()

	return UserConfigType{
		Version: *parsed.Version,
		Options: UserConfigOptionsType{
			Html: HtmlOptionsType{
				Unsafe: *parsed.Options.Html.Unsafe,
			},
			Minify: MinifyOptionsType{
				HTML: *parsed.Options.Minify.HTML,
				CSS:  *parsed.Options.Minify.CSS,
				JS:   *parsed.Options.Minify.JS,
				JSON: *parsed.Options.Minify.JSON,
				XML:  *parsed.Options.Minify.XML,
				SVG:  *parsed.Options.Minify.SVG,
			},
		},
		MarkdownPlugins: MarkdownPluginFieldType{
			Highlighting: HighlightingType{
				Enabled:        parsed.MarkdownPlugins.Highlighting.Enabled,
				Style:          parsed.MarkdownPlugins.Highlighting.Style,
				UseLineNumbers: *parsed.MarkdownPlugins.Highlighting.UseLineNumbers,
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
			GFM: GFMType{
				Enabled: parsed.MarkdownPlugins.GFM.Enabled,
			},
			CJK: CJKType{
				Enabled: parsed.MarkdownPlugins.CJK.Enabled,
			},
			DefinitionList: DefinitionListType{
				Enabled: parsed.MarkdownPlugins.DefinitionList.Enabled,
			},
			Footnote: FootnoteType{
				Enabled: parsed.MarkdownPlugins.Footnote.Enabled,
			},
			Table: TableType{
				Enabled: parsed.MarkdownPlugins.Table.Enabled,
			},
			Strikethrough: StrikethroughType{
				Enabled: parsed.MarkdownPlugins.Strikethrough.Enabled,
			},
			Typographer: TypographerType{
				Enabled: parsed.MarkdownPlugins.Typographer.Enabled,
			},
			TaskList: TaskListType{
				Enabled: parsed.MarkdownPlugins.TaskList.Enabled,
			},
			Linkify: LinkifyType{
				Enabled: parsed.MarkdownPlugins.Linkify.Enabled,
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
		Options:         mergeUserConfigOptions(defaultConfig.Options, userConfig.Options),
		MarkdownPlugins: mergeMarkdownPlugins(defaultConfig.MarkdownPlugins, userConfig.MarkdownPlugins),
	}
	return mergedConfig
}

func mergeUserConfigOptions(defaultOptions, userOptions *ParseUserConfigOptionsType) *ParseUserConfigOptionsType {
	if userOptions == nil {
		return defaultOptions
	}

	merged := &ParseUserConfigOptionsType{
		Html:   mergeHtmlOptions(defaultOptions.Html, userOptions.Html),
		Minify: mergeMinifyOptions(defaultOptions.Minify, userOptions.Minify),
	}

	return merged
}

func mergeHtmlOptions(defaultHtml, userHtml *ParseHtmlOptionsType) *ParseHtmlOptionsType {
	if userHtml == nil {
		return defaultHtml
	}

	return &ParseHtmlOptionsType{
		Unsafe: userHtml.Unsafe,
	}
}

func mergeMinifyOptions(defaultMinify, userMinify *ParseMinifyOptionsType) *ParseMinifyOptionsType {
	if userMinify == nil {
		return defaultMinify
	}

	merged := &ParseMinifyOptionsType{}

	if userMinify.HTML != nil {
		merged.HTML = userMinify.HTML
	} else {
		merged.HTML = defaultMinify.HTML
	}

	if userMinify.CSS != nil {
		merged.CSS = userMinify.CSS
	} else {
		merged.CSS = defaultMinify.CSS
	}

	if userMinify.JS != nil {
		merged.JS = userMinify.JS
	} else {
		merged.JS = defaultMinify.JS
	}

	if userMinify.JSON != nil {
		merged.JSON = userMinify.JSON
	} else {
		merged.JSON = defaultMinify.JSON
	}

	if userMinify.XML != nil {
		merged.XML = userMinify.XML
	} else {
		merged.XML = defaultMinify.XML
	}

	if userMinify.SVG != nil {
		merged.SVG = userMinify.SVG
	} else {
		merged.SVG = defaultMinify.SVG
	}

	return merged
}

func mergeMarkdownPlugins(defaultMarkdownPlugins, userMarkdownPlugins *ParseMarkdownPluginFieldType) *ParseMarkdownPluginFieldType {
	if userMarkdownPlugins == nil {
		return defaultMarkdownPlugins
	}

	merged := &ParseMarkdownPluginFieldType{
		Highlighting:   mergeHighlighting(defaultMarkdownPlugins.Highlighting, userMarkdownPlugins.Highlighting),
		Enclave:        mergeEnclave(defaultMarkdownPlugins.Enclave, userMarkdownPlugins.Enclave),
		Anchor:         mergeAnchor(defaultMarkdownPlugins.Anchor, userMarkdownPlugins.Anchor),
		GFM:            mergeGFM(defaultMarkdownPlugins.GFM, userMarkdownPlugins.GFM),
		CJK:            mergeCJK(defaultMarkdownPlugins.CJK, userMarkdownPlugins.CJK),
		DefinitionList: mergeDefinitionList(defaultMarkdownPlugins.DefinitionList, userMarkdownPlugins.DefinitionList),
		Footnote:       mergeFootnote(defaultMarkdownPlugins.Footnote, userMarkdownPlugins.Footnote),
		Table:          mergeTable(defaultMarkdownPlugins.Table, userMarkdownPlugins.Table),
		Strikethrough:  mergeStrikethrough(defaultMarkdownPlugins.Strikethrough, userMarkdownPlugins.Strikethrough),
		Typographer:    mergeTypographer(defaultMarkdownPlugins.Typographer, userMarkdownPlugins.Typographer),
		TaskList:       mergeTaskList(defaultMarkdownPlugins.TaskList, userMarkdownPlugins.TaskList),
		Linkify:        mergeLinkify(defaultMarkdownPlugins.Linkify, userMarkdownPlugins.Linkify),
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

func mergeGFM(defaultGFM, userGFM *ParseGFMType) *ParseGFMType {
	if userGFM == nil {
		return defaultGFM
	}

	return &ParseGFMType{
		Enabled: userGFM.Enabled,
	}
}

func mergeCJK(defaultCJK, userCJK *ParseCJKType) *ParseCJKType {
	if userCJK == nil {
		return defaultCJK
	}

	return &ParseCJKType{
		Enabled: userCJK.Enabled,
	}
}

func mergeDefinitionList(defaultDefinitionList, userDefinitionList *ParseDefinitionListType) *ParseDefinitionListType {
	if userDefinitionList == nil {
		return defaultDefinitionList
	}

	return &ParseDefinitionListType{
		Enabled: userDefinitionList.Enabled,
	}
}

func mergeFootnote(defaultFootnote, userFootnote *ParseFootnoteType) *ParseFootnoteType {
	if userFootnote == nil {
		return defaultFootnote
	}

	return &ParseFootnoteType{
		Enabled: userFootnote.Enabled,
	}
}

func mergeTable(defaultTable, userTable *ParseTableType) *ParseTableType {
	if userTable == nil {
		return defaultTable
	}

	return &ParseTableType{
		Enabled: userTable.Enabled,
	}
}

func mergeStrikethrough(defaultStrikethrough, userStrikethrough *ParseStrikethroughType) *ParseStrikethroughType {
	if userStrikethrough == nil {
		return defaultStrikethrough
	}

	return &ParseStrikethroughType{
		Enabled: userStrikethrough.Enabled,
	}
}

func mergeTypographer(defaultTypographer, userTypographer *ParseTypographerType) *ParseTypographerType {
	if userTypographer == nil {
		return defaultTypographer
	}

	return &ParseTypographerType{
		Enabled: userTypographer.Enabled,
	}
}

func mergeTaskList(defaultTaskList, userTaskList *ParseTaskListType) *ParseTaskListType {
	if userTaskList == nil {
		return defaultTaskList
	}

	return &ParseTaskListType{
		Enabled: userTaskList.Enabled,
	}
}

func mergeLinkify(defaultLinkify, userLinkify *ParseLinkifyType) *ParseLinkifyType {
	if userLinkify == nil {
		return defaultLinkify
	}

	return &ParseLinkifyType{
		Enabled: userLinkify.Enabled,
	}
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
