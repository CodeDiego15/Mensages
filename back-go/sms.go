// sms.go
package main

import (
	"fmt"

	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/rest/api/v2010"
)

func sendSMS(phone, code string) error {
	client := twilio.NewRestClient()

	messageParams := &api.CreateMessageParams{}
	messageParams.SetTo(phone)
	messageParams.SetFrom("your_twilio_phone_number")
	messageParams.SetBody(fmt.Sprintf("Your verification code is: %s", code))

	_, err := client.Api.CreateMessage(messageParams)
	if err != nil {
		return err
	}
	return nil
}
