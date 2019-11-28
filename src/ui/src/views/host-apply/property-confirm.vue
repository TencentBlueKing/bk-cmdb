<template>
    <div class="host-apply-confirm">
        <feature-tips
            :feature-name="'hostApply'"
            :show-tips="showFeatureTips"
            :desc="$t('冲突的属性若不解决，在应用时将被忽略，确定应用后，新转入该模块的主机将自动应用配置的属性')"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <div class="caption">
            <div class="title">请确认以下主机应用信息：</div>
            <div class="stat">
                <span class="conflict-item">
                    <span v-bk-tooltips="{ content: $t('当一台主机同属多个模块时，且模块配置了不同的属性'), width: 185 }">
                        <i class="bk-cc-icon icon-cc-tips"></i>
                    </span>
                    属性冲突<em class="conflict-num">{{conflictNum}}</em>台
                </span>
                <span>总检测<em class="check-num">{{table.total}}</em>台</span>
            </div>
        </div>
        <property-confirm-table
            ref="propertyConfirmTable"
            :list="table.list"
            :total="table.total"
        >
        </property-confirm-table>
        <div class="bottom-actionbar">
            <div class="actionbar-inner">
                <bk-button theme="primary" @click="handleApply">应用</bk-button>
                <bk-button theme="default">上一步</bk-button>
                <bk-button theme="default">取消</bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapState, mapActions } from 'vuex'
    import featureTips from '@/components/feature-tips/index'
    import propertyConfirmTable from './children/property-confirm-table'
    import { MENU_BUSINESS_HOST_APPLY } from '@/dictionary/menu-symbol'
    export default {
        components: {
            featureTips,
            propertyConfirmTable
        },
        data () {
            return {
                showFeatureTips: false,
                table: {
                    list: [],
                    total: 0
                },
                conflictNum: 0
            }
        },
        computed: {
            ...mapState('hostApply', ['propertyConfig']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters(['featureTipsParams', 'supplierAccount']),
            isBatch () {
                return this.$route.query.batch === 1
            }
        },
        watch: {

        },
        created () {
            this.showFeatureTips = this.featureTipsParams['hostApplyConfirm']
            this.setBreadcrumbs()
            this.initData()
        },
        methods: {
            ...mapActions('hostApply', [
                'getApplyPreview',
                'runApply'
            ]),
            async initData () {
                try {
                    const previewData = await this.getApplyPreview({
                        bizId: this.bizId,
                        params: this.propertyConfig,
                        config: {
                            requestId: 'getHostApplyPreview'
                        }
                    })

                    this.table.list = previewData.plans || []
                    this.table.total = previewData.count
                    this.conflictNum = previewData.unresolved_conflict_count
                } catch (e) {
                    console.error(e)
                }
            },
            setBreadcrumbs () {
                const title = this.isBatch ? '批量应用属性' : '应用属性'
                this.$store.commit('setTitle', this.$t(title))
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('主机属性自动应用'),
                    route: {
                        name: MENU_BUSINESS_HOST_APPLY
                    }
                }, {
                    label: this.$t(title)
                }])
            },
            async handleApply () {
                const conflictResolveResult = this.$refs.propertyConfirmTable.conflictResolveResult
                const conflictResolvers = []
                Object.keys(conflictResolveResult).forEach(key => {
                    const propertyList = conflictResolveResult[key]
                    propertyList.forEach(property => {
                        conflictResolvers.push({
                            bk_host_id: Number(key),
                            bk_attribute_id: property.id,
                            bk_property_value: property.__extra__.value
                        })
                    })
                })

                // 合入冲突结果数据
                const propertyConfig = { ...this.propertyConfig, ...{ conflict_resolvers: conflictResolvers } }

                try {
                    const result = await this.runApply({
                        bizId: this.bizId,
                        params: propertyConfig,
                        config: {
                            requestId: 'runHostApply'
                        }
                    })

                    // 更新属性配置
                    this.$store.commit('hostApply/setPropertyConfig', propertyConfig)
                    console.log(result)
                } catch (e) {
                    console.error(e)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-apply-confirm {
        padding: 0 20px;

        .caption {
            display: flex;
            margin-bottom: 14px;
            justify-content: space-between;
            align-items: center;

            .title {
                color: #63656e;
                font-size: 14px;
            }

            .stat {
                color: #313238;
                font-size: 12px;
                margin-right: 8px;

                .conflict-item {
                    margin-right: 12px;
                }

                .conflict-num,
                .check-num {
                    font-style: normal;
                    font-weight:bold;
                }
                .conflict-num {
                    color: #ff5656;
                }
                .check-num {
                    color: #2dcb56;
                }
            }
        }
    }

    .bottom-actionbar {
        position: absolute;
        width: 100%;
        height: 50px;
        border-top: 1px solid #dcdee5;
        bottom: 0;
        left: 0;

        .actionbar-inner {
            padding: 8px 0 0 20px;
            .bk-button {
                min-width: 86px;
            }
        }
    }
</style>
