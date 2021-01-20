<template>
    <bk-dialog
        v-model="isShow"
        :draggable="false"
        :width="730"
        @after-leave="handleClosed">
        <div class="title" slot="tools">
            <span>{{title}}</span>
            <bk-input class="filter-input" v-model.trim="filter" clearable :placeholder="$t('请输入关键字搜索')"></bk-input>
        </div>
        <div v-bkloading="{ isLoading: !ready }">
            <div class="property-tips" v-show="!isEmpty">
                <span class="tips-text">{{$t('请选择导出字段')}}:</span>
            </div>
            <bk-exception class="search-empty"
                v-if="isEmpty"
                type="search-empty"
                scene="part">
            </bk-exception>
            <section class="property-selector" v-show="!isEmpty">
                <ul class="property-list property-list-independent clearfix" v-if="independentProperties.length">
                    <li class="property-item fl"
                        v-for="property in independentProperties"
                        :key="property.bk_property_id">
                        <bk-checkbox class="property-checkbox"
                            :checked="isChecked(property)"
                            @change="handleToggleProperty(property, ...arguments)">
                            {{property.bk_property_name}}
                        </bk-checkbox>
                    </li>
                </ul>
                <div class="group"
                    v-for="group in renderGroups"
                    :key="group.id">
                    <h2 class="group-title">
                        {{group.name}}
                        <bk-checkbox class="property-all"
                            :checked="isGroupAllChecked[group.id]"
                            @change="handleGroupToggleAll(group.id)">
                            {{$t('全选')}}
                        </bk-checkbox>
                    </h2>
                    <ul class="property-list clearfix">
                        <li class="property-item fl"
                            v-for="property in group.properties"
                            :key="property.bk_property_id">
                            <bk-checkbox class="property-checkbox"
                                :checked="isChecked(property)"
                                @change="handleToggleProperty(property, ...arguments)">
                                {{property.bk_property_name}}
                            </bk-checkbox>
                        </li>
                    </ul>
                </div>
            </section>
        </div>
        <footer class="footer" slot="footer">
            <i18n path="已选个数" v-show="!!selected.length">
                <span class="count" place="count">{{selected.length}}</span>
            </i18n>
            <div class="selected-options">
                <bk-button theme="primary" :disabled="!selected.length" @click="confirm">{{$t('开始导出')}}</bk-button>
                <bk-button class="ml10" theme="default" @click="close">{{$t('取消')}}</bk-button>
            </div>
        </footer>
    </bk-dialog>
</template>

<script>
    import Throttle from 'lodash.throttle'
    export default {
        props: {
            properties: {
                type: Array,
                default: () => ([])
            },
            propertyGroups: {
                type: Array,
                default: () => ([])
            },
            invisibleProperties: {
                type: Array,
                default: () => ([
                    'bk_host_id',
                    'bk_cloud_id',
                    'bk_host_innerip',
                    '__bk_host_topology__'
                ])
            },
            handler: {
                type: Function,
                default: () => {}
            },
            title: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                ready: false,
                filter: '',
                renderGroups: [],
                isShow: false,
                selected: [],
                throttleFilter: Throttle(this.handleFilter, 500, { leading: false })
            }
        },
        computed: {
            availableProperties () {
                return this.properties.filter(property => !this.invisibleProperties.includes(property.bk_property_id))
            },
            sortedProperties () {
                return [...this.availableProperties].sort((propertyA, propertyB) => {

                })
            },
            independentProperties () {
                return this.availableProperties.filter(property => {
                    return !this.propertyGroups.some(group => {
                        return group.bk_group_id === property.bk_property_group
                            && group.bk_obj_id === property.bk_obj_id
                    })
                })
            },
            groups () {
                const sortedGroups = [...this.propertyGroups].sort((groupA, groupB) => {
                    return groupA.bk_group_index - groupB.bk_group_index
                })
                const groups = sortedGroups.map(group => {
                    return {
                        id: group.bk_group_id,
                        name: group.bk_group_name,
                        properties: this.availableProperties.filter(property => property.bk_property_group === group.bk_group_id)
                    }
                })
                return groups.filter(group => group.properties.length)
            },
            visibleProperties () {
                return this.renderGroups.reduce((accumulator, group) => {
                    return accumulator.concat(group.properties)
                }, [])
            },
            isAllChecked () {
                if (this.filter) {
                    return this.visibleProperties.every(property => this.selected.includes(property))
                }
                return !!this.selected.length && this.selected.length === this.visibleProperties.length
            },
            isEmpty () {
                return this.ready && !this.visibleProperties.length
            },
            isGroupAllChecked () {
                const isGroupAllChecked = {}
                this.renderGroups.forEach(group => {
                    isGroupAllChecked[group.id] = group.properties.every(property => this.selected.includes(property))
                })
                return isGroupAllChecked
            }
        },
        watch: {
            filter: {
                immediate: true,
                handler () {
                    this.throttleFilter()
                }
            }
        },
        methods: {
            handleFilter () {
                this.ready = true
                this.$nextTick(() => {
                    if (!this.filter.length) {
                        this.renderGroups = this.groups
                    } else {
                        const filteredGroups = []
                        const filter = this.filter.toLowerCase()
                        this.groups.forEach(group => {
                            const properties = group.properties.filter(property => {
                                const name = property.bk_property_name.toLowerCase()
                                return name.indexOf(filter) > -1
                            })
                            if (properties.length) {
                                filteredGroups.push({
                                    ...group,
                                    properties
                                })
                            }
                        })
                        this.renderGroups = filteredGroups
                    }
                })
            },
            isChecked (property) {
                return this.selected.some(target => target.id === property.id)
            },
            handleToggleProperty (property, checked) {
                if (checked) {
                    this.selected.push(property)
                } else {
                    const index = this.selected.findIndex(target => target.id === property.id)
                    index > -1 && this.selected.splice(index, 1)
                }
            },
            handleToggleAll () {
                if (this.filter) {
                    if (this.isAllChecked) {
                        this.selected = this.selected.filter(property => !this.visibleProperties.includes(property))
                    } else {
                        const newSelected = this.visibleProperties.filter(property => !this.selected.includes(property))
                        this.selected = this.selected.concat(newSelected)
                    }
                } else {
                    this.selected = this.isAllChecked ? [] : [...this.visibleProperties]
                }
            },
            handleGroupToggleAll (groupId) {
                const group = this.renderGroups.find(group => group.id === groupId)
                const isGroupAllChecked = this.isGroupAllChecked[groupId]
                if (isGroupAllChecked) {
                    this.selected = this.selected.filter(property => {
                        return !group.properties.includes(property)
                    })
                } else {
                    const newSelected = group.properties.filter(property => !this.selected.includes(property))
                    this.selected = this.selected.concat(newSelected)
                }
            },
            confirm () {
                this.handler(this.selected)
                this.close()
            },
            handleClosed () {
                this.$emit('closed')
            },
            open () {
                this.isShow = true
            },
            close () {
                this.isShow = false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .title {
        display: flex;
        justify-content: space-between;
        align-items: center;
        vertical-align: middle;
        line-height: 30px;
        font-size: 24px;
        color: #444;
        padding: 15px 0 0 24px;
        .filter-input {
            width: 240px;
            margin-right: 45px;
        }
    }
    .property-tips {
        display: flex;
        padding: 15px 0;
        justify-content: space-between;
        align-items: center;
        line-height: 20px;
    }
    .property-selector {
        margin: 0 -24px -24px 0;
        height: 350px;
        @include scrollbar-y;
    }
    .group {
        margin-bottom: 15px;
        .group-title {
            display: flex;
            justify-content: space-between;
            align-items: center;
            position: relative;
            padding: 0 24px 0 15px;
            line-height: 20px;
            font-size: 15px;
            font-weight: bold;
            color: #63656E;
            &:before {
                content: "";
                position: absolute;
                left: 0;
                top: 3px;
                width: 4px;
                height: 14px;
                background-color: #C4C6CC;
            }
            .property-all {
                font-weight: normal;
            }
        }
    }
    .property-list {
        padding: 10px 0 6px 0;
        &.property-list-independent ~ {
            .group {
                margin-top: 10px;
            }
        }
        .property-item {
            width: 33%;
        }
    }
    .property-checkbox {
        display: block;
        margin: 8px 20px 8px 0;
        @include ellipsis;
        /deep/ {
            .bk-checkbox-text {
                max-width: calc(100% - 25px);
                @include ellipsis;
            }
        }
    }
    .footer {
        display: flex;
        align-items: center;
        font-size: 14px;
        .count {
            font-weight: bold;
            color: #2DCB56;
            padding: 0 2px;
        }
        .selected-options {
            display: flex;
            align-items: center;
            margin-left: auto;
        }
    }
    .search-empty {
        height: 400px;
        margin: 0 -24px -24px 0;
        justify-content: center;
    }
</style>
