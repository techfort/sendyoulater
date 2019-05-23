<template>
    <div id="header">
        <div v-if="USER">SendYouLater | Hello {{ USER.Name }}
            <router-link to="/">browse</router-link>
            <router-link to="create">create action</router-link>
        </div>
        <div v-else>SendYouLater</div>
    </div>
</template>
<script>
import { mapGetters } from 'vuex';

export default {
    name: 'Header',
    computed: {
        ...mapGetters([
            'USER',
            'EVENTS',
        ])
    }
};
</script>

