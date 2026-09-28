package ntfy

import (
	"os"
	"testing"
)

func TestSendMsg(t *testing.T) {
	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		t.Fatal("TEST_TOKEN env var is not specified")
	}

	client, err := NewNtfyClient(WithTokenAuth(token),
		WithHTTP(),
		WithHost("127.0.0.1"),
		WithPort("8888"))
	if err != nil {
		t.Fatalf("Create ntfy client: %v", err)
	}

	err = client.SendMsg("test", "TEST", "4", []string{"warning"}, "Pay attention, please!")
	if err != nil {
		t.Fatalf("Send message: %v", err)
	}
}
