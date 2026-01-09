<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()

// Variabili reattive (stato)
const username = ref('')
const errorMessage = ref(null)

// Funzione chiamata al click del bottone
async function doLogin() {
	if (username.value.length < 3) {
		errorMessage.value = "Il nome deve avere almeno 3 caratteri."
		return
	}

	try {

		const response = await axios.post('http://localhost:3000/session', {
			name: username.value
		})

		localStorage.setItem('token', response.data.id || response.data.identifier) 
		localStorage.setItem('username', username.value)

		
		router.push('/home')
		
	} catch (e) {
		console.error(e)
		if (e.response && e.response.status === 400) {
			errorMessage.value = "Nome utente non valido."
		} else {
			errorMessage.value = "Impossibile contattare il server."
		}
	}
}
</script>

<template>
	<div class="login-container">
		<div class="login-box">
			<h2>Benvenuto in WASAText</h2>
			<p>Inserisci il tuo nome per entrare</p>
			
			<div class="input-group">
				<input 
					v-model="username" 
					type="text" 
					placeholder="Il tuo nome (es. Marco)" 
					@keyup.enter="doLogin"
				/>
				<button @click="doLogin">Entra</button>
			</div>

			<p v-if="errorMessage" class="error">{{ errorMessage }}</p>
		</div>
	</div>
</template>

<style scoped>
.login-container {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 80vh;
}

.login-box {
	background: #f0f2f5;
	padding: 2rem;
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0,0,0,0.1);
	text-align: center;
	width: 100%;
	max-width: 400px;
}

input {
	padding: 10px;
	border: 1px solid #ccc;
	border-radius: 5px;
	width: 70%;
	margin-right: 5px;
}

button {
	padding: 10px 20px;
	background-color: #008069;
	color: white;
	border: none;
	border-radius: 5px;
	cursor: pointer;
}

button:hover {
	background-color: #006d59;
}

.error {
	color: red;
	margin-top: 10px;
	font-size: 0.9em;
}
</style>