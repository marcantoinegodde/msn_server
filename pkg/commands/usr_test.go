package commands

import (
	"testing"

	"gorm.io/gorm"
)

func TestHandleUSRDispatch(t *testing.T) {
	tests := []struct {
		arguments             string
		expectedTransactionID string
		expectedAuthParams    AuthParams
		ok                    bool
	}{
		{"2 MD5 I example@passport.com", "2", AuthParams{authMethod: "MD5", authState: "I", email: "example@passport.com", password: ""}, true},
		{"MD5 I example@passport.com", "", AuthParams{}, false},
		{"2 MD5 I", "", AuthParams{}, false},
		{"2 MD5 I example@passport.com foo", "", AuthParams{}, false},
	}

	for _, tt := range tests {
		mockConn := &mockConn{}
		mockDB := &gorm.DB{}

		ap := &AuthParams{}

		gotTransactionID, gotErr := HandleReceiveUSR(mockConn, mockDB, ap, tt.arguments)

		if (gotErr == nil) != tt.ok {
			t.Errorf("Error HandleReceiveUSR(%q) = %v, want %v", tt.arguments, gotErr == nil, tt.ok)
		}

		if gotTransactionID != tt.expectedTransactionID {
			t.Errorf("TransactionID HandleReceiveUSR(%q) = %q, want %q", tt.arguments, gotTransactionID, tt.expectedTransactionID)
		}

		if ap.authMethod != tt.expectedAuthParams.authMethod {
			t.Errorf("AuthMethod HandleReceiveUSR(%q) = %q, want %q", tt.arguments, ap.authMethod, tt.expectedAuthParams.authMethod)
		}

		if ap.authState != tt.expectedAuthParams.authState {
			t.Errorf("AuthState HandleReceiveUSR(%q) = %q, want %q", tt.arguments, ap.authState, tt.expectedAuthParams.authState)
		}

		if ap.email != tt.expectedAuthParams.email {
			t.Errorf("Username HandleReceiveUSR(%q) = %q, want %q", tt.arguments, ap.email, tt.expectedAuthParams.email)
		}

	}
}
