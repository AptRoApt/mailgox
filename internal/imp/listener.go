package imp

import (
	"fmt"

	"github.com/ITC-yka/mailgox/internal/app"
	"github.com/emersion/go-imap/v2/imapclient"
)

type Client struct {
	*imapclient.Client
	UpdatesChan chan uint32
	msgCount    uint32
}

func NewClient(loginData app.LoginData) (*Client, error) {
	c := &Client{}
	c.UpdatesChan = make(chan uint32)
	options := imapclient.Options{
		UnilateralDataHandler: &imapclient.UnilateralDataHandler{
			Mailbox: func(data *imapclient.UnilateralDataMailbox) {
				if data.NumMessages != nil {
					newCount := *data.NumMessages
					if newCount > c.msgCount {
						fmt.Println("NEW UPDATE")
						c.UpdatesChan <- newCount
					}
					c.msgCount = newCount
				}
			},
		},
	}
	serverStr := fmt.Sprintf("%s:%d", loginData.Server, loginData.Port)
	var err error
	c.Client, err = imapclient.DialTLS(serverStr, &options)
	if err != nil {
		return nil, fmt.Errorf("failed to dial IMAP server: %v", err)
	}

	if err := c.Login(loginData.Login, loginData.Password).Wait(); err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}
	mb, err := c.Select("INBOX", nil).Wait()
	if err != nil {
		return nil, err
	}
	c.msgCount = mb.NumMessages
	return c, nil
}
