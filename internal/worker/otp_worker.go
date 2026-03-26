package worker

import (
	"encoding/json"
	"strings"

	"github.com/Hdeee1/go-register-login-profile/pkg/mail"
	"github.com/Hdeee1/go-register-login-profile/pkg/whatsapp"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type OTPWorker struct {
	rmqChannel *amqp091.Channel
	mailer mail.EmailSender
	waSender whatsapp.WhatsappSender
	logger *zap.Logger
}

func NewOTPWorker(rmChann *amqp091.Channel, mailer mail.EmailSender, waSender whatsapp.WhatsappSender, logger *zap.Logger) *OTPWorker {
	return &OTPWorker{
		rmqChannel: rmChann,
		mailer: mailer,
		waSender: waSender,
		logger: logger,
	}
}

func (ow *OTPWorker) Start() error {
	messages, err := ow.rmqChannel.Consume("wa_otp_queue", "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range messages {
		type otpMessage struct {
			To		string `json:"to"`
			OTPCode	string `json:"otp_code"`
		}

		var payload otpMessage
		if err := json.Unmarshal([]byte(msg.Body), &payload); err != nil {
			ow.logger.Error(err.Error())
			continue
		}

		isIdentified := strings.Contains(payload.To, "@")
		if isIdentified {
			go func() {
				if err := ow.mailer.SendMail(payload.To, payload.OTPCode); err != nil {
					ow.logger.Error("failed to send otp to email", zap.Error(err))
				}
			}()
		} else {
			go func() {
				if err := ow.waSender.OTPSender(payload.To, payload.OTPCode); err != nil {
					ow.logger.Error("failed to send otp ti whatsapp", zap.Error(err))
				}
			}()
		}
	}
	return nil
}