package actions_on_google

type SimpleResponse struct {
	TextToSpeech string `json:"textToSpeech,omitempty"`
	Ssml         string `json:"ssml,omitempty"`
	DisplayText  string `json:"displayText,omitempty"`
}

type MediaResponse struct {
	MediaType    MediaType      `json:"mediaType,omitempty"`
	MediaObjects []*MediaObject `json:"mediaObjects,omitempty"`
}

type Suggestion struct {
	Title string `json:"title,omitempty"`
}

type Item struct {
	SimpleResponse     *SimpleResponse     `json:"simpleResponse,omitempty"`
	BasicCard          *BasicCard          `json:"basicCard,omitempty"`
	StructuredResponse *StructuredResponse `json:"structuredResponse,omitempty"`
	MediaResponse      *MediaResponse      `json:"mediaResponse,omitempty"`
	CarouselBrowse     *CarouselBrowse     `json:"carouselBrowse,omitempty"`
	TableCard          *TableCard          `json:"tableCard,omitempty"`
}

type BasicCard struct {
	Title               string              `json:"title,omitempty"`
	Subtitle            string              `json:"subtitle,omitempty"`
	FormattedText       string              `json:"formattedText,omitempty"`
	Image               *Image              `json:"image,omitempty"`
	Buttons             []*Button           `json:"buttons,omitempty"`
	ImageDisplayOptions ImageDisplayOptions `json:"imageDisplayOptions,omitempty"`
}

type CarouselBrowse struct {
	Items               []*CarouselBrowseItem `json:"items,omitempty"`
	ImageDisplayOptions ImageDisplayOptions   `json:"imageDisplayOptions,omitempty"`
}

type TableCard struct {
	Title            string            `json:"title,omitempty"`
	Subtitle         string            `json:"Subtitle,omitempty"`
	Image            *Image            `json:"image,omitempty"`
	ColumnProperties []*ColumnProperty `json:"columnProperties"`
	Rows             []*Row            `json:"rows,omitempty"`
	Buttons          []*Button         `json:"buttons,omitempty"`
}

type Image struct {
	Url               string `json:"url,omitempty"`
	AccessibilityText string `json:"accessibilityText,omitempty"`
	Height            int    `json:"height,omitempty"`
	Width             int    `json:"width,omitempty"`
}

type Button struct {
	Title         string         `json:"title,omitempty"`
	OpenUrlAction *OpenUrlAction `json:"openUrlAction,omitempty"`
}

type MediaObject struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ContentUrl  string `json:"contentUrl,omitempty"`
	LargeImage  *Image `json:"largeImage,omitempty"`
	Icon        *Image `json:"icon,omitempty"`
}

type ColumnProperty struct {
	Header              string              `json:"header,omitempty"`
	HorizontalAlignment HorizontalAlignment `json:"horizontalAlignment,omitempty"`
}

type Row struct {
	Cells        []*Cell `json:"cells,omitempty"`
	DividerAfter bool    `json:"dividerAfter"`
}

type Cell struct {
	Text string `json:"text,omitempty"`
}

type LinkOutSuggestion struct {
	DestinationName string         `json:"destinationName,omitempty"`
	OpenUrlAction   *OpenUrlAction `json:"openUrlAction,omitempty"`
}

type OpenUrlAction struct {
	Url         string      `json:"url,omitempty"`
	AndroidApp  *AndroidApp `json:"androidApp,omitempty"`
	UrlTypeHint UrlTypeHint `json:"urlTypeHint,omitempty"`
}

type AndroidApp struct {
	PackageName string           `json:"packageName,omitempty"`
	Versions    []*VersionFilter `json:"versions,omitempty"`
}

type VersionFilter struct {
	MinVersion int `json:"minVersion,omitempty"`
	MaxVersion int `json:"maxVersion,omitempty"`
}

type SimpleSelect struct {
	Items []*SimpleSelectItem `json:"items,omitempty"`
}

type ListSelect struct {
	Title string      `json:"title,omitempty"`
	Items []*ListItem `json:"items,omitempty"`
}

type CarouselSelect struct {
	Items               []*CarouselItem     `json:"items,omitempty"`
	ImageDisplayOptions ImageDisplayOptions `json:"imageDisplayOptions,omitempty"`
}

type SimpleSelectItem struct {
	OptionInfo *OptionInfo `json:"optionInfo,omitempty"`
	Title      string      `json:"title,omitempty"`
}

type ListItem struct {
	OptionInfo  *OptionInfo `json:"optionInfo,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Image       *Image      `json:"image,omitempty"`
}

type OptionInfo struct {
	Key      string   `json:"key,omitempty"`
	Synonyms []string `json:"synonyms,omitempty"`
}

type CarouselItem struct {
	OptionInfo  *OptionInfo `json:"optionInfo,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Image       *Image      `json:"image,omitempty"`
}

type CarouselBrowseItem struct {
	Title         string         `json:"title,omitempty"`
	Description   string         `json:"description,omitempty"`
	Footer        string         `json:"footer,omitempty"`
	Image         *Image         `json:"image,omitempty"`
	OpenUrlAction *OpenUrlAction `json:"openUrlAction,omitempty"`
}

type Argument struct {
	/* TODO: 추후 구현
	- https://developers.google.com/actions/reference/rest/Shared.Types/Argument
	*/
}

type StructuredResponse struct {
	/* TODO: 추후 개발
	https://developers.google.com/actions/reference/rest/Shared.Types/AppResponse#StructuredResponse
	*/
}
