package commands

import (
	"testing"
)

func TestHandleUSRDispatch(t *testing.T) {
	tests := []struct {
		arguments             string
		expectedTransactionID string
		expectedSession       Session
		ok                    bool
	}{
		{"2 MD5 I example@passport.com", "2", Session{Email: "example@passport.com", authMethod: "MD5", authState: "I", password: ""}, true},
		{"MD5 I example@passport.com", "", Session{}, false},
		{"2 MD5 I", "", Session{}, false},
		{"2 MD5 I example@passport.com foo", "", Session{}, false},
	}

	for _, tt := range tests {
		s := &Session{}

		gotTransactionID, gotErr := HandleReceiveUSR(s, tt.arguments)

		if (gotErr == nil) != tt.ok {
			t.Errorf("Error HandleReceiveUSR(%q) = %v, want %v", tt.arguments, gotErr == nil, tt.ok)
		}

		if gotTransactionID != tt.expectedTransactionID {
			t.Errorf("TransactionID HandleReceiveUSR(%q) = %q, want %q", tt.arguments, gotTransactionID, tt.expectedTransactionID)
		}

		if s.authMethod != tt.expectedSession.authMethod {
			t.Errorf("AuthMethod HandleReceiveUSR(%q) = %q, want %q", tt.arguments, s.authMethod, tt.expectedSession.authMethod)
		}

		if s.authState != tt.expectedSession.authState {
			t.Errorf("AuthState HandleReceiveUSR(%q) = %q, want %q", tt.arguments, s.authState, tt.expectedSession.authState)
		}

		if s.Email != tt.expectedSession.Email {
			t.Errorf("Username HandleReceiveUSR(%q) = %q, want %q", tt.arguments, s.Email, tt.expectedSession.Email)
		}

	}
}
