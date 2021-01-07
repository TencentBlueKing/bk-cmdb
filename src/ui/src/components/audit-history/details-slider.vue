<template>
    <bk-sideslider
        :width="760"
        :title="$t('操作详情')"
        :is-show.sync="isShow"
        @hidden="handleHidden">
        <div class="details-content" slot="content"
            v-bkloading="{ isLoading: pending }">
            <component
                v-if="details"
                :is="detailsType"
                :details="details">
            </component>
        </div>
    </bk-sideslider>
</template>

<script>
    import DetailsJson from './details-json'
    import DetailsTable from './details-table'
    export default {
        components: {
            [DetailsJson.name]: DetailsJson,
            [DetailsTable.name]: DetailsTable
        },
        props: {
            id: Number
        },
        data () {
            return {
                details: null,
                isShow: false,
                pending: true
            }
        },
        computed: {
            detailsType () {
                if (!this.details) {
                    return null
                }
                const withCompare = ['host', 'module', 'set', 'mainline_instance', 'model_instance', 'business', 'cloud_area']
                return withCompare.includes(this.details.resource_type) ? DetailsTable.name : DetailsJson.name
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
                    this.details = null
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
        padding: 20px;
    }
</style>
