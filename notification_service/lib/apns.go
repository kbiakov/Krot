package push_service

import (
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"sync"
	"os"
)

const iosBadgePushCount = 1

func SendApnsMessage(subsID string, deviceToken string, text string) error {
	apnsClient := getApnsClientInstance()
	if apnsClient.err != nil {
		return apnsClient.err
	}

	newPayload := payload.NewPayload().
		Alert("Krot #" + subsID).
		Badge(iosBadgePushCount).
		Custom("message", text)

	notification := &apns.Notification{
		DeviceToken: deviceToken,
		Topic: os.Getenv("BUNDLE_ID"),
		Payload: newPayload,
	}

	_, err := apnsClient.client.Push(notification)

	return err
}

// - Singleton instance

type ApnsClient struct {
	client *apns.Client
	err error
}

var apnsClientInstance *ApnsClient

func getApnsClientInstance() *ApnsClient {
	sync.Once.Do(func() {
		apnsClientInstance = &newApnsClient()
	})

	return apnsClientInstance
}

func newApnsClient() *ApnsClient {
	certPassword := os.Getenv("CERT_PASSWORD")
	cert, err := certificate.FromPemFile("../cert.pem", certPassword)
	if err != nil {
		return &ApnsClient{nil, err}
	}

	return &ApnsClient{
		client: apns.NewClient(cert).Production(),
		err: nil,
	}
}
