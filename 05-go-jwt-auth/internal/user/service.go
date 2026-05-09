package user

import (
	"context"
	"errors"
	"fmt"
	"jwt-auth/internal/auth"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repo

	jwtSecret string
}

type RegisterInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type LoginInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthResult struct {
	Token string

	u PublicUser
}

func NewService(repo *Repo, jwtSecret string) *Service {
	return &Service {
		repo: repo,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	// clean inputs
	email := strings.ToLower(strings.TrimSpace(input.Email));
	password := strings.TrimSpace(input.Password);

	// input validation
	if email == "" || password == "" {
		return AuthResult{}, fmt.Errorf("Email and Password are required!")
	}

	if len(password) < 8 {
		return AuthResult{}, fmt.Errorf("Password must be at least 8 characters!")
	}

	// check if user already exists
	_, err := s.repo.FindByEmail(ctx, email);
	if err == nil {
		return AuthResult{}, fmt.Errorf("User already exists! If that's you, please login, else, check if email is correct.")
	}

	if !errors.Is(err, mongo.ErrNoDocuments){
		return AuthResult{}, fmt.Errorf("Failed to verify if user already exists: %w", err)
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10);
	if err != nil {
		return AuthResult{}, fmt.Errorf("Password hashing failed: %w", err);
	}

	// create and update time
	now := time.Now().UTC();

	// create our user
	u :=  User{
		Email: email,
		Password: string(hashedPassword),
		Role: input.Role,
		Name: input.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	created, err := s.repo.Create(ctx, u);
	if err != nil {
		return AuthResult{}, err;
	}

	// create token
	token, err := auth.CreateToken(s.jwtSecret, created.ID.Hex(), created.Role);
	if err != nil {
		return AuthResult{}, err;
	}

	return AuthResult{
		u: ToPublic(created),
		Token: token,
	}, nil
}


func (s *Service) Login (ctx context.Context, input LoginInput) (AuthResult, error) {
	// clean inputs
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)

	// input validation
	if email == "" || password == "" {
		return AuthResult{}, fmt.Errorf("Email and Password are required!")
	}

	if len(password) < 8 {
		return AuthResult{}, fmt.Errorf("Password must be at least 8 characters!")
	}

	// Check if there is a valid account.
	u, err := s.repo.FindByEmail(ctx, email);
	if err != nil {
		return AuthResult{}, fmt.Errorf("Invalid Credentials! %w", err)
	}

	// verify password;
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return AuthResult{}, fmt.Errorf("Invalid Credentials! %w", err)
	}


	// generate token;
	token, err := auth.CreateToken(s.jwtSecret, u.ID.Hex(), u.Role);
	if err != nil {
		return AuthResult{}, err;
	}

	return AuthResult{
		u: ToPublic(u),
		Token: token,
	}, nil
}