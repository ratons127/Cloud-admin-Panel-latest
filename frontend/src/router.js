import { createRouter, createWebHistory } from 'vue-router'
import AppShell from './components/app/AppShell.vue'
import VerifyEmail from './components/app/VerifyEmail.vue'
import InviteAccept from './components/app/InviteAccept.vue'

export default createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', component: AppShell },
    { path: '/verify', component: VerifyEmail },
    { path: '/invite', component: InviteAccept }
  ]
})
