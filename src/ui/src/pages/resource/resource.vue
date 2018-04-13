/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
   <div class="host-resource-wrapper">
        <v-index 
            ref="index"
            :isShowBiz="false" 
            :isShowCollect="false" 
            :isShowHistory="false" 
            :outerParams="index.searchParams"
            @choose="setSelectedHost" 
            @attrLoaded="search">
            <div class="button-contain clearfix" slot="btnGroup">
                <bk-select class="biz-selector fl mr10" placeholder='分配到业务空闲机池'
                    :disabled="!hasSelectedHost"
                    :selected.sync="index.bkBizId" 
                    :filterable="true"
                    @on-selected="confirmTransfer">
                    <bk-select-option v-for="(bkBiz, index) in bkBizList"
                        :key="index"
                        :value="bkBiz['bk_biz_id']"
                        :label="bkBiz['bk_biz_name']">
                    </bk-select-option>
                </bk-select>
                <form ref="exportForm" :action="exportUrl" method="POST" class="fl mr10">
                    <input type="hidden" name="bk_host_id" :value="index.selectedHost">
                    <input type="hidden" name="bk_biz_id" value="-1">
                    <button class="bk-button"
                        :disabled="!hasSelectedHost"
                        @click.prevent="exportChoose">
                        <i class="icon-cc-derivation"></i>
                        <span>导出选中</span>
                    </button>
                </form>
                <button class="bk-button del-button fl mr10" @click="showFiling">
                    <i class="icon-cc-history"></i>
                </button>
                <button class="bk-button del-button fl mr10" :disabled="!hasSelectedHost" @click="confirmDel">
                    <i class="icon-cc-del"></i>
                </button>
                <bk-button type="primary" class="fl" @click="importHostShow">导入主机</bk-button>
            </div>
        </v-index>
        <v-sideslider 
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <bk-tab :active-name="slider.tab.active" @tab-changed="tabChanged" slot="content" style="border: none;padding: 0 20px;">
                <bk-tabpanel name="import" title="批量导入">
                    <v-import 
                        :templateUrl="slider.import.templateUrl"
                        :importUrl="slider.import.importUrl"
                        @success="search()"
                        @partialSuccess="search()">
                        <span slot="download-desc" style="display: inline-block;vertical-align: top;">说明：内网IP为必填列；</span>
                    </v-import>
                </bk-tabpanel>
                <bk-tabpanel name="agent" title="自动导入">
                    <div class="automatic-import">
                        <p>说明： 当主机安装了蓝鲸Agent后，会自动录入到“主机资源池”中</p>
                        <div class="back-contain">
                            <img src="../../common/images/icon/icon-back.png">
                            <a href="javascript:void(0)" @click="openAgentApp">点此进入Agent安装APP</a>
                        </div>
                    </div>
                </bk-tabpanel>
            </bk-tab>
        </v-sideslider>
        <v-filing
            :isShow.sync="filing.isShow"
            :objId="'host'"
            :objTableHeader="index.table.header"
            :allAttr="index.table.allAttr"
        ></v-filing>
   </div>
</template>

<script type="text/javascript">
    import vIndex from '@/pages/index/index'
    import vImport from '@/components/import/import'
    import vSideslider from '@/components/slider/sideslider'
    import vFiling from '@/components/filing/filing'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        data () {
            return {
                filing: {
                    isShow: false
                },
                index: {
                    bkBizId: '',
                    selectedHost: [],
                    searchParams: {
                        'bk_biz_id': -1,
                        condition: [{
                            'bk_obj_id': 'biz',
                            fields: [],
                            condition: [{
                                field: 'default',
                                operator: '$eq',
                                value: 1
                            }]
                        }]
                    },
                    table: {
                        header: [],
                        allAttr: []
                    }
                },
                slider: {
                    title: {
                        text: '导入主机',
                        icon: 'icon-cc-import'
                    },
                    isShow: false,
                    tab: {
                        active: 'import'
                    },
                    import: {
                        templateUrl: `${window.siteUrl}importtemplate/host`,
                        importUrl: `${window.siteUrl}hosts/import`
                    }
                },
                isShowimportHost: false
            }
        },
        computed: {
            ...mapGetters(['bkBizList', 'bkSupplierAccount']),
            hasSelectedHost () {
                return this.index.selectedHost.length
            },
            exportUrl () {
                return `${window.siteUrl}hosts/export`
            }
        },
        watch: {
            'slider.isShow' (isShow) {
                if (!isShow) {
                    this.slider.tab.active = 'import'
                }
            }
        },
        methods: {
            ...mapActions(['getBkBizList']),
            tabChanged (active) {
                this.slider.tab.active = active
            },
            confirmTransfer (selected, index) {
                let h = this.$createElement
                this.$bkInfo({
                    title: '请确认是否转移',
                    content: h('p', [
                        h('span', '选中的 '),
                        h('span', {
                            style: {
                                color: '#498fe0'
                            }
                        }, this.index.selectedHost.length),
                        h('span', ' 个主机转移到 '),
                        h('span', {
                            style: {
                                color: '#498fe0'
                            }
                        }, selected.label),
                        h('span', ' 下的空闲机模块')
                    ]),
                    confirmFn: () => {
                        this.transferHost()
                    }
                })
            },
            transferHost () {
                this.$axios.post('hosts/modules/resource/idle', {
                    'bk_biz_id': this.index.bkBizId,
                    'bk_host_id': this.index.selectedHost
                }).then(res => {
                    if (res.result) {
                        this.$alertMsg('分配成功', 'success')
                        this.index.selectedHost = []
                        this.$refs.index.transferSuccess()
                        this.index.bkBizId = ''
                        this.search()
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            confirmDel () {
                this.$bkInfo({
                    title: '确定删除选中的主机？',
                    confirmFn: () => {
                        this.$axios.delete('hosts/batch', {
                            data: JSON.stringify({
                                'bk_host_id': this.index.selectedHost.join(','),
                                'bk_supplier_account': this.bkSupplierAccount
                            })
                        }).then(res => {
                            if (res.result) {
                                this.$bkInfo({
                                    statusOpts: {
                                        title: '成功删除选中的主机',
                                        subtitle: false
                                    },
                                    type: 'success'
                                })
                                this.search()
                                this.$refs.index.clearChooseId()
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        })
                    }
                })
            },
            search () {
                // 构造新对象，触发v-index中watch outerParams
                this.index.searchParams = Object.assign({}, this.index.searchParams)
            },
            setSelectedHost (chooseId) {
                this.index.selectedHost = chooseId
            },
            exportChoose () {
                this.$refs.exportForm.submit()
            },
            importHostShow () {
                this.slider.isShow = true
            },
            openAgentApp () {
                let agentAppUrl = window.agentAppUrl
                if (agentAppUrl) {
                    if (window.agentAppUrl.indexOf('paasee-g.o.qcloud.com') !== -1) {
                        window.top.postMessage(JSON.stringify({action: 'open_other_app', app_code: 'bk_agent_setup'}), '*')
                    } else {
                        window.open(agentAppUrl)
                    }
                } else {
                    this.$alertMsg('未配置Agent安装APP地址')
                }
            },
            showFiling () {
                this.index.table.header = this.$refs.index.table.tableHeader
                this.index.table.allAttr = this.$refs.index.attribute
                this.filing.isShow = true
            }
        },
        created () {
            if (!this.bkBizList.length) {
                this.getBkBizList()
            }
        },
        components: {
            vImport,
            vIndex,
            vSideslider,
            vFiling
        }
    }
</script>
<style lang="scss" scoped>
    .host-resource-wrapper{
        position: relative;
        height: 100%;
    }
    .biz-selector {
        display: inline-block;
        vertical-align: middle;
        width: 200px;
    }
    .del-button{
        width: 36px;
        padding: 0;
        &:hover{
                .icon-cc-del{
                color: #ef4c4c;
            }
        }
    }
    .automatic-import{
        padding:40px 30px 0 30px;
        >p{
            // color:#6b7baa;
        }
        .back-contain{
            margin-left:45px;
            cursor:pointer;
            img{
                margin-right: 5px;
            }
            a{
                color:#498fe0;
            }
        }
    }
</style>
