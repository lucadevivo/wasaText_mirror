package database

func (db *appdbimpl) UploadPhoto(userID int, format string) (int, error) {
	res, err := db.c.Exec("INSERT INTO photos (user_id, format) VALUES (?, ?)", userID, format)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (db *appdbimpl) GetPhoto(photoID int) (Photo, error) {
	var p Photo
	err := db.c.QueryRow("SELECT id, user_id, format FROM photos WHERE id = ?", photoID).Scan(&p.ID, &p.UserID, &p.Format)
	if err != nil {
		return p, err
	}
	return p, nil
}