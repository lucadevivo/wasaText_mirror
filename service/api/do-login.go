package api

import (
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"

)

// Questa struct definisce cosa ci aspettiamo di ricevere dal frontend (JSON)
type UserRequest struct {
	Name string `json:"name"` // Esempio JSON: {"name": "Maria"}
}

// Questa struct definisce cosa rispondiamo noi (JSON)
type UserResponse struct {
	Identifier int `json:"identifier"` // Esempio JSON: {"identifier": 123}
}

// doLogin è la funzione che viene chiamata quando qualcuno fa POST /session
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
    // 1. Leggiamo il JSON inviato dall'utente
    var userReq UserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		// Se il JSON è sbagliato (es. sintassi errata), diamo errore 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}

    // 2. Validazione minima (come da specifiche OpenAPI)
    // Esempio: il nome non può essere vuoto o troppo corto
    if len(userReq.Name) < 3 {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

	// 3. Chiamiamo il database (la funzione che abbiamo scritto prima!)
	id, err := rt.db.DoLogin(userReq.Name)
	if err != nil {
		// Se il DB fallisce, logghiamo l'errore e diamo errore 500
		ctx.Logger.WithError(err).Error("Login failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Prepariamo la risposta
	response := UserResponse{
		Identifier: id,
	}

	// 5. Inviamo la risposta al client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	_ = json.NewEncoder(w).Encode(response)
}