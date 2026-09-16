package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LOOPUU/agnos-go-backend/internal/model"
)

type mockService struct {
	createErr error
	loginErr  error
	authErr   error
	searchErr error
}

func (m *mockService) CreateStaff(ctx context.Context, hospital, username, password string) (model.Staff, error) {
	return model.Staff{ID: 1, HospitalID: 1, Username: "tester"}, m.createErr
}

func (m *mockService) Login(ctx context.Context, hospital, username, password string) (model.Token, error) {
	return model.Token{AccessToken: "token", TokenType: "Bearer", ExpiresAt: time.Now().Add(time.Hour)}, m.loginErr
}

func (m *mockService) Authenticate(ctx context.Context, token string) (model.Staff, error) {
	return model.Staff{ID: 1, HospitalID: 1}, m.authErr
}

func (m *mockService) Search(ctx context.Context, s model.Staff, f model.SearchFilters) ([]model.Patient, error) {
	return []model.Patient{{PatientHN: "HN-1"}}, m.searchErr
}

func perform(h *Handler, method, path, body, token string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	h.Router().ServeHTTP(w, r)
	return w
}

func TestCreateStaffPositive(t *testing.T) {
	w := perform(&Handler{Svc: &mockService{}}, "POST", "/staff/create", `{"username":"tester","password":"password1","hospital":"hospital-a"}`, "")
	if w.Code != 201 {
		t.Fatalf("got %d", w.Code)
	}
}

func TestCreateStaffNegative(t *testing.T) {
	w := perform(&Handler{Svc: &mockService{}}, "POST", "/staff/create", `{"username":"x","password":"1","hospital":""}`, "")
	if w.Code != 422 {
		t.Fatalf("got %d", w.Code)
	}
}

func TestLoginPositive(t *testing.T) {
	w := perform(&Handler{Svc: &mockService{}}, "POST", "/staff/login", `{"username":"tester","password":"password1","hospital":"hospital-a"}`, "")
	if w.Code != 200 {
		t.Fatalf("got %d", w.Code)
	}
}

func TestLoginNegative(t *testing.T) {
	w := perform(&Handler{Svc: &mockService{loginErr: errors.New("bad")}}, "POST", "/staff/login", `{"username":"tester","password":"password1","hospital":"hospital-a"}`, "")
	if w.Code != 401 {
		t.Fatalf("got %d", w.Code)
	}
}

func TestSearchPositive(t *testing.T) {
	w := perform(&Handler{Svc: &mockService{}}, "GET", "/patient/search?national_id=123", ``, "token")
	if w.Code != 200 {
		t.Fatalf("got %d", w.Code)
	}
	var b map[string][]model.Patient
	if json.Unmarshal(w.Body.Bytes(), &b) != nil || len(b["data"]) != 1 {
		t.Fatal("unexpected body")
	}
}

func TestSearchNegativeNoLogin(t *testing.T) {
	w := perform(&Handler{Svc: &mockService{}}, "GET", "/patient/search?national_id=123", ``, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("got %d", w.Code)
	}
}

func TestSearchNegativeNoFilters(t *testing.T) {
	w := perform(&Handler{Svc: &mockService{}}, "GET", "/patient/search", ``, "token")
	if w.Code != 422 {
		t.Fatalf("got %d", w.Code)
	}
}
