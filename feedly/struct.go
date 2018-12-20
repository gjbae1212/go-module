package feedly

type Entry struct {
	Id              string      `json:"id,omitempty"`
	Title           string      `json:"title,omitempty"`
	Author          string      `json:"author,omitempty"`
	OriginId        string      `json:"originId,omitempty"`
	Fingerprint     string      `json:"fingerprint,omitempty"`
	Sid             string      `json:"sid,omitempty"`
	Crawled         int         `json:"crawled,omitempty"`
	Recrawled       int         `json:"recrawled,omitempty"`
	Published       int         `json:"published,omitempty"`
	Updated         int         `json:"updated,omitempty"`
	Engagement      int         `json:"engagement,omitempty"`
	ActionTimestamp int         `json:"actionTimestamp,omitempty"`
	Unread          bool        `json:"unread,omitempty"`
	Keywords        []string    `json:"keywords,omitempty"`
	Origin          *Origin     `json:"origin,omitempty"`
	Content         *Content    `json:"content,omitempty"`
	Summary         *Content    `json:"summary,omitempty"`
	Alternate       []*Link     `json:"alternate,omitempty"`
	Enclosure       []*Link     `json:"enclosure,omitempty"`
	Canonical       []*Link     `json:"canonical,omitempty"`
	Visual          *Visual     `json:"visual,omitempty"`
	Thumbnail       []*Visual   `json:"thumbnail,viomitempty"`
	Tags            []*Tag      `json:"tags,omitempty"`
	Categories      []*Category `json:"categories,omitempty"`
}

type Category struct {
	Id           string          `json:"id,omitempty"`
	Label        string          `json:"label,omitempty"`
	Description  string          `json:"description,omitempty"`
	Customizable bool            `json:"customizable,omitempty"`
	Enterprise   bool            `json:"enterprise,omitempty"`
	Cover        string          `json:"cover,omitempty"`
	Created      int             `json:"created,omitempty"`
	NumFeeds     int             `json:"numFeeds,omitempty"`
	Feeds        []*Subscription `json:"feeds,omitempty"`
}

type Feed struct {
	Id            string   `json:"id,omitempty"`
	FeedId        string   `json:"feedId,omitempty"`
	Title         string   `json:"title,omitempty"`
	Description   string   `json:"description,omitempty"`
	Language      string   `json:"language,omitempty"`
	Website       string   `json:"website,omitempty"`
	Topics        []string `json:"topics,omitempty"`
	Velocity      float64  `json:"velocity,omitempty"`
	Subscribers   int      `json:"subscribers,omitempty"`
	State         string   `json:"state,omitempty"`
	LastUpdated   int      `json:"lastUpdated,omitempty"`
	IconUrl       string   `json:"iconUrl,omitempty"`
	VisualUrl     string   `json:"visualUrl,omitempty"`
	CoverUrl      string   `json:"coverUrl,omitempty"`
	Logo          string   `json:"logo,omitempty"`
	ContentType   string   `json:"contentType,omitempty"`
	CoverColor    string   `json:"coverColor,omitempty"`
	DeliciousTags []string `json:"deliciousTags,omitempty"`
	Partial       bool     `json:"partial,omitempty"`
	Featured      bool     `json:"featured,omitempty"`
}

type Profile struct {
	Id                          string           `json:"id,omitempty"`
	Email                       string           `json:"email,omitempty"`
	GivenName                   string           `json:"givenName,omitempty"`
	FamilyName                  string           `json:"familyName,omitempty"`
	FullName                    string           `json:"fullName,omitempty"`
	Picture                     string           `json:"picture,omitempty"`
	Gender                      string           `json:"gender,omitempty"`
	Locale                      string           `json:"locale,omitempty"`
	Google                      string           `json:"google,omitempty"`
	Reader                      string           `json:"reader,omitempty"`
	Wave                        string           `json:"wave,omitempty"`
	Client                      string           `json:"client,omitempty"`
	Source                      string           `json:"source,omitempty"`
	Created                     int              `json:"created,omitempty"`
	Product                     string           `json:"product,omitempty"`
	ProductExpiration           int              `json:"productExpiration,omitempty"`
	ProductRenewalAmount        int              `json:"productRenewalAmount,omitempty"`
	UpgradeDate                 int              `json:"upgradeDate,omitempty"`
	SubscriptionRenewalDate     int              `json:"subscriptionRenewalDate,omitempty"`
	SubscriptionPaymentProvider string           `json:"subscriptionPaymentProvider,omitempty"`
	SubscriptionStatus          string           `json:"subscriptionStatus,omitempty"`
	EvernoteStoreUrl            string           `json:"evernoteStoreUrl,omitempty"`
	EvernoteWebApiPrefix        string           `json:"evernoteWebApiPrefix,omitempty"`
	EvernotePartialOAuth        bool             `json:"evernotePartialOAuth,omitempty"`
	RefPage                     string           `json:"refPage,omitempty"`
	LandingPage                 string           `json:"landingPage,omitempty"`
	LoginProviders              []*LoginProvider `json:"logins,omitempty"`
	CardDetails                 *CardDetails     `json:"cardDetails,omitempty"`
	TwitterUserId               string           `json:"twitterUserId,omitempty"`
	FacebookUserId              string           `json:"facebookUserId,omitempty"`
	WordPressId                 string           `json:"wordPressId,omitempty"`
	WindowsLiveId               string           `json:"windowsLiveId,omitempty"`
	EvernoteUserId              string           `json:"evernoteUserId,omitempty"`
	EvernoteConnected           bool             `json:"evernoteConnected,omitempty"`
	PocketConnected             bool             `json:"pocketConnected,omitempty"`
	DropboxConnected            bool             `json:"dropboxConnected,omitempty"`
	TwitterConnected            bool             `json:"twitterConnected,omitempty"`
	FacebookConnected           bool             `json:"facebookConnected,omitempty"`
	WordPressConnected          bool             `json:"wordPressConnected,omitempty"`
	WindowsLiveConnected        bool             `json:"windowsLiveConnected,omitempty"`
	InstapaperConnected         bool             `json:"instapaperConnected,omitempty"`
}

type LoginProvider struct {
	Id         string `json:"id,omitempty"`
	Verified   bool   `json:"verified,omitempty"`
	Picture    string `json:"picture,omitempty"`
	Provider   string `json:"provider,omitempty"`
	ProviderId string `json:"providerId,omitempty"`
	FullName   string `json:"fullName,omitempty"`
}

type CardDetails struct {
	Brand           string `json:"brand,omitempty"`
	ExpirationMonth int    `json:"expirationMonth,omitempty"`
	ExpirationYear  int    `json:"expirationYear,omitempty"`
	Last4           string `json:"last4,omitempty"`
	Country         string `json:"country,omitempty"`
}

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

type Tag struct {
	Id          string `json:"id,omitempty"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
}

type Content struct {
	Content   string `json:"content,omitempty"`
	Direction string `json:"direction,omitempty"`
}

type Link struct {
	Href string `json:"href,omitempty"`
	Type string `json:"type,omitempty"`
}

type Origin struct {
	StreamId string `json:"streamId,omitempty"`
	Title    string `json:"title,omitempty"`
	HtmlUrl  string `json:"htmlUrl,omitempty"`
}

type Visual struct {
	Url         string `json:"url,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

type Stream struct {
	Id           string   `json:"id,omitempty"`           // stream id
	Ids          []string `json:"ids,omitempty"`          // entry ids
	Continuation string   `json:"continuation,omitempty"` // next stream id for pagination
	Title        string   `json:"title,omitempty"`
	Direction    string   `json:"direction,omitempty"`
	Updated      int      `json:"updated,omitempty"`
	Alternate    []*Link  `json:"alternate,omitempty"`
	Items        []*Entry `json:"items,omitempty"`
}

type Marker struct {
	UnreadCounts []*Unread `json:"unreadcounts,omitempty"`
}

type Unread struct {
	Id      string `json:"id,omitempty"`
	Count   int    `json:"count,omitempty"`
	Updated int    `json:"updated,omitempty"`
}

type SearchResult struct {
	Hint    string   `json:"hint,omitempty"`
	Related []string `json:"related,omitempty"`
	Feeds   []*Feed  `json:"results,omitempty"`
}
