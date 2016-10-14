package krot

const (
	ReceiverType_Email = 0
	ReceiverType_APNS = 1
	ReceiverType_GCM = 2
)

const (
	ErrReceiverNotFound = error("Receiver not found.")
	ErrInvalidReceiverType = error("Invalid receiver type.")
)

type Receiver struct {
	Name	 string	`json:"name" bson:"name"`
	Type	 uint8	`json:"type" bson:"type"`
	Endpoint string	`json:"endpoint" bson:"endpoint"`
}

// - DAO

func (r Receiver) CreateReceiver(userID string) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	user.Receivers = append(user.Receivers, r)

	return mongoUsers.UpdateId(userID, &user)
}

func GetReceivers(userID string) (*[]Receiver, error) {
	user, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &user.Receivers, nil
}

func RemoveReceiver(userID string, name string) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	for i, rec := range user.Receivers {
		if rec.Name == name {
			user.Receivers = append(user.Receivers[:i], user.Receivers[i+1:]...)
			return mongoUsers.UpdateId(userID, &user)
		}
	}

	return ErrReceiverNotFound
}

// - Job processing

func NotifyUser(userID string, subsID string, text string)  {
	user, err := GetUserByID(userID)
	if err != nil {
		Log(subsID, LogLevel_Error, "Error sending notification: " + err.Error())
		return
	}

	// TODO: check with stored result (Redis?)

	// TODO: create main interface to services

	gcmRegIDs := make([]string, 0)

	for _, rec := range user.Receivers {
		to := rec.Endpoint

		switch rec.Type {

		case ReceiverType_Email:
			if err := SendEmail(subsID, to, text); err != nil {
				Log(subsID, LogLevel_Error, CreateSmtpErrorMessage(rec.Name, to, err))
			} else {
				Log(subsID, LogLevel_Info, "Successfully sent email to " + to)
			}
		case ReceiverType_APNS:
			if err := SendApnsMessage(subsID, to, text); err != nil {
				Log(subsID, LogLevel_Error, "APNS service error: " + err.Error())
			} else {
				Log(subsID, LogLevel_Info, "Successfully sent email to " + to)
			}
		case ReceiverType_GCM:
			gcmRegIDs = append(gcmRegIDs, to)
		}
	}

	if len(gcmRegIDs) > 0 {
		if err := SendGcmMessage(subsID, gcmRegIDs, text); err != nil {
			Log(subsID, LogLevel_Error, "GCM service error: " + err.Error())
		} else {
			Log(subsID, LogLevel_Info, "Successfully sent GCM notifiocations")
		}
	}
}

// - Helpers

/*
func (r Receiver) IsEmail() bool {
	return r.Type == ReceiverType_Email
}

func (r Receiver) IsAPNS() bool {
	return r.Type == ReceiverType_APNS
}

func (r Receiver) IsGCM() bool {
	return r.Type == ReceiverType_GCM
}

func (r Receiver) GetType (string, error) {
	switch r.Type {

	case ReceiverType_Email:
		return "email", nil
	case ReceiverType_APNS:
		return "apns", nil
	case ReceiverType_GCM:
		return "gcm", nil
	default:
		return nil, ErrInvalidReceiverType
	}
}
*/
