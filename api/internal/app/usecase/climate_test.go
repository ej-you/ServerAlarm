package usecase

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"server-alarm/api/internal/app/repo"
	"server-alarm/api/internal/pkg/db"
	"server-alarm/api/internal/pkg/ntfy"
	"server-alarm/api/internal/pkg/storage"
)

var _testUC *ClimateUC

func TestCheckTemperature(t *testing.T) {
	dsn := os.Getenv("TEST_DSN")
	token := os.Getenv("TEST_TOKEN")
	require.NotZero(t, dsn, "TEST_DSN env var is not specified")
	require.NotZero(t, token, "TEST_TOKEN env var is not specified")

	location := time.FixedZone("UTC+3", 3*60*60)
	theme := "test"
	var (
		tresholdStandart float32 = 33.0
		tresholdHigh     float32 = 33.0
		tresholdUrgent   float32 = 33.0
	)

	storageInst := storage.NewKeyValueInMem()
	climateRepoCache := repo.NewClimateRepoCache(storageInst)

	dbInst, err := db.New(dsn)
	require.NoError(t, err, "create db instance")
	climateRepoDB := repo.NewClimateRepoDB(dbInst)

	ntfyClient, err := ntfy.NewClient(ntfy.WithTokenAuth(token),
		ntfy.WithHTTP(),
		ntfy.WithHost("127.0.0.1"),
		ntfy.WithPort("8888"))
	require.NoError(t, err, "create ntfy client")
	climateRepoNtfy := repo.NewClimateRepoNtfy(ntfyClient, location, theme)

	_testUC = NewClimateUC(climateRepoCache, climateRepoDB, climateRepoNtfy,
		tresholdStandart, tresholdHigh, tresholdUrgent)

	err = _testUC.CheckTemperature()
	require.NoError(t, err, "check temperature")
}

func TestCheckTemperatureAfterSleep(t *testing.T) {
	sleep := 10 * time.Second

	t.Logf("Sleep %v", sleep)
	time.Sleep(sleep)

	err := _testUC.CheckTemperature()
	require.NoError(t, err, "check temperature")
}
