<template>
    <div class="collection-layout">
        <div class="collection-options clearfix">
            <label for="searchCollection" class="options-search fl">
                <i class="bk-icon icon-search"></i>
                <input id="searchCollection" class="search-input cmdb-form-input" type="text"
                    :placeholder="`${$t('Common[\'快速查询\']')}...`"
                    v-model.trim="filter.text">
            </label>
            <div class="options-sort fr">
                <span class="sort-item"
                    v-for="(option, index) in sortOptions"
                    :key="index"
                    @click="setFilterSort(option, filter.order === 'asc' ? 'desc': 'asc')">
                    <span>{{option.name}}</span>
                    <i class="sort-angle ascing"
                        :class="{'cur-sort' : filter.sort === option.id && filter.order === 'asc'}"
                        @click.stop="setFilterSort(option, 'asc')">
                    </i>
                    <i class="sort-angle descing"
                        :class="{'cur-sort' : filter.sort === option.id && filter.order === 'desc'}"
                        @click.stop="setFilterSort(option, 'desc')">
                    </i>
                </span>
            </div>
        </div>
        <ul class="collection-list">
            <li class="collection-item clearfix"
                v-for="(collection, index) in filteredList"
                :key="index"
                :class="{'delete-confirm': collection.deleteConfirm}"
                @click="handleApplyCollection(collection)">
                <span class="collection-name fl" v-show="!collection.edit">{{collection.name}}</span>
                <input class="collection-input cmdb-form-input fl" type="text" :ref="collection.id"
                    v-show="collection.edit"
                    v-model.trim="collection.name"
                    @blur="handleUpdateCollection(collection)">
                <i class="collection-icon icon-cc-del fr" @click.stop="handleConfirmDelete(collection)"></i>
                <i class="collection-icon icon-cc-edit fr" @click.stop="handleEditCollection($event, collection)"></i>
                <div class="collection-delete-confirm"
                    v-if="collection.deleteConfirm"
                    v-click-outside="handleCancelDelete">
                    <p class="confirm-title">{{$t('Hosts[\'确认删除\']')}}</p>
                    <bk-button class="confirm-btn" size="small" type="primary"
                        :loading="$loading(`delete_deleteFavorites_${collection.id}`)"
                        @click.stop="handleDeleteCollection(collection)">
                        {{$t('Hosts[\'确认\']')}}
                    </bk-button>
                    <bk-button class="confirm-btn" size="small" type="default"
                        @click.stop="handleCancelDelete">
                        {{$t('Common[\'取消\']')}}
                    </bk-button>
                </div>
            </li>
        </ul>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        data () {
            return {
                filter: {
                    text: '',
                    sort: 'frequency',
                    order: 'desc'
                },
                sortOptions: [{
                    id: 'name',
                    name: this.$t('Hosts["名称"]')
                }, {
                    id: 'frequency',
                    name: this.$t('Hosts["频率"]')
                }],
                list: [],
                filteredList: [],
                deleteConfirmCollection: null
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['privilegeBusiness'])
        },
        watch: {
            filter: {
                deep: true,
                handler (val) {
                    this.setFilteredList()
                }
            }
        },
        created () {
            this.getCollectionList()
        },
        methods: {
            ...mapActions('hostFavorites', [
                'searchFavorites',
                'udpateFavorites',
                'deleteFavorites'
            ]),
            setFilterSort (option, order) {
                this.filter.sort = option.id
                this.filter.order = order
            },
            getCollectionList () {
                this.searchFavorites({
                    params: {},
                    config: {
                        requestId: 'searchFavorites',
                        cancelPrevious: true
                    }
                }).then(data => {
                    this.list = data.info.filter(collection => {
                        const info = JSON.parse(collection.info)
                        return this.privilegeBusiness.some(business => business['bk_biz_id'] === info['bk_biz_id'])
                    })
                    this.setFilteredList()
                })
            },
            setFilteredList () {
                const filter = this.filter
                let filteredList = this.$tools.clone(this.list)
                if (filter.text.length) {
                    const filterText = filter.text.toLowerCase()
                    filteredList = filteredList.filter(collection => collection.name.toLowerCase().indexOf(filterText) !== -1)
                }
                filteredList.sort((collectionA, collectionB) => {
                    let compareResult
                    if (filter.sort === 'frequency') {
                        compareResult = filter.order === 'desc'
                            ? (collectionB.count - collectionA.count)
                            : (collectionA.count - collectionB.count)
                    } else {
                        compareResult = filter.order === 'desc'
                            ? collectionB.name.localeCompare(collectionA.name, 'zh-Hans-CN', {sensitivity: 'accent'})
                            : collectionA.name.localeCompare(collectionB.name, 'zh-Hans-CN', {sensitivity: 'accent'})
                    }
                    return compareResult
                })
                this.filteredList = filteredList
            },
            handleEditCollection (event, collection) {
                this.$set(collection, 'edit', true)
                this.$nextTick(() => {
                    const $input = this.$refs[collection.id][0]
                    $input.focus()
                })
            },
            handleUpdateCollection (collection) {
                this.$set(collection, 'edit', false)
                const originalCollection = this.list.find(original => original.id === collection.id)
                if (originalCollection.name !== collection.name) {
                    this.udpateFavorites({
                        id: collection.id,
                        params: {
                            ...originalCollection,
                            name: collection.name
                        },
                        config: {
                            requestId: `update_collection_${collection.id}`,
                            cancelPrevious: true
                        }
                    }).then(() => {
                        originalCollection.name = collection.name
                        this.setFilteredList()
                    })
                }
            },
            handleConfirmDelete (collection) {
                this.deleteConfirmCollection = collection
                this.$set(collection, 'deleteConfirm', true)
            },
            handleDeleteCollection (deleteCollection) {
                this.deleteFavorites({
                    id: deleteCollection.id,
                    config: {
                        requestId: `delete_deleteFavorites_${deleteCollection.id}`,
                        fromCache: true
                    }
                }).then(() => {
                    this.list = this.list.filter(collection => collection.id !== deleteCollection.id)
                    this.setFilteredList()
                })
            },
            handleCancelDelete () {
                if (this.deleteConfirmCollection) {
                    const collection = this.deleteConfirmCollection
                    this.$set(collection, 'deleteConfirm', false)
                    this.deleteConfirmCollection = null
                }
            },
            handleApplyCollection (collection) {
                this.$store.commit('hostFavorites/setApplying', collection)
                this.$emit('on-apply', collection)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .collection-layout {
        height: 100%;
    }
    .collection-options {
        margin: 20px 20px 0;
        .options-search {
            position: relative;
            width: 165px;
            .search-input {
                padding: 0 28px 0 12px;
                height: 30px;
                line-height: 28px;
                font-size: 12px;
            }
            .icon-search {
                position: absolute;
                top: 7px;
                right: 10px;
                font-size: 18px;
                color: $cmdbBorderColor;
            }
        }
        .options-sort {
            color: #c3cdd7;
            line-height: 28px;
            font-size: 0;
            .sort-item {
                position: relative;
                display: inline-block;
                vertical-align: middle;
                padding: 0 24px 0 12px;
                margin: 0 0 0 10px;
                font-size: 12px;
                border: 1px solid $cmdbBorderColor;
                border-radius: 2px;
                cursor: pointer;
            }
            .sort-angle {
                position: absolute;
                right: 10px;
                width: 0;
                height: 0;
                border: 5px solid transparent;
                border-bottom-color: currentColor;
                cursor: pointer;
                &:hover {
                    border-bottom-color: #4d597d;
                }
                &.cur-sort {
                    border-bottom-color: #4b8fe0;
                }
                &.ascing {
                    top: 3px;
                }
                &.descing {
                    bottom: 3px;
                    transform: rotate(180deg);
                }
            }
        }
    }
    .collection-list {
        height: calc(100% - 60px);
        margin: 10px 0 0 0;
        font-size: 12px;
        @include scrollbar-y;
        .collection-item {
            position: relative;
            height: 40px;
            padding: 0 10px 0 20px;
            line-height: 40px;
            cursor: pointer;
            &:hover,
            &.delete-confirm {
                background-color: #ebf4ff;
                .collection-icon {
                    display: block;
                }
            }
            .collection-name {
                width: 250px;
                padding: 0 0 0 10px;
                @include ellipsis;
            }
            .collection-input {
                width: 250px;
                height: 30px;
                padding: 0 0 0 8px;
                margin: 5px 0;
                line-height: 28px;
                font-size: 12px;
            }
            .collection-icon {
                display: none;
                width: 24px;
                height: 40px;
                line-height: 40px;
                text-align: center;
                cursor: pointer;
                &:hover {
                    color: #498fe0;
                }
                &.icon-cc-del:hover {
                    color: $cmdbDangerColor;
                }
            }
        }
    }
    .collection-delete-confirm {
        position: absolute;
        right: 0;
        top: 100%;
        padding: 10px;
        border: solid 1px #ececed;
        background-color: #fff;
        text-align: center;
        border-radius: 2px;
        box-shadow: 0px 2px 5px 0px rgba(0, 0, 0, 0.13);
        z-index: 1;
        &:before,
        &:after {
            content: '';
            position: absolute;
            right: 15px;
            bottom: 100%;
            width: 0;
            height: 0;
            border-left: 6px solid transparent;
            border-right: 6px solid transparent;
            border-bottom: 10px solid rgba(0, 0, 0, 0.13);
            z-index: 1;
        }
        &:after{
            bottom: calc(100% - 1px);
            border-bottom-color: #fff;
            z-index: 2;
        }
        .confirm-title {
            text-align: center;
        }
        .confirm-btn {
            height: 28px;
            line-height: 26px;
        }
    }
    .collection-item:nth-child(n + 5) {
        .collection-delete-confirm {
            bottom: 100%;
            top: auto;
            &:before,
            &:after {
                top: 100%;
                bottom: auto;
                border-top: 10px solid rgba(0, 0, 0, 0.13);
                border-bottom: none;
            }
            &:after {
                top: calc(100% - 1px);
                border-top-color: #fff;
            }
        }
    }
</style>