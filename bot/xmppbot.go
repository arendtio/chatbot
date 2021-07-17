package bot

import (
	"log"
	"os"

	xmpp "github.com/mattn/go-xmpp"
)

type XMPPBotOptions struct {
	xmpp.Options
	Room string
}

//*************************************************

type XMPPBotMessage struct {
	body, from string
}

func (m XMPPBotMessage) Body() string {
	return m.body
}

func (m XMPPBotMessage) From() string {
	return m.from
}

//*************************************************

type XMPPBot struct {
	Options XMPPBotOptions
	Plugins []Plugin
	client  *xmpp.Client
	logger  *log.Logger
}

func (b *XMPPBot) Name() string {
	return b.Options.Room + "/" + b.Options.Resource
}

func (b *XMPPBot) AddPlugin(p Plugin) {
	b.Plugins = append(b.Plugins, p)
}

func (b *XMPPBot) Send(msg string) {
	b.client.Send(xmpp.Chat{Remote: b.Options.Room, Type: "groupchat", Text: msg})
}

func (b *XMPPBot) Connect() error {
	var err error
	b.logger.Printf("Connecting to %s:*******@%s \n", b.Options.User, b.Options.Host)
	b.client, err = b.Options.NewClient()
	if err != nil {
		b.logger.Printf("Error: %s \n", err)
		return err
	}

	b.logger.Printf("Joining %s with resource %s \n", b.Options.Room, b.Options.Resource)
	b.client.JoinMUCNoHistory(b.Options.Room+"/"+b.Options.Resource, "")
	return nil
}

func (b *XMPPBot) Listen() {
	for {
		chat, err := b.client.Recv()
		if err != nil {
			b.logger.Printf("Error: %s \n", err)
			break // receiving an error breaks the loop and allows a reconnect
		}
		switch v := chat.(type) {
		case xmpp.Chat:
			msg := XMPPBotMessage{body: v.Text, from: v.Remote}
			for _, p := range b.Plugins {
				p.Execute(msg, b)
			}
		case xmpp.Presence:
			b.logger.Printf("Presence: %+v \n", v)
		}
	}
}

func (b *XMPPBot) SetLogger(logger *log.Logger) {
	b.logger = logger
}

//*************************************************

func NewXMPPBot(host, user, password, room, name string) *XMPPBot {
	opt := XMPPBotOptions{
		xmpp.Options{
			Host:     host,
			User:     user,
			Password: password,
			Resource: name,
			NoTLS:    true,
			Debug:    false,
			Session:  true,
		},
		room,
	}

	bot := &XMPPBot{Options: opt}
	bot.SetLogger(log.New(os.Stderr, "", log.LstdFlags))
	return bot
}
