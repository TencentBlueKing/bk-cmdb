<template>
    <div class="edit-application-fields" ref="applicationFields">
        <div class="main-box">
            <div class="info">
                <label class="info-label">{{$t('已选择模块', { num: modules.length })}}</label>
                <div class="main">
                    <ul class="module-list" ref="moduleList" :style="{
                        height: listHeight()
                    }">
                        <li v-for="module in modules" :key="module.bk_module_id" ref="moduleItem" :title="'拓扑路径'">
                            <span class="module-icon">{{$i18n.locale === 'en' ? 'M' : '模'}}</span>
                            <span class="module-name">{{module.bk_module_name}}</span>
                            <i class="module-remove bk-icon icon-close" v-if="modules.length > 1" @click="handleRemoveModule(module.bk_module_id)"></i>
                        </li>
                    </ul>
                    <span class="collapse" v-if="showCollapse" @click="collapse = !collapse">
                        {{collapse ? $t('收起') : $t('展开') }}
                        <i :class="['bk-icon', 'icon-angle-down', { 'is-collapse': collapse }]"></i>
                    </span>
                </div>
            </div>
            <div class="info clearfix">
                <label class="info-label">{{$t('请勾选要删除的字段')}}</label>
                <div class="main field-table">
                    <bk-table :data="list" @selection-change="handleSelectionChange">
                        <bk-table-column type="selection" width="50"></bk-table-column>
                        <bk-table-column :label="$t('字段名称')" prop="name"></bk-table-column>
                        <bk-table-column :label="$t('所属模块')">
                            <template slot-scope="{ row }">
                                <div class="module-item"
                                    v-for="(module, index) in row.conflictList"
                                    :key="index"
                                    :title="module.path || '--'">
                                    {{module.path || '--'}}
                                </div>
                            </template>
                        </bk-table-column>
                        <bk-table-column :label="$t('当前值')">
                            <template slot-scope="{ row }">
                                <div class="module-item"
                                    v-for="(module, index) in row.conflictList"
                                    :key="index"
                                    :title="module.value | formatter(row.property.bk_property_type, row.property.option)">
                                    {{module.value | formatter(row.property.bk_property_type, row.property.option)}}
                                </div>
                            </template>
                        </bk-table-column>
                    </bk-table>
                </div>
            </div>
        </div>
        <div class="footer-btns">
            <bk-button class="mr5" theme="primary" @click="handleDeleteFields">{{$t('确认删除')}}</bk-button>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="conflict.show"
            :width="795"
            :title="$t('处理冲突', { host: conflict.host })">
            <resolve-conflict slot="content"></resolve-conflict>
        </bk-sideslider>
    </div>
</template>

<script>
    import RESIZE_EVENTS from '@/utils/resize-events'
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
                },
                collapseHeight: null,
                collapse: false,
                showCollapse: false,
                list: [],
                checked: [],
                properties: [],
                modules: [{
                    bk_module_id: 2,
                    bk_module_name: 'Gamesever1'
                }, {
                    bk_module_id: 5,
                    bk_module_name: 'Gamesever21111'
                }, {
                    bk_module_id: 3,
                    bk_module_name: 'Gamesever32222'
                }, {
                    bk_module_id: 1,
                    bk_module_name: 'Gamesever433333'
                }, {
                    bk_module_id: 4,
                    bk_module_name: 'Gamesever2444444'
                }, {
                    bk_module_id: 6,
                    bk_module_name: 'Gamesever3555555'
                }, {
                    bk_module_id: 7,
                    bk_module_name: 'Gamesever4'
                }, {
                    bk_module_id: 8,
                    bk_module_name: 'Gamesever2'
                }, {
                    bk_module_id: 9,
                    bk_module_name: 'Gamesever3'
                }, {
                    bk_module_id: 11,
                    bk_module_name: 'Gamesever4'
                }, {
                    bk_module_id: 12,
                    bk_module_name: 'Gamesever2'
                }, {
                    bk_module_id: 13,
                    bk_module_name: 'Gamesever3'
                }, {
                    bk_module_id: 14,
                    bk_module_name: 'Gamesever4'
                }, {
                    bk_module_id: 15,
                    bk_module_name: 'Gamesever3'
                }, {
                    bk_module_id: 16,
                    bk_module_name: 'Gamesever4'
                }, {
                    bk_module_id: 17,
                    bk_module_name: 'Gamesever2'
                }, {
                    bk_module_id: 18,
                    bk_module_name: 'Gamesever3'
                }, {
                    bk_module_id: 19,
                    bk_module_name: 'Gamesever4'
                }, {
                    bk_module_id: 20,
                    bk_module_name: 'Gamesever3'
                }, {
                    bk_module_id: 21,
                    bk_module_name: 'Gamesever4'
                }]
            }
        },
        watch: {
            modules () {
                this.$nextTick(this.listHeight)
            }
        },
        async created () {
            this.setBreadcrumbs()
            const [list] = await Promise.all([
                this.getApplicationData(),
                this.getProperties()
            ])
            this.getPropertyInfo(list)
        },
        mounted () {
            RESIZE_EVENTS.addResizeListener(this.$refs.applicationFields, this.listHeight)
        },
        beforeDestroy () {
            RESIZE_EVENTS.removeResizeListener(this.$refs.applicationFields, this.listHeight)
        },
        methods: {
            setBreadcrumbs () {
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('主机属性自动应用'),
                    route: {
                        name: ''
                    }
                }, {
                    label: this.$t('批量删除')
                }])
            },
            listHeight () {
                const moduleItem = this.$refs.moduleItem
                const moduleList = this.$refs.moduleList
                if (!moduleItem || !moduleList) return 'auto'
                let totalWidth = null
                const rowWidth = moduleList.clientWidth
                moduleItem.forEach(item => {
                    totalWidth += (item.offsetWidth + 10)
                })
                const isOverflow = totalWidth > (rowWidth * 2)
                const listHeight = isOverflow ? 72 + 'px' : 'auto'
                this.showCollapse = isOverflow
                return this.collapse && isOverflow ? moduleList.scrollHeight + 'px' : listHeight
            },
            getPropertyInfo (list) {
                list.forEach((item, index) => {
                    const property = this.properties.find(property => property.bk_property_id === item.bk_property_id)
                    item.property = property
                })
                this.list = list
            },
            async getProperties () {
                try {
                    const properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                        params: this.$injectMetadata({
                            bk_obj_id: 'host'
                        })
                    })
                    this.properties = properties
                } catch (e) {
                    console.error(e)
                }
            },
            getApplicationData () {
                return [{
                    name: '主机名称',
                    bk_property_id: 'bk_host_name',
                    conflictList: [{
                        path: '当前主机值',
                        value: 'Quelly'
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: 'Selina,Quelly'
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块dddddddd',
                        value: 'Selina,Quelly'
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块ccccc',
                        value: 'Selina,Quelly'
                    }]
                }, {
                    name: '所属运营商',
                    bk_property_id: 'bk_isp_name',
                    conflictList: [{
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: '0'
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: '1'
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块dddddddd',
                        value: '2'
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块ccccc',
                        value: '3'
                    }]
                }, {
                    name: 'CPU逻辑核心数',
                    bk_property_id: 'bk_cpu',
                    conflictList: [{
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: 1
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: 2
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块dddddddd',
                        value: 3
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块ccccc',
                        value: 4
                    }]
                }, {
                    name: '备份维护人',
                    bk_property_id: 'operator',
                    conflictList: [{
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: 'admin'
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块aaaa',
                        value: 'peibiaoyang'
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块dddddddd',
                        value: 'yunyaoyang'
                    }, {
                        path: '广东区 / 大区一 / Apache业务模块ccccc',
                        value: 'v_dgzheng'
                    }]
                }]
            },
            handleSelectionChange (selection) {
                this.checked = selection.map(item => item.bk_property_id)
            },
            handleRemoveModule (id) {
                const moduleIndex = this.modules.findIndex(module => module.bk_module_id === id)
                moduleIndex > -1 && this.modules.splice(moduleIndex, 1)
            },
            handleDeleteFields () {
                return new Promise((resolve, reject) => {
                    this.$bkInfo({
                        title: this.$t('确认删除自动应用字段'),
                        subTitle: this.$t('自动应用字段删除提示'),
                        extCls: 'bk-dialog-sub-header-center',
                        confirmFn: () => { }
                    })
                })
            },
            handleCancel () {}
        }
    }
</script>

<style lang="scss" scoped>
    @mixin moduleIcon ($size) {
        width: $size;
        height: $size;
        line-height: $size;
        text-align: center;
        background-color: #C4C6CC;
        border-radius: 50%;
    }
    .edit-application-fields {
        padding: 0 20px;
        .info {
            display: flex;
            font-size: 14px;
            color: #63656E;
            margin-bottom: 10px;
        }
        .info-label {
            width: 150px;
            margin-right: 20px;
            text-align: right;
            font-weight: bold;
        }
        .main {
            position: relative;
            flex: 1;
            width: 0;
        }
        .module-list {
            display: flex;
            flex-wrap: wrap;
            padding-right: 50px;
            will-change: height;
            overflow: hidden;
            transition: all .2s;
            li {
                display: flex;
                align-items: center;
                height: 26px;
                line-height: 26px;
                font-size: 12px;
                border: 1px solid #C4C6CC;
                border-radius: 60px;
                margin: 0 10px 10px 0;
                &:hover {
                    border-color: #3A84FF;
                    .module-name {
                        color: #3A84FF;
                    }
                    .module-icon {
                        background-color: #3A84FF;
                    }
                    .module-remove {
                        background-color: #3A84FF;
                    }
                }
            }
            .module-icon {
                @include moduleIcon(20px);
                color: #FFFFFF;
                margin-left: 5px;
            }
            .module-name {
                color: #63656E;
                padding: 0 6px;
            }
            .module-remove {
                @include moduleIcon(16px);
                color: #FFFFFF;
                line-height: 12px;
                margin-right: 5px;
                cursor: pointer;
                &::before {
                    @include inlineBlock;
                    font-size: 20px;
                    transform: translateX(-2px) scale(.5);
                }
            }
        }
        .collapse {
            position: absolute;
            right: 4px;
            bottom: 12px;
            font-size: 14px;
            color: #3A84FF;
            cursor: pointer;
            .icon-angle-down {
                font-size: 12px;
                font-weight: bold;
                transition: all 0.2s;
                &.is-collapse {
                    transform: rotateZ(180deg);
                }
            }
        }
        .field-table {
            .module-item {
                height: 22px;
                line-height: 22px;
            }
        }
        .footer-btns {
            position: sticky;
            bottom: 0;
            left: 0;
            padding: 10px 0 10px 170px;
            background-color: #FAFBFD;
            z-index: 5;
        }
    }
</style>
