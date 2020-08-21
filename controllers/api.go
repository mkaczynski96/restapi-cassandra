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

type Api struct {
	connection *configs.DBConnection
	config     configs.Config
}

func RunApi(config configs.Config, connection *configs.DBConnection) {
	conn := &Api{
		connection: connection,
		config:     config,
	}
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc(apiPath+"/messages/{emailValue}", conn.getMessages).Methods("GET")
	muxRouter.HandleFunc(apiPath+"/message", conn.addMessage).Methods("POST")
	muxRouter.HandleFunc(apiPath+"/send", conn.sendMessage).Methods("POST")

	log.Printf("Api runs on port %v", config.Api.Port)
	log.Fatal(http.ListenAndServe(config.Api.Port, muxRouter))
}

// Get message by email value
func (api Api) getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	params := mux.Vars(r)
	emailParameter := params["emailValue"]

	args := make(map[string]interface{})
	args["email"] = emailParameter

	query := helpers.CreateSelectQuery(api.config, args)
	if messages := helpers.ExtractMessagesFromSelectResponse(api.connection, query); messages != nil {
		errEncode := json.NewEncoder(w).Encode(messages)
		if errEncode != nil {
			log.Fatal(errEncode)
			return
		}
	}
	return
}

// Add new message
func (api Api) addMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	var newMessage models.EmailMessage
	_ = json.NewDecoder(r.Body).Decode(&newMessage)
	if err := helpers.AddRecordToDatabase(api.connection, api.config, newMessage); err != nil {
		_ = json.NewEncoder(w).Encode(err)
	}
	errEncode := json.NewEncoder(w).Encode("200 OK")
	if errEncode != nil {
		log.Fatal(errEncode)
		return
	}
}

// Send email and remove message from db
func (api Api) sendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	var message models.EmailMessage
	_ = json.NewDecoder(r.Body).Decode(&message)

	args := make(map[string]interface{})
	args["magic_number"] = message.MagicNumber

	query := helpers.CreateSelectQuery(api.config, args)
	returnedMessages := helpers.ExtractMessagesFromSelectResponse(api.connection, query)
	if returnedMessages != nil {
		for _, message := range returnedMessages {
			errSend := helpers.SendMail(api.config, message.Email, message.Title, message.Content)
			if errSend != nil {
				log.Fatal(errSend)
			}
			_ = helpers.RemoveRecordFromDatabase(api.connection, api.config, message)
		}
		errEncode := json.NewEncoder(w).Encode("200 OK")
		if errEncode != nil {
			log.Fatal(errEncode)
			return
		}
	}
	return
}
