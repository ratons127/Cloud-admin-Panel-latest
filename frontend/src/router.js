import { createRouter, createWebHistory } from 'vue-router'
import AppShell from './components/app/AppShell.vue'
import VerifyEmail from './components/app/VerifyEmail.vue'
import InviteAccept from './components/app/InviteAccept.vue'
<<<<<<< HEAD
import ForgotPassword from './components/app/ForgotPassword.vue'
import ResetPassword from './components/app/ResetPassword.vue'
import Welcome from './components/app/Welcome.vue'
import OAuthCallback from './components/app/OAuthCallback.vue'
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af

export default createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', component: AppShell },
    { path: '/verify', component: VerifyEmail },
<<<<<<< HEAD
    { path: '/invite', component: InviteAccept },
    { path: '/forgot', component: ForgotPassword },
    { path: '/reset', component: ResetPassword },
    { path: '/welcome', component: Welcome },
    { path: '/oauth', component: OAuthCallback }
=======
    { path: '/invite', component: InviteAccept }
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
  ]
})
