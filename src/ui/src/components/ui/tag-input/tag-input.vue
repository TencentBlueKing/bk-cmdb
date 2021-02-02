<template>
    <div class="tag-input"
        :style="{
            height: fixedHeight ? inputHeight + 'px' : 'auto'
        }"
        @mousedown="() => (shouldUpdate = false)"
        @click="focus">
        <div class="tag-input-layout">
            <div class="tag-input-container"
                ref="container"
                :class="{
                    focus: isFocus,
                    disabled: disabled,
                    placeholder: !localValue.length && !isFocus,
                    'is-fast-clear': fastClear,
                    'is-loading': loading,
                    'is-flex-height': !fixedHeight
                }"
                :style="containerStyle"
                :data-placeholder="placeholder"
                @mousewheel="handleContainerScroll">
                <template v-if="multiple || !isFocus">
                    <span class="tag-input-selected"
                        v-for="(tag, index) in localValueData"
                        ref="selected"
                        :key="tag.value"
                        @click.stop
                        @mousedown.left.stop="handleSelectedMousedown($event, index)"
                        @mouseup.left.stop="handleSelectedMouseup($event, index)"
                        @mouseenter="handleSelectedMouseenter($event, tag)"
                        @mouseleave="handleSelectedMouseleave($event, tag)">
                        <template v-if="renderTag">
                            <render-tag :tag="tag" :index="index"></render-tag>
                        </template>
                        <template v-else>
                            <span class="tag-input-selected-value">
                                {{getDisplayText(tag)}}
                            </span>
                        </template>
                        <i class="tag-input-selected-clear bk-icon icon-close" v-if="tagClearable"
                            @mouseup.left.stop
                            @mousedown.left.stop="handleRemoveMouseDown"
                            @click.stop.prevent="handleRemoveSelected(tag, index)">
                        </i>
                    </span>
                </template>
                <span ref="input" class="tag-input-input"
                    spellcheck="false"
                    contenteditable
                    v-show="isFocus"
                    @click.stop
                    @input="handleInput($event)"
                    @blur="handleBlur"
                    @paste.prevent.stop="handlePaste($event)"
                    @keydown="handleKeydown($event)">
                </span>
            </div>
            <i class="tag-input-clear bk-icon icon-close-circle-shape"
                v-if="fastClear && !disabled && localValue.length"
                @click.stop="handleFastClear">
            </i>
        </div>
    </div>
</template>

<script>
    import Throttle from 'lodash.throttle'
    import RenderTag from './render-tag.js'
    import AlternateList from './alternate-list'
    import Vue from 'vue'
    import Tippy from 'bk-magic-vue/lib/utils/tippy'
    export default {
        name: 'cmdb-tag-input',
        components: {
            RenderTag
        },
        props: {
            value: {
                type: Array,
                default: () => ([])
            },
            placeholder: {
                type: String,
                default: '请输入'
            },
            disabled: {
                type: Boolean,
                default: false
            },
            multiple: {
                type: Boolean,
                default: true
            },
            allowCreate: {
                type: Boolean,
                default: true
            },
            createOnly: Boolean,
            focusRowLimit: {
                type: Number,
                default: 6
            },
            defaultAlternate: {
                type: [String, Array, Function],
                validator (value) {
                    return value === 'history' || typeof value === 'function' || value instanceof Array
                }
            },
            searchFromDefaultAlternate: {
                type: Boolean,
                default: true
            },
            historyKey: String,
            historyLabel: {
                type: String,
                default: '最近选择'
            },
            historyRecord: {
                type: Number,
                default: 5
            },
            fuzzySearchMethod: Function,
            exactSearchMethod: Function,
            emptyText: {
                type: String,
                default: '暂无数据'
            },
            tagClearable: {
                type: Boolean,
                default: true
            },
            fastClear: Boolean,
            renderList: Function,
            renderTag: Function,
            displayTagTips: Boolean,
            tagTipsContent: Function,
            tagTipsDelay: {
                type: Number,
                default: 300
            },
            fixedHeight: {
                type: Boolean,
                default: true
            },
            disabledData: {
                type: Array,
                default: () => ([])
            },
            listScrollHeight: [Number, String],
            searchLimit: {
                type: Number,
                default: 20
            },
            pasteFormatter: {
                type: Function,
                default (value) {
                    return value.trim()
                }
            },
            pasteValidator: {
                type: Function,
                default: values => values
            },
            panelWidth: {
                type: [Number, String],
                validator (value) {
                    const pixel = parseInt(value)
                    return pixel >= 190
                }
            }
        },
        data () {
            return {
                inputHeight: 32,
                singleRowHeight: 30,
                inputValue: '',
                inputIndex: 0,
                highlightIndex: -1,
                shouldUpdate: true,
                isFocus: false,
                overflowTagIndex: null,
                currentData: [],
                matchedData: [],
                flattenedData: [],
                scheduleSearch: Throttle(this.search, 800, { leading: false }),
                popoverInstance: null,
                alternateContent: null,
                selectedTipsTimer: {},
                overflowTagNode: null,
                loading: false
            }
        },
        computed: {
            containerStyle () {
                const style = {}
                if (this.isFocus) {
                    style.maxHeight = this.fixedHeight ? this.focusRowLimit * this.singleRowHeight + 'px' : 'auto'
                } else if (this.fixedHeight) {
                    style.height = this.singleRowHeight + 'px'
                }
                return style
            },
            localValue: {
                get () {
                    return [...this.value]
                },
                set (value) {
                    this.$emit('input', value)
                    this.$emit('change', value)
                }
            },
            localValueData () {
                return this.localValue.map(value => {
                    const tag = this.currentData.find(tag => tag.value === value)
                    return tag || { value }
                })
            }
        },
        watch: {
            inputValue (value) {
                if (value.length) {
                    this.highlightIndex = -1
                    this.updateScroller()
                    this.scheduleSearch(value)
                } else if (this.isFocus) {
                    this.search()
                }
            },
            isFocus (isFocus) {
                if (isFocus) {
                    this.search()
                    this.$emit('focus')
                } else {
                    this.reset()
                    this.$emit('blur')
                }
                this.calcOverflow()
            },
            localValue (localValue) {
                this.calcOverflow()
                this.getCurrentData()
            },
            highlightIndex () {
                this.updateScroller()
            }
        },
        created () {
            this.getCurrentData()
        },
        mounted () {
            this.calcOverflow()
        },
        methods: {
            async getCurrentData () {
                try {
                    if (this.exactSearchMethod) {
                        this.currentData = await this.exactSearchMethod(this.localValue)
                    } else {
                        console.warn('No exact search method has been set')
                    }
                } catch (error) {
                    console.error(error)
                }
            },
            // 搜索，滚动加载时next可设置为页码，无更多设置为false
            async search (value, next) {
                if (this.createOnly) return
                try {
                    const popoverInstance = this.getPopoverInstance()
                    const alternateContent = this.getAlternateContent()
                    popoverInstance.setContent(alternateContent.$el)
                    this.showPopover()
                    alternateContent.loading = !!value
                    const { results: data, next: nextPage } = await new Promise(async (resolve, reject) => {
                        if (value) {
                            const promise = [(this.fuzzySearchMethod || this.defaultFuzzySearchMethod)(value, next)]
                            if (this.searchFromDefaultAlternate) {
                                promise.push(this.getDefaultAlternateData(value))
                            }
                            const [fuzzySearchData, defaultAlternateData] = await Promise.all(promise)
                            if (defaultAlternateData) {
                                fuzzySearchData.results.unshift(...defaultAlternateData.results)
                            }
                            resolve(fuzzySearchData)
                        } else {
                            const defaultAlternateData = this.getDefaultAlternateData()
                            resolve(defaultAlternateData)
                        }
                    })
                    if (!this.isFocus) {
                        return
                    }
                    const { matched, flattened } = this.filterData(data)
                    if (!value && !flattened.length) {
                        this.hidePopover()
                        return
                    }
                    this.matchedData = next ? [...this.matchedData, ...matched] : matched
                    this.flattenedData = next ? [...this.flattenedData, ...flattened] : flattened
                    this.highlightIndex = flattened.length && !!this.inputValue ? 0 : -1
                    alternateContent.next = nextPage
                    alternateContent.keyword = value
                    alternateContent.matchedData = this.matchedData
                    alternateContent.loading = false
                } catch (e) {
                    if (e.type === 'reset') {
                        return
                    }
                    this.matchedData = []
                    this.flattenedData = []
                    console.error(e)
                }
            },
            async getDefaultAlternateData (keyword) {
                let data = []
                const isMatch = (tag, keyword) => tag.value.toLowerCase().indexOf(keyword.toString().toLowerCase()) > -1
                if (this.defaultAlternate === 'history') {
                    data = [{ display_name: this.historyLabel, children: this.getHistoryData() }]
                } else if (this.defaultAlternate instanceof Array) {
                    data = this.defaultAlternate
                } else if (typeof this.defaultAlternate === 'function') {
                    data = await this.defaultAlternate()
                }
                if (keyword) {
                    const filterResult = []
                    data.forEach(tag => {
                        if (tag.hasOwnProperty('children')) {
                            const children = tag.children.filter(child => isMatch(child, keyword))
                            if (children.length) {
                                filterResult.push({
                                    ...tag,
                                    children
                                })
                            }
                        } else if (isMatch(tag, keyword)) {
                            filterResult.push(tag)
                        }
                    })
                    data = filterResult
                }
                return Promise.resolve({ results: data, next: false })
            },
            // 过滤已被选择的
            filterData (data) {
                const matched = []
                const flattened = []
                data.forEach(tag => {
                    if (tag.hasOwnProperty('children')) {
                        const children = tag.children.filter(child => !flattened.some(flattenedTag => flattenedTag.value === child.value))
                        if (this.multiple) {
                            const unexistTag = children.filter(child => !this.localValue.includes(child.value))
                            if (unexistTag.length) {
                                tag.children = unexistTag
                                matched.push(tag)
                                flattened.push(...unexistTag)
                            }
                        } else {
                            matched.push(tag)
                            flattened.push(...children)
                        }
                        return
                    }
                    const exist = this.localValue.includes(tag.value)
                    const repeat = flattened.some(flattenedTag => flattenedTag.value === tag.value)
                    if ((!this.multiple || !exist) && !repeat) {
                        matched.push(tag)
                        flattened.push(tag)
                    }
                })
                return {
                    matched,
                    flattened
                }
            },
            // 默认模糊搜索方法，可通过fuzzySearchMethod自定义
            async defaultFuzzySearchMethod (value, next) {
                if (!this.fuzzySearchMethod) {
                    console.warn('No fuzzy search method has been set')
                    return Promise.resolve({ next: false, results: [] })
                }
                return this.fuzzySearchMethod(value, next)
            },
            // 已选择的标签tootips内容，可通过tagTipsContent自定义
            async getTagTips (instance, value) {
                try {
                    const contentElement = document.createElement('span')
                    if (typeof this.tagTipsContent === 'function') {
                        const content = await this.tagTipsContent(value)
                        contentElement.innerHTML = content
                    } else {
                        const tag = await (this.exactSearchMethod || this.defaultExactSearchMethod)(value)
                        contentElement.innerHTML = tag ? tag.text : 'Non existing tag'
                    }
                    instance.setContent(contentElement)
                } catch (e) {
                    console.error(e)
                    instance.setContent(e.message)
                }
            },
            // 默认精确搜索方法，通过此方法获取单个信息
            defaultExactSearchMethod (value) {
                if (!this.exactSearchMethod) {
                    console.warn('No exact search method has been set')
                    return Promise.resolve({})
                }
                return this.exactSearchMethod(value)
            },
            // 创建/获取备选面板popover实例
            getPopoverInstance () {
                if (!this.popoverInstance) {
                    this.popoverInstance = Tippy(this.$refs.input, {
                        theme: 'light tag-input-popover',
                        appendTo: document.body,
                        trigger: 'manual',
                        placement: 'bottom-start',
                        // distance: 5,
                        arrow: false,
                        hideOnClick: false,
                        content: '',
                        interactive: true,
                        boundary: 'window',
                        onHide: () => {
                            this.handlePopoverHide()
                        }
                    })
                }
                return this.popoverInstance
            },
            // 创建/获取备选面板内容
            getAlternateContent () {
                if (!this.alternateContent) {
                    this.alternateContent = new Vue(AlternateList)
                    this.alternateContent.tagInput = this
                    this.alternateContent.$mount()
                }
                return this.alternateContent
            },
            // 获取最近输入
            getHistoryData () {
                try {
                    if (this.historyKey) {
                        const data = JSON.parse(window.localStorage.getItem(this.historyKey)) || []
                        return data.filter(tag => !this.disabledData.includes(tag.value))
                    }
                    throw new Error('History key not provide')
                } catch (e) {
                    console.error(e)
                    return []
                }
            },
            // 更新最近输入
            updateHistoryData (tag) {
                if (this.historyKey) {
                    try {
                        const histories = this.getHistoryData()
                        const exist = histories.findIndex(history => history.value === tag.value)
                        if (exist > -1) {
                            histories.splice(exist, 1)
                        }
                        Array.isArray(tag) ? histories.unshift(...tag) : histories.unshift(tag)
                        const newHistories = histories.filter(history => !this.disabledData.includes(history.value)).slice(0, this.historyRecord)
                        window.localStorage.setItem(this.historyKey, JSON.stringify(newHistories))
                    } catch (e) {
                        console.error(e)
                    }
                }
            },
            // 更新面板定位
            updatePopover () {
                this.popoverInstance && this.popoverInstance.popperInstance && this.popoverInstance.popperInstance.update()
            },
            showPopover () {
                this.updatePopover()
                this.popoverInstance && this.popoverInstance.show(0)
            },
            hidePopover () {
                this.popoverInstance && this.popoverInstance.hide(0)
            },
            // 面板隐藏后重置匹配上的数据
            handlePopoverHide () {
                this.$nextTick(() => {
                    this.matchedData = []
                    this.flattenedData = []
                    this.alternateContent.matchedData = []
                })
            },
            // 显示的文本
            getDisplayText (tag) {
                const isObject = typeof tag === 'object'
                return isObject ? tag.text || tag.value : tag
            },
            // 元素整体点击, 获得焦点，初始化输入相关内容
            focus () {
                if (this.disabled) {
                    return false
                }
                this.clearOverflowTimer()
                this.inputIndex = this.localValue.length
                if (!this.multiple && this.value.length) {
                    this.inputValue = this.getDisplayText(this.value[0])
                    this.$refs.input.innerHTML = this.inputValue
                    this.moveInput(0, { selectRange: true })
                } else {
                    this.moveInput(0)
                }
            },
            // 选择面板显示时，禁止已选容器的滚动，防止面板滚动后错位
            handleContainerScroll (event) {
                this.popoverInstance && this.popoverInstance.state.isVisible && event.preventDefault()
            },
            // 已选择click事件拆分为mousedown,mouseup, 因为click事件晚于input的blur事件，但mousedown早于blur
            // shouldUpdate = false, 阻止input blur后续处理
            handleSelectedMousedown (event, index) {
                if (this.disabled) {
                    return false
                }
                this.shouldUpdate = false
            },
            // 根据点击的位置判断在点击的前方还是后方进行输入
            handleSelectedMouseup (event, index) {
                if (this.disabled) {
                    return false
                }
                if (this.multiple) {
                    const $referenceTarget = event.target
                    const offsetWidth = $referenceTarget.offsetWidth
                    const eventX = event.offsetX
                    // const $input = this.$refs.input
                    this.inputIndex = eventX > (offsetWidth / 2) ? index + 1 : index
                    this.moveInput(0)
                } else {
                    this.inputValue = this.getDisplayText(this.value[0])
                    this.$refs.input.innerHTML = this.inputValue
                    this.moveInput(0, { selectRange: true })
                }
            },
            // 已选的鼠标划过事件，如果配置了显示tips，则在一定延时后显示
            handleSelectedMouseenter (event, { value }) {
                if (!this.displayTagTips) return
                const target = event.currentTarget
                if (target._tag_tips_) {
                    return false
                }
                this.selectedTipsTimer[value] = setTimeout(() => {
                    target._tag_tips_ = Tippy(target, {
                        theme: 'light small-arrow tag-selected-tips',
                        boundary: 'window',
                        appendTo: document.body,
                        arrow: true,
                        content: 'loading...',
                        placement: 'top',
                        interactive: true,
                        onShow: instance => this.getTagTips(instance, value)
                    })
                    target._user_tips_.show()
                    delete this.selectedTipsTimer[value]
                }, this.tagTipsDelay)
            },
            // 在tips延时执行前离开，放弃执行
            handleSelectedMouseleave (event, { value }) {
                if (this.displayTagTips) {
                    this.selectedTipsTimer[value] && clearTimeout(this.selectedTipsTimer[value])
                }
            },
            // 删除单个mousedown,阻止input blur
            handleRemoveMouseDown () {
                this.shouldUpdate = false
            },
            // 删除单个
            handleRemoveSelected ({ value }, index) {
                if (this.disabled) {
                    return false
                }
                const localValue = [...this.localValue]
                localValue.splice(index, 1)
                this.localValue = localValue
                this.reset()
                if (this.isFocus) {
                    this.moveInput(index >= this.inputIndex ? 0 : -1)
                } else {
                    this.handleBlur()
                }
                this.$emit('remove-selected', value)
            },
            // 选择列表mousedown, 阻止input blur
            handleTagMousedown (tag, disabled) {
                this.shouldUpdate = false
            },
            // 选择列表，区分单选多选
            handleTagMouseup (tag, disabled) {
                if (disabled || this.disabled) {
                    this.moveInput(0)
                    return false
                }
                this.updateHistoryData(tag)
                this.currentData.push(tag)
                if (this.multiple) {
                    const localValue = [...this.localValue]
                    localValue.splice(this.inputIndex, 0, tag.value)
                    this.localValue = localValue
                    setTimeout(() => {
                        this.moveInput(1)
                        this.setSelection({ reset: true })
                        this.search()
                    }, 0)
                } else {
                    this.localValue = [tag.value]
                    this.reset()
                    this.handleBlur()
                }
                this.$emit('select-tag', tag)
            },
            // 选择列表分组 mousedown, 阻止input blur
            handleGroupMousedown () {
                this.shouldUpdate = false
            },
            // 点击选择列表分组会导致input 失去焦点，mouseup后重新获得焦点
            handleGroupMouseup () {
                this.moveInput(0)
            },
            // input 通过contenteditable模拟，已输入情况下，再次点击同一位置，光标会在文字最前方
            // 此方法多选时，将光标移至文字末尾，单选时，选中文本
            setSelection (option = {}) {
                if (option.reset) {
                    this.reset()
                }
                this.isFocus = true
                this.shouldUpdate = true
                this.$nextTick(() => {
                    const $input = this.$refs.input
                    $input.focus()
                    const range = window.getSelection()
                    range.selectAllChildren($input)
                    !option.selectRange && range.collapseToEnd()
                })
            },
            // 集中分配input 的键盘事件
            handleKeydown (event) {
                if (this.loading) {
                    event.preventDefault()
                    event.stopPropagation()
                    return
                }
                const key = event.key
                const keyMap = {
                    'Enter': this.handleEnter,
                    'Backspace': this.handleBackspace,
                    'Delete': this.handleBackspace,
                    'ArrowLeft': this.handleArrow,
                    'ArrowRight': this.handleArrow,
                    'ArrowUp': this.handleArrow,
                    'ArrowDown': this.handleArrow
                }
                if (keyMap.hasOwnProperty(key)) {
                    keyMap[key](event)
                }
                this.$emit('keydown', event)
            },
            // 通过键盘enter进行确认
            // 1.如果有备选高亮，则直接选中该项
            // 2.如果配置允许创建，允许输入不在备选列表中的，则进行相关值的设定
            handleEnter (event) {
                event.preventDefault()
                this.shouldUpdate = false
                if (this.highlightIndex !== -1) {
                    const value = this.flattenedData[this.highlightIndex]['value']
                    const disabled = this.disabledData.includes(value)
                    if (disabled) {
                        return false
                    }
                    if (this.multiple) {
                        const localValue = [...this.localValue]
                        localValue.splice(this.inputIndex, 0, value)
                        this.localValue = localValue
                        this.moveInput(1, { reset: true })
                    } else {
                        this.localValue = [value]
                        this.reset()
                        this.handleBlur()
                    }
                } else if (this.inputValue) {
                    const exist = this.localValue.includes(this.inputValue)
                    if (exist || !this.allowCreate) {
                        this.reset()
                    } else if (this.multiple) {
                        const localValue = [...this.localValue]
                        localValue.splice(this.inputIndex, 0, this.inputValue)
                        this.localValue = localValue
                        this.moveInput(1, { reset: true })
                    } else {
                        this.localValue = [this.inputValue]
                        this.reset()
                        this.handleBlur()
                    }
                } else {
                    this.reset()
                    this.handleBlur()
                }
                this.hidePopover()
            },
            // 删除事件，无文本输入时删除前一个
            handleBackspace (event) {
                if (this.inputValue || !this.localValue.length || !this.inputIndex) {
                    return true
                }
                this.shouldUpdate = false
                const localValue = [...this.localValue]
                localValue.splice(this.inputIndex - 1, 1)
                this.localValue = localValue
                this.moveInput(-1)
                this.search()
            },
            // 箭头事件
            // 1.左右事件处理为无文本时移动光标在已选中的位置
            // 2.上下事件处理备选高亮状态及备选列表的滚动位置
            handleArrow (event) {
                const arrow = event.key
                if (['ArrowLeft', 'ArrowRight'].includes(arrow)) {
                    if (this.inputValue || !this.localValue.length) {
                        return true
                    }
                    if (arrow === 'ArrowLeft' && this.inputIndex !== 0) {
                        this.moveInput(-1)
                    } else if (arrow === 'ArrowRight' && this.inputIndex !== this.localValue.length) {
                        this.moveInput(1)
                    }
                } else if (this.flattenedData.length) {
                    event.preventDefault()
                    if (arrow === 'ArrowDown') {
                        if (this.highlightIndex < (this.flattenedData.length - 1)) {
                            this.highlightIndex++
                        } else if (this.alternateContent && this.alternateContent.next) {
                            this.alternateContent.$refs.alternateList.scrollTop = this.alternateContent.$refs.alternateList.scrollTop + 32
                            this.alternateContent.handleScroll()
                        } else {
                            this.highlightIndex = 0
                        }
                    } else if (arrow === 'ArrowUp' && this.highlightIndex !== -1) {
                        this.highlightIndex--
                    }
                }
            },
            // 处理input输入
            handleInput (event) {
                if (this.loading) {
                    event.preventDefault()
                    event.stopPropagation()
                    return
                }
                this.inputValue = this.$refs.input.textContent.trim()
            },
            // 处理input blur事件，在其他元素的mousedown事件中控制shouldUpdate，达到阻止后续逻辑的能力
            // blur时，如果允许输入不存在备选列表中的值，则以当前值填充，否则匹配备选列表
            handleBlur () {
                if (!this.shouldUpdate) {
                    return true
                }
                this.isFocus = false
                this.hidePopover()
            },
            // 获取匹配当前输入文本的数据
            getMatchedTag (nameToMatch) {
                const tag = this.flattenedData.find(tag => {
                    const value = tag.value
                    const text = tag.text
                    const isMatch = [value, text].some(name => name.toLowerCase() === nameToMatch.toLowerCase())
                    const isSelected = this.localValue.includes(value)
                    return isMatch && !isSelected
                })
                return tag
            },
            // 粘贴事件，逗号或者分号分隔
            async handlePaste (event) {
                this.hidePopover()
                if (this.loading) {
                    event.preventDefault()
                    event.stopPropagation()
                }
                try {
                    this.loading = true
                    const pasteStr = event.clipboardData.getData('text').trim()
                    const values = pasteStr.split(/,|;|\n/).map(value => this.pasteFormatter(value)).filter(value => value.length)
                    const uniqueValues = [...new Set(values)]
                    if (!uniqueValues.length) {
                        return
                    }
                    const validValues = await this.pasteValidator(uniqueValues)
                    const newValues = validValues.filter(value => !this.localValue.includes(value))
                    if (!validValues.length) {
                        return
                    }
                    const localValue = [...this.localValue]
                    localValue.splice(this.inputIndex, 0, ...newValues)
                    this.localValue = localValue
                    this.isFocus && this.moveInput(newValues.length, { reset: true })
                } catch (error) {
                    console.error(error)
                } finally {
                    this.loading = false
                }
            },
            getSelectedDOM () {
                return Array.from(this.$refs.container.querySelectorAll('.tag-input-selected'))
            },
            // 移动光标位置
            moveInput (step, option = {}) {
                this.inputIndex = this.inputIndex + step
                this.$nextTick(() => {
                    const selected = this.getSelectedDOM()
                    const $referenceTarget = selected[this.inputIndex] || null
                    this.$refs.container.insertBefore(this.$refs.input, $referenceTarget)
                    this.setSelection(option)
                    this.updatePopover()
                })
            },
            // 设置备选高亮时，当前备选项不在视图内时，进行滚动，到底/顶后循环滚动
            updateScroller () {
                if (!this.alternateContent) {
                    return false
                }
                this.$nextTick(() => {
                    const highlightIndex = this.highlightIndex
                    const $alternateList = this.alternateContent.$refs.alternateList
                    if (!$alternateList) {
                        return false
                    }
                    if (highlightIndex !== -1) {
                        const $alternateItem = this.alternateContent.$refs.alternateItem[highlightIndex].$el
                        const listClientHeight = $alternateList.clientHeight
                        const listScrollTop = $alternateList.scrollTop
                        const itemOffsetTop = $alternateItem.offsetTop
                        const itemOffsetHeight = $alternateItem.offsetHeight
                        if (itemOffsetTop >= listScrollTop && (itemOffsetTop + itemOffsetHeight) <= (listScrollTop + listClientHeight)) {
                            return false
                        } else if (itemOffsetTop <= listScrollTop) {
                            $alternateList.scrollTop = itemOffsetTop
                        } else if ((itemOffsetTop + itemOffsetHeight) > (listScrollTop + listClientHeight)) {
                            $alternateList.scrollTop = itemOffsetTop + itemOffsetHeight - listClientHeight
                        }
                    } else {
                        $alternateList.scrollTop = 0
                    }
                })
            },
            // 计算第二行第一个的index，在其前方插入overflow tag
            calcOverflow () {
                this.removeOverflowTagNode()
                if (!this._isMounted || !this.fixedHeight || this.isFocus || this.localValue.length < 2) {
                    return false
                }
                this.clearOverflowTimer()
                this.overflowTimer = setTimeout(() => {
                    const selectedData = this.getSelectedDOM()
                    const userIndexInSecondRow = selectedData.findIndex((currentTag, index) => {
                        if (!index) return
                        const previousTag = selectedData[index - 1]
                        return previousTag.offsetTop !== currentTag.offsetTop
                    })
                    if (userIndexInSecondRow > -1) {
                        this.overflowTagIndex = userIndexInSecondRow
                    } else {
                        this.overflowTagIndex = null
                    }
                    this.$refs.container.scrollTop = 0
                    this.insertOverflowTag()
                }, 0)
            },
            clearOverflowTimer () {
                this.overflowTimer && clearTimeout(this.overflowTimer)
            },
            // 根据计算的overflow index，插入tag并进行校正
            insertOverflowTag () {
                if (!this.overflowTagIndex) return
                const overflowTagNode = this.getOverflowTagNode()
                const selectedTag = this.getSelectedDOM()
                const referenceTag = selectedTag[this.overflowTagIndex]
                if (referenceTag) {
                    overflowTagNode.textContent = `+${this.localValue.length - this.overflowTagIndex}`
                    this.$refs.container.insertBefore(overflowTagNode, referenceTag)
                } else {
                    this.overflowTagIndex = null
                    return
                }
                setTimeout(() => {
                    const previousTag = selectedTag[this.overflowTagIndex - 1]
                    if (overflowTagNode.offsetTop !== previousTag.offsetTop) {
                        this.overflowTagIndex--
                        this.$refs.container.insertBefore(overflowTagNode, overflowTagNode.previousSibling)
                        overflowTagNode.textContent = `+${this.localValue.length - this.overflowTagIndex}`
                    }
                }, 0)
            },
            // 创建/获取溢出数字节点
            getOverflowTagNode () {
                if (this.overflowTagNode) {
                    return this.overflowTagNode
                }
                const overflowTagNode = document.createElement('span')
                overflowTagNode.setAttribute(this.$options._scopeId, '')
                overflowTagNode.className = 'tag-input-overflow-tag'
                this.overflowTagNode = overflowTagNode
                return overflowTagNode
            },
            // 从容器中移除溢出数字节点
            removeOverflowTagNode () {
                if (this.overflowTagNode && this.overflowTagNode.parentNode === this.$refs.container) {
                    this.$refs.container.removeChild(this.overflowTagNode)
                }
            },
            // 一键清空
            handleFastClear () {
                this.localValue = []
                this.$emit('clear')
            },
            // 重置选择器输入状态
            reset () {
                this.shouldUpdate = true
                this.highlightIndex = -1
                this.inputValue = ''
                this.$refs.input.innerHTML = ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    @import './style.scss';
</style>

<style lang="scss">
    .tippy-tooltip.tag-input-popover-theme {
        border: 1px solid #dcdee5;
        border-radius: 2px;
        box-shadow: 0px 2px 6px 0px rgba(0, 0, 0, 0.1);
        color: #63656E;
        font-size: 12px;
        line-height: 24px;
    }
</style>
