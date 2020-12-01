<template>
    <bk-dialog
        v-model="isShow"
        :draggable="false"
        :width="730"
        @after-leave="handleClosed">
        <div class="title" slot="tools">
            <span>{{$t('选择导出字段')}}</span>
            <bk-input class="filter-input" v-model.trim="filter" clearable :placeholder="$t('请输入关键字搜索')"></bk-input>
        </div>
        <section class="property-selector">
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
        <footer class="footer" slot="footer">
            <bk-checkbox class="property-all"
                :checked="isAllChecked"
                @change="handleToggleAll">
                {{$t('全选')}}
            </bk-checkbox>
            <div class="selected-options">
                <bk-button theme="primary" :disabled="!selected.length" @click="confirm">{{$t('确定')}}</bk-button>
                <bk-button theme="default" @click="close">{{$t('取消')}}</bk-button>
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
                    'bk_cloud_id',
                    'bk_host_innerip'
                ])
            },
            handler: {
                type: Function,
                default: () => {}
            }
        },
        data () {
            return {
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
            isAllChecked () {
                return !!this.selected.length && this.selected.length === this.availableProperties.length
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
                this.selected = this.isAllChecked ? [] : [...this.availableProperties]
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
        line-height: 31px;
        font-size: 24px;
        color: #444;
        padding: 15px 0 0 24px;
        .filter-input {
            width: 240px;
            margin-right: 45px;
        }
    }
    .property-selector {
        margin: 0 -24px -24px 0;
        height: 350px;
        @include scrollbar-y;
    }
    .group {
        margin-top: 15px;
        .group-title {
            position: relative;
            padding: 0 0 0 15px;
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
        .selected-options {
            margin-left: auto;
        }
    }
</style>
