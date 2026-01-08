package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// Send
func (rt *_router) sendGroupMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Auth
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	myID, _ := strconv.Atoi(strings.TrimPrefix(authHeader, "Bearer "))

	// Params
	groupID, _ := strconv.Atoi(ps.ByName("id"))

	// Body
	var body struct {
		Content string `json:"content"`
		PhotoID int    `json:"photo_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// DB Call
	msgID, err := rt.db.SendMessageToGroup(groupID, myID, body.Content, body.PhotoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't send group message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{ID: msgID})
}

func (rt *_router) getGroupMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Auth: Estraiamo il myID dall'header
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	myID, err := strconv.Atoi(strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. Parametri URL
	groupID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. Chiamata al DB
	messages, err := rt.db.GetGroupMessages(groupID, myID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't get group messages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Risposta
	w.Header().Set("Content-Type", "application/json")
	if messages == nil {
		messages = []database.Message{}
	}
	_ = json.NewEncoder(w).Encode(messages)
}
