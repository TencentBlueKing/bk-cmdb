<template>
    <div class="batch-edit">
        <multi-module-config></multi-module-config>
    </div>
</template>

<script>
    import multiModuleConfig from './children/multi-module-config'
    import { MENU_BUSINESS_HOST_APPLY } from '@/dictionary/menu-symbol'
    export default {
        components: {
            multiModuleConfig
        },
        data () {
            return {
            }
        },
        created () {
            this.setBreadcrumbs()
        },
        methods: {
            setBreadcrumbs () {
                this.$store.commit('setTitle', this.$t('批量编辑'))
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('主机属性自动应用'),
                    route: {
                        name: MENU_BUSINESS_HOST_APPLY
                    }
                }, {
                    label: this.$t('批量编辑')
                }])
            }
        }
    }
</script>

<style lang="scss" scoped>
    .batch-edit {
        padding: 0 20px;
    }
</style>
