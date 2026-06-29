package service

import (
	"errors"
	"os"
	"spotSync-golang-Project/internal/dto" 
	"spotSync-golang-Project/internal/models" 
	"spotSync-golang-Project/internal/repository" 
	"time" 

	"github.com/golang-jwt/jwt/v5" 
	"golang.org/x/crypto/bcrypt" 
	"gorm.io/gorm" 
)

type AuthService interface { // Interface for auth service (অথ সার্ভিসের ইন্টারফেস)
	Register(req dto.RegisterRequest) (*dto.AuthUserResponse, error) // Register method (রেজিস্টার মেথড)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error) // Login method (লগইন মেথড)
}

type authService struct { // Struct implementing interface (স্ট্রাকচার)
	userRepo repository.UserRepository // User repository (ইউজার রিপোজিটরি)
}

func NewAuthService(repo repository.UserRepository) AuthService { // Constructor (কনস্ট্রাক্টর)
	return &authService{repo} // Return instance (ইন্সট্যান্স রিটার্ন)
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthUserResponse, error) { // Register implementation (রেজিস্টার এর কাজ)
	// check if user already exists (আগে থেকে ইউজার আছে কিনা চেক করা)
	_, err := s.userRepo.FindByEmail(req.Email) // Find by email (ইমেইল দিয়ে খোঁজা)
	if err == nil { // If no error, means user exists (এরর না থাকা মানে ইউজার আছে)
		return nil, errors.New("email already in use") // Return error (এরর রিটার্ন)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) { // If other database error (অন্য কোনো ডাটাবেস এরর হলে)
		return nil, err // Return error (এরর রিটার্ন)
	}

	// Hash password (পাসওয়ার্ড হ্যাশ করা)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12) // Hash with cost 12 (১২ কস্টে হ্যাশ তৈরি)
	if err != nil { // If hash fails (হ্যাশ করতে ব্যর্থ হলে)
		return nil, err // Return error (এরর রিটার্ন)
	}

	role := "driver" // Default role (ডিফল্ট রোল)
	if req.Role != "" { // If role provided (যদি রোল দেওয়া থাকে)
		role = req.Role // Set role (রোল সেট করা)
	}

	user := models.User{ // Create user model instance (ইউজার মডেল তৈরি)
		Name:     req.Name, // Set name (নাম সেট করা)
		Email:    req.Email, // Set email (ইমেইল সেট করা)
		Password: string(hashedPassword), // Set hashed password (হ্যাশ করা পাসওয়ার্ড সেট করা)
		Role:     role, // Set role (রোল সেট করা)
	}

	if err := s.userRepo.CreateUser(&user); err != nil { // Save to database (ডাটাবেসে সেভ করা)
		return nil, err // Return if error (এরর হলে রিটার্ন)
	}

	return &dto.AuthUserResponse{ // Return response DTO (রেসপন্স ডিটিও রিটার্ন)
		ID:        user.ID, // User ID (ইউজার আইডি)
		Name:      user.Name, // Name (নাম)
		Email:     user.Email, // Email (ইমেইল)
		Role:      user.Role, // Role (রোল)
		CreatedAt: user.CreatedAt, // Creation time (তৈরির সময়)
		UpdatedAt: user.UpdatedAt, // Update time (আপডেটের সময়)
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) { // Login implementation (লগইন এর কাজ)
	user, err := s.userRepo.FindByEmail(req.Email) // Find user by email (ইমেইল দিয়ে ইউজার খোঁজা)
	if err != nil { // If error (যদি এরর হয়)
		if errors.Is(err, gorm.ErrRecordNotFound) { // If not found (না পাওয়া গেলে)
			return nil, errors.New("invalid email or password") // Invalid credentials (ভুল ইমেইল বা পাসওয়ার্ড)
		}
		return nil, err // Return other errors (অন্য এরর রিটার্ন)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil { // Compare password hashes (পাসওয়ার্ড মিলিয়ে দেখা)
		return nil, errors.New("invalid email or password") // Incorrect password (ভুল পাসওয়ার্ড)
	}

	// Generate JWT (JWT তৈরি করা)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // Use HS256 algorithm (HS256 অ্যালগরিদম ব্যবহার করা)
		"id":   user.ID, // Embed user ID (ইউজার আইডি রাখা)
		"role": user.Role, // Embed user role (ইউজারের রোল রাখা)
		"exp":  time.Now().Add(time.Hour * 72).Unix(), // Set expiration time to 72 hours (৭২ ঘণ্টা মেয়াদ দেওয়া)
	})

	secret := os.Getenv("JWT_SECRET") // Read secret key (সিক্রেট কি পড়া)
	if secret == "" { // Check if empty (খালি কিনা চেক করা)
		return nil, errors.New("JWT_SECRET not set") // Error if missing (না থাকলে এরর)
	}

	t, err := token.SignedString([]byte(secret)) // Sign the token with secret (সিক্রেট দিয়ে টোকেন সাইন করা)
	if err != nil { // If error signing (সাইন করতে সমস্যা হলে)
		return nil, err // Return error (এরর রিটার্ন)
	}

	return &dto.LoginResponse{ // Return response (রেসপন্স রিটার্ন)
		Token: t, // The JWT string (JWT স্ট্রিং)
		User: dto.AuthUserResponse{ // User details (ইউজারের তথ্য)
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
