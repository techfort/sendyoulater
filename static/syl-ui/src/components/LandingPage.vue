<template>
    <div id="mainview">
        <div class="menu">
            <div>Upcoming actions: {{ EMAILS.length }}</div>
            <div class="syl-action" v-for="e in EMAILS" v-bind:key="e.ID">
                    {{ e.To }}
            </div>

        </div>
        <div class="main">
            <div v-if="USER">
                Click on any action to view details.
            </div>
            <div v-else>
                <div id="gSignIn"></div>
            </div>
        </div>
    </div>
</template>

<style>

#mainview {
  grid-column: span 12;
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  grid-gap:10px;
}
.menu {
  background: #ccc;
  grid-column: span 4;
}
.main {
  background: #eee;
  grid-column: span 8;
}
.syl-action {
    font-size: 0.9em;
}
</style>


<script>
// import Vue from 'vue'
// import GSignInButton from 'vue-google-signin-button'
// Vue.use(GSignInButton)

import { mapGetters } from 'vuex';
import to from '../helpers';

export default {
    name: 'LandingPage',
    computed: {
        ...mapGetters([
            'USER',
            'EMAILS',
        ]),
    },
    watch: {
        EMAILS(n, o) {
            console.log(n, o);
        },
    },
    async mounted() {
        const that = this;
        const onSuccess = async (googleUser) => {
            const user = googleUser.getBasicProfile();
            const profile = {
                Name: user.getName(),
                Email: user.getEmail(),
                Id: user.getId(),
                Avatar: user.getImageUrl(),
            };
            const { data, error } = await to(this.$store.dispatch('setUser', profile));
            if (error) {
                console.log(`Error setting user: ${error}`)
                return;
            }
            console.log('emitting signin event');
            that.$emit('signin', profile);
        };
        console.log('EMAILS', this.$store.getters.EMAILS);
        const onFailure = (err) => {
            console.error(err);
        };
        gapi.signin2.render('gSignIn', {
            'scope': 'profile email openid',
            'width': 240,
            'height': 50,
            'longtitle': true,
            'theme': 'dark',
            'onsuccess': onSuccess,
            'onfailure': onFailure,
        });
    },
    data () {
        return {
        };
    },
    methods: {
    },
}
</script>