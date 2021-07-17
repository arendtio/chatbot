package main

import (
	"flag"
	"log"
	"time"

	"github.com/arendtio/chatbot/bot"
)

func main() {
	var host, user, pass, room, name string
	flag.StringVar(&host, "host", "", "Hostname:port of the XMPP server")
	flag.StringVar(&user, "user", "", "Username of XMPP server (i.e.: foo@hostname.com")
	flag.StringVar(&pass, "pass", "", "Password for XMPP server")
	flag.StringVar(&room, "room", "", "Room to join (i.e.: #myroom@hostname.com")
	flag.StringVar(&name, "name", "CrazyBot", "Name of the bot")
	flag.Parse()

	//TODO:Add some validation...but whatever for now

	chatbot := bot.NewXMPPBot(host, user, pass, room, name)
	chatbot.AddPlugin(bot.PluginEcho{})

	for {
		err := chatbot.Connect()
		if err != nil {
			log.Println(err)
			log.Println("Connect failed, waiting 60 seconds...")
			time.Sleep(60 * time.Second)
			continue
		}

		chatbot.Listen()
		time.Sleep(10 * time.Second)
	}
}
