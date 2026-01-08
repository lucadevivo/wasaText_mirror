package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

// DoLogin crea un nuovo utente se non esiste, o restituisce l'ID se esiste
func (db *appdbimpl) DoLogin(username string) (int, error) {
	// 1. Proviamo a inserire l'utente (se non esiste)
	res, err := db.c.Exec("INSERT OR IGNORE INTO users (username) VALUES (?)", username)
	if err != nil {
		return 0, err
	}

	// 2. Se abbiamo inserito una riga, prendiamo l'ID generato
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if affected > 0 {
		id, err := res.LastInsertId()
		return int(id), err
	}

	// 3. Se non abbiamo inserito nulla (l'utente esisteva già), cerchiamo il suo ID
	var id int
	err = db.c.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SetUserName cambia il nome utente
func (db *appdbimpl) SetUserName(id int, name string) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id = ?", name, id)
	return err
}

// SearchUsers cerca utenti che contengono la stringa nella query
func (db *appdbimpl) SearchUsers(query string) ([]User, error) {
	var users []User

	// Usiamo i wildcard % per cercare "contiene"
	searchTerm := "%" + query + "%"
	
	// Nota: includiamo photo_url nella SELECT
	rows, err := db.c.Query("SELECT id, username, photo_url FROM users WHERE username LIKE ?", searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		// Gestiamo il caso in cui photo_url sia NULL
		var photoURL sql.NullString
		
		err = rows.Scan(&u.ID, &u.Username, &photoURL)
		if err != nil {
			return nil, err
		}
		

		if photoURL.Valid {
			u.PhotoURL = photoURL.String
		}
		
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// SetUserPhoto imposta la foto profilo dell'utente
func (db *appdbimpl) SetUserPhoto(userID int, photoID int) error {
	photoURL := "/photos/" + strconv.Itoa(photoID)
	

	res, err := db.c.Exec("UPDATE users SET photo_url = ? WHERE id = ?", photoURL, userID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	

	if affected == 0 {
		return errors.New("nessun utente trovato con questo ID") 
	}

	return nil
}