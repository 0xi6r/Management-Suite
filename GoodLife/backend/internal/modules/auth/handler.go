package auth

import (
	"encoding/json"
	"net/http"
	"time"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"

	"github.com/0xi6r/Management-Suite/GoodLife/backend/internal/middleware"
)

// register section
type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type registerResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Basic validation
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.FullName = strings.TrimSpace(req.FullName)

	if req.Email == "" || req.Password == "" || req.FullName == "" {
		http.Error(w, `{"error":"email, password, and full_name are required"}`, http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, `{"error":"password must be at least 8 characters"}`, http.StatusBadRequest)
		return
	}
	// Simple email format check (not exhaustive)
	if !strings.Contains(req.Email, "@") {
		http.Error(w, `{"error":"invalid email format"}`, http.StatusBadRequest)
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("password hashing failed", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Insert user
	var userID string
	err = h.pool.QueryRow(r.Context(),
		`INSERT INTO auth.users (email, password_hash, full_name)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		req.Email, string(hash), req.FullName,
	).Scan(&userID)

	if err != nil {
		// Check for duplicate email (PostgreSQL error code 23505)
		if strings.Contains(err.Error(), "23505") {
			http.Error(w, `{"error":"email already registered"}`, http.StatusConflict)
			return
		}
		h.logger.Error("user insertion failed", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

		// Assign the "patient" role to the new user
	_, err = h.pool.Exec(r.Context(),
		`INSERT INTO auth.user_roles (user_id, role_id)
		SELECT $1, id FROM auth.roles WHERE name = 'patient'`,
		userID,
	)
	if err != nil {
		// Log the error but don't fail the registration (the user is still created).
		// In production you'd want to alert, but for now we log and continue.
		h.logger.Error("failed to assign patient role",
			zap.String("user_id", userID),
			zap.Error(err),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registerResponse{
		Message: "registration successful",
		UserID:  userID,
	})
}


// login and other
type Handler struct {
	pool      *pgxpool.Pool
	jwtSecret []byte
	logger    *zap.Logger
}

func NewHandler(pool *pgxpool.Pool, jwtSecret string, logger *zap.Logger) *Handler {
	return &Handler{
		pool:      pool,
		jwtSecret: []byte(jwtSecret),
		logger:    logger,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"email and password are required"}`, http.StatusBadRequest)
		return
	}

	// Look up user in the database
	var userID string
	var passwordHash string
	err := h.pool.QueryRow(r.Context(),
		`SELECT id, password_hash FROM auth.users WHERE email = $1 AND is_active = true`,
		req.Email,
	).Scan(&userID, &passwordHash)

	if err != nil {
		h.logger.Warn("login failed: user not found", zap.String("email", req.Email))
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		h.logger.Warn("login failed: wrong password", zap.String("email", req.Email))
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Generate JWT
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(h.jwtSecret)
	if err != nil {
		h.logger.Error("failed to sign JWT", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Optional: insert a session record into auth.sessions
	// (skipped for now, but you can add it later)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{Token: signedToken})
}


func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string) // safe because middleware enforced
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"user_id": userID,
	})
}
