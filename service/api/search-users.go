package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Autenticazione (Opzionale, ma consigliata: solo gli utenti loggati possono cercare)
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. Leggi il parametro "search" dall'URL (es. /users?search=marco)
	query := r.URL.Query().Get("search")

	// Qui restituiamo lista vuota se la query è vuota
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]database.User{})
		return
	}

	// 3. Chiama il Database
	users, err := rt.db.SearchUsers(query)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't search users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Rispondi
	w.Header().Set("Content-Type", "application/json")
	if users == nil {
		users = []database.User{}
	}
	_ = json.NewEncoder(w).Encode(users)
}
