package ids

import (
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hiragi-gkuth/bitris-analyzer/pkg/threshold"
)

// ReconstructIdsModel は検知モデルをDBから再取得します
func ReconstructIdsModel(serverID, dbHost, dbUser string, dbPass string) {

	repository := threshold.NewRepository(serverID, mysql.Config{
		Addr:   dbHost,
		Net:    "tcp",
		DBName: "bitris",
		User:   dbUser,
		Passwd: dbPass,
	})
	newModel := repository.FetchModel(16, 24*time.Hour, 24)

	GetModel().Write(newModel)
}
