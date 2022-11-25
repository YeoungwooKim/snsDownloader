package email

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/gomail.v2"
)

type Msg struct {
	Subject    string
	Body       string
	BodyType   string
	Attachment string
}

var EMAIL_TOKEN = "THIS IS DUMMY VALUE. MUST CHANGE!!!"
var defaultMsg = &Msg{
	Subject:    "",
	Body:       "",
	BodyType:   "text/html",
	Attachment: "",
}

// this function is made for loading EMAIL_TOKEN which came from token.txt
// its okay to remove this function. not so much effect other codes.
// the only thing to do if you remove this function is just replace variable value which name is "EMAIL_TOKEN"
func init() {
	token_info_file := "token.txt"
	file, err := os.Open(token_info_file)
	defer file.Close()
	if os.IsNotExist(err) {
		// fmt.Printf("Not exist \"token.txt\"\ndeclared value -> \"EMAIL_TOKEN:%v\"\n", EMAIL_TOKEN)
		return
	}
	b, _ := ioutil.ReadAll(file)
	EMAIL_TOKEN = string(b)
	// println("EMAIL_TOKEN : " + EMAIL_TOKEN)
}

func InitMsg(subject, body, attachment string) {
	defaultMsg.Subject = subject
	defaultMsg.Body = body
	defaultMsg.Attachment = attachment
}

func SendMail(sender string, receiver string) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", sender)
	mail.SetHeader("To", receiver)
	mail.SetHeader("Subject", defaultMsg.Subject)
	mail.SetBody(defaultMsg.BodyType, defaultMsg.Body)
	if len(defaultMsg.Attachment) != 0 {
		mail.Attach(defaultMsg.Attachment)
	}

	dial := gomail.NewDialer("smtp.gmail.com", 587, sender, EMAIL_TOKEN)
	if err := dial.DialAndSend(mail); err != nil {
		fmt.Printf("err %v", err)
		return
	}
}
