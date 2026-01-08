package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Qui registriamo il Login
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Esempio: PUT /settings/username
	rt.router.PUT("/settings/username", rt.wrap(rt.setMyUserName))

	// Invio messaggio a un utente specifico
	// Esempio URL: POST /conversations/2/messages
	rt.router.POST("/conversations/:id/messages", rt.wrap(rt.sendMessage))

	// Leggi i messaggi di una conversazione
	rt.router.GET("/conversations/:id/messages", rt.wrap(rt.getConversation))

	// Cancella messaggio (Unsend)
	rt.router.DELETE("/conversations/:id/messages/:messageId", rt.wrap(rt.unsendMessage))

	// Lista delle conversazioni (Inbox)
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))

	// Crea un nuovo gruppo
	rt.router.POST("/groups", rt.wrap(rt.createGroup))

	// Cambia nome gruppo
	rt.router.PUT("/groups/:id/name", rt.wrap(rt.setGroupName))

	// Aggiungi membro al gruppo
	rt.router.PUT("/groups/:id/members/:userId", rt.wrap(rt.addGroupMember))

	// Messaggi di Gruppo
	rt.router.POST("/groups/:id/messages", rt.wrap(rt.sendGroupMessage))
	rt.router.GET("/groups/:id/messages", rt.wrap(rt.getGroupMessages))

	// Rimuovi membro (o esci dal gruppo)
	rt.router.DELETE("/groups/:id/members/:userId", rt.wrap(rt.removeGroupMember))

	// Cerca utenti
	rt.router.GET("/users", rt.wrap(rt.searchUsers))

	// Foto
	rt.router.POST("/photos", rt.wrap(rt.uploadPhoto))
	rt.router.GET("/photos/:id", rt.wrap(rt.getPhoto))

	// Segna conversazione come letta
	rt.router.PUT("/conversations/:id/seen", rt.wrap(rt.setConversationSeen))

	// Segna conversazione come ricevuta (consegnata al client)
	rt.router.PUT("/conversations/:id/received", rt.wrap(rt.setConversationReceived))

	rt.router.PUT("/groups/:id/received", rt.wrap(rt.setGroupReceived))
	rt.router.PUT("/groups/:id/seen", rt.wrap(rt.setGroupSeen))

	// Foto Profilo
	rt.router.PUT("/users/me/photo", rt.wrap(rt.setUserPhoto))
	rt.router.PUT("/groups/:id/photo", rt.wrap(rt.setGroupPhoto))

	// Reazioni (Emoji) ai messaggi
	rt.router.PUT("/messages/:id/reaction", rt.wrap(rt.reactToMessage))
	rt.router.DELETE("/messages/:id/reaction", rt.wrap(rt.unreactToMessage))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
