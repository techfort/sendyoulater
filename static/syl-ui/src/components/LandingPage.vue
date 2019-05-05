<template>
    <div id="landingpage">
        <h1>SendYouLater</h1>
        <div v-if="USER">
            Logged in as {{ USER.Name }}
        </div>
        <div v-else>
            <div id="gSignIn"></div>
        </div>
    </div>
</template>

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
            'EVENTS',
        ])
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