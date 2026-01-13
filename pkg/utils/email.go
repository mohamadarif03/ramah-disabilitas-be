package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendVerificationEmail(toEmail, token string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	baseURL := os.Getenv("BASE_URL") // URL of the frontend or backend depending on where the verify link points.
	// Currently pointing to backend for simplicity as requested, but standard is frontend.
	// The user request: "berikan di response ketika abis regist user disuruh mengecek email verifikasi"
	// and "buatkan fitur verifikasi email".
	// Usually one clicks a link like http://api.com/verify?token=xyz

	if baseURL == "" {
		baseURL = "http://localhost:8080" // Default dev
	}

	verificationLink := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", baseURL, token)

	if smtpHost == "" || smtpUser == "" {
		// Mock email sending for dev environment without SMTP credentials
		log.Printf("==================================================\n")
		log.Printf("[MOCK EMAIL] To: %s\n", toEmail)
		log.Printf("[MOCK EMAIL] Verification Link: %s\n", verificationLink)
		log.Printf("==================================================\n")
		return nil
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Verifikasi Email Anda\r\n"+
		"\r\n"+
		"Halo,\r\n\r\n"+
		"Terima kasih telah mendaftar. Silakan klik tautan di bawah ini untuk memverifikasi email Anda:\r\n"+
		"%s\r\n\r\n"+
		"Jika Anda tidak merasa mendaftar, abaikan email ini.\r\n", toEmail, verificationLink))

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, smtpUser, []string{toEmail}, msg)
	if err != nil {
		return err
	}

	return nil
}
