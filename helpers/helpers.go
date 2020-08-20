package helpers

import (
	"bytes"
	"fmt"
	"net/smtp"
	"reflect"
	"restapi-cassandra/configs"
	"restapi-cassandra/models"
	"strconv"
)

func CreateSelectQuery(config configs.Config, args map[string]interface{}) string {
	var buffer bytes.Buffer

	buffer.WriteString("SELECT * FROM " + config.Database.Keyspace + "." + config.Database.TableName)
	// if args exist add WHERE statement
	if args != nil {
		buffer.WriteString(" WHERE ")
	}

	index := 0
	for name, value := range args {
		index++
		// add AND statement before every argument, except first argument
		if index > 1 {
			buffer.WriteString(" AND ")
		}
		// if arg is int, add it without single quotes
		if reflect.TypeOf(value).Kind() == reflect.Int {
			buffer.WriteString(fmt.Sprintf("%s=%v", name, value))
		} else {
			buffer.WriteString(fmt.Sprintf("%s='%v'", name, value))
		}
	}

	buffer.WriteString(" ALLOW FILTERING;")
	return buffer.String()
}

func ExtractMessagesFromSelectResponse(connection *configs.DBConnection, query string) []models.EmailMessage {
	var messages []models.EmailMessage
	m := map[string]interface{}{}

	iter := connection.ExecuteSelectQuery(query)
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

func AddRecordToDatabase(connection *configs.DBConnection, config configs.Config, message models.EmailMessage) error {
	query := `INSERT INTO ` + config.Database.Keyspace + `.` + config.Database.TableName + `(email, title, content, magic_number) VALUES ('` + message.Email + `', '` + message.Title + `', '` + message.Content + `', ` + strconv.Itoa(message.MagicNumber) + `) USING TTL ` + strconv.Itoa(config.Mail.MessageExpirationSeconds) + `;`
	if err := connection.ExecuteQuery(query); err != nil {
		return err
	}
	return nil
}

func RemoveRecordFromDatabase(connection *configs.DBConnection, config configs.Config, message models.EmailMessage) error {
	query := `DELETE FROM ` + config.Database.Keyspace + `.` + config.Database.TableName + ` WHERE email='` + message.Email + `' AND magic_number=` + strconv.Itoa(message.MagicNumber) + `;`
	if err := connection.ExecuteQuery(query); err != nil {
		return err
	}
	return nil
}

// Send mail function with args: config, email of receiver, title of email, content of email
func SendMail(config configs.Config, email, title, content string) error {
	auth := smtp.PlainAuth("", config.Mail.Username, config.Mail.Password, config.Mail.Host)

	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + title + "\r\n" +
		"\r\n" +
		"" + content + "\r\n")
	hostPort := config.Mail.Host + ":" + strconv.Itoa(config.Mail.Port)
	err := smtp.SendMail(hostPort, auth, config.Mail.Username, to, msg)
	if err != nil {
		return err
	}
	return nil
}
