package auth

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/palantir/stacktrace"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
)

type AuthDomain interface {
	Register(ctx context.Context, input *model.UserInput) (*model.User, error)
	Login(ctx context.Context, input *model.LoginInput) (*model.LoginResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*Claims, error)
	GetCurrentUser(ctx context.Context, userID string) (*model.User, error)
	GenerateToken(user *model.User) (string, int64, error)
}

type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type authDomain struct {
	databasePort outbound_port.DatabasePort
	messagePort  outbound_port.MessagePort
}

func NewAuthDomain(databasePort outbound_port.DatabasePort, messagePort outbound_port.MessagePort) AuthDomain {
	return &authDomain{
		databasePort: databasePort,
		messagePort:  messagePort,
	}
}

func (d *authDomain) Register(ctx context.Context, input *model.UserInput) (*model.User, error) {
	// Check if email already exists
	exists, err := d.databasePort.User().EmailExists(input.Email)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to check email")
	}
	if exists {
		return nil, stacktrace.NewError("email already registered")
	}

	user, err := model.UserPrepare(input)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to prepare user")
	}

	err = d.databasePort.User().Create(user)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to create user")
	}

	// Send Welcome WhatsApp Notification
	if d.messagePort != nil && user.WhatsApp != "" {
		message := "Halo " + user.Name + "!\n\nSelamat datang di EduVera. Akun admin Anda telah berhasil dibuat.\n\n" +
			"Email: " + user.Email + "\n" +
			"Silakan lanjutkan proses pendaftaran dengan melengkapi data institusi Anda.\n\n" +
			"Terima kasih!"
		_ = d.messagePort.WhatsApp().Send(user.WhatsApp, message)
	}

	return user, nil
}

func (d *authDomain) Login(ctx context.Context, input *model.LoginInput) (*model.LoginResponse, error) {
	user, err := d.databasePort.User().FindByEmail(input.Email)
	if err != nil {
		return nil, stacktrace.NewError("invalid email or password")
	}

	if !user.CheckPassword(input.Password) {
		return nil, stacktrace.NewError("invalid email or password")
	}

	if !user.IsActive {
		return nil, stacktrace.NewError("account not activated")
	}

	// Update last login
	_ = d.databasePort.User().UpdateLastLogin(user.ID)

	// Generate JWT
	signedToken, expiresAt, err := d.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		User:        *user,
		AccessToken: signedToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func (d *authDomain) GenerateToken(user *model.User) (string, int64, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "eduvera",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(getJWTSecret()))
	if err != nil {
		return "", 0, stacktrace.Propagate(err, "failed to sign token")
	}

	return signedToken, expiresAt.Unix(), nil
}

func (d *authDomain) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, stacktrace.NewError("unexpected signing method")
		}
		return []byte(getJWTSecret()), nil
	})

	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to parse token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, stacktrace.NewError("invalid token")
	}

	return claims, nil
}

func (d *authDomain) GetCurrentUser(ctx context.Context, userID string) (*model.User, error) {
	user, err := d.databasePort.User().FindByID(userID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to find user")
	}
	return user, nil
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "eduvera-default-secret-change-in-production"
	}
	return secret
}
