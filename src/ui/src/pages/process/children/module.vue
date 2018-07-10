/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template lang="html">
    <v-table class="module-table"
        :header="table.header"
        :list="table.list"
        :loading="table.isLoading"
        :width="754"
        :wrapperMinusHeight="150"
        :sortable="false">
        <template slot="is_bind" slot-scope="{ item }">
            <bk-button :type="item['is_bind'] ? 'primary' : 'default'" :loading="$loading(`${item['bk_module_name']}Bind`)" @click="changeBinding(item)" 
                :class="item['is_bind'] ? 'main-btn' : 'vice-btn'">
                {{item['is_bind'] ? $t("ProcessManagement['已绑定']") : $t("ProcessManagement['未绑定']")}}
            </bk-button>
        </template>
    </v-table>
</template>

<script>
    import {mapGetters} from 'vuex'
    import vTable from '@/components/table/table'
    export default {
        props: {
            bkProcessId: {
                required: true
            },
            bkBizId: {
                required: true
            }
        },
        data () {
            return {
                table: {
                    header: [{
                        id: 'bk_module_name',
                        name: this.$t("ProcessManagement['模块名']")
                    }, {
                        id: 'set_num',
                        name: this.$t("ProcessManagement['所属集群数']")
                    }, {
                        id: 'is_bind',
                        name: this.$t("ProcessManagement['状态']")
                    }],
                    list: [],
                    isLoading: false,
                    maxHeight: 0
                }
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount'])
        },
        watch: {
            bkProcessId (bkProcessId) {
                if (bkProcessId) {
                    this.getModuleList()
                }
            }
        },
        methods: {
            changeBinding (item) {
                let moduleName = item['bk_module_name'].replace(' ', '')
                if (item['is_bind'] === 0) {
                    this.$axios.put(`proc/module/${this.bkSupplierAccount}/${this.bkBizId}/${this.bkProcessId}/${moduleName}`, {}, {id: `${item['bk_module_name']}Bind`}).then((res) => {
                        if (res.result) {
                            this.$alertMsg(this.$t("ProcessManagement['绑定进程到该模块成功']"), 'success')
                            item['is_bind'] = 1
                        } else {
                            this.$alertMsg(this.$t("ProcessManagement['绑定进程到该模块失败']"))
                        }
                    })
                } else {
                    this.$axios.delete(`proc/module/${this.bkSupplierAccount}/${this.bkBizId}/${this.bkProcessId}/${moduleName}`, {id: `${item['bk_module_name']}Bind`}).then(res => {
                        if (res.result) {
                            this.$alertMsg(this.$t("ProcessManagement['解绑进程模块成功']"), 'success')
                            item['is_bind'] = 0
                        } else {
                            this.$alertMsg(this.$t("ProcessManagement['解绑进程模块失败']"))
                        }
                    })
                }
            },
            getModuleList () {
                this.isLoading = true
                this.$axios.get(`proc/module/${this.bkSupplierAccount}/${this.bkBizId}/${this.bkProcessId}`).then((res) => {
                    if (res.result) {
                        this.table.list = this.sortModule(res.data)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.calcMaxHeight()
                    this.isLoading = false
                }).catch(() => {
                    this.isLoading = false
                })
            },
            sortModule (data) {
                let bindedModule = []
                let unbindModule = []
                data.forEach(module => {
                    module['is_bind'] ? bindedModule.push(module) : unbindModule.push(module)
                })
                bindedModule.sort((moduleA, moduleB) => {
                    return moduleA['bk_module_name'].localeCompare(moduleB['bk_module_name'])
                })
                unbindModule.sort((moduleA, moduleB) => {
                    return moduleA['bk_module_name'].localeCompare(moduleB['bk_module_name'])
                })
                return [...bindedModule, ...unbindModule]
            },
            calcMaxHeight () {
                this.table.maxHeight = document.body.getBoundingClientRect().height - 160
            }
        },
        components: {
            vTable
        }
    }
</script>
<style lang="scss" scoped>
    .module-table{
        height: calc(100% - 20px);
        margin: 20px 0 0 0;
        .btn{
            width: 52px;
            height: 25px;
        }
        .bk-button {
            padding: 1px 7px 2px;
            height: 22px;
            line-height: 20px;
            font-size: 12px;
        }
    }
</style>