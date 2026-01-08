package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Definiamo dove salvare le foto.
const photoFolderPath = "/tmp/wasa-photos"

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Auth
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, _ := strconv.Atoi(strings.TrimPrefix(authHeader, "Bearer "))

	// 2. Parsare il form (max 10 MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. Recuperare il file dal form (il campo si deve chiamare "file")
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 4. Capire il formato (estensione)
	// fileHeader.Filename ci dà il nome originale (es. "vacanza.png")
	ext := filepath.Ext(fileHeader.Filename) // ".png"
	if ext == "" {
		ext = ".jpg" // Default
	}
	// Rimuoviamo il punto iniziale per pulizia ("png" invece di ".png")
	format := strings.TrimPrefix(ext, ".")

	// 5. Registrare nel Database per avere un ID
	photoID, err := rt.db.UploadPhoto(userID, format)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't insert photo in DB")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 6. Salvare il file su disco
	// Creiamo la cartella se non esiste
	err = os.MkdirAll(photoFolderPath, 0755)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't create folder")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Costruiamo il percorso: /tmp/wasa-photos/1.png
	fileName := strconv.Itoa(photoID) + "." + format
	filePath := filepath.Join(photoFolderPath, fileName)

	// Creiamo il file vuoto
	dst, err := os.Create(filePath)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't create file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copiamo i byte dal caricamento al disco
	_, err = io.Copy(dst, file)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't save file content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 7. Risposta
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{ID: photoID})
}
