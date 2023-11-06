package ses

import (
	"errors"
	"olx-clone/conf"
	"olx-clone/constants/messages"
	"olx-clone/functions/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func SendEmail(sender, recipient, subject, htmlBody, textBody, charSet string) (*ses.SendEmailOutput, error) {
	var result *ses.SendEmailOutput

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(conf.AWS_REGION),
			Credentials: credentials.NewStaticCredentials(conf.AWS_ACCESS_KEY, conf.AWS_SECRET_ACCESS_KEY, conf.AWS_TOKEN),
		},
	})

	if err != nil {
		return result, errors.New(messages.SomethingWentWrongMessage)
	}

	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	result, err = svc.SendEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
			case ses.ErrCodeMailFromDomainNotVerifiedException:
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				logger.Log.Errorln(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
				return result, errors.New(messages.SomethingWentWrongMessage)
			default:
				logger.Log.Errorln(aerr.Error())
			}
		}
		return result, errors.New(messages.SomethingWentWrongMessage)
	}
	logger.Log.Println("Email Sent to address: " + recipient)
	return result, err
}
