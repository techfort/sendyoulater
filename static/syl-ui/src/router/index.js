import VueRouter from 'vue-router';
import LandingPage from '../components/LandingPage.vue';
import Browser from '../components/Browser.vue';
import Vue from 'vue';

Vue.use(VueRouter);

const routes = [
  { path: '/', name: 'home', component:  LandingPage, },
  { path: '/app', name: 'app', component: Browser, }
];

export default new VueRouter({
  routes,
});