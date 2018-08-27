<template>
    <div class="field-wrapper">
       <v-base-info></v-base-info> 
    </div>
</template>

<script>
    import vBaseInfo from './base-info'
    export default {
        components: {
            vBaseInfo
        }
    }
</script>

<style lang="scss" scoped>
    .field-wrapper {
        padding: 0 30px;
    }
</style>
