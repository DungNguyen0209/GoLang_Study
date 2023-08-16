package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func TestSendEmailWithGmail(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "A test mail"
	content := `
	<h1> Hello World </h1>
	<p> This is test message from <a href="https://www.facebook.com/">Jackie</a></p>
	`
	to := []string{"nguyenhoangminhdungbk18@gmail.com"}
	attachFiles := []string{"../README.md"}
	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
