// package middleware

// import (
//     "fmt"
//     "net/http"
//     "os"
//     "strings"

//     "github.com/gin-gonic/gin"
//     "github.com/golang-jwt/jwt/v5"
// )

// func jwtSecret() []byte {
//     s := os.Getenv("JWT_SECRET")
//     if s == "" {
//         s = "replace-with-secure-secret"
//     }
//     return []byte(s)
// }

// func AuthMiddleware() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         auth := c.GetHeader("Authorization")
//         if auth == "" {
//             c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
//             return
//         }
//         parts := strings.SplitN(auth, " ", 2)
//         if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
//             c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
//             return
//         }
//         tokenString := parts[1]

//         claims := jwt.MapClaims{}
//         token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
//             if t.Method != jwt.SigningMethodHS256 {
//                 return nil, fmt.Errorf("unexpected signing method")
//             }
//             return jwtSecret(), nil
//         })
//         if err != nil || !token.Valid {
//             c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
//             return
//         }

//         c.Set("claims", claims)
//         c.Next()
//     }
// }

package middleware

import (
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func jwtSecret() []byte {
    s := os.Getenv("JWT_SECRET")
    if s == "" {
        s = "replace-with-secure-secret"
    }
    return []byte(s)
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            return
        }

        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
            return
        }
        tokenString := parts[1]

        claims := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
            if t.Method != jwt.SigningMethodHS256 {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return jwtSecret(), nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        // 🔹 ดึง sut_id จาก claims
        sutID, _ := claims["sut_id"].(string)
        if sutID == "" {
            // กันเคส token ไม่มี sut_id จริง ๆ
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing sut_id in token"})
            return
        }

        // เซ็ตลง context ให้ controller ดึงไปใช้ได้
        c.Set("sut_id", sutID)

        // ถ้าอยากใช้ user_id / role ที่อื่น ก็เซ็ตเพิ่มได้
        if uid, ok := claims["user_id"]; ok {
            c.Set("user_id", uid)
        }
        if role, ok := claims["role"]; ok {
            c.Set("role", role)
        }

        // เผื่ออยากอ่าน claims ดิบ ๆ ที่อื่น
        c.Set("claims", claims)

        c.Next()
    }
}
