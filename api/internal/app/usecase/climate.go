package usecase

import (
	"fmt"
	"log/slog"

	"server-alarm/api/internal/app/entity"
	"server-alarm/api/internal/app/repo"
)

// ClimateUC represents a usecases for entity.Climate.
type ClimateUC struct {
	climateRepoCache *repo.ClimateRepoCache
	climateRepoDB    *repo.ClimateRepoDB
	climateRepoNtfy  *repo.ClimateRepoNtfy
	tempTreshold     float32
}

// NewClimateUC returns a new instance of ClimateUC.
func NewClimateUC(climateRepoCache *repo.ClimateRepoCache, climateRepoDB *repo.ClimateRepoDB,
	climateRepoNtfy *repo.ClimateRepoNtfy, tempTreshold float32) *ClimateUC {

	return &ClimateUC{
		climateRepoCache: climateRepoCache,
		climateRepoDB:    climateRepoDB,
		climateRepoNtfy:  climateRepoNtfy,
		tempTreshold:     tempTreshold,
	}
}

// CheckTemperature compares actual temperature with treshold.
// Sends ntfy message if treshold exceeded.
func (u *ClimateUC) CheckTemperature() error {
	// from cache
	lastOld, err := u.climateRepoCache.Get()
	if err != nil {
		return fmt.Errorf("get last record from cache: %w", err)
	}
	// from db
	lastNew, err := u.climateRepoDB.GetLastRecord()
	if err != nil {
		return fmt.Errorf("get last record from db: %w", err)
	}
	// set new one to cache
	if err := u.climateRepoCache.Set(lastNew); err != nil {
		slog.Warn("check temperature: set last record to cache", "err", err)
	}

	// skip cases
	if lastOld != nil {
		if lastOld.Datetime.Equal(lastNew.Datetime) {
			slog.Debug("check temperature: no new data received from db")
			return nil
		}
		if lastOld.Temperature >= u.tempTreshold && lastNew.Temperature >= u.tempTreshold {
			sendMsgLog(lastNew, true, false)
			return nil
		}
	}
	if lastNew.Temperature < u.tempTreshold {
		sendMsgLog(lastNew, false, false)
		return nil
	}
	// send message
	sendMsgLog(lastNew, true, true)
	if err := u.climateRepoNtfy.SendTempTresholdMsg(lastNew); err != nil {
		return fmt.Errorf("send ntfy message: %w", err)
	}
	return nil
}

// sendMsgLog create info log about tempreture checking.
func sendMsgLog(climate *entity.Climate, tresholdExceeded, sendMsg bool) {
	slog.Info("check temperature",
		"temperature", climate.Temperature,
		"datetime", climate.Datetime,
		"treshold exceeded", tresholdExceeded,
		"send message", sendMsg)
}
