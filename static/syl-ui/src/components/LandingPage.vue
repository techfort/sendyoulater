<template>
    <div id="landingpage">
        <h1>SendYouLater</h1>
        <div v-if="USER">
            Logged in as {{ USER.Name }}
        </div>
        <div v-else>
            <button class="google-signin-button" @click="signIn">
                Sign in with Google
            </button>
        </div>
    </div>
</template>

<script>
import Vue from 'vue'
import GSignInButton from 'vue-google-signin-button'
Vue.use(GSignInButton)

import { mapGetters } from 'vuex';
import { post } from 'axios';

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
                client_id: "541640626027-uls0c5go44hag7oe1b4n78f6heqlqad4.apps.googleusercontent.com",
                client_secret: "5jbcSzmUBPjFKww6BsoEKpC8",
                scope: 'profile email openid',
            }
        }
    },
    methods: {
        
        async signIn() {
            gapi.load('auth2', async () => {
                console.log("gapi auth2", gapi.auth2)
                let auth2 = await gapi.auth2.init(this.googleSignInParams);
                auth2.grantOfflineAccess().then(this.signInCallback);
            });
        },
        signInCallback(authResult) {
            if (authResult['code']) {
                const { code } = authResult;
                console.log("CODE", code)
                // Hide the sign-in button now that the user is authorized, for example:
                $('.google-signin-button').attr('style', 'display: none');

                // Send the code to the server
                this.$http({
                    method: 'POST',
                    url: 'http://localhost:1323/token',
                    // Always include an `X-Requested-With` header in every AJAX request,
                    // to protect against CSRF attacks.
                    headers: {
                        'X-Requested-With': 'XMLHttpRequest',
                        'Content-Type': 'application/json'
                    },
                    success: function(result) {
                        console.log(result)    // Handle or verify the server response.
                    },
                    data: { code }
                });
                
            } else {
                // There was an error.
                console.log("There was an error")
            }
        },
    },
}
</script>

<style>
.google-signin-button {
  /* This is where you control how the button looks. Be creative! */
  display: inline-block;
  padding: 4px 8px;
  border-radius: 3px;
  background-color: #3c82f7;
  color: #fff;
  box-shadow: 0 3px 0 #0f69ff;
}
</style> 
