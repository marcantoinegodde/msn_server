package commands

import (
	"testing"

	"gorm.io/gorm"
)

func TestHandleUSRDispatch(t *testing.T) {
	tests := []struct {
		arguments             string
		expectedTransactionID string
		expectedSession       Session
		ok                    bool
	}{
		{"2 MD5 I example@passport.com", "2", Session{authMethod: "MD5", authState: "I", email: "example@passport.com", password: ""}, true},
		{"MD5 I example@passport.com", "", Session{}, false},
		{"2 MD5 I", "", Session{}, false},
		{"2 MD5 I example@passport.com foo", "", Session{}, false},
	}

	for _, tt := range tests {
		mockConn := &mockConn{}
		mockDB := &gorm.DB{}

		s := &Session{}

		gotTransactionID, gotErr := HandleReceiveUSR(mockConn, mockDB, s, tt.arguments)

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

		if s.email != tt.expectedSession.email {
			t.Errorf("Username HandleReceiveUSR(%q) = %q, want %q", tt.arguments, s.email, tt.expectedSession.email)
		}

	}
}
