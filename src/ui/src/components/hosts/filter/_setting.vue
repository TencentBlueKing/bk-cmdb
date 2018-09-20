<template>
    <div class="setting-layout">
        <i class="setting-icon icon-cc-collection" v-tooltip="$t('Hosts[\'收藏\']')"
            v-if="activeSetting.includes('collection')"
            :class="{active: collection.show}"
            @click="handleCollection">
        </i>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
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
                        requestId: `post_batchSearchObjectAttribute_${Object.keys(this.filterConfig.properties).join('_')}`,
                        requestGroup: Object.keys(this.filterConfig.properties).map(id => `post_searchObjectAttribute_${id}`),
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