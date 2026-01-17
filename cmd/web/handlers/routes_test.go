package handlers

import (
	"echohub-community/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomePage(t *testing.T) {
	// Since WebApp depends on models which we've tested separately,
	// we can do a basic test here. Complex integration tests would need
	// a setup similar to common_test.go but accessible here.

	// For now, let's just test that the router is correctly configured
	// and doesn't panic on a basic request.

	db := models.SetupTestDB(t)
	defer db.Close()

	app := &WebApp{
		Users: &models.UserModel{DB: db},
	}
	mux := app.Router()

	req, _ := http.NewRequest("GET", "/login", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	// We expect a 200 OK for the login page
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
