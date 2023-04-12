package mocking

import (
	"database/sql"
	"fmt"
	"strings"
)

/*
Function to test.
Notice this the function to test, should accept the "DbHandlerIf" so that we can pass the mocked object.
If the function to test does not accept the handler, make changes to the function.
So that we can pass the mocked DB handler or mocked one.
*/

func handleDatabaseCreation(dbHandle DbHandlerIf, name string, dbPassword string, userCreated *bool) error {
	// Create the database
	result, err := dbHandle.CreateDatabase(name)
	if err != nil && strings.Contains(err.Error(), "already exists") {
		println("already exist")
	} else if result != nil && err == nil {
		println("Created database " + name)
	} else {
		println(err, "Failed to create database "+name)
		return err
	}

	// Create the database user
	var tmpBool bool
	result, err = dbHandle.CreateDatabaseUser(name)
	if err != nil && strings.Contains(err.Error(), "already exists") {
		println("Database user " + name + " already present")
	} else if result != nil && err == nil {
		println("Created database user " + name)
		tmpBool = true
		userCreated = &tmpBool
	} else {
		println(err, "Failed to create database user "+name)
		return err
	}

	// Create the database password only when the database user is created
	if *userCreated {
		result, err = dbHandle.AssignPassword(name, dbPassword)
		if err != nil {
			println(err, "Failed to create database user password")
			return err
		}
		println("Assigned database user password")
	}

	// Grant the privileges on the database to the user
	result, err = dbHandle.GrantPrivileges(name, name)
	if result != nil && err == nil {
		println("Grant privileges on database " + name + " to user " + name)
	} else {
		println(err, "Failed to grant privileges on database "+name+" to user "+name)
		return err
	}

	return nil
}

/*
First of all encapsulate all the function that you need to mock under an interface.
Mock generator will generate the NewMockDbHandlerIf() which we will use in the testing.
*/

//go:generate mockgen -package mocks -destination=./mocks/mock_dbhandler.go -source=db.go
type DbHandlerIf interface {
	CreateDatabase(dbName string) (sql.Result, error)
	CreateDatabaseUser(userName string) (sql.Result, error)
	AssignPassword(userName string, password string) (sql.Result, error)
	GrantPrivileges(dbName string, userName string) (sql.Result, error)
	TerminateOpenDatabaseConnection(dbName string) (sql.Result, error)
	DropDatabase(dbName string) (sql.Result, error)
	DropDatabaseUser(userName string) (sql.Result, error)
	Close()
}

/*
Create a struct which will implement the interface
*/
type DbHandler struct {
	db *sql.DB
}

/*
Create a New() which should return the pointer to the above struct
*/
func NewDbHandler() (*DbHandler, error) {
	// Get the connection to database cluster
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		"Hostname", "DbPort", "Username", "Password", "DbName", "SslMode")
	dbConnection, err := sql.Open("postgres", dbUrl)

	if err != nil {
		return nil, err
	}
	h := DbHandler{
		db: dbConnection,
	}
	return &h, nil
}

/*
Provide the minimal definition to them.
*/
func (h *DbHandler) CreateDatabase(dbName string) (sql.Result, error) {
	return h.db.Exec("CREATE DATABASE " + dbName)
}

func (h *DbHandler) CreateDatabaseUser(userName string) (sql.Result, error) {
	return h.db.Exec("CREATE USER " + userName)
}

func (h *DbHandler) AssignPassword(userName string, password string) (sql.Result, error) {
	return h.db.Exec("ALTER ROLE " + userName + " WITH PASSWORD '" + password + "'")
}

func (h *DbHandler) GrantPrivileges(dbName string, userName string) (sql.Result, error) {
	return h.db.Exec("GRANT ALL PRIVILEGES ON DATABASE " + dbName + " TO " + userName)
}

func (h *DbHandler) TerminateOpenDatabaseConnection(dbName string) (sql.Result, error) {
	return h.db.Exec("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '" + dbName + "';")
}

func (h *DbHandler) DropDatabase(dbName string) (sql.Result, error) {
	return h.db.Exec("DROP DATABASE " + dbName)
}

func (h *DbHandler) DropDatabaseUser(userName string) (sql.Result, error) {
	return h.db.Exec("DROP USER " + userName)
}

func (h *DbHandler) Close() {
	h.db.Close()
}

/*
Make the call to the NewDbHandler() and get the handler
Pass that handler to the function that need to be tested
*/
func handleDatabaseArtifacts(name string, dbPassword string, userCreated *bool) error {
	// Get the connection to database cluster
	dbHandler, err := NewDbHandler()
	if err != nil {
		return err
	}
	defer dbHandler.Close()
	if err := handleDatabaseCreation(dbHandler, name, dbPassword, userCreated); err != nil {
		return err
	}
	return nil
}
