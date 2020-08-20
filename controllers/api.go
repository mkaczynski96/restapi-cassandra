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

var cfg configs.Config

func RunApi(config configs.Config) {
	cfg = config
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc(apiPath+"/messages/{emailValue}", getMessages).Methods("GET")
	muxRouter.HandleFunc(apiPath+"/message", createMessage).Methods("POST")
	muxRouter.HandleFunc(apiPath+"/send", sendMessage).Methods("POST")

	log.Printf("Api runs on port %v", config.Api.Port)
	log.Fatal(http.ListenAndServe(config.Api.Port, muxRouter))
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	params := mux.Vars(r)
	emailParameter := params["emailValue"]

	args := make(map[string]interface{})
	args["email"] = emailParameter

	query := helpers.CreateSelectQuery(cfg, args)
	if messages := helpers.GetMessagesFromSelect(query); messages != nil {
		_ = json.NewEncoder(w).Encode(messages)
	}
	return
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	var newMessage models.EmailMessage
	_ = json.NewDecoder(r.Body).Decode(&newMessage)
	if err := helpers.AddRecordToDatabase(cfg, newMessage); err != nil {
		_ = json.NewEncoder(w).Encode(err)
	}
	_ = json.NewEncoder(w).Encode("200 OK")
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	var message models.EmailMessage
	_ = json.NewDecoder(r.Body).Decode(&message)

	args := make(map[string]interface{})
	args["magic_number"] = message.MagicNumber

	query := helpers.CreateSelectQuery(cfg, args)
	returnedMessages := helpers.GetMessagesFromSelect(query)
	if returnedMessages != nil {
		for _, message := range returnedMessages {
			_ = helpers.RemoveRecordFromDatabase(cfg, message)
			helpers.SendMail(cfg,message.Email, message.Title, message.Content)
		}
		_ = json.NewEncoder(w).Encode("200 OK")
	}
	return
}
