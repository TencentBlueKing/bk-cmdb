<template>
    <div class="setting-layout">
        <i class="setting-icon icon-cc-broom" v-tooltip="$t('HostResourcePool[\'清空查询条件\']')"
            v-if="activeSetting.includes('reset')"
            @click="handleReset">
        </i>
        <i class="setting-icon icon-cc-collection" v-tooltip="$t('Hosts[\'收藏\']')"
            v-if="activeSetting.includes('collection')"
            :class="{active: collection.show}"
            @click="handleCollection">
        </i>
        <i class="setting-icon icon-cc-funnel" v-tooltip="$t('HostResourcePool[\'设置筛选项\']')"
            v-if="activeSetting.includes('filter-config')"
            @click="filterConfig.show = true">
        </i>
        <cmdb-slider :is-show.sync="filterConfig.show" :title="$t('HostResourcePool[\'主机筛选项设置\']')" :width="600">
            <cmdb-filter-config slot="content"
                :properties="filterConfig.properties"
                :selected="customFilterFields"
                @on-cancel="filterConfig.show = false"
                @on-apply="handleApplyFilterConfig">
            </cmdb-filter-config>
        </cmdb-slider>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import cmdbFilterConfig from './_filter-config.vue'
    export default {
        components: {
            cmdbFilterConfig
        },
        props: {
            activeSetting: {
                type: Array,
                default () {
                    return ['reset', 'collection', 'filter-config']
                }
            },
            filterConfigKey: {
                type: String,
                required: true
            }
        },
        data () {
            return {
                filterConfig: {
                    show: false,
                    properties: {
                        'biz': [],
                        'host': [],
                        'set': [],
                        'module': []
                    }
                },
                collection: {
                    show: false
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('hostFavorites', ['applyingProperties']),
            customFilterFields () {
                return this.applyingProperties.length ? this.applyingProperties : (this.usercustom[this.filterConfigKey] || [])
            }
        },
        watch: {
            applyingProperties (properties) {
                let hasUnloadObj = false
                properties.forEach(property => {
                    if (!this.filterConfig.properties.hasOwnProperty(property['bk_obj_id'])) {
                        hasUnloadObj = true
                        this.$set(this.filterConfig.properties, property['bk_obj_id'], [])
                    }
                })
                if (hasUnloadObj) {
                    this.$http.cancel('hostsAttribute')
                    this.getProperties()
                }
            }
        },
        async created () {
            await this.getProperties()
        },
        methods: {
            ...mapActions('objectModelProperty', ['batchSearchObjectAttribute']),
            getProperties () {
                return this.batchSearchObjectAttribute({
                    params: {
                        bk_obj_id: {'$in': Object.keys(this.filterConfig.properties)},
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: 'hostsAttribute',
                        fromCache: true
                    }
                }).then(result => {
                    Object.keys(this.filterConfig.properties).forEach(objId => {
                        this.filterConfig.properties[objId] = result[objId]
                    })
                    this.$delete(this.filterConfig.properties, 'biz')
                    return result
                })
            },
            handleApplyFilterConfig (properties) {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.filterConfigKey]: properties.map(property => {
                        return {
                            'bk_property_id': property['bk_property_id'],
                            'bk_obj_id': property['bk_obj_id']
                        }
                    })
                }).then(() => {
                    this.$store.commit('hostFavorites/setApplying', null)
                })
                this.filterConfig.show = false
            },
            handleReset () {
                this.$emit('on-reset')
            },
            handleCollection () {
                this.$emit('on-collection')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .setting-icon{
        font-size: 16px;
        margin: 0 -20px 0 30px;
        cursor: pointer;
        color: #c3cdd7;
        &.active {
            color: #ffb400;
        }
    }
</style>