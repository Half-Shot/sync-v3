package sync3

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"testing"
)

var postgresConnectionString = "user=xxxxx dbname=syncv3_test sslmode=disable"

func createLocalDB() string {
	fmt.Println("Note: tests require a postgres install accessible to the current user")
	dbName := "syncv3_sync3_test"
	exec.Command("dropdb", dbName).Run()
	if err := exec.Command("createdb", dbName).Run(); err != nil {
		fmt.Println("createdb failed: ", err)
		os.Exit(2)
	}
	return dbName
}

func currentUser() string {
	user, err := user.Current()
	if err != nil {
		fmt.Println("cannot get current user: ", err)
		os.Exit(2)
	}
	return user.Username
}

func prepareDBConnectionString() (connStr string) {
	// Required vars: user and db
	// We'll try to infer from the local env if they are missing
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = currentUser()
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		dbName = createLocalDB()
	}
	connStr = fmt.Sprintf(
		"user=%s dbname=%s sslmode=disable",
		user, dbName,
	)
	// optional vars, used in CI
	password := os.Getenv("POSTGRES_PASSWORD")
	if password != "" {
		connStr += fmt.Sprintf(" password=%s", password)
	}
	host := os.Getenv("POSTGRES_HOST")
	if host != "" {
		connStr += fmt.Sprintf(" host=%s", host)
	}
	return
}

func TestMain(m *testing.M) {
	postgresConnectionString = prepareDBConnectionString()
	exitCode := m.Run()
	os.Exit(exitCode)
}