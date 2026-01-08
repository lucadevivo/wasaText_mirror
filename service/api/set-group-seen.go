package api

import (
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Mark Received
func (rt *_router) setGroupReceived(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") { w.WriteHeader(http.StatusUnauthorized); return }
	myID, _ := strconv.Atoi(strings.TrimPrefix(authHeader, "Bearer "))

	groupID, _ := strconv.Atoi(ps.ByName("id"))

	err := rt.db.MarkGroupAsReceived(groupID, myID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Mark Read (Seen)
func (rt *_router) setGroupSeen(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") { w.WriteHeader(http.StatusUnauthorized); return }
	myID, _ := strconv.Atoi(strings.TrimPrefix(authHeader, "Bearer "))

	groupID, _ := strconv.Atoi(ps.ByName("id"))

	err := rt.db.MarkGroupAsRead(groupID, myID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}