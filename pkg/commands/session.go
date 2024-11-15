package commands

type Session struct {
	authMethod string
	authState  string
	email      string
	password   string
	connected  bool
}
