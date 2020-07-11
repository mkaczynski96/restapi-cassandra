package controllers

import (
	"acaisoft-mkaczynski-api/helpers"
	"acaisoft-mkaczynski-api/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	apiPath             = "/api"
	apiContentTypeKey   = "Content-Type"
	apiContentTypeValue = "application/json"
)

func RunApi() {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc(apiPath+"/messages/{emailValue}", getMessages).Methods("GET")
	muxRouter.HandleFunc(apiPath+"/message", createMessage).Methods("POST")
	muxRouter.HandleFunc(apiPath+"/send", sendMessage).Methods("POST")

	log.Printf("API runs on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", muxRouter))
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	params := mux.Vars(r)
	emailParameter := params["emailValue"]

	args := make(map[string]interface{})
	args["email"] = emailParameter

	query := helpers.CreateSelectQuery(args)
	if messages := helpers.GetMessagesFromSelect(query); messages != nil {
		_ = json.NewEncoder(w).Encode(messages)
	}
	return
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(apiContentTypeKey, apiContentTypeValue)

	var newMessage models.EmailMessage
	_ = json.NewDecoder(r.Body).Decode(&newMessage)
	if err := helpers.AddRecordToDatabase(newMessage); err != nil {
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

	query := helpers.CreateSelectQuery(args)
	returnedMessages := helpers.GetMessagesFromSelect(query)
	if returnedMessages != nil {
		for _, message := range returnedMessages {
			// There are mail sending - disabled due to the lack of smtp server :)
			//	helpers.SendMail(message.Email, message.Title, message.Content)
			_ = helpers.RemoveRecordFromDatabase(message)
		}
		_ = json.NewEncoder(w).Encode("Done!")
	}
	return
}
