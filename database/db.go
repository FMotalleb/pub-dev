package database

import (
	"gorm.io/gorm"

	"github.com/fmotalleb/north_outage/database/driver"
)

var rootDB *gorm.DB

func Connect(connection string) (*gorm.DB, error) {
	if rootDB != nil {
		return rootDB, nil
	}
	var conn gorm.Dialector
	var db *gorm.DB
	var err error
	if conn, err = driver.MakeConnection(connection); err != nil {
		return nil, err
	}
	if db, err = gorm.Open(conn, &gorm.Config{}); err != nil {
		return nil, err
	}
	rootDB = db
	return db, nil
}

// Get requires db to be [Connect]ed first.
// if called before a successful [Connect] will panic.
func Get() *gorm.DB {
	if rootDB == nil {
		panic("database is not initialized")
	}
	return rootDB
}
