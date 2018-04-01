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
        :tableHeader="table.header"
        :tableList="table.list"
        :isLoading="table.isLoading"
        :maxHeight="table.maxHeight"
        :sortable="false">
        <td slot="is_bind" slot-scope="{ item }">
            <button @click="changeBinding(item)" 
                :class="item['is_bind'] ? 'main-btn' : 'vice-btn'">
                {{item['is_bind'] ? '已绑定' : '未绑定'}}
            </button>
        </td>    
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
                        name: '模块名'
                    }, {
                        id: 'set_num',
                        name: '所属集群数'
                    }, {
                        id: 'is_bind',
                        name: '状态'
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
                this.getModuleList()
            }
        },
        methods: {
            changeBinding (item) {
                let moduleName = item['bk_module_name'].replace(' ', '')
                if (item['is_bind'] === 0) {
                    this.$axios.put(`proc/module/${this.bkSupplierAccount}/${this.bkBizId}/${this.bkProcessId}/${moduleName}`).then((res) => {
                        if (res.result) {
                            this.$alertMsg('绑定进程到该模块成功', 'success')
                            this.getModuleList()
                        } else {
                            this.$alertMsg('绑定进程到该模块失败')
                        }
                    })
                } else {
                    this.$axios.delete(`proc/module/${this.bkSupplierAccount}/${this.bkBizId}/${this.bkProcessId}/${moduleName}`).then(res => {
                        if (res.result) {
                            this.$alertMsg('解绑进程模块成功', 'success')
                            this.getModuleList()
                        } else {
                            this.$alertMsg('解绑进程模块失败')
                        }
                    })
                }
            },
            getModuleList () {
                this.isLoading = true
                this.$axios.get(`proc/module/${this.bkSupplierAccount}/${this.bkBizId}/${this.bkProcessId}`).then((res) => {
                    if (res.result) {
                        this.table.list = res.data
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.calcMaxHeight()
                    this.isLoading = false
                }).catch(() => {
                    this.isLoading = false
                })
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
        padding: 20px 0 0 0;
        .btn{
            width: 52px;
            height: 25px;
        }
    }
</style>