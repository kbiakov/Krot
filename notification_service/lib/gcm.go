package push_service

import (
	"errors"
	"github.com/alexjlockwood/gcm"
	"os"
)

const RetriesCount = 2

func SendGcmMessage(subsID string, regIDs []string, text string) error {
	data := map[string]interface{}{
		"subscription_id": subsID,
		"text":            text,
	}

	msg := gcm.NewMessage(data, regIDs...)

	sender := &gcm.Sender{ApiKey: os.Getenv("GCM_API_KEY")}

	if response, err := sender.Send(msg, RetriesCount); err != nil {
		return err
	} else if response.Failure == len(regIDs) {
		return errors.New("None of the messages were not delivered")
	}
	return nil
}
