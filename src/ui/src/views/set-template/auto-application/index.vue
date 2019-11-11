<template>
    <bk-sideslider
        v-transfer-dom
        :is-show.sync="conflict.show"
        :width="795"
        :title="$t('处理冲突', conflict.host)">
        <resolve-conflict slot="content"></resolve-conflict>
    </bk-sideslider>
</template>

<script>
    import resolveConflict from './conflict.vue'
    export default {
        components: {
            resolveConflict
        },
        data () {
            return {
                conflict: {
                    show: true,
                    host: '192.168.0.1'
                }
            }
        }
    }
</script>

<style lang="scss" scoped>

</style>
