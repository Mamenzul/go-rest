package mailgun

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/mamenzul/go-rest/configs"
)

func SendSimpleMessage(subject string, body string, recipient string) (string, error) {
	sender := configs.Envs.MAILGUN_SENDER
	mg := mailgun.NewMailgun(configs.Envs.MAILGUN_DOMAIN, configs.Envs.MAILGUN_API_KEY)
	m := mg.NewMessage(sender, subject, body, recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, m)
	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return id, err
}
