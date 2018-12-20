package feedly

type Subscription struct {
	Id          string      `json:"id,omitempty"`
	FeedId      string      `json:"feedId,omitempty"`
	SortId      string      `json:"sortid,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	ContentType string      `json:"contentType,omitempty"`
	Language    string      `json:"language,omitempty"`
	Website     string      `json:"website,omitempty"`
	IconUrl     string      `json:"iconUrl,omitempty"`
	CoverUrl    string      `json:"coverUrl,omitempty"`
	VisualUrl   string      `json:"visualUrl,omitempty"`
	CoverColor  string      `json:"coverColor,omitempty"`
	Subscribers int         `json:"subscribers,omitempty"`
	Added       int         `json:"added,omitempty"`
	Updated     int         `json:"updated,omitempty"`
	Velocity    float64     `json:"velocity,omitempty"`
	Partial     bool        `json:"partial,omitempty"`
	Categories  []*Category `json:"categories,omitempty"`
	Topics      []string    `json:"topics,omitempty"`
}
