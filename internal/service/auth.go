package service

import (
    "errors"
    "time"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "reflection-app/internal/models"
    "reflection-app/internal/repository"
)

type AuthService struct {
    userRepo  *repository.UserRepository
    jwtSecret []byte
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
    // In production, get this from environment variable
    jwtSecret := []byte("your-secret-key") // TODO: Move to environment variable
    return &AuthService{
        userRepo:  userRepo,
        jwtSecret: jwtSecret,
    }
}

// Register creates a new user account
func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
    // Check if user already exists
    existingUser, _ := s.userRepo.GetByUsername(req.Username)
    if existingUser != nil {
        return nil, errors.New("username already exists")
    }
    
    // Check if email already exists
    existingEmail, _ := s.userRepo.GetByEmail(req.Email)
    if existingEmail != nil {
        return nil, errors.New("email already exists")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    // Create user
    user := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        FirstName:    req.FirstName,
        LastName:     req.LastName,
        PasswordHash: string(hashedPassword),
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    // Generate JWT token
    token, err := s.generateJWT(user)
    if err != nil {
        return nil, err
    }
    
    return &models.AuthResponse{
        Token: token,
        User:  *user,
    }, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
    // Get user by username
    user, err := s.userRepo.GetByUsername(req.Username)
    if err != nil {
        return nil, errors.New("invalid username or password")
    }
    
    // Check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid username or password")
    }
    
    // Generate JWT token
    token, err := s.generateJWT(user)
    if err != nil {
        return nil, err
    }
    
    return &models.AuthResponse{
        Token: token,
        User:  *user,
    }, nil
}

// generateJWT creates a JWT token for the user
func (s *AuthService) generateJWT(user *models.User) (string, error) {
    claims := &models.JWTClaims{
        UserID:   user.ID,
        Username: user.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   string(rune(user.ID)),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return s.jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}

// GetUserByID retrieves a user by their ID
func (s *AuthService) GetUserByID(userID int) (*models.User, error) {
    return s.userRepo.GetByID(userID)
}