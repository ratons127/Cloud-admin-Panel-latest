import { createApp } from 'vue'
import App from './components/app/App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import 'roboto-fontface/css/roboto/roboto-fontface.css'

const app = createApp(App)

app.directive('blur', {
  mounted(el) {
    el.onfocus = (ev) => ev.target.blur()
  }
})

app.use(vuetify)
app.use(router)
app.use(store)
app.mount('#app')
