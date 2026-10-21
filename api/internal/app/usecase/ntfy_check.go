package usecase

import (
	"server-alarm/api/internal/app/repo"
)

// NtfyCheckUC represents a usecases to check ntfy availability.
type NtfyCheckUC struct {
	repoNtfy *repo.AnyMsgRepoNtfy
}

// NewNtfyCheckUC returns a new instance of NtfyCheckUC.
func NewNtfyCheckUC(repoNtfy *repo.AnyMsgRepoNtfy) *NtfyCheckUC {
	return &NtfyCheckUC{
		repoNtfy: repoNtfy,
	}
}

// SendCheckMsg sends message to check ntfy availability.
func (u *NtfyCheckUC) SendCheckMsg() error {
	return u.repoNtfy.SendCheckMsg()
}
