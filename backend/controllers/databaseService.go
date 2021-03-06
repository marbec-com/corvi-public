package controllers

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type DatabaseService interface {
	Open() error
	Connection() *sql.DB
	Close() error
	Destroy() error
}

type SQLiteDBService struct {
	databasePath string
	connection   *sql.DB
}

func NewSQLiteDBService(path string) (*SQLiteDBService, error) {
	c := &SQLiteDBService{
		databasePath: path,
		connection:   nil,
	}

	// Open database, create if not exists
	err := c.Open()
	if err != nil {
		return nil, err
	}

	// Test database connection
	err = c.connection.Ping()
	if err != nil {
		return nil, err
	}

	// Enable foreign key constraints
	_, err = c.connection.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *SQLiteDBService) Open() error {
	con, err := sql.Open("sqlite3", c.databasePath)
	if err != nil {
		return err
	}
	c.connection = con
	return nil
}

func (c *SQLiteDBService) Connection() *sql.DB {
	return c.connection
}

func (c *SQLiteDBService) Close() error {
	return c.Connection().Close()
}

func (c *SQLiteDBService) Destroy() error {
	return os.Remove(c.databasePath)
}
