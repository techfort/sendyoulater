<template>
    <div id="header">
        <div v-if="USER">SendYouLater | Hello {{ USER.Name }}</div>
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

