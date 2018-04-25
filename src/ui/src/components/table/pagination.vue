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
    <div class="table-page-contain clearfix" :class="{'empty': !tableList.length}">
        <div v-show="tableList.length" class="clearfix page-info-box">
            <div class="page-info fl">
                <span class="mr20" v-if="hasCheckbox">{{$tc('Common["已选N行"]', chooseId.length, {N: chooseId.length})}}</span>
                <span>{{$tc('Common[\'页码\']', pagination.current, {current: pagination.current, total: totalPage})}}</span>
                <span class="ml20 mr20">
                    <i18n path="Common['每页显示']" tag="span">
                        <span class="select_page_setting mr5" place="page">
                            <bk-select class="select-box" ref="sizeSelector"
                                :selected.sync="defaultSize" 
                                :list="pagelist">
                                <bk-select-option v-for="(option, index) in pagelist" 
                                    :key="index" 
                                    :value="option.value" 
                                    :label="option.label">
                                </bk-select-option>
                            </bk-select>
                        </span>
                    </i18n>
                </span>
            </div>
            <div class="bk-page bk-page-compact fr">
                <ul class="pagination">
                    <li :title="$t('Common[\'首页\']')" :class="['page-item', {'disabled': pagination.current === 1}]" v-show="pagination.current !== 1">
                        <button class="page-button" style="font-size: 12px;" 
                            :disabled="pagination.current === 1" 
                            @click="turnToPage(1)">
                            <i class="icon-cc-backward"></i>
                        </button>
                    </li>
                    <li :title="$t('Common[\'上一页\']')" :class="['page-item', {'disabled': pagination.current === 1}]" v-show="pagination.current !== 1">
                        <button class="page-button" 
                            :disabled="pagination.current === 1" 
                            @click="turnToPage(pagination.current - 1)">
                            <i class="icon-cc-angle-left"></i>
                        </button>
                    </li>
                    <li class="page-item"
                        v-for="(page, index) in pageNum"
                        :title="`${page}`"
                        :class="{'cur-page':pagination.current === page}"  
                        :key="index">
                        <button class="page-button" 
                            @click="turnToPage(page)">
                            {{page}}
                        </button>
                    </li>
                    <li :title="$t('Common[\'下一页\']')" :class="['page-item', {'disabled': pagination.current === totalPage}]" v-show="pagination.current !== totalPage">
                        <button class="page-button"
                            :disabled="pagination.current === totalPage"
                            @click="turnToPage(pagination.current + 1)">
                            <i class="icon-cc-angle-right"></i>
                        </button>
                    </li>
                    <li :title="$t('Common[\'尾页\']')" :class="['page-item', {'disabled': pagination.current === totalPage}]" v-show="pagination.current !== totalPage">
                        <button class="page-button" style="font-size: 12px;" 
                            :disabled="pagination.current === totalPage" 
                            @click="turnToPage(totalPage)">
                            <i class="icon-cc-forward"></i>
                        </button>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            tableList: {
                type: Array,
                required: true
            },
            hasCheckbox: {
                type: Boolean,
                default: false
            },
            pagination: {
                type: Object,
                required: true
            },
            chooseId: {
                type: Array,
                default () {
                    return []
                }
            }
        },
        computed: {
            ...mapGetters([
                'language'
            ]),
            totalPage () {
                return Math.ceil(this.pagination.count / this.pagination.size)
            },
            pageNum () {
                let pageNum = []
                if (this.pagination.current < this.pagination.size) { // 如果当前的激活的项 小于要显示的条数
                    // 总页数和要显示的条数那个大就显示多少条
                    let i = Math.min(this.pagination.size, this.totalPage)
                    while (i) {
                        pageNum.unshift(i--)
                    }
                } else { // 当前页数大于显示页数了
                    let middle = this.pagination.current - Math.floor(this.pagination.size / 2) // 从哪里开始
                    let i = this.pagination.size
                    if (middle > (this.totalPage - this.pagination.size)) {
                        middle = (this.totalPage - this.pagination.size) + 1
                    }
                    while (i--) {
                        pageNum.push(middle++)
                    }
                }
                return pageNum
            }
        },
        data () {
            return {
                defaultSize: 10,
                pagelist: [{
                    value: 10,
                    label: 10
                }, {
                    value: 20,
                    label: 20
                }, {
                    value: 50,
                    label: 50
                }, {
                    value: 100,
                    label: 100
                }]
            }
        },
        watch: {
            defaultSize (newSize) {
                this.$emit('onPageSizeChange', newSize)
            }
        },
        created () {
            this.defaultSize = this.pagination.size || 10
        },
        mounted () {
            this.$watch('$refs.sizeSelector.open', (isOpen) => {
                this.$emit('handleSizeToggle', isOpen)
            })
        },
        methods: {
            turnToPage (pageNum) {
                if (pageNum !== this.pagination.current) {
                    this.$emit('onPageTurning', pageNum)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table-page-contain{
        width:100%;
        padding:5px 20px;
        background:#f9f9f9;
        &.empty{
            padding: 0;
        }
        .bk-page{
            height: 32px;
            >ul{
                height: 32px;
            }
            .page-item{
                min-width: 32px;
                height: 32px;
                line-height: 32px;
                button{
                    display: block;
                    width: 100%;
                    height: 100%;
                    margin: 0;
                    padding: 0;
                    border: none;
                    outline: 0;
                    &:disabled{
                        border-color: #e7e9ef;
                        background-color: #fafafa;
                        color: #bec6de;
                    }
                }
            }
        }
        .bk-page-compact{
            float:right;
        }
        .page-info-box{
            .page-info{
                padding: 4px 0;
                font-size:12px;
                color:#c3cdd7;
            }
        }
    }
</style>
