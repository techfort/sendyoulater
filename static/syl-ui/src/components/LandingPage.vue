<template>
    <div id="landingpage">
        <h1>SendYouLater</h1>
        <div v-if="USER">
            Logged in as {{ USER.U3 }}
        </div>
        <div v-else>
            <g-signin-button
                :params="googleSignInParams"
                @success="onSignInSuccess"
                @error="onSignInError">
                Sign in with Google
            </g-signin-button>
        </div>
    </div>
</template>

<script>
import Vue from 'vue'
import GSignInButton from 'vue-google-signin-button'
Vue.use(GSignInButton)

import { mapGetters } from 'vuex';

export default {
    name: 'LandingPage',
    computed: {
        ...mapGetters([
            'USER',
            'EVENTS',
        ])
    },
    data () {
        return {
        /**
         * The Auth2 parameters, as seen on
         * https://developers.google.com/identity/sign-in/web/reference#gapiauth2initparams.
         * As the very least, a valid client_id must present.
         * @type {Object} 
         */
            googleSignInParams: {
                client_id: "541640626027-l7s3mcv05cbdhqsq0vf54tcvpprb6s63.apps.googleusercontent.com",
                client_secret: "5jbcSzmUBPjFKww6BsoEKpC8",
                project_id: "spry-surf-230621",
            }
        }
    },
    methods: {
        onSignInSuccess (googleUser) {
            // `googleUser` is the GoogleUser object that represents the just-signed-in user.
            // See https://developers.google.com/identity/sign-in/web/reference#users
            const profile = googleUser.getBasicProfile()
            this.$store.dispatch('setUser', profile);
        },
        onSignInError (error) {
            // `error` contains any error occurred.
            console.log('OH NOES', error)
        }
    },
}
</script>

<style>
.g-signin-button {
  /* This is where you control how the button looks. Be creative! */
  display: inline-block;
  padding: 4px 8px;
  border-radius: 3px;
  background-color: #3c82f7;
  color: #fff;
  box-shadow: 0 3px 0 #0f69ff;
}
</style> 
