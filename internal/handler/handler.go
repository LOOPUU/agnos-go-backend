package handler

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/LOOPUU/agnos-go-backend/internal/model"
	"github.com/gin-gonic/gin"
)

type AppService interface {
	CreateStaff(context.Context, string, string, string) (model.Staff, error)
	Login(context.Context, string, string, string) (model.Token, error)
	Authenticate(context.Context, string) (model.Staff, error)
	Search(context.Context, model.Staff, model.SearchFilters) ([]model.Patient, error)
}

type Handler struct {
	Svc AppService
}

type credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{3,100}$`)

func fail(c *gin.Context, status int, code, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"error": gin.H{"code": code, "message": msg}})
}

func validate(v credentials) error {
	if !usernamePattern.MatchString(v.Username) {
		return errors.New("invalid username")
	}
	if len(v.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func (h *Handler) CreateStaff(c *gin.Context) {
	var v credentials
	if c.ShouldBindJSON(&v) != nil || validate(v) != nil {
		fail(c, 422, "VALIDATION_ERROR", "invalid username, password, or hospital")
		return
	}

	s, err := h.Svc.CreateStaff(c, v.Hospital, v.Username, v.Password)
	if err != nil {
		fail(c, 409, "CREATE_STAFF_FAILED", err.Error())
		return
	}

	c.JSON(201, gin.H{"data": s})
}

func (h *Handler) Login(c *gin.Context) {
	var v credentials
	if c.ShouldBindJSON(&v) != nil || validate(v) != nil {
		fail(c, 422, "VALIDATION_ERROR", "invalid username, password, or hospital")
		return
	}

	t, err := h.Svc.Login(c, v.Hospital, v.Username, v.Password)
	if err != nil {
		fail(c, 401, "INVALID_CREDENTIALS", "invalid username, password, or hospital")
		return
	}

	c.JSON(200, gin.H{"data": t})
}

func (h *Handler) Search(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	parts := strings.Fields(auth)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		fail(c, 401, "UNAUTHORIZED", "valid Bearer token required")
		return
	}

	staff, err := h.Svc.Authenticate(c, parts[1])
	if err != nil {
		fail(c, 401, "UNAUTHORIZED", "invalid or expired token")
		return
	}

	var f model.SearchFilters
	if c.ShouldBindQuery(&f) != nil || f.Empty() {
		fail(c, 422, "VALIDATION_ERROR", "at least one search field is required")
		return
	}

	patients, err := h.Svc.Search(c, staff, f)
	if err != nil {
		fail(c, 502, "HOSPITAL_SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(200, gin.H{"data": patients})
}

func (h *Handler) Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/staff/create", h.CreateStaff)
	r.POST("/staff/login", h.Login)
	r.GET("/patient/search", h.Search)

	return r
}
