package helpers

import (
	"acaisoft-mkaczynski-api/configs"
	"acaisoft-mkaczynski-api/models"
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"reflect"
)

const (
	messageExpirationSeconds = 300
	usernameMail             = "test@example.com"
	passwordMail             = "test123"
	hostMail                 = "smtp.test.com"
	portMail                 = "25"
)

func CreateSelectQuery(args map[string]interface{}) string {
	var buffer bytes.Buffer

	buffer.WriteString("SELECT * FROM " + configs.Keyspace + "." + configs.TableName)
	if args != nil {
		buffer.WriteString(" WHERE ")
	}

	index := 0
	for name, value := range args {
		index++
		if index > 1 {
			buffer.WriteString(" AND ")
		}
		if reflect.TypeOf(value).Kind() == reflect.Int {
			buffer.WriteString(fmt.Sprintf("%s=%v", name, value))
		} else {
			buffer.WriteString(fmt.Sprintf("%s='%v'", name, value))
		}
	}

	buffer.WriteString(" ALLOW FILTERING;")
	return buffer.String()
}

func GetMessagesFromSelect(query string) []models.EmailMessage {
	var messages []models.EmailMessage
	m := map[string]interface{}{}

	iter := configs.ExecuteSelectQuery(query)
	for iter.MapScan(m) {
		messages = append(messages, models.EmailMessage{
			Email:       m["email"].(string),
			Title:       m["title"].(string),
			Content:     m["content"].(string),
			MagicNumber: m["magic_number"].(int),
		})
		m = map[string]interface{}{}
	}
	return messages
}

func AddRecordToDatabase(message models.EmailMessage) error {
	query := fmt.Sprintf("INSERT INTO %s.%s(email, title, content, magic_number) VALUES ('%s', '%s', '%s', %v) USING TTL %v;", configs.Keyspace, configs.TableName, message.Email, message.Title, message.Content, message.MagicNumber, messageExpirationSeconds)
	if err := configs.ExecuteQuery(query); err != nil {
		return err
	}
	return nil
}

func RemoveRecordFromDatabase(message models.EmailMessage) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE email='%s' AND magic_number=%v;", configs.Keyspace, configs.TableName, message.Email, message.MagicNumber)
	if err := configs.ExecuteQuery(query); err != nil {
		return err
	}
	return nil
}

func SendMail(email, title, content string) {
	auth := smtp.PlainAuth("", usernameMail, passwordMail, hostMail)

	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + title + "?\r\n" +
		"\r\n" +
		"" + content + "\r\n")
	err := smtp.SendMail(hostMail+portMail, auth, usernameMail, to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
