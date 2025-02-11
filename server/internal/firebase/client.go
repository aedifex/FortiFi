package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/aedifex/FortiFi/config"

	"google.golang.org/api/option"
)

type FcmClient struct {
	messagingClient *messaging.Client
}

func NewFirebaseMessagingClient(config *config.Config) (*FcmClient, error) {
	opt := option.WithCredentialsFile(config.FcmKeyPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
	  return nil, fmt.Errorf("error initializing app: %v", err)
	}
	messaging, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}
	return &FcmClient{
		messagingClient: messaging,
	}, nil
}

func (client *FcmClient) SendMessage(registrationToken string) (string,error) {
	ctx := context.Background()
	notification := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Intrusion Detection",
			Body: "An anomaly has been detected on your network",
		},
		Token: registrationToken,
	}
	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.messagingClient.Send(ctx, notification)
	if err != nil {
		return "", err
	}
	// Response is a message ID string.
	return response, nil
}