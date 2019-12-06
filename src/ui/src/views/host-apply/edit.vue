<template>
    <div class="apply-edit">
        <component
            :is="currentView"
            :module-ids="moduleIds"
            :action="action"
        >
        </component>
    </div>
</template>

<script>
    import multiModuleConfig from './children/multi-module-config'
    import singleModuleConfig from './children/single-module-config'
    import { MENU_BUSINESS_HOST_APPLY } from '@/dictionary/menu-symbol'
    export default {
        components: {
            multiModuleConfig,
            singleModuleConfig
        },
        data () {
            return {
                currentView: ''
            }
        },
        computed: {
            moduleIds () {
                const mid = this.$route.query.mid
                let moduleIds = []
                if (mid) {
                    moduleIds = String(mid).split(',').map(id => Number(id))
                }
                return moduleIds
            },
            isBatch () {
                return this.$route.query.hasOwnProperty('batch')
            },
            action () {
                return this.$route.query.action
            },
            title () {
                let title = '编辑'
                if (this.isBatch) {
                    title = this.action === 'batch-del' ? '批量删除' : '批量编辑'
                }
                return title
            }
        },
        created () {
            this.setBreadcrumbs()
            this.currentView = this.isBatch ? multiModuleConfig.name : singleModuleConfig.name
        },
        methods: {
            setBreadcrumbs () {
                this.$store.commit('setTitle', this.$t(this.title))
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('主机属性自动应用'),
                    route: {
                        name: MENU_BUSINESS_HOST_APPLY
                    }
                }, {
                    label: this.$t(this.title)
                }])
            }
        }
    }
</script>

<style lang="scss" scoped>
    .apply-edit {
        padding: 0 20px;
    }
</style>
