package api

import (
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) removeGroupMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Auth
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. Parametri URL
	groupID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	targetUserID, err := strconv.Atoi(ps.ByName("userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. Esecuzione
	err = rt.db.RemoveGroupMember(groupID, targetUserID)
	if err != nil {

		ctx.Logger.WithError(err).Error("Can't remove member")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Successo
	w.WriteHeader(http.StatusNoContent)
}
