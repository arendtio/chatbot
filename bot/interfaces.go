package bot

type Bot interface {
	Name() string
	Send(msg string)
	AddPlugin(p Plugin)
	Connect() error
	Listen()
}

//*************************************************

type Message interface {
	Body() string
	From() string
}

//*************************************************

type Plugin interface {
	Name() string
	Execute(msg Message, bot Bot) error
}
