package handlers

import (
	"net/http"
	"strings"
	"time"
	"todo_api/internal/config"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type user_input struct {
	Mail     string `json:"mail" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Mail     string `json:"mail" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func TodoUsersHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input user_input

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(input.Password) < 8 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "password must be greater than 8 characters"})
			return
		}

		hashedPassword, err := hashPassword(input.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
			return
		}

		user, err := repository.CreateUser(pool, input.Mail, hashedPassword)

		if err != nil {
			if strings.Contains(err.Error(), "SQLSTATE 23505") {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "mail already registered"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, user)

	}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var Request LoginRequest
		if err := ctx.BindJSON(&Request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repository.GetUserByMail(pool, Request.Mail)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Request.Password))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Mail,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, LoginResponse{Token: tokenString})

	}
}
