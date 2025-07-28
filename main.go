package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"net/smtp"

	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a user in the system
// It will be stored in the MongoDB users collection
// Passwords are stored as SHA-256 hashes (Go stdlib only)
type User struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty"`
	Email               string             `bson:"email"`
	PasswordHash        string             `bson:"password_hash"`
	Verified            bool               `bson:"verified"`
	VerificationToken   string             `bson:"verification_token"`
	VerificationExpiry  time.Time          `bson:"verification_expiry"`
	PasswordResetToken  string             `bson:"password_reset_token"`
	PasswordResetExpiry time.Time          `bson:"password_reset_expiry"`
}

var verifyTmpl = template.Must(template.ParseFiles("templates/verify.html"))
var forgotTmpl = template.Must(template.ParseFiles("templates/forgot_password.html"))
var resetTmpl = template.Must(template.ParseFiles("templates/reset_password.html"))

func sendMail(to, subject, body string) error {
	from := os.Getenv("GMAIL_USER")
	pass := os.Getenv("GMAIL_PASS")
	smtpHost := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}

func hashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func verifyRequestHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			verifyTmpl.Execute(w, nil)
			return
		}
		if r.Method == http.MethodPost {
			var req struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"message": "Email and password required.", "status": "error"})
				return
			}
			email := req.Email
			password := req.Password
			coll := client.Database("practiceproject").Collection("users")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var user User
			err := coll.FindOne(ctx, map[string]interface{}{"email": email}).Decode(&user)
			if err != nil {
				// New user: create and send verification
				hash := hashPassword(password)
				token, err := generateToken()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"message": "Failed to generate token.", "status": "error"})
					return
				}
				expiry := time.Now().Add(3 * time.Minute)
				user := User{
					Email:              email,
					PasswordHash:       hash,
					Verified:           false,
					VerificationToken:  token,
					VerificationExpiry: expiry,
				}
				_, err = coll.InsertOne(ctx, user)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"message": "Failed to register user.", "status": "error"})
					return
				}
				verifyLink := "localhost:3000/verify?token=" + token
				subject := "Verify your email address"
				body := "Click the following link to verify your email: http://" + verifyLink
				err = sendMail(email, subject, body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"message": "Could not send verification email. Please try again later.", "status": "error"})
					return
				}
				json.NewEncoder(w).Encode(map[string]string{"message": "A verification email has been sent. Please check your inbox.", "status": "success"})
				return
			}
			if user.PasswordHash != hashPassword(password) {
				json.NewEncoder(w).Encode(map[string]string{"message": "Incorrect password.", "status": "error"})
				return
			}
			if user.Verified {
				subject := "You are already verified"
				body := "Hello,\n\nYour email is already verified. You do not need to verify again.\n\nThank you!"
				_ = sendMail(email, subject, body)
				json.NewEncoder(w).Encode(map[string]string{"message": "Your email is already verified!", "status": "success"})
				return
			}
			token, err := generateToken()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Failed to generate token.", "status": "error"})
				return
			}
			expiry := time.Now().Add(3 * time.Minute)
			_, err = coll.UpdateOne(ctx, map[string]interface{}{"_id": user.ID}, map[string]interface{}{"$set": map[string]interface{}{"verification_token": token, "verification_expiry": expiry, "verified": false}})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Failed to update user.", "status": "error"})
				return
			}
			verifyLink := "localhost:3000/verify?token=" + token
			subject := "Verify your email address"
			body := "Click the following link to verify your email: http://" + verifyLink
			err = sendMail(email, subject, body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Could not send verification email. Please try again later.", "status": "error"})
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "A verification email has been sent. Please check your inbox.", "status": "success"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed.", "status": "error"})
	}
}

func verifyHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			verifyTmpl.Execute(w, map[string]string{"Title": "Invalid Link", "Message": "Verification token is missing."})
			return
		}
		coll := client.Database("practiceproject").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var user User
		err := coll.FindOne(ctx, map[string]interface{}{"verification_token": token}).Decode(&user)
		if err != nil {
			verifyTmpl.Execute(w, map[string]string{"Title": "Invalid Link", "Message": "Invalid or expired verification link."})
			return
		}
		if user.Verified {
			verifyTmpl.Execute(w, map[string]string{"Title": "Already Verified", "Message": "Your email is already verified!"})
			return
		}
		if time.Now().After(user.VerificationExpiry) {
			verifyTmpl.Execute(w, map[string]string{"Title": "Verification Failed", "Message": "Email verification failed. The link has expired."})
			return
		}
		// Mark user as verified
		_, err = coll.UpdateOne(ctx, map[string]interface{}{"_id": user.ID}, map[string]interface{}{"$set": map[string]interface{}{"verified": true}})
		if err != nil {
			verifyTmpl.Execute(w, map[string]string{"Title": "Error", "Message": "Could not verify your email. Please try again later."})
			return
		}
		verifyTmpl.Execute(w, map[string]string{"Title": "Email Verified", "Message": "Your email is verified successfully!"})
	}
}

func checkVerifiedHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]bool{"verified": false, "expired": false})
			return
		}
		coll := client.Database("practiceproject").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var user User
		err := coll.FindOne(ctx, map[string]interface{}{"email": email}).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]bool{"verified": false, "expired": false})
			return
		}
		if user.Verified {
			json.NewEncoder(w).Encode(map[string]bool{"verified": true, "expired": false})
			return
		}
		if time.Now().After(user.VerificationExpiry) {
			json.NewEncoder(w).Encode(map[string]bool{"verified": false, "expired": true})
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"verified": false, "expired": false})
	}
}

func forgotPasswordHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			forgotTmpl.Execute(w, nil)
			return
		}
		if r.Method == http.MethodPost {
			var req struct {
				Email string `json:"email"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"message": "Email is required."})
				return
			}
			email := req.Email
			coll := client.Database("practiceproject").Collection("users")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var user User
			err := coll.FindOne(ctx, map[string]interface{}{"email": email}).Decode(&user)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"message": "If the email exists, a reset link will be sent."})
				return
			}
			token, err := generateToken()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Failed to generate reset token."})
				return
			}
			expiry := time.Now().Add(15 * time.Minute)
			_, err = coll.UpdateOne(ctx, map[string]interface{}{"_id": user.ID}, map[string]interface{}{"$set": map[string]interface{}{"password_reset_token": token, "password_reset_expiry": expiry}})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Failed to set reset token."})
				return
			}
			resetLink := "localhost:3000/reset-password?token=" + token
			subject := "Password Reset Request"
			body := "Click the following link to reset your password (valid for 15 minutes): http://" + resetLink
			_ = sendMail(email, subject, body)
			json.NewEncoder(w).Encode(map[string]string{"message": "If the email exists, a reset link will be sent."})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed."})
	}
}

func resetPasswordHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Token    string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" || req.Token == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Email, new password, and token are required."})
			return
		}
		token := req.Token
		coll := client.Database("practiceproject").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var user User
		err := coll.FindOne(ctx, map[string]interface{}{"password_reset_token": token}).Decode(&user)
		if err != nil || token == "" || time.Now().After(user.PasswordResetExpiry) {
			resetTmpl.Execute(w, map[string]string{"Message": "Invalid or expired reset link."})
			return
		}
		if r.Method == http.MethodGet {
			resetTmpl.Execute(w, map[string]interface{}{"Token": token})
			return
		}
		if r.Method == http.MethodPost {
			if req.Email != user.Email {
				json.NewEncoder(w).Encode(map[string]string{"message": "Email does not match reset request."})
				return
			}
			hash := hashPassword(req.Password)
			_, err = coll.UpdateOne(ctx, map[string]interface{}{"_id": user.ID}, map[string]interface{}{"$set": map[string]interface{}{"password_hash": hash}, "$unset": map[string]interface{}{"password_reset_token": "", "password_reset_expiry": ""}})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Failed to reset password."})
				return
			}
			subject := "Password Reset Successful"
			body := "Hello,\n\nYour password has been reset successfully. If you did not perform this action, please contact support immediately.\n\nThank you!"
			_ = sendMail(req.Email, subject, body)
			json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successful! You can now log in.", "status": "success"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed."})
	}
}

// Add CORS middleware
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env file")
	}

	// Get the MongoDB URI from environment variable
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable.")
	}

	// Set a timeout for connecting
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}
	fmt.Println("Connected to MongoDB!")

	http.HandleFunc("/", withCORS(verifyRequestHandler(client)))
	http.HandleFunc("/verify", withCORS(verifyHandler(client)))
	http.HandleFunc("/check-verified", withCORS(checkVerifiedHandler(client)))
	http.HandleFunc("/forgot-password", withCORS(forgotPasswordHandler(client)))
	http.HandleFunc("/reset-password", withCORS(resetPasswordHandler(client)))
	http.HandleFunc("/resend-verification", withCORS(resendVerificationHandler(client)))
	http.Handle("/animations/", http.StripPrefix("/animations/", http.FileServer(http.Dir("animations"))))
	log.Println("Server started at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

// Add resendVerificationHandler if not present
func resendVerificationHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed."})
			return
		}
		var req struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Email is required."})
			return
		}
		coll := client.Database("practiceproject").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var user User
		err := coll.FindOne(ctx, map[string]interface{}{"email": req.Email}).Decode(&user)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"message": "If the email exists, a verification email will be resent."})
			return
		}
		if user.Verified {
			json.NewEncoder(w).Encode(map[string]string{"message": "Your email is already verified!"})
			return
		}
		token, err := generateToken()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Failed to generate token."})
			return
		}
		expiry := time.Now().Add(3 * time.Minute)
		_, err = coll.UpdateOne(ctx, map[string]interface{}{"_id": user.ID}, map[string]interface{}{"$set": map[string]interface{}{"verification_token": token, "verification_expiry": expiry, "verified": false}})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Failed to update user."})
			return
		}
		verifyLink := "localhost:3000/verify?token=" + token
		subject := "Verify your email address"
		body := "Click the following link to verify your email: http://" + verifyLink
		_ = sendMail(req.Email, subject, body)
		json.NewEncoder(w).Encode(map[string]string{"message": "If the email exists, a verification email will be resent."})
	}
}
