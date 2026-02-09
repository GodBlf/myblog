package test

import (
	"myblog/database"
	"testing"
)

func TestGetBlogDBConnection(t *testing.T) {
	blogDBConnection := database.GetBlogDBConnection()
	db, err := blogDBConnection.DB()
	if err != nil {
		t.Fatalf("Failed to get DB connection: %v", err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping DB: %v", err)
	}
}
