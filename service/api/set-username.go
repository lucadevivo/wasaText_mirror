package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
)

type SetUsernameRequest struct {
	Name string `json:"name"`
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. AUTENTICAZIONE (Bearer Token)
	authHeader := r.Header.Get("Authorization")
	// L'header deve essere formato da "Bearer {id}"
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Estraiamo l'ID (togliamo "Bearer " dall'inizio)
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := strconv.Atoi(tokenString) // Convertiamo in intero
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. Parsing della richiesta
	var body SetUsernameRequest
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. Validazione (stesse regole del login)
	if len(body.Name) < 3 || len(body.Name) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 4. Aggiornamento Database
	err = rt.db.SetUserName(userID, body.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't update username")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 5. Successo (204 No Content è standard per gli aggiornamenti senza risposta)
	w.WriteHeader(http.StatusNoContent)
}