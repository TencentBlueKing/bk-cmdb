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
    <div class="history-wrapper">
        <ul class="history-list">
            <li class="history-item clearfix" v-for="(history, index) in historyList" :key="index" @click="apply(history)">
                <span class="history-item-params fl" :title="history.paramsStr">{{history.paramsStr}}</span>
                <span class="history-item-time fr">{{history.time}}</span>
            </li>
        </ul>
    </div>
</template>
<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            isShow: {
                type: Boolean,
                default: false
            },
            start: {
                type: Number,
                default: 0
            },
            limit: {
                type: Number,
                default: 20
            }
        },
        computed: {
            ...mapGetters(['bkBizList'])
        },
        data () {
            return {
                historyList: [],
                operatorMap: {
                    '$eq': '=',
                    '$ne': '!=',
                    '$in': '~',
                    '$regex': '~'
                }
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.getHistory()
                }
            }
        },
        methods: {
            getHistory () {
                this.$axios.get(`hosts/history/${this.start}/${this.limit}`).then((res) => {
                    if (res.result) {
                        if (res.data.count) {
                            let historyList = res.data.info.map(history => {
                                let content = JSON.parse(history['content'])
                                let time = this.$formatTime(history['bk_create_time'], 'YYYY-MM-DD HH:mm')
                                let paramsStr = []
                                let queryParams = []
                                let filters = {
                                    'biz': [],
                                    'host': [],
                                    'module': [],
                                    'set': []
                                }
                                let info = {
                                    'bk_biz_id': content['bk_biz_id'],
                                    'exact_search': content.ip.exact,
                                    'inner_ip': content.ip.flag.indexOf('bk_host_innerip') !== -1,
                                    'outer_ip': content.ip.flag.indexOf('bk_host_outerip') !== -1,
                                    'ip_list': content.ip.data
                                }
                                paramsStr.push(`bk_biz_id:${content['bk_biz_id']}`)
                                if (content.ip.data.length) {
                                    paramsStr.push(`IP:${content.ip.data.join(',')}`)
                                }
                                content.condition.map((condition, index) => {
                                    condition.condition.map(filter => {
                                        if (filter.field !== 'default' && filters.hasOwnProperty(condition['bk_obj_id'])) {
                                            filters[condition['bk_obj_id']].push(`${filter.field}${this.operatorMap[filter.operator]}${filter.value}`)
                                            queryParams.push({
                                                'bk_obj_id': condition['bk_obj_id'],
                                                'field': filter.field,
                                                'operator': filter.operator,
                                                'value': filter.value
                                            })
                                        }
                                    })
                                })
                                Object.keys(filters).map(filterType => {
                                    if (filters[filterType].length) {
                                        paramsStr.push(`${filterType}:${filters[filterType].join('; ')}`)
                                    }
                                })
                                return {
                                    time: time,
                                    paramsStr: paramsStr.join(' | '),
                                    info: info,
                                    queryParams: queryParams
                                }
                            })
                            this.historyList = historyList.reverse()
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            apply ({info, queryParams}) {
                let isAppExist = false
                this.bkBizList.map(({bk_biz_id: bkBizId}) => {
                    if (bkBizId === info['bk_biz_id']) {
                        isAppExist = true
                    }
                })
                if (isAppExist) {
                    // 构造的数据结构跟收藏列表类似
                    this.$emit('apply', {info: JSON.stringify(info), query_params: JSON.stringify(queryParams)})
                } else {
                    this.$alertMsg(this.$t('Common[\'该查询条件对应的业务不存在\']'))
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .history-wrapper{
        padding: 20px 0 0 0;
        font-size: 12px;
        overflow: auto;
        height: 100%;
        @include scrollbar;
    }
    .history-list{
        line-height: 30px;
        .history-item{
            cursor: pointer;
            padding: 0 10px;
            &:hover{
                background-color: #f1f7ff;
            }
            .history-item-params{
                width: 180px;
                @include ellipsis;
            }
        }
    }
</style>