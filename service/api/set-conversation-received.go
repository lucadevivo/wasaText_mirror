package api

import (
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setConversationReceived(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Auth
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	myID, _ := strconv.Atoi(strings.TrimPrefix(authHeader, "Bearer "))

	// 2. Parametro ID (l'altro utente che mi ha mandato i messaggi)
	otherUserID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. Database Call
	err = rt.db.MarkConversationAsReceived(myID, otherUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't mark conversation as received")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Risposta
	w.WriteHeader(http.StatusNoContent)
}