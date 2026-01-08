package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	// Importiamo il pacchetto database per usare la struct Message
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Identifica CHI sta facendo la richiesta
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	myID, err := strconv.Atoi(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. Identifica CON CHI vuole vedere la chat
	otherIDString := ps.ByName("id")
	otherID, err := strconv.Atoi(otherIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. Chiama il database
	messages, err := rt.db.GetConversation(myID, otherID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't retrieve conversation")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Rispondi col JSON
	w.Header().Set("Content-Type", "application/json")

	if messages == nil {
		messages = []database.Message{}
	}

	_ = json.NewEncoder(w).Encode(messages)
}
