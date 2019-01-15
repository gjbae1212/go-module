package actions_on_google

type WebhookResponse struct {
	ExpectUserResponse bool          `json:"expectUserResponse"`
	UserStorage        string        `json:"userStorage,omitempty"`
	RichResponse       *RichResponse `json:"richResponse,omitempty"`
	SystemIntent       *SystemIntent `json:"systemIntent,omitempty"`
}

type RichResponse struct {
	Items             []*Item            `json:"items,omitempty"`
	Suggestions       []*Suggestion      `json:"suggestions,omitempty"`
	LinkOutSuggestion *LinkOutSuggestion `json:"linkOutSuggestion,omitempty"`
}

type SystemIntent struct {
	Intent string            `json:"intent,omitempty"`
	Data   *SystemIntentData `json:"data,omitempty"`
}
