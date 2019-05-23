<template>
    <div id="createemail">
        <div class="formrow">
            <div class="formlabel">
                to
            </div>
            <div class="forminput">
               <input type="text" class="syl-input" v-model="to" placeholder="addresses of recipients, separated by cmomas" /> 
            </div>
        </div>
        <div class="formrow">
            <div class="formlabel">
               subject 
            </div>
            <div class="forminput">
               <input type="text" class="syl-input" v-model="subject" placeholder="addresses of recipients, separated by cmomas" /> 
            </div>
        </div>
        <div class="formrow">
            <div class="formlabel">
                body
            </div>
            <div class="forminput">
               <textarea class="syl-input" rows="20" v-model="body" placeholder="addresses of recipients, separated by cmomas"></textarea>
            </div>
        </div>
        <div class="formrow">
            <div class="formlabel">
                when (interval)
            </div>
            <div class="forminput">
               <div class="input-alt">
                   <input type="text" class="syl-input" v-model="whenInterval" placeholder="interval here" />
               </div>
               <div class="input-alt">
                   Format uses h for hours, m for minutes and s for seconds, without spaces. Eg, 10 hours, 25 minutes and 30 seconds is 10h25m30s
               </div>
            </div>
        </div>
        <div class="formrow">
            <div class="formlabel">or when (date-time)</div>
            <div class="forminput">
               <div class="input-alt">
                   <datetime type="datetime" v-model="whenDate" class="syl-input"></datetime>
               </div>
            </div>
        </div>
        <div class="formrow">
            <button class="createactionbtn" @click="createAction">create action</button>
        </div>
    </div>
</template>

<style>
#createemail {
    width: 100vh;
    display: grid;
    grid-gap: 10px;
    grid-template-columns: repeat(2, 1fr);
}

.formrow {
    width: 100vh;
    grid-column: span 2;
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    grid-gap: 10px;
}

.formlabel {
    grid-column: span 1;
    font-variant: small-caps;
    font-weight: bolder;
}
.forminput {
    grid-column: span 3;
}
.syl-input {
    width: 100%;
    padding: 4px;
}
</style>

<script>
import services from '../services';
import Settings from '../config';
import { mapGetters } from 'vuex';
const { APIUrl } = Settings;
const { session } = services(APIUrl);

export default {
    name: 'CreateEmail',
    data() {
        return {
            whenInterval: '10s',
            whenDate: '',
            to: 'joe.minichino@gmail.com',
            body: 'testing syl',
            subject: 'testing syl subject',
        };
    },
    methods: {
      async createAction() {
          const data = {
              userId: this.$store.getters.USER.Email,
              ex: this.whenInterval,
              to: this.to,
              body: this.body,
              subject: this.subject,
          };
          const { error } = await session.saveEmailAction(data);
          if (error) {
              this.$toasted.error('Ooops! Something went wrong with your request. Please retry', {
                  duration: 6000,
              });
              return;
          }
          this.$toasted.success('Email action created', {
              duration: 3000,
          })
          this.$emit('emailcreated');
      },
      ...mapGetters([
          'USER',
      ]),
    },
}
</script>
