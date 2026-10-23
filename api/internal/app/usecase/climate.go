package usecase

import (
	"fmt"
	"log/slog"

	"server-alarm/api/internal/app/entity"
	"server-alarm/api/internal/app/repo"
	"server-alarm/api/internal/pkg/utils/float"
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
	history, lastOld, lastNew, err := u.getClimates()
	if err != nil {
		return err
	}

	// send warning message if no new data
	if lastOld != nil && lastOld.Datetime.Equal(lastNew.Datetime) {
		err := u.climateRepoNtfy.SendNoDataMsg(lastNew)
		if err != nil {
			return fmt.Errorf("send ntfy message: %w", err)
		}
		return nil
	}
	// send warning message if data is repeated
	if u.isBrokenSensor(history, lastNew) {
		err := u.climateRepoNtfy.SendBrokenSensorMsg(lastNew)
		if err != nil {
			return fmt.Errorf("send ntfy message: %w", err)
		}
		return nil
	}

	lvl := u.tempThreshold.Level(lastNew.Temperature)
	lvlString := threshold.LevelString(lvl)
	// skip
	if u.skipCases(lastOld, lastNew, lvl) {
		return nil
	}
	// send message
	sendMsgLog(lastNew, lvl, true)
	if err := u.climateRepoNtfy.SendTempTresholdMsg(lastNew, lvlString); err != nil {
		return fmt.Errorf("send ntfy message: %w", err)
	}
	return nil
}

// getClimates get history and lastOld climate objects from cache and lastNew object from DB.
func (u *ClimateUC) getClimates() (history []*entity.Climate,
	lastOld, lastNew *entity.Climate, err error) {

	// from cache
	history, err = u.climateRepoCache.GetHistory()
	if err != nil {
		err = fmt.Errorf("get history from cache: %w", err)
		return
	}
	if len(history) != 0 {
		lastOld = history[0]
	}
	// from db
	lastNew, err = u.climateRepoDB.GetLastRecord()
	if err != nil {
		err = fmt.Errorf("get last record from db: %w", err)
		return
	}
	// set new one to cache history
	if err := u.climateRepoCache.SetHistory(lastNew); err != nil {
		slog.Warn("check temperature: set last record to cache history", "err", err)
	}
	return
}

// isBrokenSensor returns true if data from sensor is repeated.
func (u *ClimateUC) isBrokenSensor(history []*entity.Climate, lastNew *entity.Climate) bool {
	// skip if history is not full
	if len(history) < u.climateRepoCache.HistoryLen {
		return false
	}
	// check that all temperatures are equal to each other
	for _, elem := range history {
		if !float.Equal(lastNew.Temperature, elem.Temperature) {
			return false
		}
	}
	return true
}

// skipCases returns true on any skip case.
func (u *ClimateUC) skipCases(lastOld, lastNew *entity.Climate, lvl int) bool {
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
