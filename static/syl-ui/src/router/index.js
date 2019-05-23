import VueRouter from 'vue-router';
import LandingPage from '../components/LandingPage.vue';
import Browser from '../components/Browser.vue';
import CreateEmail from '../components/CreateEmail.vue';
import Vue from 'vue';

Vue.use(VueRouter);

const routes = [
  { path: '/', name: 'home', component:  LandingPage, },
  { path: '/app', name: 'app', component: Browser, },
  { path: '/create', name: 'create', component: CreateEmail, }
];

export default new VueRouter({
  routes,
});