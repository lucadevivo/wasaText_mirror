package database

func (db *appdbimpl) GetMyConversations(myId int) ([]Conversation, error) {
	var conversations []Conversation

	// Questa query seleziona ID e Nome degli utenti con cui ho interagito.
	query := `
		SELECT u.id, u.username
		FROM users u
		WHERE u.id IN (
			SELECT recipient_id FROM messages WHERE sender_id = ?
			UNION
			SELECT sender_id FROM messages WHERE recipient_id = ?
		)`

	rows, err := db.c.Query(query, myId, myId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Conversation
		err = rows.Scan(&c.OtherUserID, &c.OtherUserName)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}