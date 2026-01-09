import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: LoginView
    },
    // 2. AGGIUNGI QUESTA ROTTA:
    {
      path: '/home',
      name: 'home',
      component: HomeView
    }
  ]
})

export default router