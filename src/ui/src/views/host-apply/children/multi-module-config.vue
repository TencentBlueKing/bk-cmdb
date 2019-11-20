<template>
    <div class="multi-module-config">
        <div class="config-bd">
            <div class="config-item">
                <div class="item-label">已选择 {{checkedModuleList.length}} 个模块：</div>
                <div class="item-content">
                    <div :class="['module-list', { 'show-more': showMore.isMoreModuleShowed }]" ref="moduleList">
                        <div class="module-item" :title="'a / b / c'" v-for="(m, index) in checkedModuleList" :key="index">
                            <span class="module-icon">{{$i18n.locale === 'en' ? 'M' : '模'}}</span>
                            GameSevers
                        </div>
                        <div
                            :class="['module-item', 'more', { 'opened': showMore.isMoreModuleShowed }]"
                            :style="{ left: `${showMore.linkLeft}px` }"
                            v-show="showMore.showLink" @click="handleShowMore"
                        >
                            {{showMore.isMoreModuleShowed ? '收起' : '展开更多'}}<i class="bk-cc-icon icon-cc-arrow-down"></i>
                        </div>
                    </div>
                </div>
            </div>
            <div class="config-item">
                <div class="item-label">
                    {{$t(isDel ? '请勾选要删除的字段' : '请添加自动应用的字段')}}
                </div>
                <div class="item-content">
                    <bk-button theme="default" icon="plus" @click="handleChooseField" class="choose-button" v-if="!isDel">添加字段</bk-button>
                    <div class="config-table" v-if="checkedPropertyIdList.length">
                        <property-config-table
                            ref="configEditTable"
                            :multiple="true"
                            :readonly="isDel"
                            :deletable="isDel"
                            :checked-property-id-list="checkedPropertyIdList"
                            :init-data="initConfigData"
                            @property-value-change="handlePropertyValueChange"
                            @selection-change="handlePropertySelectionChange"
                        >
                        </property-config-table>
                    </div>
                </div>
            </div>
        </div>
        <div class="config-ft">
            <bk-button theme="primary" :disabled="nextButtonDisabled" @click="handleNextStep" v-if="!isDel">下一步</bk-button>
            <bk-button theme="primary" :disabled="delButtonDisabled" @click="handleDel" v-else>确定删除</bk-button>
            <bk-button theme="default">取消</bk-button>
        </div>

        <host-property-modal
            :visiable.sync="propertyModalVisiable"
            :checked-list.sync="checkedPropertyIdList"
        >
        </host-property-modal>
    </div>
</template>
<script>
    import { mapGetters } from 'vuex'
    import hostPropertyModal from './host-property-modal'
    import propertyConfigTable from './property-config-table'
    export default {
        components: {
            hostPropertyModal,
            propertyConfigTable
        },
        data () {
            return {
                propertyModalVisiable: false,
                nextButtonDisabled: true,
                delButtonDisabled: true,
                initConfigData: [],
                checkedPropertyIdList: [],
                checkedModuleList: [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
                showMore: {
                    moduleListMaxRow: 2,
                    showLink: false,
                    isMoreModuleShowed: false,
                    linkLeft: 0
                },
                selectedPropertyRow: []
            }
        },
        computed: {
            ...mapGetters('hosts', ['configPropertyList']),
            isDel () {
                return this.$route.query.action === 'batch-del'
            }
        },
        watch: {
            configPropertyList: {
                handler () {
                    this.initData()
                },
                immediate: true
            },
            checkedPropertyIdList () {
                this.$nextTick(() => {
                    this.toggleNexButtonDisabled()
                })
            }
        },
        created () {
            this.initData()
        },
        mounted () {
            this.setShowMoreLinkStatus()
            window.addEventListener('resize', this.setShowMoreLinkStatus)
        },
        beforeDestroy () {
            window.removeEventListener('resize', this.setShowMoreLinkStatus)
        },
        methods: {
            async initData () {
                // 请求接口获取配置数据
                this.initConfigData = await Promise.resolve([
                    {
                        bk_property_id: 'bk_asset_id',
                        module_list: [
                            {
                                bk_inst_id: 141,
                                path: '广东区 / 大区一 / Apache模块1',
                                value: 'No13878667'
                            },
                            {
                                bk_inst_id: 142,
                                path: '广东区 / 大区一 / Apache模块2',
                                value: 'No13878645'
                            }
                        ]
                    },
                    {
                        bk_property_id: 'bk_state_name',
                        module_list: [
                            {
                                bk_inst_id: 143,
                                path: '广东区 / 大区一 / Apache模块1',
                                value: 'BR'
                            },
                            {
                                bk_inst_id: 132,
                                path: '广东区 / 大区一 / Apache模块2',
                                value: 'CN'
                            },
                            {
                                bk_inst_id: 112,
                                path: '广东区 / 大区一 / Apache模块2',
                                value: 'US'
                            }
                        ]
                    }
                ])

                this.checkedPropertyIdList = this.initConfigData.map(item => item.bk_property_id)
            },
            setShowMoreLinkStatus () {
                const moduleList = this.$refs.moduleList
                const moduleItemEl = moduleList.getElementsByClassName('module-item')[0]
                const moduleItemStyle = getComputedStyle(moduleItemEl)
                const moduleItemWidth = moduleItemEl.offsetWidth + parseInt(moduleItemStyle.marginLeft, 10) + parseInt(moduleItemStyle.marginRight, 10)
                const moduleListWidth = moduleList.clientWidth
                const maxCountInRow = Math.floor(moduleListWidth / moduleItemWidth)
                const rowCount = Math.ceil(this.checkedModuleList.length / maxCountInRow)
                this.showMore.showLink = rowCount > this.showMore.moduleListMaxRow
                this.showMore.linkLeft = moduleItemWidth * (maxCountInRow - 1)
            },
            toggleNexButtonDisabled () {
                const modulePropertyList = this.$refs.configEditTable.modulePropertyList
                this.nextButtonDisabled = !this.checkedPropertyIdList.length || !modulePropertyList.every(property => property.__extra__.value)
            },
            handleNextStep () {
                const modulePropertyList = this.$refs.configEditTable.modulePropertyList
                const configData = modulePropertyList.map(property => ({
                    key: property.bk_property_id,
                    value: property.__extra__.value
                }))

                this.$router.push({
                    name: 'hostApplyConfirm',
                    query: {
                        batch: 1
                    }
                })
                console.log(configData)
            },
            handleDel () {
                this.$bkInfo({
                    title: this.$t('确认删除自动应用字段？'),
                    subTitle: this.$t('删除后，将会移除字段在对应模块中的配置'),
                    confirmFn: () => {
                    }
                })
            },
            handlePropertySelectionChange (value) {
                this.selectedPropertyRow = value
                this.delButtonDisabled = this.selectedPropertyRow.length <= 0
            },
            handlePropertyValueChange () {
                this.toggleNexButtonDisabled()
            },
            handleChooseField () {
                this.propertyModalVisiable = true
            },
            handleShowMore () {
                this.showMore.isMoreModuleShowed = !this.showMore.isMoreModuleShowed
            }
        }
    }
</script>
<style lang="scss" scoped>
    .multi-module-config {
        --labelWidth: 180px;
        .config-item {
            display: flex;
            margin: 8px 0;

            .item-label {
                flex: none;
                width: var(--labelWidth);
                font-size: 14px;
                font-weight: bold;
                color: #63656e;
                text-align: right;
                margin-right: 12px;
            }
            .item-content {
                flex: auto;
            }

            .choose-button {
                margin-bottom: 18px;;
            }
        }
        .config-ft {
            margin-left: calc(var(--labelWidth) + 12px);
            margin-top: 20px;
            .bk-button {
                min-width: 86px;
            }
        }
    }

    .module-list {
        position: relative;
        max-height: 72px;
        overflow: hidden;
        transition: all .2s ease-out;

        &.show-more {
            max-height: 100%;
        }
    }
    .module-item {
        position: relative;
        display: inline-block;
        vertical-align: middle;
        height: 26px;
        width: 120px;
        margin: 0 10px 10px 0;
        line-height: 24px;
        padding: 0 20px 0 25px;
        border: 1px solid #c4c6cc;
        border-radius: 13px;
        color: $textColor;
        font-size: 12px;
        cursor: default;
        @include ellipsis;
        .module-icon {
            position: absolute;
            left: 2px;
            top: 2px;
            width: 20px;
            height: 20px;
            border-radius: 50%;
            line-height: 20px;
            text-align: center;
            color: #fff;
            font-size: 12px;
            background-color: #c4c6cc;
        }

        &.more {
            position: absolute;
            left: 0;
            bottom: 0;
            background: #fafbfd;
            border: 0 none;
            border-radius: unset;
            cursor: pointer;
            color: #3a84ff;
            font-size: 14px;
            text-align: center;
            padding: 0;
            .bk-cc-icon {
                font-size: 22px;
            }

            &.opened {
                position: static;
                .bk-cc-icon {
                    transform: rotate(180deg);
                }
            }
        }
    }
</style>
