package controller

import (
    "net/http"
    "os"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "backend/internal/service/users"
)

type AuthController struct {
    AuthService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
    return &AuthController{AuthService: authService}
}

type LoginRequest struct {
    SutID    string `json:"sut_id" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func jwtSecret() []byte {
    s := os.Getenv("JWT_SECRET")
    if s == "" {
        s = "replace-with-secure-secret"
    }
    return []byte(s)
}

// เพิ่ม struct claims แบบชัดเจน
type CustomClaims struct {
    UserID uint   `json:"user_id"`
    SutID  string `json:"sut_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func (ctr *AuthController) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    user, err := ctr.AuthService.Login(req.SutID, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
        return
    }

    now := time.Now()
    claims := CustomClaims{
        UserID: user.ID,
        SutID:  user.SutId,
        Role:   user.Role.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
            Issuer:    "your-app", // เปลี่ยนตามต้องการ
            Subject:   "auth-token",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(jwtSecret())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": signedToken,
        "user": gin.H{
            "id":    user.ID,
            "sutId": user.SutId,
            "role":  user.Role.Role,
        },
    })
}