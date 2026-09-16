package handler

import (
	"context"
	"errors"
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

func (m *mockService) CreateStaff(
	ctx context.Context,
	hospital, username, password string,
) (model.Staff, error) {
	return model.Staff{
		ID:         1,
		HospitalID: 1,
		Username:   "tester",
	}, m.createErr
}

func (m *mockService) Login(
	ctx context.Context,
	hospital, username, password string,
) (model.Token, error) {
	return model.Token{
		AccessToken: "token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Hour),
	}, m.loginErr
}

func (m *mockService) Authenticate(
	ctx context.Context,
	token string,
) (model.Staff, error) {
	return model.Staff{
		ID:         1,
		HospitalID: 1,
	}, m.authErr
}

func (m *mockService) Search(
	ctx context.Context,
	staff model.Staff,
	filters model.SearchFilters,
) ([]model.Patient, error) {
	return []model.Patient{
		{PatientHN: "HN-1"},
	}, m.searchErr
}

func perform(
	h *Handler,
	method, path, body, authorization string,
) *httptest.ResponseRecorder {

	r := httptest.NewRequest(
		method,
		path,
		strings.NewReader(body),
	)

	r.Header.Set("Content-Type", "application/json")

	if authorization != "" {
		r.Header.Set("Authorization", authorization)
	}

	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	return w
}

func check(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()

	if w.Code != want {
		t.Fatalf("got %d, want %d", w.Code, want)
	}
}

// Health
func TestHealth(t *testing.T) {
	w := perform(
		&Handler{Svc: &mockService{}},
		"GET",
		"/health",
		"",
		"",
	)

	check(t, w, 200)
}

// Create Staff
func TestCreateStaff(t *testing.T) {
	valid := `{"username":"tester","password":"password1","hospital":"hospital-a"}`

	tests := []struct {
		name string
		body string
		svc  *mockService
		want int
	}{
		{
			name: "success",
			body: valid,
			svc:  &mockService{},
			want: 201,
		},
		{
			name: "invalid json",
			body: `{`,
			svc:  &mockService{},
			want: 422,
		},
		{
			name: "invalid username",
			body: `{"username":"x","password":"password1","hospital":"hospital-a"}`,
			svc:  &mockService{},
			want: 422,
		},
		{
			name: "short password",
			body: `{"username":"tester","password":"123","hospital":"hospital-a"}`,
			svc:  &mockService{},
			want: 422,
		},
		{
			name: "create error",
			body: valid,
			svc: &mockService{
				createErr: errors.New("create error"),
			},
			want: 409,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := perform(
				&Handler{Svc: tt.svc},
				"POST",
				"/staff/create",
				tt.body,
				"",
			)

			check(t, w, tt.want)
		})
	}
}

// Login
func TestLogin(t *testing.T) {
	valid := `{"username":"tester","password":"password1","hospital":"hospital-a"}`

	tests := []struct {
		name string
		body string
		svc  *mockService
		want int
	}{
		{
			name: "success",
			body: valid,
			svc:  &mockService{},
			want: 200,
		},
		{
			name: "invalid json",
			body: `{`,
			svc:  &mockService{},
			want: 422,
		},
		{
			name: "invalid username",
			body: `{"username":"x","password":"password1","hospital":"hospital-a"}`,
			svc:  &mockService{},
			want: 422,
		},
		{
			name: "short password",
			body: `{"username":"tester","password":"123","hospital":"hospital-a"}`,
			svc:  &mockService{},
			want: 422,
		},
		{
			name: "login error",
			body: valid,
			svc: &mockService{
				loginErr: errors.New("bad login"),
			},
			want: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := perform(
				&Handler{Svc: tt.svc},
				"POST",
				"/staff/login",
				tt.body,
				"",
			)

			check(t, w, tt.want)
		})
	}
}

// Patient Search
func TestSearch(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		auth  string
		svc   *mockService
		want  int
	}{
		{
			name: "success",
			path: "/patient/search?national_id=123",
			auth: "Bearer token",
			svc:  &mockService{},
			want: 200,
		},
		{
			name: "no token",
			path: "/patient/search?national_id=123",
			auth: "",
			svc:  &mockService{},
			want: 401,
		},
		{
			name: "wrong authorization",
			path: "/patient/search?national_id=123",
			auth: "Basic token",
			svc:  &mockService{},
			want: 401,
		},
		{
			name: "invalid token",
			path: "/patient/search?national_id=123",
			auth: "Bearer bad",
			svc: &mockService{
				authErr: errors.New("invalid token"),
			},
			want: 401,
		},
		{
			name: "no filter",
			path: "/patient/search",
			auth: "Bearer token",
			svc:  &mockService{},
			want: 422,
		},
		{
			name: "hospital error",
			path: "/patient/search?national_id=123",
			auth: "Bearer token",
			svc: &mockService{
				searchErr: errors.New("hospital error"),
			},
			want: 502,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := perform(
				&Handler{Svc: tt.svc},
				"GET",
				tt.path,
				"",
				tt.auth,
			)

			check(t, w, tt.want)
		})
	}
}