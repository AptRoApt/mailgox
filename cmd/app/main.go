package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ITC-yka/mailgox/internal/app"
	"github.com/ITC-yka/mailgox/internal/imp"
	"github.com/emersion/go-imap/v2"
	"github.com/jhillyerd/enmime"
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
		//создаём reader, т.к. функция, парсящая сообщение, работает с ним.
		//Т.к. reader'у нужен string, преобразуем байты в string
		//Я уверен, что тут лучше как-то через указатели это сделать, чтоб лишнюю память не тратить.
		//но я пока плохо знаю язык.
		r := strings.NewReader(string(raw))
		// https://pkg.go.dev/github.com/jhillyerd/enmime#Envelope - устройство структуры.
		env, err := enmime.ReadEnvelope(r)
		if err != nil {
			panic(err)
		}
		fmt.Println(env.Text)
		// TODO: Обработать содержимое письма

	}

}
