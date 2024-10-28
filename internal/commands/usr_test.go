package commands

import (
	"testing"
)

func TestHandleUSRDispatch(t *testing.T) {
	tests := []struct {
		arguments             string
		expectedTransactionID string
		expectedAuthParams    authParams
		ok                    bool
	}{
		{"2 MD5 I example@passport.com", "2", authParams{authMethod: "MD5", authState: "I", username: "example@passport.com"}, true},
		{"MD5 I example@passport.com", "", authParams{}, false},
		{"2 MD5 I", "", authParams{}, false},
		{"2 MD5 I example@passport.com foo", "", authParams{}, false},
	}

	for _, tt := range tests {
		mock := &mockConn{}
		gotTransactionID, gotAuthParams, gotErr := HandleReceiveUSR(mock, tt.arguments)

		if (gotErr == nil) != tt.ok {
			t.Errorf("Error HandleReceiveUSR(%q) = %v, want %v", tt.arguments, gotErr == nil, tt.ok)
		}

		if gotTransactionID != tt.expectedTransactionID {
			t.Errorf("TransactionID HandleReceiveUSR(%q) = %q, want %q", tt.arguments, gotTransactionID, tt.expectedTransactionID)
		}

		if gotAuthParams.authMethod != tt.expectedAuthParams.authMethod {
			t.Errorf("AuthMethod HandleReceiveUSR(%q) = %q, want %q", tt.arguments, gotAuthParams.authMethod, tt.expectedAuthParams.authMethod)
		}

		if gotAuthParams.authState != tt.expectedAuthParams.authState {
			t.Errorf("AuthState HandleReceiveUSR(%q) = %q, want %q", tt.arguments, gotAuthParams.authState, tt.expectedAuthParams.authState)
		}

		if gotAuthParams.username != tt.expectedAuthParams.username {
			t.Errorf("Username HandleReceiveUSR(%q) = %q, want %q", tt.arguments, gotAuthParams.username, tt.expectedAuthParams.username)
		}

	}
}
