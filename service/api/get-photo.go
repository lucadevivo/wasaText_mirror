package api

import (
	"net/http"
	"path/filepath"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Recupera ID
	photoID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 2. Chiedi al DB i dettagli (soprattutto il formato)
	photo, err := rt.db.GetPhoto(photoID)
	if err != nil {
		// Se non la trova nel DB
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 3. Costruisci il percorso del file
	fileName := strconv.Itoa(photo.ID) + "." + photo.Format
	filePath := filepath.Join(photoFolderPath, fileName) // photoFolderPath è definita nell'altro file

	// 4. Servi il file
	http.ServeFile(w, r, filePath)
}
