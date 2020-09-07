<template>
    <bk-sideslider
        :width="800"
        :title="$t('操作详情')"
        :is-show.sync="isShow"
        @hidden="handleHidden">
        <div class="details-content" slot="content"
            v-bkloading="{ isLoading: pending }">
            <details-json :details="details"></details-json>
        </div>
    </bk-sideslider>
</template>

<script>
    import DetailsJson from './details-json'
    export default {
        components: {
            DetailsJson
        },
        props: {
            id: Number
        },
        data () {
            return {
                details: {},
                isShow: false,
                pending: true
            }
        },
        async created () {
            this.getDetails()
        },
        methods: {
            show () {
                this.isShow = true
            },
            handleHidden () {
                this.$emit('close')
            },
            async getDetails () {
                try {
                    this.pending = true
                    this.details = await this.$store.dispatch('audit/getDetails', { id: this.id })
                } catch (error) {
                    console.error(error)
                    this.details = {}
                } finally {
                    this.pending = false
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-content {
        height: calc(100vh - 60px);
        padding: 20px 40px;
    }
</style>
