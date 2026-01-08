package database

import (
	"errors"
	"strconv"
)

// Crea un gruppo (accetta lista membri)
func (db *appdbimpl) CreateGroup(name string, members []int) (int, error) {
	// 1. Crea il gruppo
	res, err := db.c.Exec("INSERT INTO groups (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}
	groupID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 2. Aggiungi i membri iniziali
	for _, memberID := range members {
		_, err = db.c.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, memberID)
		if err != nil {
			return 0, err
		}
	}

	return int(groupID), nil
}

// Cambia nome al gruppo
func (db *appdbimpl) SetGroupName(groupID int, userID int, name string) error {
	_, err := db.c.Exec("UPDATE groups SET name = ? WHERE id = ?", name, groupID)
	return err
}
// Aggiungi membro
func (db *appdbimpl) AddGroupMember(groupID int, userID int) error {
	_, err := db.c.Exec("INSERT OR IGNORE INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, userID)
	return err
}

// Rimuovi membro
func (db *appdbimpl) RemoveGroupMember(groupID int, userID int) error {
	res, err := db.c.Exec("DELETE FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("member not found in group")
	}
	return nil
}

// INVIO MESSAGGIO 
func (db *appdbimpl) SendMessageToGroup(groupID int, senderID int, content string, photoID int) (int, error) {
	// 1. Inserisci il messaggio nella tabella principale
	var photoURL *string
	if photoID > 0 {
		url := "/photos/" + strconv.Itoa(photoID)
		photoURL = &url
	}

	res, err := db.c.Exec(`
		INSERT INTO messages (sender_id, group_id, content, recipient_id, photo_url) 
		VALUES (?, ?, ?, NULL, ?)`, 
		senderID, groupID, content, photoURL)
	if err != nil {
		return 0, err
	}

	msgID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 2. STEP SICURO: Prima raccogliamo tutti gli ID in una slice (lista)
	var memberIDs []int
	
	rows, err := db.c.Query("SELECT user_id FROM group_members WHERE group_id = ? AND user_id != ?", groupID, senderID)
	if err != nil {
		return int(msgID), nil 
	}
	
	for rows.Next() {
		var memberID int
		if err := rows.Scan(&memberID); err == nil {
			memberIDs = append(memberIDs, memberID)
		}
	}
	rows.Close() 

	// 3. Ora facciamo le Insert per le ricevute
	for _, memberID := range memberIDs {
		_, _ = db.c.Exec(`INSERT INTO group_receipts (message_id, user_id, received, read) VALUES (?, ?, FALSE, FALSE)`, msgID, memberID)
	}

	return int(msgID), nil
}

// LETTURA MESSAGGI
func (db *appdbimpl) GetGroupMessages(groupID int, requestorID int) ([]Message, error) {
	var messages []Message

	query := `
		SELECT 
			m.id, m.sender_id, 0, m.content, m.photo_url, m.timestamp,
			-- Calcolo RECEIVED
			CASE 
				WHEN m.sender_id = ? THEN 
					(SELECT COUNT(*) FROM group_receipts WHERE message_id = m.id AND received = 0) = 0
				ELSE 
					COALESCE(gr.received, 0)
			END as received,
			-- Calcolo READ
			CASE 
				WHEN m.sender_id = ? THEN 
					(SELECT COUNT(*) FROM group_receipts WHERE message_id = m.id AND read = 0) = 0
				ELSE 
					COALESCE(gr.read, 0)
			END as read
		FROM messages m
		LEFT JOIN group_receipts gr ON m.id = gr.message_id AND gr.user_id = ?
		WHERE m.group_id = ?
		ORDER BY m.timestamp ASC`

	rows, err := db.c.Query(query, requestorID, requestorID, requestorID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m Message
		err = rows.Scan(&m.ID, &m.SenderID, &m.RecipientID, &m.Content, &m.PhotoURL, &m.Timestamp, &m.Received, &m.Read)
		if err != nil {
			return nil, err
		}

		reactions, err := db.getMessageReactions(m.ID)
		if err != nil {
		} else {
			m.Reactions = reactions
		}
		// -------------------------------------------------------------

		messages = append(messages, m)
	}

	return messages, nil
}

// Segna messaggi del gruppo come RICEVUTI (per l'utente corrente)
func (db *appdbimpl) MarkGroupAsReceived(groupID int, userID int) error {
	_, err := db.c.Exec(`
		UPDATE group_receipts 
		SET received = TRUE 
		WHERE user_id = ? 
		  AND received = FALSE
		  AND message_id IN (SELECT id FROM messages WHERE group_id = ?)`,
		userID, groupID)
	return err
}

// Segna messaggi del gruppo come LETTI (per l'utente corrente)
func (db *appdbimpl) MarkGroupAsRead(groupID int, userID int) error {
	_, err := db.c.Exec(`
		UPDATE group_receipts 
		SET read = TRUE, received = TRUE 
		WHERE user_id = ? 
		  AND read = FALSE
		  AND message_id IN (SELECT id FROM messages WHERE group_id = ?)`,
		userID, groupID)
	return err
}

func (db *appdbimpl) SetGroupPhoto(groupID int, photoID int) error {
	photoURL := "/photos/" + strconv.Itoa(photoID)
	_, err := db.c.Exec("UPDATE groups SET photo_url = ? WHERE id = ?", photoURL, groupID)
	return err
}