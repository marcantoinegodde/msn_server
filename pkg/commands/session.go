package commands

type Session struct {
	Email      string
	authMethod string
	authState  string
	password   string
	connected  bool
}
