package commands

import (
	"msnserver/pkg/clients"
	"testing"
)

func TestHandleUSRDispatch(t *testing.T) {
	tests := []struct {
		arguments             string
		expectedTransactionID string
		expectedSession       clients.Session
		ok                    bool
	}{
		{"2 MD5 I example@passport.com", "2", clients.Session{Email: "example@passport.com", AuthMethod: "MD5", AuthState: "I", Password: ""}, true},
		{"MD5 I example@passport.com", "", clients.Session{}, false},
		{"2 MD5 I", "", clients.Session{}, false},
		{"2 MD5 I example@passport.com foo", "", clients.Session{}, false},
	}

	for _, tt := range tests {
		s := &clients.Session{}

		gotTransactionID, gotErr := HandleReceiveUSR(s, tt.arguments)

		if (gotErr == nil) != tt.ok {
			t.Errorf("Error HandleReceiveUSR(%q) = %v, want %v", tt.arguments, gotErr == nil, tt.ok)
		}

		if gotTransactionID != tt.expectedTransactionID {
			t.Errorf("TransactionID HandleReceiveUSR(%q) = %q, want %q", tt.arguments, gotTransactionID, tt.expectedTransactionID)
		}

		if s.AuthMethod != tt.expectedSession.AuthMethod {
			t.Errorf("AuthMethod HandleReceiveUSR(%q) = %q, want %q", tt.arguments, s.AuthMethod, tt.expectedSession.AuthMethod)
		}

		if s.AuthState != tt.expectedSession.AuthState {
			t.Errorf("AuthState HandleReceiveUSR(%q) = %q, want %q", tt.arguments, s.AuthState, tt.expectedSession.AuthState)
		}

		if s.Email != tt.expectedSession.Email {
			t.Errorf("Username HandleReceiveUSR(%q) = %q, want %q", tt.arguments, s.Email, tt.expectedSession.Email)
		}

	}
}
