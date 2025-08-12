package mail

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)






func SendMessage(sender string, subject string, body string, recipient string) (string,error){
		// automatic loading of env files
	if err := godotenv.Load();err != nil {
		panic("Error loading .env file")
	}
	 
	var domain = os.Getenv("MAIL_DOMAIN")
	var privateApiKey = os.Getenv("MAIL_API_KEY")
	var MailInstance = mailgun.NewMailgun(domain,privateApiKey)
	fmt.Print("the domain is this",domain)
	message := mailgun.NewMessage(sender,subject,body,recipient)
	message.SetHTML(body)
	ctx,cancel := context.WithTimeout(context.Background(),time.Second * 10)
	defer cancel()

	resp, _, err := MailInstance.Send(ctx, message)

	if err != nil{
		return "",err
	}
	return resp,nil
}