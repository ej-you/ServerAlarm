package ntfy

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendMsg(t *testing.T) {
	token := os.Getenv("TEST_TOKEN")
	require.NotZero(t, token, "TEST_TOKEN env var is not specified")

	client, err := NewClient(WithTokenAuth(token),
		WithHTTP(),
		WithHost("127.0.0.1"),
		WithPort("8888"))
	if err != nil {
		t.Fatalf("Create ntfy client: %v", err)
	}

	msg, err := client.SendMsg("test", "TEST", "4", []string{"warning"}, "Pay attention, please!")
	if err != nil {
		t.Fatalf("Send message: %v", err)
	}
	t.Logf("Answer: %s", msg)
}
