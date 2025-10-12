// import '@yuelioi/components/the-footer.css'
// import '@yuelioi/toast/dist/toast.css'

import './assets/main.css'
import { createPinia } from 'pinia'

import { createApp } from 'vue'
const pinia = createPinia()
import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(router)
app.use(pinia)

app.mount('#app')
