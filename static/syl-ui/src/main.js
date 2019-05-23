import Vue from 'vue'
import axios from 'axios';
import Datetime from 'vue-datetime'
import Toasted from 'vue-toasted';
import App from './App.vue'
import store from './store'
import router from './router';
// You need a specific loader for CSS files
import 'vue-datetime/dist/vue-datetime.css';

Vue.use(Datetime)
Vue.use(Toasted);
Vue.http = Vue.prototype.$http = axios;
Vue.config.productionTip = false

new Vue({
  store,
  router,
  render: h => h(App),
}).$mount('#app')
