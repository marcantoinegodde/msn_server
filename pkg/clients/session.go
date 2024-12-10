package clients

type Session struct {
	Email                       string
	AuthMethod                  string
	AuthState                   string
	Password                    string
	Connected                   bool
	InitialPresenceNotification bool
}
