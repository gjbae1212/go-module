package actions_on_google

type SystemIntentData struct {
	Type IntentValueType `json:"@type,omitempty"`
	*ConfirmationData
	*DateTimeData
	*DeliveryAddressData
	*OptionData
	*PermissionData
	*SignInData
	*LinkData
}

type ConfirmationData struct {
	DialogSpec *ConfirmationDialogSpec `json:"dialogSpec,omitempty"`
}

type DateTimeData struct {
	DialogSpec *DateTimeDialogSpec `json:"dialogSpec,omitempty"`
}

type DeliveryAddressData struct {
	AddressOptions *AddressOptions `json:"addressOptions,omitempty"`
}

type OptionData struct {
	SimpleSelect   *SimpleSelect   `json:"simpleSelect,omitempty"`
	ListSelect     *ListSelect     `json:"listSelect,omitempty"`
	CarouselSelect *CarouselSelect `json:"carouselSelect,omitempty"`
}

type PermissionData struct {
	PermissionValueSpec
}

type SignInData struct {
	SignInValueSpec
}

type LinkData struct {
	LinkValueSpec
}

type ConfirmationDialogSpec struct {
	RequestConfirmationText string `json:"requestConfirmationText,omitempty"`
}

type DateTimeDialogSpec struct {
	RequestDatetimeText string `json:"requestDatetimeText,omitempty"`
	RequestDateText     string `json:"requestDateText,omitempty"`
	RequestTimeText     string `json:"requestTimeText,omitempty"`
}

type AddressOptions struct {
	Reason string `json:"reason,omitempty"`
}

type UpdatePermissionValueSpec struct {
	Intent    string      `json:"intent,omitempty"`
	Arguments []*Argument `json:"arguments,omitempty"`
}

type PermissionValueSpec struct {
	OptContext                string                     `json:"optContext,omitempty"`
	Permissions               []Permission               `json:"permissions,omitempty"`
	UpdatePermissionValueSpec *UpdatePermissionValueSpec `json:"updatePermissionValueSpec,omitempty"`
}

type SignInValueSpec struct {
	OptContext string `json:"optContext,omitempty"`
}

type LinkValueSpec struct {
	/* TODO: 추후 개발
	https://developers.google.com/actions/reference/rest/Shared.Types/LinkValueSpec
	*/
}
