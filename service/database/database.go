package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Reaction struct {
	UserID int    `json:"user_id"`
	Emoji  string `json:"emoji"`
}

type Photo struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Format string `json:"format"` // es: "png", "jpeg"
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	PhotoURL string `json:"photo_url"` // Aggiornato
}

// Struct per i Gruppi
type Group struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	PhotoURL  string `json:"photo_url"`
	CreatedAt string `json:"created_at"`
}

// Struct per la lista conversazioni
type Conversation struct {
	OtherUserID   int    `json:"user_id"`
	OtherUserName string `json:"username"`
}

type Message struct {
	ID          int        `json:"id"`
	SenderID    int        `json:"sender_id"`
	RecipientID int        `json:"recipient_id"`
	Content     string     `json:"content"`
	PhotoURL    *string    `json:"photo_url"`
	Timestamp   string     `json:"timestamp"`
	Received    bool       `json:"received"`
	Read        bool       `json:"read"`
	Reactions   []Reaction `json:"reactions"` 
}

// AppDatabase è l'interfaccia principale
type AppDatabase interface {
	DoLogin(username string) (int, error)
	SetUserName(id int, name string) error

	// Messaggi 1-a-1
	SendMessage(senderId int, recipientId int, content string, photoID int) (int, error)
	GetConversation(myId int, otherId int) ([]Message, error)
	DeleteMessage(messageId int, senderId int) error
	GetMyConversations(myId int) ([]Conversation, error)

	// --- GRUPPI ---
	// Crea un gruppo con un nome e una lista iniziale di membri
	CreateGroup(name string, memberIDs []int) (int, error)

	SetGroupName(groupID int, userID int, newName string) error
	AddGroupMember(groupID int, userID int) error

	// Messaggi di Gruppo
	SendMessageToGroup(groupID int, senderID int, content string, photoID int) (int, error)
	GetGroupMessages(groupID int, requestorID int) ([]Message, error)

	RemoveGroupMember(groupID int, userID int) error

	SearchUsers(query string) ([]User, error)

	// Foto
	UploadPhoto(userID int, format string) (int, error)
	GetPhoto(photoID int) (Photo, error)

	// Spunte (Read Receipts)
	MarkConversationAsRead(myID int, otherUserID int) error
	MarkConversationAsReceived(myID int, otherUserID int) error

	MarkGroupAsRead(groupID int, userID int) error
	MarkGroupAsReceived(groupID int, userID int) error
	
	// Foto Profilo (Utenti e Gruppi)
	SetUserPhoto(userID int, photoID int) error
	SetGroupPhoto(groupID int, photoID int) error

	// Reazioni (Emoji)
	ReactToMessage(messageID int, userID int, emoji string) error
	UnreactToMessage(messageID int, userID int) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New restituisce una nuova istanza di AppDatabase
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// 1. Tabella Users (AGGIORNATA con photo_url)
	queryUsers := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        photo_url TEXT
    );`
	if _, err := db.Exec(queryUsers); err != nil {
		return nil, fmt.Errorf("error creating database structure (users): %w", err)
	}

	// 2. Tabella Messages
	queryMessages := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender_id INTEGER NOT NULL,
        recipient_id INTEGER,  
        group_id INTEGER,      
        content TEXT,
        photo_url TEXT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        received BOOLEAN DEFAULT FALSE,
        read BOOLEAN DEFAULT FALSE,
        FOREIGN KEY (sender_id) REFERENCES users(id),
        FOREIGN KEY (recipient_id) REFERENCES users(id),
        FOREIGN KEY (group_id) REFERENCES groups(id)
    );`
	if _, err := db.Exec(queryMessages); err != nil {
		return nil, fmt.Errorf("error creating database structure (messages): %w", err)
	}

	// 3. Tabella Gruppi (AGGIORNATA con photo_url se non c'era)
	queryGroups := `
    CREATE TABLE IF NOT EXISTS groups (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        photo_url TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
	if _, err := db.Exec(queryGroups); err != nil {
		return nil, fmt.Errorf("error creating database structure (groups): %w", err)
	}

	// 4. Tabella Membri del Gruppo
	queryGroupMembers := `
    CREATE TABLE IF NOT EXISTS group_members (
        group_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        FOREIGN KEY (group_id) REFERENCES groups(id),
        FOREIGN KEY (user_id) REFERENCES users(id),
        PRIMARY KEY (group_id, user_id)
    );`
	if _, err := db.Exec(queryGroupMembers); err != nil {
		return nil, fmt.Errorf("error creating database structure (group_members): %w", err)
	}

	// 5. Tabella Foto
	queryPhotos := `
    CREATE TABLE IF NOT EXISTS photos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        format TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
	if _, err := db.Exec(queryPhotos); err != nil {
		return nil, fmt.Errorf("error creating database structure (photos): %w", err)
	}

	// 6. Tabella Group Receipts (Stati messaggi di gruppo)
	queryGroupReceipts := `
    CREATE TABLE IF NOT EXISTS group_receipts (
        message_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        received BOOLEAN DEFAULT FALSE,
        read BOOLEAN DEFAULT FALSE,
        PRIMARY KEY (message_id, user_id),
        FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`
	if _, err := db.Exec(queryGroupReceipts); err != nil {
		return nil, fmt.Errorf("error creating database structure (group_receipts): %w", err)
	}

	// 7. Tabella Reactions (NUOVA)
	queryReactions := `
	CREATE TABLE IF NOT EXISTS reactions (
		message_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		emoji TEXT NOT NULL,
		PRIMARY KEY (message_id, user_id),
		FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	if _, err := db.Exec(queryReactions); err != nil {
		return nil, fmt.Errorf("error creating database structure (reactions): %w", err)
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}