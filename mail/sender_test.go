package mail

import (
	"fmt"
	"testing"

	"github.com/arya2004/xyfin/utils"
)



func TestSendEmailWithGmail(t *testing.T) {	

	if testing.Short() {
		t.Skip()
	}

	config, err := utils.LoadConfig("..")

	if err != nil {
		fmt.Println("Failed to load config:", err)
		
	}

	sender := NewGmailSender("Xphyrus", config.EmailSenderAddress, config.EmailSenderPassword)
	to := []string{"arya.pathak22@vit.edu"}
	cc := []string{}
	bcc := []string{}
	attachFiles := []string{"../sqlc.md"} // Add file paths if you have attachments

	data := map[string]string{
		"Name": "John Doe",
	}

	err = sender.SendEmail("Welcome to Our Service!", "./welcome_template.html", to, cc, bcc, attachFiles, data)
	if err != nil {
		fmt.Println("Failed to send email:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
	
}