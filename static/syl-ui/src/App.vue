<template>
  <div id="app">
    <div class="header">
      <Header />
    </div>
    <div class="menu">MENU</div>
    <div class="main">MAIN</div>
    <div class="footer">FOOTER</div>
  </div>
</template>

<script>
import services from './services/';
import Header from './components/Header.vue';

const { session } = services('http://localhost:1323');

export default {
  name: 'app',
  components: {
    Header,
  },
  mounted() {
    setInterval(async () => {
      const res = await session.check();
      console.log(res);
    }, 5000);
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
.menu {
  background: #ccc;
  grid-column: span 4;
}
.main {
  background: #eee;
  grid-column: span 8;
}
.footer {
  background: #ddd;
  grid-column: span 12;
}
</style>
