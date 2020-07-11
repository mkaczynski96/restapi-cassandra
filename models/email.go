package models

type EmailMessage struct {
	Email       string `json:"email"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MagicNumber int    `json:"magic_number"`
}
