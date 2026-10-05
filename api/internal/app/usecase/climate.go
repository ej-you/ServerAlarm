package usecase

import (
	"fmt"
	"log/slog"

	"server-alarm/api/internal/app/entity"
	"server-alarm/api/internal/app/repo"
	"server-alarm/api/internal/pkg/utils/threshold"
)

// ClimateUC represents a usecases for entity.Climate.
type ClimateUC struct {
	climateRepoCache *repo.ClimateRepoCache
	climateRepoDB    *repo.ClimateRepoDB
	climateRepoNtfy  *repo.ClimateRepoNtfy
	tempThreshold    *threshold.Threshold[float32]
}

// NewClimateUC returns a new instance of ClimateUC.
func NewClimateUC(climateRepoCache *repo.ClimateRepoCache,
	climateRepoDB *repo.ClimateRepoDB, climateRepoNtfy *repo.ClimateRepoNtfy,
	tempThresholdStandart, tempThresholdHigh, tempThresholdUrgent float32) *ClimateUC {

	return &ClimateUC{
		climateRepoCache: climateRepoCache,
		climateRepoDB:    climateRepoDB,
		climateRepoNtfy:  climateRepoNtfy,
		tempThreshold: threshold.New(
			tempThresholdStandart, tempThresholdHigh, tempThresholdUrgent,
		),
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

	lvl := u.tempThreshold.Level(lastNew.Temperature)
	lvlString := threshold.LevelString(lvl)
	// skip
	if u.skipCases(lastOld, lastNew, lvl) {
		return nil
	}
	// send message
	sendMsgLog(lastNew, lvl, true)
	if err = u.climateRepoNtfy.SendTempTresholdMsg(lastNew, lvlString); err != nil {
		return fmt.Errorf("send ntfy message: %w", err)
	}
	return nil
}

// skipCases returns true on any skip case.
func (u *ClimateUC) skipCases(lastOld, lastNew *entity.Climate, lvl int) bool {
	// no new data (skip if level is not Urgent)
	if lvl != threshold.Urgent && lastOld != nil && lastOld.Datetime.Equal(lastNew.Datetime) {
		slog.Debug("check temperature: no new data received from db")
		return true
	}
	// skip if threshold was not exceeded
	if lvl == threshold.Zero {
		sendMsgLog(lastNew, threshold.Zero, false)
		return true
	}

	if lastOld == nil {
		return false
	}
	// skip if last old level is not Zero
	if lvl == threshold.Standart && u.tempThreshold.Level(lastOld.Temperature) != threshold.Zero {
		sendMsgLog(lastNew, threshold.Standart, false)
		return true
	}
	// skip if last old level is not Standart
	if lvl == threshold.High && u.tempThreshold.Level(lastOld.Temperature) != threshold.Standart {
		sendMsgLog(lastNew, threshold.High, false)
		return true
	}
	return false
}

// sendMsgLog create info log about tempreture checking.
func sendMsgLog(climate *entity.Climate, thresholdLevel int, sendMsg bool) {
	slog.Info("check temperature",
		"temperature", climate.TemperatureString(),
		"datetime", climate.Datetime,
		"threshold exceeded", threshold.LevelString(thresholdLevel),
		"send message", sendMsg)
}
