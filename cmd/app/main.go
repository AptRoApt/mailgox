package main

import (
	"fmt"
	"log"

	"github.com/ITC-yka/mailgox/internal/app"
	"github.com/ITC-yka/mailgox/internal/imp"
	"github.com/emersion/go-imap/v2"
)

func main() {
	// конфиг сделан для тестирования. В реалной версии LoginData будет собираться из запроса
	cfg, err := app.ParseConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	c, err := imp.NewClient(app.LoginData{
		Server:   cfg.Server,
		Port:     cfg.ImapPort,
		Login:    cfg.Login,
		Password: cfg.Password,
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	for {
		// start idling
		idleCmd, err := c.Idle()
		if err != nil {
			panic(err)
		}
		// wait update
		fmt.Println("waiting for update")
		upd := <-c.UpdatesChan
		// stop idling
		err = idleCmd.Close()
		if err != nil {
			panic(err)
		}
		fetchOptions := &imap.FetchOptions{
			BodySection: []*imap.FetchItemBodySection{
				{Part: []int{}}, // так надо :)
			},
		}
		seqSet := imap.SeqSetNum(upd)
		msgs, err := c.Fetch(seqSet, fetchOptions).Collect()
		if err != nil {
			panic(err)
		}
		msg := msgs[0]
		var raw []byte
		for _, v := range msg.BodySection {
			raw = v
		}
		fmt.Println(string(raw))
		// TODO: сюда встаить MIME парсер

	}

}
