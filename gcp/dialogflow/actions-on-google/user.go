package actions_on_google

import (
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type User struct {
	UserId      string
	LastSeen    string
	UserStorage string
	Locale      string
}

func (u *User) IsEmpty() bool {
	return u.UserId == ""
}

func GetUserFromRequest(req *dialogflowpb.WebhookRequest) (*User, error) {
	user := &User{}
	if req.OriginalDetectIntentRequest == nil {
		return user, nil
	}

	if req.OriginalDetectIntentRequest.Payload == nil {
		return user, nil
	}
	s := req.OriginalDetectIntentRequest.Payload.GetFields()["user"]
	if s == nil {
		return user, nil
	}
	fields := s.GetStructValue().GetFields()
	user.UserId = fields["userId"].GetStringValue()
	user.LastSeen = fields["lastSeen"].GetStringValue()
	user.UserStorage = fields["userStorage"].GetStringValue()
	user.Locale = fields["locale"].GetStringValue()
	return user, nil
}
