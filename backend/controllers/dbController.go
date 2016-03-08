package controllers

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DBController struct {
	databasePath string
	connection   *sql.DB
}

func NewDBController(path string) (*DBController, error) {
	c := &DBController{
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

func (c *DBController) Open() error {
	con, err := sql.Open("sqlite3", c.databasePath)
	if err != nil {
		return err
	}
	c.connection = con
	return nil
}

func (c *DBController) Connection() *sql.DB {
	return c.connection
}

func (c *DBController) Close() error {
	return c.Connection().Close()
}
