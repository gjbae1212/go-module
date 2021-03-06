package actions_on_google

type UrlTypeHint string

type ImageDisplayOptions string

type MediaType string

type HorizontalAlignment string

type ActionType string

type PriceType string

type IntentValueType string

type Permission string

const (
	UTH_URL_TYPE_HINT_UNSPECIFIED UrlTypeHint = "URL_TYPE_HINT_UNSPECIFIED"
	UTH_AMP_CONTENT               UrlTypeHint = "AMP_CONTENT"
)

const (
	IDO_DEFAULT ImageDisplayOptions = "DEFAULT"
	IDO_WHITE   ImageDisplayOptions = "WHITE"
	IDO_CROPPED ImageDisplayOptions = "CROPPED"
)

const (
	MT_MEDIA_TYPE_UNSPECIFIED MediaType = "MEDIA_TYPE_UNSPECIFIED"
	MT_AUDIO                  MediaType = "AUDIO"
)

const (
	HA_LEADING  HorizontalAlignment = "LEADING"
	HA_CENTER   HorizontalAlignment = "CENTER"
	HA_TRAILING HorizontalAlignment = "TRAILING"
)

const (
	AT_UNKNOWN          ActionType = "UNKNOWN"
	AT_VIEW_DETAILS     ActionType = "VIEW_DETAILS"
	AT_MODIFY           ActionType = "MODIFY"
	AT_CANCEL           ActionType = "CANCEL"
	AT_RETURN           ActionType = "RETURN"
	AT_EXCHANGE         ActionType = "EXCHANGE"
	AT_EMAIL            ActionType = "EMAIL"
	AT_CALL             ActionType = "CALL"
	AT_REORDER          ActionType = "REORDER"
	AT_REVIEW           ActionType = "REVIEW"
	AT_CUSTOMER_SERVICE ActionType = "CUSTOMER_SERVICE"
	AT_FIX_ISSUE        ActionType = "FIX_ISSUE"
)

const (
	PT_UNKNOWN  PriceType = "UNKNOWN"
	PT_ESTIMATE PriceType = "ESTIMATE"
	PT_ACTUAL   PriceType = "ACTUAL"
)

const (
	IVT_CONFIRMATION     IntentValueType = "type.googleapis.com/google.actions.v2.ConfirmationValueSpec"
	IVT_DATETIME         IntentValueType = "type.googleapis.com/google.actions.v2.DateTimeValueSpec"
	IVT_DELIVERY_ADDRESS IntentValueType = "type.googleapis.com/google.actions.v2.DeliveryAddressValueSpec"
	IVT_LINK             IntentValueType = "type.googleapis.com/google.actions.v2.LinkValueSpec"
	IVT_OPTION           IntentValueType = "type.googleapis.com/google.actions.v2.OptionValueSpec"
	IVT_PERMISSION       IntentValueType = "type.googleapis.com/google.actions.v2.PermissionValueSpec"
	IVT_SIGN_IN          IntentValueType = "type.googleapis.com/google.actions.v2.SignInValueSpec"
)

const (
	PM_UNSPECIFIED_PERMISSION  Permission = "UNSPECIFIED_PERMISSION"
	PM_NAME                    Permission = "NAME"
	PM_DEVICE_PRECISE_LOCATION Permission = "DEVICE_PRECISE_LOCATION"
	PM_DEVICE_COARSE_LOCATION  Permission = "DEVICE_COARSE_LOCATION"
	PM_UPDATE                  Permission = "UPDATE"
)

func (ivt IntentValueType) Intent() string {
	switch ivt {
	case IVT_SIGN_IN:
		return "actions.intent.SIGN_IN"
	case IVT_PERMISSION:
		return "actions.intent.PERMISSION"
	case IVT_OPTION:
		return "actions.intent.OPTION"
	case IVT_LINK:
		return "actions.intent.LINK"
	case IVT_DELIVERY_ADDRESS:
		return "actions.intent.DELIVERY_ADDRESS"
	case IVT_DATETIME:
		return "actions.intent.DATETIME"
	case IVT_CONFIRMATION:
		return "actions.intent.CONFIRMATION"
	}
	return ""
}

func (ivt IntentValueType) Event() string {
	switch ivt {
	case IVT_SIGN_IN:
		return "actions_intent_SIGN_IN"
	case IVT_PERMISSION:
		return "actions_intent_PERMISSION"
	case IVT_OPTION:
		return "actions_intent_OPTION"
	case IVT_LINK:
		return "actions_intent_LINK"
	case IVT_DELIVERY_ADDRESS:
		return "actions_intent_DELIVERY_ADDRESS"
	case IVT_DATETIME:
		return "actions_intent_DATETIME"
	case IVT_CONFIRMATION:
		return "actions_intent_CONFIRMATION"
	}
	return ""
}

func (ivt IntentValueType) Context() string {
	switch ivt {
	case IVT_SIGN_IN:
		return "actions_intent_sign_in"
	case IVT_PERMISSION:
		return "actions_intent_permission"
	case IVT_OPTION:
		return "actions_intent_option"
	case IVT_LINK:
		return "actions_intent_link"
	case IVT_DELIVERY_ADDRESS:
		return "actions_intent_delivery_address"
	case IVT_DATETIME:
		return "actions_intent_datetime"
	case IVT_CONFIRMATION:
		return "actions_intent_confirmation"
	}
	return ""
}
