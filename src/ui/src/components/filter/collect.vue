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
    <div class="collect-wrapper">
        <div class="collect-filter clearfix">
            <div class="filter-group filter-group-search fl">
                <input class="bk-form-input" type="text" :placeholder="`${$t('Common[\'快速查询\']')}...`" v-model.trim="filter.text">
                <i class="bk-icon icon-search"></i>
            </div>
            <div class="filter-group filter-group-sort fr">
                <span class="sort-box" v-for="(option, index) in sortOptions" :key="index">
                    <span>{{option.name}}</span>
                    <i class="sort-angle ascing" :class="{'cur-sort' : filter.sort.type === option.id && filter.sort.order === 'asc'}" @click="sortFavoriteList(option, 'asc')"></i>
                    <i class="sort-angle descing" :class="{'cur-sort' : filter.sort.type === option.id && filter.sort.order === 'desc'}" @click="sortFavoriteList(option, 'desc')"></i>
                </span>
            </div>
        </div>
        <div class="collect-list-wrapper">
            <ul class="collect-list">
                <li ref="collectItem" class="collect-item"
                    v-for="(item, index) in localFavoriteList"
                    v-if="item.isShow"
                    :key="index"
                    :class="{'editing': item.edit || item.isShowDeleteConfirm}"
                    @click="apply(item)">
                    <span class="collect-item-name" v-show="!item.edit">{{item['name']}}</span>
                    <input ref="collectNameInput" class="bk-form-input" type="text"
                        :itemId="item['id']"
                        v-show="item['edit']"
                        v-model.trim="item['name']"
                        @click.stop
                        @blur="updateCollectName(item)">
                    <i class="icon-cc-edit" @click.stop="editCollectName(item)"></i>
                    <i class="icon-cc-del" @click.stop="item.isShowDeleteConfirm = true"></i>
                    <div class="collect-delete-pop" v-if="item.isShowDeleteConfirm">
                        <div class="btn-content">
                            <p>{{$t('Hosts[\'确认删除\']')}}</p>
                            <button class="main-btn" @click.stop="deleteCollect(item, index)">{{$t('Hosts[\'确认\']')}}</button>
                            <button class="vice-btn" @click.stop="item.isShowDeleteConfirm = false">{{$t('Common[\'取消\']')}}</button>
                        </div>
                        <div class="collect-delete-mask" @click.stop="item.isShowDeleteConfirm = false">
                        </div>
                    </div>
                </li>
            </ul>
        </div>
    </div>
</template>
<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            favoriteList: {
                type: Array,
                required: true
            },
            active: Boolean
        },
        data () {
            return {
                filter: {
                    text: '',
                    sort: {
                        type: '',
                        order: ''
                    }
                },
                sortOptions: [{
                    id: 'name',
                    name: this.$t('Hosts[\'名称\']')
                }, {
                    id: 'frequency',
                    name: this.$t('Hosts[\'频率\']')
                }],
                localFavoriteList: []
            }
        },
        computed: {
            ...mapGetters(['bkBizList'])
        },
        watch: {
            active () {
                this.filter.text = ''
            },
            favoriteList (favoriteList) {
                this.localFavoriteList = this.favoriteList.map(favorite => {
                    return Object.assign({
                        edit: false,
                        isShow: favorite.name.toLowerCase().indexOf(this.filter.text.toLowerCase()) !== -1,
                        isShowDeleteConfirm: false
                    }, favorite)
                })
            },
            'filter.text' (filterText) {
                this.localFavoriteList.map((item, index) => {
                    item.isShow = item.name.toLowerCase().indexOf(filterText.toLowerCase()) !== -1
                })
            },
            'filter.sort' (sort) {
                let favoriteList = this.favoriteList.slice(0)
                if (sort.type === 'name') {
                    favoriteList.sort((collectA, collectB) => {
                        return collectA['name'].localeCompare(collectB['name'])
                    })
                } else if (sort.type === 'frequency') {
                    favoriteList.sort((collectA, collectB) => {
                        return collectA['count'] - collectB['count']
                    })
                }
                if (sort.order === 'desc') {
                    favoriteList.reverse()
                }
                this.localFavoriteList = favoriteList.map(favorite => {
                    return Object.assign({
                        edit: false,
                        isShow: true,
                        isShowDeleteConfirm: false
                    }, favorite)
                })
            }
        },
        methods: {
            sortFavoriteList (option, order) {
                this.filter.sort = {
                    type: option.id,
                    order: order
                }
            },
            editCollectName (item) {
                item.edit = true
                this.$nextTick(() => {
                    this.$refs.collectNameInput.map(input => {
                        if (input.getAttribute('itemId') === item['id']) {
                            input.focus()
                        }
                    })
                })
            },
            updateCollectName (item) {
                let updateItem = Object.assign({}, item)
                delete updateItem.isShow
                delete updateItem.edit
                delete updateItem.isShowDeleteConfirm
                this.$axios.put(`hosts/favorites/${updateItem['id']}`, updateItem).then(res => {
                    if (res.result) {
                        item.edit = false
                        this.$emit('update', updateItem)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            deleteCollect (item, index) {
                this.$axios.delete(`hosts/favorites/${item['id']}`).then(res => {
                    if (res.result) {
                        this.localFavoriteList.splice(index, 1)
                        this.$alertMsg(this.$t('Common[\'删除成功\']'), 'success')
                        this.$emit('delete', item, index)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            apply (collect) {
                let isAppExist = false
                let collectBkBizId = JSON.parse(collect['info'])['bk_biz_id']
                this.bkBizList.map(({bk_biz_id: bkBizId}) => {
                    if (bkBizId === collectBkBizId) {
                        isAppExist = true
                    }
                })
                if (isAppExist) {
                    this.$emit('apply', collect)
                } else {
                    this.$alertMsg(this.$t('Common[\'该查询条件对应的业务不存在\']'))
                }
            }
        }
    }
</script>
<style lang="scss" scoped>
    .collect-wrapper{
        padding: 20px 0;
        height: 100%;
    }
    .collect-filter{
        .filter-group{
            position: relative;
            &.filter-group-search{
                width: 175px;
                .bk-form-input{
                    padding: 0 10px 0 28px;
                    height: 30px;
                    line-height: 28px;
                    border: 1px solid transparent;
                    font-size: 12px;
                    &:focus{
                        border-color: #6b7baa;
                    }
                }
                .icon-search{
                    position: absolute;
                    left: 7px;
                    top: 7px;
                    font-size: 16px;
                    color:#bec6de;
                }
            }
            &.filter-group-sort{
                .sort-box{
                    position: relative;
                    padding: 0 15px 0 0;
                    font-size: 12px;
                    color: #c3cdd7;
                    // color: #bfc7de;
                    height: 30px;
                    line-height: 30px;
                    display: inline-block;
                    vertical-align: middle;
                    margin: 0 0 0 5px;
                    &:hover{
                        color: #6b7baa;
                    }
                    .sort-angle{
                        position: absolute;
                        right: 0;
                        width: 0;
                        height: 0;
                        border: 5px solid transparent;
                        cursor: pointer;
                        &.ascing{
                            border-bottom-color: #bfc7de;
                            top: 4px;
                            &:hover{
                                border-bottom-color: #4d597d;
                            }
                            &.cur-sort{
                                border-bottom-color: #4b8fe0;
                            }
                        }
                        &.descing{
                            border-top-color: #bfc7de;
                            bottom: 4px;
                            &:hover{
                                border-top-color: #4d597d;
                            }
                            &.cur-sort{
                                border-top-color: #4b8fe0;
                            }
                        }
                    }
                }
            }
        }
    }
    .collect-list-wrapper{
        height: calc(100% - 20px);
        padding: 10px 0;
        font-size: 12px;
        overflow-y: auto;
        @include scrollbar;
        .collect-list{
            .collect-item{
                position: relative;
                height: 30px;
                line-height: 30px;
                padding: 0 50px 0 0;
                cursor: pointer;
                &:hover,
                &.editing{
                    background-color: #f9f9f9;
                    .icon-cc-edit,
                    .icon-cc-del{
                        display: block;
                    }
                }
                .collect-item-name{
                    display: inline-block;
                    vertical-align: top;
                    max-width: 250px;
                    padding: 0 0 0 10px;
                    @include ellipsis;
                }
                .bk-form-input{
                    height: 30px;
                    line-height: 28px;
                    vertical-align: top;
                }
                .icon-cc-edit,
                .icon-cc-del{
                    display: none;
                    position: absolute;
                    font-size: 12px;
                    top: 9px;
                    width: 24px;
                    text-align: center;
                }
                .icon-cc-edit{
                    right: 30px;
                }
                .icon-cc-del{
                    right: 5px;
                }
                .collect-delete-pop {
                    position: absolute;
                    width: 181px;
                    height: 106px;
                    background-color: #ffffff;
                    box-shadow: 0px 2px 5px 0px rgba(0, 0, 0, 0.13);
                    border: solid 1px #ececed;
                    right: 2px;
                    text-align: center;
                    padding-top: 18px;
                    z-index: 1;
                    &:before{
                        content: '';
                        right: 9px;
                        bottom: 105px;
                        width: 0;
                        height: 0;
                        border-left: 6px solid transparent;
                        border-right: 6px solid transparent;
                        border-bottom: 10px solid #e7e9ef;
                        position: absolute;
                    }
                    &:after{
                        content: "";
                        right: 9px;
                        bottom: 104px;
                        width: 0px;
                        height: 0px;
                        border-left: 6px solid transparent;
                        border-right: 6px solid transparent;
                        border-bottom: 10px solid #fff;
                        position: absolute;
                    }
                    .btn-content{
                        position: relative;
                        z-index: 2;
                        p{
                            font-size: 16px;
                            line-height: 0px;
                            letter-spacing: 0px;
                            color: #6b7baa;
                        }
                        button{
                            width: 55px;
                            height: 28px;
                            line-height: 26px;
                            display: inline-block;
                            border-radius: 2px;
                            margin-top: 10px;
                        }
                    }
                    .collect-delete-mask{
                        position: fixed;
                        left: 0;
                        top: 0;
                        right: 0;
                        bottom: 0;
                        z-index: 1;
                        cursor: default;
                    }
                }
                &:nth-child(n+5){
                    .collect-delete-pop{
                        bottom: 100%;
                        margin-bottom: 5px;
                        &:before{
                            border-top: 10px solid #e7e9ef;
                            border-bottom: none;
                            bottom: -13px;
                        }
                        &:after{
                            border-top: 10px solid #fff;
                            border-bottom: none;
                            bottom: -10px;
                        }
                    }
                }
            }
        }
    }
</style>