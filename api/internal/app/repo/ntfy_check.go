package repo

import (
	"log/slog"

	"server-alarm/api/internal/pkg/ntfy"
)

// AnyMsgRepoNtfy represents an ntfy repo to send any messages.
type AnyMsgRepoNtfy struct {
	client *ntfy.Client
	theme  string
}

// NewAnyMsgRepoNtfy returns a new instance of AnyMsgRepoNtfy.
func NewAnyMsgRepoNtfy(client *ntfy.Client, theme string) *AnyMsgRepoNtfy {
	return &AnyMsgRepoNtfy{
		client: client,
		theme:  theme,
	}
}

// SendCheckMsg sends message with the lowest priority to check ntfy availability.
func (r *AnyMsgRepoNtfy) SendCheckMsg() error {
	title := "Проверка"
	tags := []string{"paperclip"}
	priority := "min"
	text := "Регулярная проверка доступности сервиса. Не обращайте внимания"

	answer, err := r.client.SendMsg(r.theme, title, priority, tags, text)
	if err != nil {
		return err
	}
	slog.Info("send ntfy check message", "answer", string(answer))
	return nil
}
