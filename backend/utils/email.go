package utils

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

type EmailConfig struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

func LoadEmailConfig() (*EmailConfig, error) {
	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		portStr = "587"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	cfg := &EmailConfig{
		Host: os.Getenv("SMTP_HOST"),
		Port: port,
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
		From: os.Getenv("SMTP_FROM"),
	}
	if cfg.Host == "" || cfg.User == "" || cfg.Pass == "" || cfg.From == "" {
		return nil, errors.New("SMTP_HOST, SMTP_USER, SMTP_PASS, SMTP_FROM are required")
	}
	return cfg, nil
}

func SendEmail(cfg *EmailConfig, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s", cfg.From, to, subject, body))
	return smtp.SendMail(addr, auth, cfg.From, []string{to}, msg)
}
