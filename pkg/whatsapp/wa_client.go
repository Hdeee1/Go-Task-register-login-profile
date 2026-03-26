package whatsapp

import (
	"context"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
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
	if clientWa.Store.ID == nil {
		qrChan, _ := clientWa.GetQRChannel(context.Background())
		if err := clientWa.Connect(); err != nil {
			return nil, err
		}

		for event := range qrChan {
			if event.Event == "code" {
				qrterminal.GenerateHalfBlock(event.Code, qrterminal.L, os.Stdout)
				fmt.Println("scan for connect to whatsapp")
			}
		}
	} else {
		clientWa.Connect()
	}
	return &waClient{client: clientWa}, nil
}

func (w *waClient) OTPSender(to, otpCode string) error {
	if strings.HasPrefix(to, "0") {
		to = "62" + to[1:]
	}
	to = fmt.Sprintf("%s@s.whatsapp.net", to)

	jid, err := types.ParseJID(to)
	if err != nil {
		return err
	}
	protoMess := &waE2E.Message{Conversation: proto.String("Your OTP Code: " + otpCode)}

	_, err = w.client.SendMessage(context.Background(), jid, protoMess)
	if err != nil {
		return err
	}
	return nil
}