<template>
    <div class="json-details-layout">
        <pre>{{details}}</pre>
    </div>
</template>

<script>
    export default {
        name: 'details-json',
        props: {
            details: Object
        }
    }
</script>

<style lang="scss" scoped>
    .json-details-layout {
        background: #455070;
        color: #bfc6e0;
        border-radius: 2px;
        font-size: 12px;
        line-height: 20px;
        height: 100%;
        overflow: auto;
    }
</style>
