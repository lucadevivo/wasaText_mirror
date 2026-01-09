<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()

// Stato
const conversations = ref([])
const foundUsers = ref([])
const searchQuery = ref('')
const loading = ref(false)

const activeChat = ref(null) 
const activeChatMessages = ref([]) 
const newMessage = ref('')        

const showProfileModal = ref(false)
const newUsername = ref('')

// Timer per l'aggiornamento automatico
let pollingInterval = null

const myUsername = localStorage.getItem('username')
const token = localStorage.getItem('token')
const myID = parseInt(localStorage.getItem('token')) 

const api = axios.create({
	baseURL: 'http://localhost:3000',
	headers: { Authorization: `Bearer ${token}` }
})

// Funzione per aprire il profilo
function openProfile() {
    newUsername.value = myUsername 
    showProfileModal.value = true
}

// Funzione per salvare il nome
async function updateUsername() {
    if (!newUsername.value || newUsername.value.length < 3) {
        alert("Il nome deve essere di almeno 3 caratteri")
        return
    }

    try {
        await api.put('/settings/username', { username: newUsername.value })
        
        localStorage.setItem('username', newUsername.value)
        window.location.reload()
        
    } catch (e) {
        console.error("Errore cambio nome:", e)
        alert("Errore durante il cambio nome (il backend è pronto?)")
    }
}

async function refreshConversations() {
    try {
        // 1. Scarica la lista delle chat
        const response = await api.get(`/conversations?t=${Date.now()}`)
        conversations.value = response.data || []

        // 2. LOOP FORZATO PER LE SPUNTE "RICEVUTO"
        if (conversations.value.length > 0) {
            conversations.value.forEach(chat => {

                api.get(`/conversations/${chat.user_id}/messages?anti_cache=${Date.now()}`)
                   .catch(() => {}) 
            })
        }
    } catch (e) {
        console.error("Errore caricamento chat:", e)
    }
}

async function searchUsers() {
	if (searchQuery.value.length < 2) {
		foundUsers.value = []
		return
	}
	loading.value = true
	try {
		const response = await api.get('/users', { params: { search: searchQuery.value } })
		foundUsers.value = response.data || []
	} catch (e) {
		console.error("Errore ricerca:", e)
	} finally {
		loading.value = false
	}
}

// --- GESTIONE CHAT ---

// 1. Scarica messaggi (e NON resetta lo scroll se stiamo solo aggiornando)
async function getChatMessages(userId, isBackgroundUpdate = false) {
	try {
		const response = await api.get(`/conversations/${userId}/messages`)
		activeChatMessages.value = response.data || []
		
		// Scrolla giù solo se NON è un aggiornamento automatico (altrimenti da fastidio mentre leggi)
		if (!isBackgroundUpdate) {
			scrollToBottom()
		}
	} catch (e) {
		console.error("Errore scaricamento messaggi:", e)
	}
}

// 2. Diciamo al server: "Ho letto tutto!"
async function markAsRead(userId) {
	try {
		// Chiama l'endpoint PUT /conversations/:id/seen
		await api.put(`/conversations/${userId}/seen`)
	} catch (e) {
		console.error("Errore markAsRead:", e)
	}
}

// 3. Seleziona chat
function selectChat(userOrChat, isNewSearch = false) {
	let targetId, targetName

	if (isNewSearch) {
		targetId = userOrChat.id
		targetName = userOrChat.username
	} else {
		targetId = userOrChat.user_id
		targetName = userOrChat.username
	}

	activeChat.value = {
		userId: targetId,
		username: targetName
	}
	
	activeChatMessages.value = [] 
	
	// Appena apro la chat:
	getChatMessages(targetId) // Scarico i messaggi
	markAsRead(targetId)      // Dico che li ho letti 
	startPolling()            // Inizio a controllare aggiornamenti

	if (isNewSearch) {
		searchQuery.value = ''
		foundUsers.value = []
	}
}

// 4. Invia messaggio
async function sendMessage() {
	if (!newMessage.value.trim() || !activeChat.value) return

	try {
		await api.post(`/conversations/${activeChat.value.userId}/messages`, {
			content: newMessage.value 
		})
		newMessage.value = ''
		
		// Aggiorniamo subito
		await getChatMessages(activeChat.value.userId)
		refreshConversations()
	} catch (e) {
		console.error("Errore invio:", e)
		alert("Errore invio messaggio")
	}
}

// --- POLLING AUTOMATICO ---
function startPolling() {
	// Ferma eventuali timer precedenti
	stopPolling()
	
	// Avvia un timer che ogni 3 secondi scarica i messaggi nuovi
	pollingInterval = setInterval(() => {
		if (activeChat.value) {
			getChatMessages(activeChat.value.userId, true) // true = non scrollare
			refreshConversations() // Aggiorna anche la lista a sinistra (per l'anteprima)
		}
	}, 3000)
}

function stopPolling() {
	if (pollingInterval) {
		clearInterval(pollingInterval)
		pollingInterval = null
	}
}

// Helper scroll
function scrollToBottom() {
	nextTick(() => {
		const container = document.querySelector('.messages-container')
		if (container) container.scrollTop = container.scrollHeight
	})
}

// Helper tempo
function formatTime(isoString) {
	if (!isoString) return ''
	const date = new Date(isoString)
	return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function doLogout() {
	stopPolling()
	localStorage.removeItem('token')
	localStorage.removeItem('username')
	router.push('/')
}

onMounted(() => {
	if (!token) router.push('/')
	else refreshConversations()
})

// Quando chiudo la pagina o cambio view, uccido il timer
onUnmounted(() => {
	stopPolling()
})
</script>

<template>
	<div class="main-container">
		<aside class="sidebar">
			<div class="header">
				<div class="my-info">
					<div class="avatar-small" @click="openProfile">
						{{ myUsername.charAt(0).toUpperCase() }}
					</div>
					<h3>{{ myUsername }}</h3>
				</div>
				<div class="actions">
					<button @click="openProfile" class="btn-icon" title="Profilo">⚙️</button>
					<button @click="doLogout" class="btn-small">Esci</button>
				</div>
			</div>
			
			<div class="search-box">
				<input 
					v-model="searchQuery" 
					@input="searchUsers"
					type="text" 
					placeholder="Cerca utente..."
				/>
			</div>

			<div v-if="foundUsers.length > 0" class="search-results">
				<p class="section-title">Risultati Ricerca:</p>
				<div 
					v-for="user in foundUsers" 
					:key="user.id" 
					class="chat-item"
					@click="selectChat(user, true)"
				>
					<div class="avatar-circle">{{ user.username.charAt(0) }}</div>
					<div class="chat-info">
						<span class="chat-name">{{ user.username }}</span>
						<span class="chat-preview">Clicca per chattare</span>
					</div>
				</div>
			</div>

			<div class="conversation-list">
				<p class="section-title">Le tue chat:</p>
				<div v-if="conversations.length === 0" class="empty-msg">
					Nessuna chat attiva.
				</div>

				<div 
					v-for="chat in conversations" 
					:key="chat.user_id" 
					class="chat-item"
					@click="selectChat(chat)"
				>
					<div class="avatar-circle">{{ chat.username ? chat.username.charAt(0) : '?' }}</div>
					<div class="chat-info">
						<span class="chat-name">{{ chat.username }}</span>
						<span class="chat-preview">{{ chat.last_message || 'Apri per leggere' }}</span>
					</div>
				</div>
			</div>
		</aside>

        <section class="chat-area">
			<div v-if="activeChat" class="chat-window">
				<div class="chat-header">
					<h3>Chat con {{ activeChat.username }}</h3>
				</div>

				<div class="messages-container">
					<div 
						v-for="msg in activeChatMessages" 
						:key="msg.id" 
						class="message-bubble"
						:class="{ 'my-message': msg.sender_id === myID, 'other-message': msg.sender_id !== myID }"
					>
						<div class="msg-text">{{ msg.content }}</div>
						
						<div class="msg-info">
                            <span class="msg-time">{{ formatTime(msg.timestamp) }}</span>
                            
                            <span v-if="msg.sender_id === myID" class="msg-status">
                                <span v-if="msg.read" class="read-icon">✓✓</span>
                                
                                <span v-else-if="msg.received">✓✓</span>
                                
                                <span v-else>✓</span>
                            </span>
                        </div>
					</div>
				</div>

				<div class="input-area">
					<input 
						v-model="newMessage" 
						@keyup.enter="sendMessage"
						type="text" 
						placeholder="Scrivi un messaggio..." 
						autofocus
					/>
					<button @click="sendMessage" :disabled="!newMessage.trim()">
						Invia
					</button>
				</div>
			</div>

			<div v-else class="welcome-screen">
				<h2>Benvenuto su WASAText</h2>
				<p>Seleziona una chat o cerca un utente per iniziare.</p>
			</div>
		</section>
	</div>
	<div v-if="showProfileModal" class="modal-overlay">
        <div class="modal-box">
            <h2>Il tuo Profilo</h2>
            
            <div class="form-group">
                <label>Cambia Nome Utente:</label>
                <input v-model="newUsername" type="text" placeholder="Nuovo nome..." />
            </div>

            <div class="photo-placeholder">
                (Presto potrai caricare una foto qui)
            </div>

            <div class="modal-actions">
                <button @click="showProfileModal = false" class="btn-cancel">Annulla</button>
                <button @click="updateUsername" class="btn-save">Salva</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
* {
    box-sizing: border-box;
}

.main-container {
    display: flex;
    height: 100%;
    width: 100%;
    background-color: #f0f2f5;
}

/* --- SIDEBAR --- */
.sidebar {
    width: 30%;
    min-width: 300px;
    background: white;
    border-right: 1px solid #ddd;
    display: flex;
    flex-direction: column;
    height: 100%;
}

.header {
    background-color: #f0f2f5;
    padding: 10px 15px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #ddd;
    height: 60px;
    flex-shrink: 0;
}

.search-box {
    padding: 10px;
    border-bottom: 1px solid #eee;
    flex-shrink: 0;
}
.search-box input {
    width: 100%;
    padding: 8px;
    border-radius: 8px;
    border: 1px solid #ddd;
    background: #f0f2f5;
}

.conversation-list, .search-results {
    flex: 1;
    overflow-y: auto;
}

.chat-item {
    display: flex;
    padding: 10px 15px;
    cursor: pointer;
    align-items: center;
    border-bottom: 1px solid #f9f9f9;
    height: 72px;
}
.chat-item:hover { background-color: #f5f5f5; }

.avatar-circle {
    width: 45px; height: 45px;
    background-color: #ddd; border-radius: 50%;
    display: flex; align-items: center; justify-content: center;
    margin-right: 15px; font-weight: bold; color: #555;
    flex-shrink: 0;
}

.chat-info {
    display: flex; flex-direction: column;
    overflow: hidden; width: 100%;
}
.chat-name { font-weight: bold; font-size: 1rem; }
.chat-preview {
    font-size: 0.85rem; color: #777;
    white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}

/* --- CHAT AREA --- */
.chat-area {
    flex: 1;
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: #e5ddd5;
}

.chat-window {
    display: flex;
    flex-direction: column;
    height: 100%;
    width: 100%;
}

.chat-header {
    background: #f0f2f5;
    padding: 10px 20px;
    border-bottom: 1px solid #ddd;
    display: flex;
    align-items: center;
    height: 60px;
    flex-shrink: 0;
}

/* MESSAGGI */
.messages-container {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
}

/* INPUT AREA: Ancorata in basso */
.input-area {
    background: #f0f2f5;
    padding: 10px;
    display: flex;
    gap: 10px;
    border-top: 1px solid #ddd;
    min-height: 60px;
    flex-shrink: 0;
}
.input-area input {
    flex: 1; padding: 12px; border-radius: 20px;
    border: none; outline: none;
}
.input-area button {
    padding: 10px 20px; border: none; border-radius: 20px;
    background-color: #008069; color: white; cursor: pointer;
}
.input-area button:disabled { background-color: #ccc; cursor: not-allowed; }

/* NUVOLETTE */
.message-bubble {
    max-width: 60%; padding: 8px 12px;
    border-radius: 8px; margin-bottom: 10px;
    display: flex; flex-direction: column;
    font-size: 0.95rem; line-height: 1.4; word-wrap: break-word;
}
.my-message { align-self: flex-end; background-color: #d9fdd3; margin-left: auto; }
.other-message { align-self: flex-start; background-color: white; margin-right: auto; }

.msg-info {
    display: flex; justify-content: flex-end; align-items: center;
    font-size: 0.7rem; color: #667781; gap: 5px; margin-top: 2px;
}
.read-icon { color: #53bdeb; }
.msg-status span:not(.read-icon) { color: #999; } /* Grigio per Inviato/Ricevuto */

/* Stili Header Aggiornati */
.my-info {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
}
.avatar-small {
    width: 35px; height: 35px;
    background-color: #008069; color: white;
    border-radius: 50%;
    display: flex; align-items: center; justify-content: center;
    font-weight: bold;
}
.actions {
    display: flex; gap: 10px;
}
.btn-icon {
    background: none; border: none; cursor: pointer; font-size: 1.2rem;
}

/* Stili Modale */
.modal-overlay {
    position: fixed;
    top: 0; left: 0; width: 100%; height: 100%;
    background: rgba(0,0,0,0.5);
    display: flex; align-items: center; justify-content: center;
    z-index: 1000;
}
.modal-box {
    background: white;
    padding: 20px;
    border-radius: 10px;
    width: 90%; max-width: 400px;
    box-shadow: 0 4px 10px rgba(0,0,0,0.2);
}
.form-group {
    margin: 20px 0;
}
.form-group input {
    width: 100%; padding: 10px;
    border: 1px solid #ddd; border-radius: 5px;
}
.modal-actions {
    display: flex; justify-content: flex-end; gap: 10px;
}
.btn-save {
    background-color: #008069; color: white;
    border: none; padding: 8px 15px; border-radius: 5px; cursor: pointer;
}
.btn-cancel {
    background-color: #ddd; color: black;
    border: none; padding: 8px 15px; border-radius: 5px; cursor: pointer;
}
</style>