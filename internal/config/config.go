package config

import (
	"fmt"
)

// HTTPServer - minimal config to start Go HTTP server
type HTTPServer struct {
	Host    string
	Port    int
	BaseURI string
}

// Database - db configuration
type Database struct {
	// Type - db type, only mysql supported now
	Type     string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	// Options - db connections options
	Options string
	// TablePrefix - helps avoid conflicts with existing tables
	TablePrefix string
}

// Application - params and preferences for all applications
type Application struct {
	CounterREST HTTPServer
	CounterDB   Database
	// CounterID - maintained counter ID
	CounterID int
}

// DSN - formats connection string based on configuration.
func (db Database) DSN() string {
	// only mysql is supported
	// mysql://aura:aura@tcp(127.0.0.1:3306)/aura?parseTime=true&timeout=3m
	// mysql-prefix is not needed for driver, it used for uniformity
	return fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%d)/%s?%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.Options,
	)
}
