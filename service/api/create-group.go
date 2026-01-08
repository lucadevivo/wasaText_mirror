package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Struttura del JSON che ci invia l'utente
type CreateGroupRequest struct {
	Name    string `json:"name"`
	Members []int  `json:"members"` // Lista degli ID degli utenti da aggiungere
}

// Struttura della risposta
type GroupResponse struct {
	ID int `json:"id"`
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Autenticazione (Chi sta creando il gruppo?)
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

	// 2. Parsing della richiesta
	var body CreateGroupRequest
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validazione: Il nome non può essere vuoto
	if len(body.Name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. Logica di Business: Aggiungi il creatore alla lista dei membri!
	// Un gruppo deve contenere almeno chi lo ha creato.
	finalMembers := append(body.Members, myID)

	// 4. Chiama il Database
	groupID, err := rt.db.CreateGroup(body.Name, finalMembers)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't create group")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 5. Risposta (201 Created)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(GroupResponse{ID: groupID})
}