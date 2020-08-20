package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"restapi-cassandra/configs"
	"restapi-cassandra/helpers"
	"restapi-cassandra/models"
)

const (
	apiPath             = "/api"
	apiContentTypeKey   = "Content-Type"
	apiContentTypeValue = "application/json"
)

type Database struct {
	connection *configs.DBConnection
}

var cfg configs.Config

func RunApi(config configs.Config, connection *configs.DBConnection) {
	cfg = config
	conn := &Database{
		connection: connection,
	}
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc(apiPath+"/messages/{emailValue}", conn.getMessages).Methods("GET")
	muxRouter.HandleFunc(apiPath+"/message", conn.addMessage).Methods("POST")
	muxRouter.HandleFunc(apiPath+"/send", conn.sendMessage).Methods("POST")

	log.Printf("Api runs on port %v", config.Api.Port)
	log.Fatal(http.ListenAndServe(config.Api.Port, muxRouter))
}

// Get message by email value
func (connection Database) getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	params := mux.Vars(r)
	emailParameter := params["emailValue"]

	args := make(map[string]interface{})
	args["email"] = emailParameter

	query := helpers.CreateSelectQuery(cfg, args)
	if messages := helpers.ExtractMessagesFromSelectResponse(connection.connection, query); messages != nil {
		_ = json.NewEncoder(w).Encode(messages)
	}
	return
}

// Add new message
func (connection Database) addMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	var newMessage models.EmailMessage
	_ = json.NewDecoder(r.Body).Decode(&newMessage)
	if err := helpers.AddRecordToDatabase(connection.connection, cfg, newMessage); err != nil {
		_ = json.NewEncoder(w).Encode(err)
	}
	_ = json.NewEncoder(w).Encode("200 OK")
}

// Send email and remove message from db
func (connection Database) sendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	var message models.EmailMessage
	_ = json.NewDecoder(r.Body).Decode(&message)

	args := make(map[string]interface{})
	args["magic_number"] = message.MagicNumber

	query := helpers.CreateSelectQuery(cfg, args)
	returnedMessages := helpers.ExtractMessagesFromSelectResponse(connection.connection, query)
	if returnedMessages != nil {
		for _, message := range returnedMessages {
			errSend := helpers.SendMail(cfg, message.Email, message.Title, message.Content)
			if errSend != nil {
				log.Fatal(errSend)
			}
			_ = helpers.RemoveRecordFromDatabase(connection.connection, cfg, message)
		}
		_ = json.NewEncoder(w).Encode("200 OK")
	}
	return
}
