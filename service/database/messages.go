package database

import (
	"strconv"
)

// SendMessage invia un messaggio (supporta foto opzionale)
func (db *appdbimpl) SendMessage(senderId int, recipientId int, content string, photoID int) (int, error) {
	var photoURL *string
	
	// Se c'è un ID foto, costruiamo l'URL
	if photoID > 0 {
		url := "/photos/" + strconv.Itoa(photoID)
		photoURL = &url
	}

	// Unica esecuzione della query corretta
	res, err := db.c.Exec(`INSERT INTO messages (sender_id, recipient_id, content, photo_url) VALUES (?, ?, ?, ?)`,
		senderId, recipientId, content, photoURL)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func (db *appdbimpl) GetConversation(myId int, otherId int) ([]Message, error) {
	var messages []Message

	// AGGIUNTO photo_url nella SELECT
	query := `
		SELECT id, sender_id, recipient_id, content, photo_url, timestamp, received, read 
		FROM messages 
		WHERE (sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)
		ORDER BY timestamp ASC`

	rows, err := db.c.Query(query, myId, otherId, otherId, myId)
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

		messages = append(messages, m)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (db *appdbimpl) DeleteMessage(messageId int, senderId int) error {
	// Eseguiamo la DELETE con una condizione WHERE rigorosa:
	// "Cancella il messaggio X SOLO SE il mittente è Y"
	res, err := db.c.Exec("DELETE FROM messages WHERE id = ? AND sender_id = ?", messageId, senderId)
	if err != nil {
		return err
	}

	
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return nil 
	}
	return nil
}

func (db *appdbimpl) MarkConversationAsReceived(myID int, otherUserID int) error {
	_, err := db.c.Exec(`
		UPDATE messages 
		SET received = TRUE 
		WHERE recipient_id = ? AND sender_id = ? AND received = FALSE`,
		myID, otherUserID)
	return err
}


func (db *appdbimpl) MarkConversationAsRead(myID int, otherUserID int) error {
	_, err := db.c.Exec(`
		UPDATE messages 
		SET read = TRUE, received = TRUE 
		WHERE recipient_id = ? AND sender_id = ? AND read = FALSE`,
		myID, otherUserID)
	return err
}

func (db *appdbimpl) ReactToMessage(messageID int, userID int, emoji string) error {

	_, err := db.c.Exec(`
		INSERT INTO reactions (message_id, user_id, emoji) 
		VALUES (?, ?, ?) 
		ON CONFLICT(message_id, user_id) DO UPDATE SET emoji = excluded.emoji`,
		messageID, userID, emoji)
	return err
}

func (db *appdbimpl) UnreactToMessage(messageID int, userID int) error {
	_, err := db.c.Exec("DELETE FROM reactions WHERE message_id = ? AND user_id = ?", messageID, userID)
	return err
}

// Funzione helper interna per recuperare le reazioni di un messaggio
func (db *appdbimpl) getMessageReactions(messageID int) ([]Reaction, error) {
	rows, err := db.c.Query("SELECT user_id, emoji FROM reactions WHERE message_id = ?", messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []Reaction
	for rows.Next() {
		var r Reaction
		if err := rows.Scan(&r.UserID, &r.Emoji); err == nil {
			reactions = append(reactions, r)
		}
	}
	// Se non ci sono reazioni, restituiamo una slice vuota invece di nil (più bello nel JSON)
	if reactions == nil {
		reactions = []Reaction{}
	}
	return reactions, nil
}