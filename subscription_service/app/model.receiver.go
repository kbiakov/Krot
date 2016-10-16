package main

const (
	ReceiverType_Email = iota
	ReceiverType_APNS
	ReceiverType_GCM
)

type Receiver struct {
	Name	 string	`json:"name" bson:"name"`
	Type	 uint8	`json:"type" bson:"type"`
	Endpoint string	`json:"endpoint" bson:"endpoint"`
}

func NotifyUser(userID string, subsID string, text string) {
	user, err := GetUserByID(userID)
	if err != nil {
		Log(subsID, LogLevel_Error, "Error sending notification: " + err.Error())
		return
	}

	// TODO: check with stored result (Redis?)

	gcmRegIDs := make([]string, 0)

	for _, rec := range user.Receivers {
		to := rec.Endpoint

		switch rec.Type {

		case ReceiverType_Email:
			PushToQueue("email", CreateNotification(subsID, to))
		case ReceiverType_APNS:
			PushToQueue("apns", CreateNotification(subsID, to))
		case ReceiverType_GCM:
			gcmRegIDs = append(gcmRegIDs, to)
		}
	}

	if len(gcmRegIDs) > 0 {
		PushToQueue("gcm", &Notification{subsID, &gcmRegIDs, text})
	}
}
