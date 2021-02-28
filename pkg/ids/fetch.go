package ids

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

// FetchIdsModel は検知モデルをDBから再取得します
func FetchIdsModel() {
	// TODO: impelement
	log.Print("Fetching ids model...")
	time.Sleep(3 * time.Second)
	log.Print("done")
}

func fetchAnalysisResult(dbHost, serverID string) string {
	config := mysql.NewConfig()

	config.User = ""
	config.Passwd = ""
	config.DBName = ""
	config.Net = ""
	config.Addr = ""

	db, e := sql.Open("mysql", config.FormatDSN())
	if e != nil {
		log.Print("failed to connect db server. cannot use ids features")
	}

	tableName := fmt.Sprintf("%s_%s", serverID, "threshold")

	sql := fmt.Sprintf("SELECT * FROM %s WHERE `updateAt` < %d",
		tableName, time.Now().Truncate(24*time.Hour).Unix())

	rows, e := db.Query(sql)
	if e != nil {
		log.Print("failed to querying. disable ids features")
	}

	// TODO: implement rows to threshold.Threshold orm mapper

	for rows.Next() {

	}

	return sql
}
