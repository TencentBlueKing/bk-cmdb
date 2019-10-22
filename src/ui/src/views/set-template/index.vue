<template>
    <div class="template-layout">
        <div class="options clearfix">
            <div class="fl">
                <span class="fl" v-cursor="{
                    active: !$isAuthorized($OPERATION.C_SET_TEMPLATE),
                    auth: [$OPERATION.C_SET_TEMPLATE]
                }">
                    <bk-button
                        theme="primary"
                        :disabled="!$isAuthorized($OPERATION.C_SET_TEMPLATE)"
                        @click="handleCreate"
                    >
                        {{$t('新建')}}
                    </bk-button>
                </span>
            </div>
            <div class="fr">
                <bk-input :placeholder="$t('模板名称搜索')"
                    right-icon="icon-search"
                    v-model="searchName"
                    @enter="handleFilterTemplate">
                </bk-input>
            </div>
        </div>
        <bk-table class="template-table"
            :data="list"
            :row-style="{ cursor: 'pointer' }"
            @row-click="handleRowClick">
            <bk-table-column :label="$t('模板名称')" prop="name" class-name="is-highlight"></bk-table-column>
            <bk-table-column :label="$t('应用数量')" prop="set_instance_count"></bk-table-column>
            <bk-table-column :label="$t('修改人')" prop="modifier"></bk-table-column>
            <bk-table-column :label="$t('修改时间')" prop="last_time">
                <template slot-scope="{ row }">
                    <span>{{$tools.formatTime(row.last_time, 'YYYY-MM-DD HH:mm')}}</span>
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('操作')" width="180">
                <template slot-scope="{ row }">
                    <span
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.U_SET_TEMPLATE),
                            auth: [$OPERATION.U_SET_TEMPLATE]
                        }">
                        <bk-button
                            text
                            :disabled="!$isAuthorized($OPERATION.U_SET_TEMPLATE)"
                            @click="handleEdit(row)"
                        >
                            {{$t('编辑')}}
                        </bk-button>
                    </span>
                    <span class="text-primary ml15"
                        style="color: #dcdee5 !important; cursor: not-allowed;"
                        v-if="row.set_instance_count && $isAuthorized($OPERATION.D_SET_TEMPLATE)"
                        v-bk-tooltips.top="$t('不可删除')">
                        {{$t('删除')}}
                    </span>
                    <span v-else
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.D_SET_TEMPLATE),
                            auth: [$OPERATION.D_SET_TEMPLATE]
                        }">
                        <bk-button text class="ml15"
                            :disabled="!$isAuthorized($OPERATION.D_SET_TEMPLATE)"
                            @click="handleDelete(row)"
                        >
                            {{$t('删除')}}
                        </bk-button>
                    </span>
                </template>
            </bk-table-column>
            <template slot="empty">
                <i class="bk-table-empty-icon bk-icon icon-empty"></i>
                <i18n path="空集群模板提示" tag="div">
                    <span
                        place="link"
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.C_SET_TEMPLATE),
                            auth: [$OPERATION.C_SET_TEMPLATE]
                        }">
                        <bk-button
                            text
                            :disabled="!$isAuthorized($OPERATION.C_SET_TEMPLATE)"
                            @click="handleCreate"
                        >
                            {{$t('立即创建')}}
                        </bk-button>
                    </span>
                </i18n>
            </template>
        </bk-table>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                list: [],
                originList: [],
                searchName: ''
            }
        },
        computed: {
            business () {
                return this.$store.state.objectBiz.bizId
            }
        },
        async created () {
            await this.getSetTemplates()
        },
        methods: {
            async getSetTemplates () {
                const data = await this.$store.dispatch('setTemplate/getSetTemplates', {
                    bizId: this.business,
                    params: {},
                    config: {
                        requestId: 'getSetTemplates'
                    }
                })
                const list = (data.info || []).map(item => ({
                    set_instance_count: item.set_instance_count,
                    ...item.set_template
                }))
                this.list = list
                this.originList = list
            },
            handleCreate () {
                this.$router.push({
                    name: 'setTemplateConfig',
                    params: {
                        mode: 'create'
                    }
                })
            },
            handleEdit (row) {
                this.$router.push({
                    name: 'setTemplateConfig',
                    params: {
                        mode: 'edit',
                        templateId: row.id
                    }
                })
            },
            async handleDelete (row) {
                this.$bkInfo({
                    title: this.$t('确认删除xx集群模板', { name: row.name }),
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('setTemplate/deleteSetTemplate', {
                                bizId: this.$store.getters['objectBiz/bizId'],
                                config: {
                                    data: {
                                        set_template_ids: [row.id]
                                    }
                                }
                            })
                            this.getSetTemplates()
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            handleFilterTemplate () {
                const originList = this.$tools.clone(this.originList)
                this.list = this.searchName
                    ? originList.filter(template => template.name.indexOf(this.searchName) !== -1)
                    : originList
            },
            handleSelectable (row) {
                return !row.set_instance_count
            },
            handleRowClick (row, event, column) {
                if (!column.property) {
                    return false
                }
                this.$router.push({
                    name: 'setTemplateConfig',
                    params: {
                        mode: 'view',
                        templateId: row.id
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .template-layout {
        padding: 0 20px;
    }
    .options {
        font-size: 0;
    }
    .template-table {
        margin-top: 16px;
    }
</style>
