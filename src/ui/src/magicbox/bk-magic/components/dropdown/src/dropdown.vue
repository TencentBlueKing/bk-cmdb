<template>
    <div class="bk-dropdown"
        :class="[extCls, {'open': open}]"
        @click="openFn"
        v-clickoutside="close">
        <div class="bk-dropdown-wrapper">
            <input class="bk-dropdown-input" readonly="readonly"
                :class="{placeholder: localSelected === -1, active: open}"
                :value="multiSelect ? selectedList.join(',') : selectedText"
                :placeholder="placeholder"
                :disabled="disabled"
                @mouseover="showClearFn"
                @mouseleave="showClear = false">
            <i class="bk-icon icon-angle-down bk-dropdown-icon" v-show="!showClear"></i>
            <i class="bk-icon icon-close bk-dropdown-icon clear-icon"
                v-show="showClear"
                @mouseover="showClearFn"
                @mouseleave="showClear = false"
                @click="clearSelected($event)">
            </i>
        </div>
        <transition name="toggle-slide">
            <div class="bk-dropdown-list" v-show="open">
                <!-- 搜索栏 -->
                <div class="bk-dropdown-search-item"
                    @click="$event.stopPropagation()"
                    v-if="searchable">
                    <i class="bk-icon icon-search"></i>
                    <input type="text" v-model="condition" @input="inputFn">
                </div>
                <ul>
                    <!-- 新增项 -->
                    <li class="bk-dropdown-list-item bk-dropdown-create-item"
                        v-if="hasCreateItem"
                        @click.stop="createFn">
                        <div class="text">
                            {{ createText }}
                        </div>
                    </li>
                    <li class="bk-dropdown-list-item"
                        :class="{'bk-dropdown-selected': !multiSelect && index === localSelected}"
                        v-if="localList.length !== 0"
                        v-for="(item, index) in localList"
                        @click.stop="selectItem(index)">
                        <div class="text">
                            <label class="bk-form-checkbox bk-checkbox-small mr0 bk-dropdown-multi-label" v-if="multiSelect">
                                <input type="checkbox"
                                :name="'multiSelect' + +new Date()"
                                :value="index"
                                v-model="localSelected">
                                {{ item[displayKey] }}
                            </label>
                            <template v-else>
                                {{ item[displayKey] }}
                            </template>
                        </div>
                        <div class="bk-dropdown-tools" v-if="tools !== false">
                            <i class="bk-icon icon-edit2 bk-dropdown-list-icon"
                                v-if="tools.edit !== false"
                                @click.stop="editFn(index)"></i>
                            <i class="bk-icon icon-close bk-dropdown-list-icon"
                                v-if="tools.del !== false"
                                @click.stop="delFn(index)"></i>
                        </div>
                    </li>
                    <li class="bk-dropdown-list-item" v-if="localList.length === 0">
                        <div class="text">
                            {{ emptyText }}
                        </div>
                    </li>
                </ul>
            </div>
        </transition>
    </div>
</template>

<script>
    import clickoutside from './../../../utils/clickoutside'

    /**
     *  bk-dropdown
     *  @module components/dropdown
     *  @desc 下拉选框组件，模拟原生select
     *  @param extCls {String} - 自定义的样式
     *  @param hasCreateItem {Boolean} - 下拉菜单中是否有新增项，默认为true
     *  @param createText {String} - 下拉菜单中新增项的文字
     *  @param tools {Object, Boolean} - 待选项右侧的工具按钮，有两个可配置的key：edit和del，默认为两者都不显示。
     *  @param list {Array} - 必选，下拉菜单所需的数据列表
     *  @param selected {Number} - 必选，选中的项的index值，支持.sync修饰符
     *  @param placeholder {String, Boolean} - 是否显示占位行，默认为显示，且文字为“请选择”
     *  @param displayKey {String} - 循环list时，显示字段的key值，默认为name
     *  @param disabled {Boolean} - 是否禁用组件，默认为false
     *  @param multiSelect {Boolean} - 是否支持多选，默认为false
     *  @param searchable {Boolean} - 是否支持筛选，默认为false
     *  @param searchKey {Boolean} - 筛选时，搜索的key值，默认为'name'
     *  @param allowClear {Boolean} - 是否可以清除单选时选中的项
     *  @param settingKey {String} - 根据配置这个字段，自定义在单选时，选中某项之后的回调函数的第一个返回值的内容
     *  @example
        <bk-dropdown
            :list="list"
            :tools="tools"
            :selected.sync="selected"
            :placeholder="placeholder"
            :displayKey="displayKey"
            :has-create-item="hasCreateItem"
            :create-text="createText"
            :ext-cls="extCls"></bk-dropdown>
    */
    export default {
        name: 'bk-dropdown',
        props: {
            extCls: {
                type: String
            },
            hasCreateItem: {
                type: Boolean,
                default: false
            },
            createText: {
                type: String,
                default: '新增'
            },
            tools: {
                type: [Object, Boolean],
                default: false
            },
            list: {
                type: Array,
                required: true
            },
            selected: {
                type: [Number, Array],
                required: true
            },
            placeholder: {
                type: [String, Boolean],
                default: '请选择'
            },
            displayKey: {
                type: String,
                default: 'name'
            },
            disabled: {
                type: Boolean,
                default: false
            },
            multiSelect: {
                type: Boolean,
                default: false
            },
            searchable: {
                type: Boolean,
                default: false
            },
            searchKey: {
                type: String,
                default: 'name'
            },
            allowClear: {
                type: Boolean,
                default: false
            },
            settingKey: {
                type: String,
                default: 'id'
            }
        },
        data () {
            return {
                open: false,
                selectedList: this.multiSelect ? this.calcSelected() : [],
                condition: '',
                localList: this.list,
                localSelected: this.selected,
                emptyText: this.list.length ? '无匹配数据' : '暂无数据',
                showClear: false
            }
        },
        watch: {
            selected (newVal) {
                // 多选时，将index转成文字赋给控制显示的数组
                if (this.multiSelect) {
                    this.selectedList.splice(0, this.selectedList.length, this.calcSelected(newVal))
                } else {
                    this.selectedList = []
                }

                this.localSelected = newVal
            },
            condition () {
                this.calcList()
            }
        },
        directives: {
            clickoutside
        },
        computed: {
            currentItem () {
                return this.list[this.localSelected]
            },
            selectedText () {
                if (!this.multiSelect) {
                    return this.localSelected === -1 ? this.placeholder : (this.list.length ? this.currentItem[this.displayKey] : '暂无数据')
                }
            }
        },
        methods: {
            // 多选时，根据传入的index值，算出应该显示和选中的文字
            calcSelected () {
                let list = this.list
                let displayKey = this.displayKey
                let selectedTextArr = []

                for (let _arr of this.selected) {
                    for (let [index, item] of list.entries()) {
                        if (_arr === index) {
                            selectedTextArr.push(item[displayKey])
                        }
                    }
                }

                return selectedTextArr
            },
            close () {
                this.open = false
                this.$emit('visible-toggle', this.open)
            },
            openFn () {
                if (!this.disabled) {
                    this.open = !this.open
                    this.$emit('visible-toggle', this.open)
                }
            },
            /**
             *  计算返回渲染的数组
             */
            calcList () {
                if (this.searchable) {
                    let arr = []
                    let key = this.searchKey

                    for (let item of this.list) {
                        if (item[key].includes(this.condition)) {
                            arr.push(item)
                        }
                    }

                    this.localList = arr
                } else {
                    this.localList = this.list
                }
            },
            /**
             *  是否显示清除当前选择的icon
             */
            showClearFn () {
                if (this.allowClear && !this.multiSelect && this.localSelected !== -1) {
                    this.showClear = true
                }
            },
            /**
             *  清除选择
             */
            clearSelected (e) {
                this.$emit('clear', this.localSelected)
                this.localSelected = -1
                this.showClear = false
                e.stopPropagation()
            },
            /**
             *  选中列表中的项
             */
            selectItem (index) {
                let item, $selected, $index
                let settingKey = this.settingKey
                let isMultiSelect = this.multiSelect
                let list = this.localList

                if (index !== undefined) {
                    if (!isMultiSelect) {
                        $selected = index
                    } else {
                        $index = []
                        $selected = []

                        for (let item of this.localSelected) {
                            $selected.push(list[item])
                            $index.push(item)
                        }
                    }

                    item = this.list[index]
                    this.$emit('update:selected', isMultiSelect ? $index : $selected)
                } else {
                    item = null
                    this.$emit('update:selected', -1)
                }

                // 单选时，返回的两个参数是选中项的id（或通过settingKey配置的值）和选中项的数据
                // 多选时，返回的是选中项的索引数组和选中项的数据数组
                this.$emit('item-selected', isMultiSelect ? $index : (settingKey === 'id' ? index : list[index][settingKey]), isMultiSelect ? $selected : item)

                if (!this.multiSelect) {
                    this.openFn()
                }

                this.condition = ''
            },
            editFn (e, index) {
                this.$emit('edit', index)
                this.openFn()
                e.stopPropagation()
            },
            delFn (e, index) {
                this.$emit('del', index)
                this.openFn()
                e.stopPropagation()
            },
            createFn (e) {
                this.$emit('create')
                this.openFn()
                e.stopPropagation()
            },
            inputFn () {
                this.$emit('typing', this.condition)
            }
        },
        mounted () {
            this.popup = this.$el
        }
    }
</script>
