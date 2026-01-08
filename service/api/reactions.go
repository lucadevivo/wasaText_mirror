package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Mette (o aggiorna) una reazione
func (rt *_router) reactToMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Auth
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") { w.WriteHeader(http.StatusUnauthorized); return }
	myID, _ := strconv.Atoi(strings.TrimPrefix(authHeader, "Bearer "))

	// Parametri
	msgID, _ := strconv.Atoi(ps.ByName("id"))

	// Body {"emoji": "😎"}
	var body struct {
		Emoji string `json:"emoji"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Emoji) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := rt.db.ReactToMessage(msgID, myID, body.Emoji)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Toglie la reazione
func (rt *_router) unreactToMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") { w.WriteHeader(http.StatusUnauthorized); return }
	myID, _ := strconv.Atoi(strings.TrimPrefix(authHeader, "Bearer "))

	msgID, _ := strconv.Atoi(ps.ByName("id"))

	err := rt.db.UnreactToMessage(msgID, myID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}