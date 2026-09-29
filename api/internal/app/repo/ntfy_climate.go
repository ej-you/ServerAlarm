package repo

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"server-alarm/api/internal/app/entity"
	"server-alarm/api/internal/pkg/ntfy"
)

const _datetimeLayout = "02-01-2006 15:04:05"

// ClimateRepoNtfy represents an ntfy repo for entity.Climate.
type ClimateRepoNtfy struct {
	client *ntfy.Client
	locale *time.Location
	theme  string
}

// NewClimateRepoNtfy returns a new instance of ClimateRepoNtfy.
func NewClimateRepoNtfy(client *ntfy.Client,
	locale *time.Location, theme string) *ClimateRepoNtfy {

	return &ClimateRepoNtfy{
		client: client,
		locale: locale,
		theme:  theme,
	}
}

// SendTempTresholdMsg sends message about critical temperature value.
func (r *ClimateRepoNtfy) SendTempTresholdMsg(climate *entity.Climate) error {
	title := "Критическая температура"
	priority := "high"
	tags := []string{"warning"}
	datetime := climate.Datetime.In(r.locale).Format(_datetimeLayout)

	var textBuilder strings.Builder
	textBuilder.WriteString("ВНИМАНИЕ! Температура в серверной достигла критической отметки.\n")
	textBuilder.WriteString("Значение: ")
	fmt.Fprintf(&textBuilder, "%.1f°C", climate.Temperature)
	textBuilder.WriteString(" | Дата: ")
	textBuilder.WriteString(datetime)

	answer, err := r.client.SendMsg(r.theme, title, priority, tags, textBuilder.String())
	if err != nil {
		return err
	}
	slog.Info("send temperature treshold message", "answer", string(answer))
	return nil
}
