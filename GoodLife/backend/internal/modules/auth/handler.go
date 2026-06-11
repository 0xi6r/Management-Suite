package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"

	"github.com/0xi6r/Management-Suite/GoodLife/backend/internal/middleware"
)

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
