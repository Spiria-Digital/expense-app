package storage

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func NewBunDB(filePath ...string) (*bun.DB, error) {
	dataSource := "file::memory:?cache=shared&_fk=1"
	if len(filePath) > 0 {
		dataSource = fmt.Sprintf("file:%s?cache=shared&_fk=1", filePath[0])
	}
	conn, err := sql.Open(sqliteshim.ShimName, dataSource)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(1000)
	conn.SetConnMaxLifetime(0)

	return bun.NewDB(conn, sqlitedialect.New()), nil
}

func GetRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	directory := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	// check if the directory exists
	_, err := os.Stat(directory)
	if err != nil && os.IsNotExist(err) {
		directory, _ = os.Getwd()
		log.Debug().Msgf("defaulting to current working directory: %s", directory)
	} else {
		log.Debug().Msgf("using root directory: %s", directory)
	}
	return directory
}
