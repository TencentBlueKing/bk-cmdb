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
    <div class="filter-wrapper" :class="{'close': !isOpen}">
        <a href="javascript:void(0)" class="bk-icon icon-angle-right filter-toggle" 
            :class="{'filter-toggle-close': !isOpen}" 
            @click="filterToggle">
        </a>
        <bk-tab class="filter-tab" :active-name="tab.active" @tab-changed="tabChanged" v-show="isOpen">
            <bk-tabpanel name="screening" :title="$t('HostResourcePool[\'筛选\']')" ref="screeningTabpanel">
                <v-screening ref="screening" 
                    :queryColumns="queryColumns"
                    :queryColumnData="queryColumnData"
                    :attribute="attribute"
                    :isShowBiz="isShowBiz"
                    :isShowScope="isShowScope"
                    @bkBizSelected="bkBizSelected"
                    @refresh="refresh"
                    @filterChange="filterChange">
                </v-screening>
            </bk-tabpanel>
            <bk-tabpanel name="collect" :title="$t('Hosts[\'收藏\']')" :show="isShowCollect">
                <v-collect :favoriteList="tab.favorite.list" :active="tab.active === 'collect'" @delete="getFavoriteList" @update="updateFavoriteList" @apply="applyCollect"></v-collect>
            </bk-tabpanel>
            <bk-tabpanel name="history" :title="$t('Hosts[\'历史\']')" :show="false">
                <v-history :isShow="tab.active === 'history'" @apply="applyHistory"></v-history>
            </bk-tabpanel>
            <template slot="setting">
                <div class="filter-operate" v-show="tab.active === 'screening'">
                    <i class="icon-cc icon-cc-broom" @click="emptiedField" v-tooltip="$t('HostResourcePool[\'清空查询条件\']')"></i>
                    <i class="icon-cc icon-cc-collection" :class="{'collecting': tab.screening.isCollecting}" @click.stop="showCollectBox" v-tooltip="$t('Hosts[\'收藏\']')" v-if="isShowCollect"></i>
                    <i class="icon-cc icon-cc-funnel" @click="showField" v-tooltip="$t('HostResourcePool[\'设置筛选项\']')"></i>
                </div>
                <div class="collect-box" v-if="isShowCollect" v-show="tab.screening.isCollecting" @click.stop v-click-outside="hideCollectBox">
                    <div class="title tl">{{$t('Hosts[\'收藏此查询\']')}}</div>
                    <form id="validate-form">
                        <div class="input-box tl">
                            <input type="text" :placeholder="$t('Hosts[\'请填写名称\']')" v-model.trim="tab.screening.collectName"
                            :data-vv-name="$t('Hosts[\'名称\']')"
                            v-validate="'required|name'">
                            <span v-show="errors.has($t('Hosts[\'名称\']'))" class="help is-danger">{{ errors.first($t('Hosts[\'名称\']')) }}</span>
                        </div>
                    </form>
                    <div class="collect-list">
                        <ul>
                            <li>
                                <span :title="favoriteList" style="word-break: break-all;">{{favoriteList}}</span>
                            </li>
                        </ul>
                    </div>
                    <div class="footer">
                        <label class="bk-form-checkbox bk-checkbox-small p0 fl mt10" hidden>
                            <input type="checkbox" name="checkbox" 
                                v-model="tab.screening.default" 
                                :true-value="1" 
                                :false-value="2">
                            <span class="acquiescence vm pl5">{{$t('Hosts[\'默认\']')}}</span>
                        </label>
                        <div class="btn-wrapper">
                            <bk-button type="primary" class="mr10 main-btn"
                                :loading="$loading('collect')"
                                :disabled="!tab.screening.collectName"
                                @click="makeSureCollect">
                                {{$t('Hosts[\'确认\']')}}
                            </bk-button>
                            <bk-button type="default" class="cancel-btn vice-btn" @click="hideCollectBox">
                                {{$t('Common[\'取消\']')}}
                            </bk-button>
                        </div>
                    </div>
                </div>
            </template>
        </bk-tab>
    </div>
</template>
<script>
    import vScreening from '@/components/filter/screening'
    import vCollect from '@/components/filter/collect'
    import vHistory from '@/components/filter/history'
    export default {
        props: {
            queryColumns: {
                type: Array,
                required: true
            },
            queryColumnData: {
                type: Object,
                default () {
                    return {}
                }
            },
            attribute: {
                type: Array,
                required: true
            },
            isShowBiz: {
                type: Boolean,
                default: true
            },
            isShowCollect: {
                type: Boolean,
                default: true
            },
            isShowHistory: {
                type: Boolean,
                default: true
            },
            isShowScope: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                isOpen: true,
                tab: {
                    active: 'screening',
                    screening: {
                        isCollecting: false,
                        filter: {},
                        collectName: '',
                        default: 2
                    },
                    favorite: {
                        list: [],
                        filter: '',
                        loaded: false
                    }
                },
                operatorMap: {
                    '$ne': '!=',
                    '$eq': '=',
                    '$regex': '~',
                    '$in': '~'
                }
            }
        },
        watch: {
            'tab.active' (activeName) {
                if (activeName === 'collect') {
                    this.getFavoriteList()
                }
            }
        },
        computed: {
            favoriteList () {
                let favoriteList = []
                let filter = this.tab.screening.filter
                if (filter['bk_biz_id']) {
                    favoriteList.push(`bk_biz_id:${filter['bk_biz_id']}`)
                }
                if (filter.ip && filter.ip.data.length) {
                    favoriteList.push(filter.ip.data.join(','))
                }
                if (filter.condition) {
                    filter.condition.map(({condition, bk_obj_id: bkObjId}) => {
                        if (bkObjId !== 'app' && condition.length) {
                            let objectFilter = []
                            condition.map(({field, operator, value}) => {
                                objectFilter.push(`${field}${this.operatorMap[operator]}${Array.isArray(value) ? value.join(',') : value}`)
                            })
                            favoriteList.push(`${bkObjId}: ${objectFilter.join('|')}`)
                        }
                    })
                }
                return favoriteList.join('|')
            }
        },
        methods: {
            bkBizSelected (bkBizId) {
                this.$emit('bkBizSelected', bkBizId)
            },
            filterChange (filter) {
                this.tab.screening.filter = filter
                this.$emit('filterChange', filter)
            },
            showCollectBox () {
                this.tab.screening.isCollecting = true
            },
            hideCollectBox () {
                this.tab.screening.isCollecting = false
                this.tab.screening.collectName = ''
            },
            makeSureCollect () {
                this.$validator.validateAll().then(res => {
                    if (res) {
                        this.$axios.post('hosts/favorites', this.getCollectParams(), {id: 'collect'}).then(res => {
                            if (res.result) {
                                this.$alertMsg(this.$t('Common[\'收藏成功\']'), 'success')
                                this.hideCollectBox()
                                this.updateFavoriteCount(res.data.id)
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        })
                    }
                })
            },
            getFavoriteList () {
                this.$axios.post('hosts/favorites/search', {}).then(res => {
                    if (res.result) {
                        let list = res.data.info
                        this.tab.favorite.list = list
                        this.tab.favorite.loaded = true
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            updateFavoriteCount (favoriteId) {
                this.$axios.put(`hosts/favorites/${favoriteId}/incr`).then(res => {
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            updateFavoriteList (updateItem) {
                let originItem = this.tab.favorite.list.find(({id}) => id === updateItem.id)
                if (originItem) {
                    const originItemIndex = this.tab.favorite.list.indexOf(originItem)
                    this.tab.favorite.list.splice(originItemIndex, 1, Object.assign({}, originItem, updateItem))
                }
            },
            getCollectParams () {
                let filter = this.tab.screening.filter
                let info = {
                    bk_biz_id: filter['bk_biz_id'],
                    exact_search: filter.ip.exact,
                    'bk_host_innerip': filter.ip.flag.split('|').indexOf('bk_host_innerip') !== -1,
                    'bk_host_outerip': filter.ip.flag.split('|').indexOf('bk_host_outerip') !== -1,
                    ip_list: filter.ip.data
                }
                let queryParams = []
                filter.condition.map(({condition, bk_obj_id: bkObjId}) => {
                    condition.map(({field, operator, value}) => {
                        queryParams.push({
                            bk_obj_id: bkObjId,
                            field,
                            operator,
                            value
                        })
                    })
                })
                return {
                    info: JSON.stringify(info),
                    query_params: JSON.stringify(queryParams),
                    is_default: this.tab.screening.default,
                    name: this.tab.screening.collectName
                }
            },
            applyCollect (collect) {
                this.tab.active = 'screening'
                this.updateFavoriteCount(collect['id'])
                this.$emit('applyCollect', collect)
            },
            applyHistory (history) {
                this.tab.active = 'screening'
                this.$emit('applyHistory', history)
            },
            emptiedField () {
                this.$refs.screening.resetQueryColumnData()
                this.$emit('emptyField')
            },
            showField () {
                this.$emit('showField')
            },
            tabChanged (tabName) {
                this.tab.active = tabName
            },
            filterToggle () {
                this.isOpen = !this.isOpen
                this.$emit('filterToggle', this.isOpen)
            },
            refresh () {
                this.$emit('refresh')
            }
        },
        components: {
            vScreening,
            vCollect,
            vHistory
        }
    }
</script>

<style lang="scss" scoped>
.filter-wrapper{
    position: relative;
    height: 100%;
    border-left: 1px solid #e7e9ef;
    overflow: visible;
    &.close{
        border: none;
    }
    .filter-toggle{
        position: absolute;
        right: 100%;
        top: 50%;
        margin-top: -50px;
        width: 14px;
        height: 100px;
        line-height: 100px;
        color: #fff;
        font-size: 12px;
        text-align: center;
        border-top-left-radius: 12px;
        border-bottom-left-radius: 12px;
        background-color: #c3cdd7;
        &.filter-toggle-close{
            transform: rotate(180deg);
            border-top-right-radius: 12px;
            border-bottom-right-radius: 12px;
            border-top-left-radius: 0;
            border-bottom-left-radius: 0;
        }
        &:hover{
            background-color: #6b7baa;
        }
    }
}
.filter-tab{
    width: 357px;
    height: 100%;
    border: none;
    padding: 19px 20px 20px;
    .filter-operate{
        height: 36px;
        line-height: 36px;
        color: #c3cdd7;
        font-size: 0;
        margin-right: -20px;
        .icon-cc{
            font-size: 16px;
            margin: 0 0 0 10px;
            cursor: pointer;
            &:hover{
                color: #3c96ff;
            }
        }
        .icon-cc-collection.collecting{
            color: #ffb400;
        }
    }
}
.collect-box{
    position: absolute;
    z-index: 1;
    right: -19px;
    background: #fff;
    padding: 18px 20px 30px 20px;
    line-height: 1;
    width: 317px;
    /*height: 360px;*/
    box-shadow: 0px 2px 10px 4px rgba(12, 34, 59, .13);
    font-size: 14px;
    top: 50px;
    z-index: 100;
    &:before{
        content: '';
        right: 29px;
        bottom: 100%;
        width: 0;
        height: 0;
        border-left: 6px solid transparent;
        border-right: 6px solid transparent;
        border-bottom: 10px solid #e7e9ef;
        position: absolute;
        margin-bottom: 2px;
    }
    &:after{
        content: "";
        right: 29px;
        bottom: 100%;
        width: 0px;
        height: 0px;
        border-left: 6px solid transparent;
        border-right: 6px solid transparent;
        border-bottom: 10px solid #fff;
        position: absolute;
    }
    .no-data{
        text-align: center;
        height: 200px;
        line-height: 200px;
        .close-del{
            position: absolute;
            top: 5px;
            padding: 6px;
            right: 5px;
            cursor: pointer;
            color: #bec6de;
            &:hover{
                border-radius: 50%;
                background: #f3f3f3;
            }
        }
    }
    .collect-list{
        margin-top: 8px;
        width: 100%;
        min-height: 100px;
        overflow-y: auto;
        background:#f9f9f9;
        padding: 10px;
        &::-webkit-scrollbar{
            width: 6px;
            height: 5px;
        }
        &::-webkit-scrollbar-thumb{
            border-radius: 5%;
            background: #a5a5a5;
        }
        ul{
            padding: 0;
            margin: 0;
            li{
                line-height: 30px;
                text-align: left;
                width: 100%;
                word-wrap: break-word;
            }
        }
    }
    .acquiescence{
        color: #6b7baa;
    }
    .title{
        border-left: 2px solid #6b7baa;
        margin-bottom: 18px;
        padding-left: 5px;
    }
    .input-box{
        input{
            margin-bottom: 2px;
            width: 100%;
            border-radius: 2px;
            padding: 0 10px;
            height: 32px;
            outline: none;
            border: 1px solid #e7e9ef;
        }
    }
    .footer{
        margin-top: 20px;
        line-height: 34px;
        .btn-wrapper{
            font-size: 0;
            float: right;
        }
        .bk-form-checkbox{
            cursor: pointer;
        }
    }
}
</style>
<style lang="scss">
.filter-wrapper{
    .filter-tab{
        .bk-tab2-head{
            height: 37px;
        }
        .bk-tab2-content{
            height: calc(100% - 98px);
            section{
                height: 100%;
            }
        }
        .bk-tab2-nav{
            .tab2-nav-item{
                height: 36px;
                line-height: 36px;
                padding: 0 18px;
            }
        }
    }
}
</style>