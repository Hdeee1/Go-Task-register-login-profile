package whatsapp

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	_ "github.com/mattn/go-sqlite3"
)

type WhatsappSender interface {
	OTPSender(to, otpCode string) error
}

type waClient struct {
	client *whatsmeow.Client
}

func NewWaClient() (*waClient, error) {
	container, err := sqlstore.New(context.Background(), "sqlite3", "file:wastorage.db?_foreign_keys=on", nil)
	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, err
	}

	clientWa := whatsmeow.NewClient(deviceStore, nil)
	return &waClient{client: clientWa}, nil
}

func (w *waClient) OTPSender(to, otpCode string) error {
	return nil
}