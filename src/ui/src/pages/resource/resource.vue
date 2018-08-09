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
        <v-hosts 
            ref="hosts"
            :isShowBiz="false" 
            :isShowCollect="false" 
            :isShowHistory="false"
            :isShowTransfer="false"
            :isShowScope="true"
            :outerParams="hosts.searchParams"
            @choose="setSelectedHost" 
            @attrLoaded="search">
            <div class="button-contain clearfix" slot="btnGroup">
                <bk-select class="biz-selector fl mr10" :placeholder="$t('HostResourcePool[\'分配到业务空闲机池\']')"
                    :disabled="!hasSelectedHost"
                    :selected.sync="hosts.bkBizId" 
                    :filterable="true"
                    @on-selected="confirmTransfer">
                    <bk-select-option v-for="(bkBiz, index) in bkBizList"
                        :key="index"
                        :value="bkBiz['bk_biz_id']"
                        :label="bkBiz['bk_biz_name']">
                    </bk-select-option>
                </bk-select>
                <div class="fl bk-group bk-button-group">
                    <div class="btn-tooltip-wrapper" v-tooltip="$t('HostResourcePool[\'导出选中\']')">
                        <form ref="exportForm" :action="exportUrl" method="POST" class="fl">
                            <input type="hidden" name="bk_host_id" :value="hosts.selectedHost">
                            <input type="hidden" name="bk_biz_id" value="-1">
                            <bk-button type="default" class="first"
                                :disabled="!hasSelectedHost"
                                @click.prevent="exportChoose">
                                <i class="icon-cc-derivation"></i>
                            </bk-button>
                        </form>
                    </div>
                    <div class="btn-tooltip-wrapper" v-tooltip="$t('Common[\'删除\']')">
                        <bk-button type="default" class="delete-button fl" :class="{'disabled': !hasSelectedHost || $loading('deleteHosts')}" :disabled="!hasSelectedHost || $loading('deleteHosts')" @click="confirmDel">
                            <i class="icon-cc-del"></i>
                        </bk-button>
                    </div>
                </div>
                <div class="fr">
                    <bk-button type="primary" class="fl" @click="importHostShow">{{$t('HostResourcePool[\'导入主机\']')}}</bk-button>
                    <bk-button class="icon-btn icon-history fl ml10" @click="showFiling" v-tooltip="$t('Common[\'查看删除历史\']')">
                        <i class="icon-cc-history"></i>
                    </bk-button>
                </div>
            </div>
        </v-hosts>
        <v-sideslider 
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <bk-tab :active-name="slider.tab.active" @tab-changed="tabChanged" slot="content" style="border: none;padding: 0 20px;">
                <bk-tabpanel name="import" :title="$t('HostResourcePool[\'批量导入\']')">
                    <v-import v-if="slider.isShow"
                        :templateUrl="slider.import.templateUrl"
                        :importUrl="slider.import.importUrl"
                        @success="search()"
                        @partialSuccess="search()">
                        <span slot="download-desc" style="display: inline-block;vertical-align: top;">{{$t('HostResourcePool[\'说明：内网IP为必填列\']')}}</span>
                    </v-import>
                </bk-tabpanel>
                <bk-tabpanel name="agent" :title="$t('HostResourcePool[\'自动导入\']')">
                    <div class="automatic-import">
                        <p>{{$t("HostResourcePool['agent安装说明']")}}</p>
                        <div class="back-contain">
                            <i class="icon-cc-skip"></i>
                            <a href="javascript:void(0)" @click="openAgentApp">{{$t("HostResourcePool['点此进入节点管理']")}}</a>
                        </div>
                    </div>
                </bk-tabpanel>
            </bk-tab>
        </v-sideslider>
        <v-delete-history
            :isShow.sync="filing.isShow"
            :objId="'host'"
            :objTableHeader="hosts.table.header"
        ></v-delete-history>
   </div>
</template>

<script type="text/javascript">
    import vHosts from '@/pages/hosts/hosts'
    import vImport from '@/components/import/import'
    import vSideslider from '@/components/slider/sideslider'
    import vDeleteHistory from '@/components/history/delete'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        data () {
            return {
                isDropdownShow: false,
                selectedList: [],
                forSelectedList: [],
                filing: {
                    isShow: false
                },
                hosts: {
                    bkBizId: '',
                    selectedHost: [],
                    searchParams: {
                        'bk_biz_id': -1
                    },
                    table: {
                        header: [],
                        allAttr: []
                    }
                },
                slider: {
                    title: {
                        text: this.$t("HostResourcePool['导入主机']"),
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
                isShowimportHost: false,
                prevHistoryCount: 0
            }
        },
        computed: {
            ...mapGetters([
                'bkBizList',
                'bkSupplierAccount',
                'language'
            ]),
            hasSelectedHost () {
                return this.hosts.selectedHost.length
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
            },
            'filing.isShow' (isShow) {
                this.updateHistoryCount(isShow)
            }
        },
        methods: {
            ...mapActions(['getBkBizList']),
            tabChanged (active) {
                this.slider.tab.active = active
            },
            hasAssignedHosts () {
                return this.$refs.hosts.selectedList.find(host => !!host['biz'].find(biz => biz['default'] !== 1))
            },
            confirmTransfer (selected, index) {
                if (this.hasAssignedHosts()) {
                    this.$alertMsg(this.$t('Hosts["请勿选择已分配主机"]'))
                    this.$nextTick(() => {
                        this.hosts.bkBizId = ''
                    })
                    return
                }
                let h = this.$createElement
                let content = ''
                if (this.language === 'en') {
                    content = h('p', [
                        h('span', 'Selected '),
                        h('span', {
                            style: {
                                color: '#3c96ff'
                            }
                        }, this.hosts.selectedHost.length),
                        h('span', ' Hosts Transfer to Idle machine under '),
                        h('span', {
                            style: {
                                color: '#3c96ff'
                            }
                        }, selected.label)
                    ])
                } else {
                    content = h('p', [
                        h('span', '选中的 '),
                        h('span', {
                            style: {
                                color: '#3c96ff'
                            }
                        }, this.hosts.selectedHost.length),
                        h('span', ' 个主机转移到 '),
                        h('span', {
                            style: {
                                color: '#3c96ff'
                            }
                        }, selected.label),
                        h('span', ' 下的空闲机模块')
                    ])
                }
                this.$bkInfo({
                    title: this.$t("HostResourcePool['请确认是否转移']"),
                    content,
                    confirmFn: () => {
                        this.transferHost()
                    },
                    cancelFn: () => {
                        this.hosts.bkBizId = ''
                    }
                })
            },
            transferHost () {
                this.$axios.post('hosts/modules/resource/idle', {
                    'bk_biz_id': this.hosts.bkBizId,
                    'bk_host_id': this.hosts.selectedHost
                }).then(res => {
                    if (res.result) {
                        this.$alertMsg(this.$t("HostResourcePool['分配成功']"), 'success')
                        this.hosts.selectedHost = []
                        this.$refs.hosts.transferSuccess()
                        this.hosts.bkBizId = ''
                        this.search()
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            confirmDel () {
                this.$bkInfo({
                    title: `${this.$t("HostResourcePool['确定删除选中的主机']")}？`,
                    confirmFn: () => {
                        this.$axios.delete('hosts/batch', {
                            data: JSON.stringify({
                                'bk_host_id': this.hosts.selectedHost.join(','),
                                'bk_supplier_account': this.bkSupplierAccount
                            }),
                            id: 'deleteHosts'
                        }).then(res => {
                            if (res.result) {
                                this.$bkInfo({
                                    statusOpts: {
                                        title: this.$t("HostResourcePool['成功删除选中的主机']"),
                                        subtitle: false
                                    },
                                    type: 'success'
                                })
                                this.search()
                                this.$refs.hosts.clearChooseId()
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        })
                    }
                })
            },
            search () {
                // 构造新对象，触发v-index中watch outerParams
                this.hosts.searchParams = Object.assign({}, this.hosts.searchParams)
            },
            setSelectedHost (chooseId) {
                this.hosts.selectedHost = chooseId
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
                        window.top.postMessage(JSON.stringify({action: 'open_other_app', app_code: 'bk_nodeman'}), '*')
                    } else {
                        window.open(agentAppUrl)
                    }
                } else {
                    this.$alertMsg(this.$t("HostResourcePool['未配置Agent安装APP地址']"))
                }
            },
            showFiling () {
                this.hosts.table.header = this.$refs.hosts.table.tableHeader
                this.hosts.table.allAttr = this.$refs.hosts.attribute
                this.filing.isShow = true
            },
            updateHistoryCount (isShow) {
                if (this.prevHistoryCount) {
                    this.$store.commit('navigation/updateHistoryCount', isShow ? -1 : 1)
                }
            }
        },
        created () {
            this.prevHistoryCount = this.$store.state.navigation.historyCount
            if (!this.bkBizList.length) {
                this.getBkBizList()
            }
        },
        beforeRouteLeave (to, from, next) {
            this.$store.commit('resetHostSearch')
            next()
        },
        components: {
            vImport,
            vHosts,
            vSideslider,
            vDeleteHistory
        }
    }
</script>
<style lang="scss" scoped>
    .host-resource-wrapper{
        position: relative;
        height: 100%;
        .button-contain{
            display: inline-block;
            width: calc(100% - 195px);
            vertical-align: middle;
        }
    }
    .biz-selector {
        display: inline-block;
        vertical-align: middle;
        width: 200px;
    }
    .icon-cc-history{
        position: relative;
        top: -1px;
        font-size: 16px;
    }
    .icon-btn{
        width: 36px;
        padding: 0;
        &:not(:disabled):hover{
            border-color: #ef4c4c;
            &.icon-history {
                border-color: #c3cdd7;
            }
            .icon-cc-del{
                color: #ef4c4c;
            }
        }
        &.disabled{
            .icon-cc-del{
                color: #ccc;
            }
        }
    }
    .automatic-import{
        padding:40px 30px 0 30px;
        >p{
            // color:#6b7baa;
        }
        .back-contain{
            cursor:pointer;
            color: #3c96ff;
            img{
                margin-right: 5px;
            }
            a{
                color:#3c96ff;
            }
        }
    }
</style>
