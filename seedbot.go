package main

import (
	"flag"
	"github.com/dpatrie/chatbot/bot"
	"log"
)

//TODO: Command-line params
func main() {
	var host, user, pass, room, name string
	flag.StringVar(&host, "host", "", "Hostname:port of the XMPP server")
	flag.StringVar(&user, "user", "", "Username of XMPP server (i.e.: foo@hostname.com")
	flag.StringVar(&pass, "pass", "", "Password for XMPP server")
	flag.StringVar(&room, "room", "", "Room to join (i.e.: #myroom@hostname.com")
	flag.StringVar(&name, "name", "CrazyBot", "Name of the bot")
	flag.Parse()

	//TODO:Add some validation...but whatever for now

	seedbot := Seedbot{
		bot.NewXMPPBot(host, user, pass, room, name),
		[]bot.Plugin{
			bot.PluginEcho{},
		},
	}
	err := seedbot.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	recv := make(chan bot.Message)
	seedbot.Listen(recv)

	for msg := range recv {
		for _, p := range seedbot.Plugins {
			p.Execute(msg, seedbot)
		}
	}
}

type Seedbot struct {
	*bot.XMPPBot
	Plugins []bot.Plugin
}
