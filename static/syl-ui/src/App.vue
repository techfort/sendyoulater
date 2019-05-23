<template>
  <div id="app">
    <div class="header">
      <Header />
    </div>
    <router-view @signin="login" @emailcreated="loadData"></router-view>
    <div class="footer">FOOTER</div>
  </div>
</template>

<script>
import services from './services/';
import Header from './components/Header.vue';
import { mapGetters } from 'vuex';
import to from './helpers';
import Settings from './config';
const { APIUrl } = Settings;
const { session } = services(APIUrl);

export default {
  name: 'app',
  components: {
    Header,
  },
  methods: {
    async getUserData () { return to(session.getUserData()); },
    async login (profile) {
      console.log('Reacting to signin', profile);
      const { data, error } = await session.login(profile);
      console.log('Result of user data', data);
      this.loadData();
    },
    async loadData() {
      const { error, data } = await session.loadData(this.$store.getters.USER.Email);
      if (error) {
        console.log(error);
        return;
      }
      console.log('received emails', data.data);
      await this.$store.dispatch('setEmails', data.data);
    },
    ...mapGetters([
      'USER',
    ]),
  },
  async mounted() {
    
  },
  async created() {
  },
}
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=Titillium+Web:300');
*{
  box-sizing: border-box;
  padding: 0;
  margin: 0;
}
html, body {
  font-family: 'Titillium Web', Helvetica, sans-serif;
  font-size: 24px;
  height: 100vh;
  width: 100vw;
}

#app {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  grid-template-rows: 40px auto 40px;
  grid-auto-rows: 100px;
  grid-gap: 10px;
  height: 100vh;
  width: 100vw;

}
.header {
  background: #ddd;
  grid-column: span 12;
}

.footer {
  background: #ddd;
  grid-column: span 12;
}
</style>
