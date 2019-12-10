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
    import { mapGetters } from 'vuex'
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
                currentView: '',
                moduleMap: {}
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
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
                let title
                if (this.isBatch) {
                    title = this.$t(this.action === 'batch-del' ? '批量删除' : '批量编辑')
                } else {
                    title = `${this.$t('编辑')} ${this.getModuleName(this.moduleIds[0])}`
                }
                return title
            }
        },
        created () {
            this.initData()
            this.currentView = this.isBatch ? multiModuleConfig.name : singleModuleConfig.name
        },
        beforeRouteLeave (to, from, next) {
            if (to.name !== 'hostApplyConfirm') {
                this.$store.commit('hostApply/clearRuleDraft')
            }
            next()
        },
        methods: {
            async initData () {
                try {
                    const topopath = await this.getTopopath()
                    const moduleMap = {}
                    topopath.nodes.forEach(node => {
                        moduleMap[node.topo_node.bk_inst_id] = node.topo_path
                    })
                    this.moduleMap = Object.freeze(moduleMap)

                    this.setBreadcrumbs()
                } catch (e) {
                    console.log(e)
                }
            },
            setBreadcrumbs () {
                this.$store.commit('setTitle', this.title)
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('主机属性自动应用'),
                    route: {
                        name: MENU_BUSINESS_HOST_APPLY
                    }
                }, {
                    label: this.title
                }])
            },
            getTopopath () {
                return this.$store.dispatch('hostApply/getTopopath', {
                    bizId: this.bizId,
                    params: {
                        topo_nodes: this.moduleIds.map(id => ({ bk_obj_id: 'module', bk_inst_id: id }))
                    }
                })
            },
            getModulePath (id) {
                const info = this.moduleMap[id] || []
                const path = info.map(node => node.bk_inst_name).reverse().join(' / ')
                return path
            },
            getModuleName (id) {
                const topoInfo = this.moduleMap[id] || []
                const target = topoInfo.find(target => target.bk_obj_id === 'module' && target.bk_inst_id === id) || {}
                return target.bk_inst_name
            }
        }
    }
</script>

<style lang="scss" scoped>
    .apply-edit {
        padding: 0 20px;
    }
</style>
