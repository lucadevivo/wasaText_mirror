<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="app-container">
    <header class="main-header">
      <h1>WASAText</h1>
    </header>

    <main>
      <RouterView />
    </main>
  </div>
</template>

<style>
/* RESET GLOBALE: Blocca lo scroll della pagina intera */
html, body, #app {
  height: 100%;
  margin: 0;
  padding: 0;
  overflow: hidden; /* Nessuno scrollbar sul body! */
  font-family: sans-serif;
}

.app-container {
  display: flex;
  flex-direction: column; /* Impila Header sopra, Main sotto */
  height: 100%;
}

.main-header {
  background-color: #008069;
  color: white;
  padding: 15px 20px;
  flex-shrink: 0; /* L'header non si rimpicciolisce mai */
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
  z-index: 10;
}

main {
  flex: 1; /* Occupa tutto lo spazio rimanente */
  position: relative;
  overflow: hidden; /* Gestisce il contenuto interno */
  display: flex;
  flex-direction: column;
}
</style>