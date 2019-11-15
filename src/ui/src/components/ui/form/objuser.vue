<template>
    <!-- eslint-disable vue/space-infix-ops -->
    <div class="cmdb-form form-objuser"
        v-click-outside="handleClickOutside"
        @mousedown="shouldUpdate = false"
        @click="handleClick">
        <!--eslint-enable-->
        <div class="objuser-layout"
            @contextmenu="handleContextmenu($event)">
            <i class="objuser-loading" v-if="loading"></i>
            <div class="objuser-container"
                ref="container"
                :class="{
                    focus,
                    ellipsis,
                    disabled: localDisabled,
                    placeholder: !localValue.length && !focus
                }"
                :data-placeholder="localPlaceholder">
                <span class="objuser-selected"
                    v-for="(user, index) in selectedUsers"
                    ref="selected"
                    :key="index"
                    @click.stop
                    @mousedown.left.stop="handleSelectedMousedown($event, index)"
                    @mouseup.left.stop="handleSelectedMouseup($event, index)">
                    {{getLable(user)}}
                </span>
                <span ref="input" class="objuser-input"
                    spellcheck="false"
                    contenteditable
                    v-show="focus"
                    @click.stop
                    @input="handleInput"
                    @blur="handleBlur"
                    @paste="handlePaste"
                    @keydown="handleKeydown($event)">
                </span>
            </div>
            <ul class="suggestion-list" ref="suggestionList"
                v-show="focus && matchedUsers.length"
                :class="suggestionListPostion">
                <li class="suggestion-item"
                    v-for="(user, index) in matchedUsers"
                    :key="index"
                    ref="suggestionItem"
                    :title="getLable(user)"
                    :class="{ highlight: index === highlightIndex }"
                    @click.stop
                    @mousedown.left.stop="handleUserMousedown(user, index)"
                    @mouseup.left.stop="handleUserMouseup(user, index)">
                    {{getLable(user)}}
                </li>
            </ul>
            <a href="javascript:void(0)" class="objuser-menu"
                ref="contextmenu"
                v-show="contextmenu"
                @click.stop="handleCopy">{{$t('复制')}}</a>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-form-objuser',
        props: {
            value: {
                type: String,
                default: ''
            },
            placeholder: {
                type: String,
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            multiple: {
                type: Boolean,
                default: true
            },
            exclude: {
                type: Boolean,
                default: true
            }
        },
        data () {
            return {
                users: [],
                localValue: [],
                inputValue: '',
                inputIndex: 0,
                highlightIndex: -1,
                shouldUpdate: true,
                focus: false,
                ellipsis: false,
                contextmenu: false,
                exception: false,
                updateTimer: null,
                suggestionListPostion: 'bottom'
            }
        },
        computed: {
            selectedUsers () {
                return this.localValue.map(enName => this.users.find(user => user['english_name'] === enName) || { 'english_name': enName })
            },
            matchedUsers () {
                const inputValue = this.inputValue.toLowerCase()
                if (inputValue) {
                    return this.users.filter(user => {
                        const enName = user['english_name']
                        const cnName = user['chinese_name']
                        const isEnMatch = enName.toLowerCase().indexOf(inputValue) !== -1
                        const isCnMatch = cnName.toLowerCase().indexOf(inputValue) !== -1
                        const isSelected = this.localValue.includes(enName)
                        return (isEnMatch || isCnMatch) && !isSelected
                    })
                }
                return []
            },
            loading () {
                return this.$loading('get_user_list')
            },
            localDisabled () {
                return this.disabled || this.exception || this.loading
            },
            localPlaceholder () {
                if (this.exception) {
                    return this.$t('获取人员列表失败')
                } else if (this.placeholder) {
                    return this.placeholder
                }
                return this.$t('请输入用户')
            }
        },
        watch: {
            focus (focus) {
                if (this.focus) {
                    this.ellipsis = false
                } else {
                    this.reset()
                    this.calcEllipsis()
                }
            },
            value (value) {
                if (this.localValue.join(',') !== value) {
                    this.setLocalValue()
                }
            },
            localValue (localValue, oldValue) {
                const localValueStr = localValue.join(',')
                if (localValueStr !== this.value) {
                    this.$emit('input', localValueStr)
                    this.$emit('on-change', localValueStr, oldValue.join(','))
                }
            },
            matchedUsers (matchedUsers) {
                this.highlightIndex = -1
                if (matchedUsers.length) {
                    if (this.exclude) {
                        this.highlightIndex = 0
                    }
                    this.updateSuggestionPosition()
                }
            },
            highlightIndex () {
                this.updateScroller()
            }
        },
        async created () {
            try {
                this.users = await this.$store.dispatch('getUserList')
            } catch (e) {
                this.exception = true
                this.users = []
            }
            this.setLocalValue()
        },
        mounted () {
            this.calcEllipsis()
        },
        methods: {
            setLocalValue () {
                const value = (this.value || '').split(',')
                const localValue = []
                value.forEach(enName => {
                    if (this.users.some(user => user['english_name'] === enName)) {
                        localValue.push(enName)
                    }
                })
                this.localValue = localValue
            },
            calcEllipsis () {
                this.$nextTick(() => {
                    const $selected = this.$refs.selected
                    if ($selected && $selected.length) {
                        const $container = this.$refs.container
                        const $lastSelected = this.$refs.selected[$selected.length - 1]
                        const lastSelectedWidth = $lastSelected.offsetWidth
                        const lastSelectedLeft = $lastSelected.offsetLeft
                        const containerWidth = $container.offsetWidth
                        this.ellipsis = (lastSelectedWidth + lastSelectedLeft) > containerWidth
                    }
                })
            },
            getLable (user) {
                const enName = user['english_name']
                const cnName = user['chinese_name']
                if (enName && cnName) {
                    return `${user['english_name']}(${user['chinese_name']})`
                }
                return enName
            },
            handleClickOutside () {
                this.contextmenu = false
            },
            handleClick () {
                if (this.localDisabled) {
                    return false
                }
                if (!this.multiple) {
                    this.localValue = []
                }
                this.inputIndex = this.localValue.length
                this.moveInput(0)
            },
            handleSelectedMousedown (event, index) {
                if (this.localDisabled) {
                    return false
                }
                this.shouldUpdate = false
            },
            handleSelectedMouseup (event, index) {
                if (this.localDisabled) {
                    return false
                }
                if (this.multiple) {
                    const $refrenceTarget = event.target
                    const offsetWidth = $refrenceTarget.offsetWidth
                    const eventX = event.offsetX
                    // const $input = this.$refs.input
                    this.inputIndex = eventX > (offsetWidth / 2) ? index + 1 : index
                } else {
                    this.localValue = []
                    this.inputIndex = 0
                }
                this.moveInput(0)
            },
            handleUserMousedown (user, index) {
                if (this.localDisabled) {
                    return false
                }
                this.shouldUpdate = false
            },
            handleUserMouseup (user, index) {
                if (this.localDisabled) {
                    return false
                }
                if (this.multiple) {
                    this.localValue.splice(this.inputIndex, 0, user['english_name'])
                    this.$nextTick(() => {
                        this.moveInput(1)
                        this.setSelection({ reset: true })
                    })
                } else {
                    this.localValue = [user['english_name']]
                    this.reset()
                    this.handleBlur()
                }
            },
            setSelection (option = {}) {
                if (option.reset) {
                    this.reset()
                }
                this.focus = true
                this.shouldUpdate = true
                this.$nextTick(() => {
                    const $input = this.$refs.input
                    $input.focus()
                    if (window.getSelection) {
                        const range = window.getSelection()
                        range.selectAllChildren($input)
                        range.collapseToEnd()
                    } else if (document.selection) {
                        const range = document.selection.createRange()
                        range.moveToElementText($input)
                        range.collapse(false)
                        range.select()
                    }
                })
            },
            handleKeydown (event) {
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
            },
            handleEnter (event) {
                event.preventDefault()
                this.shouldUpdate = false
                if (this.highlightIndex !== -1) {
                    if (this.multiple) {
                        this.localValue.splice(this.inputIndex, 0, this.matchedUsers[this.highlightIndex]['english_name'])
                        this.moveInput(1, { reset: true })
                    } else {
                        this.localValue = [this.matchedUsers[this.highlightIndex]['english_name']]
                        this.reset()
                        this.handleBlur()
                    }
                } else if (this.inputValue) {
                    if (!this.exclude && !this.localValue.includes(this.inputValue)) {
                        if (this.multiple) {
                            this.localValue.splice(this.inputIndex, 0, this.inputValue)
                            this.moveInput(1, { reset: true })
                        } else {
                            this.localValue = [this.inputValue]
                            this.reset()
                            this.handleBlur()
                        }
                    } else {
                        this.reset()
                    }
                } else {
                    this.handleBlur()
                }
            },
            handleBackspace (event) {
                if (this.inputValue || !this.localValue.length || !this.inputIndex) {
                    return true
                }
                this.localValue.splice(this.inputIndex - 1, 1)
                this.moveInput(-1)
            },
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
                } else if (this.matchedUsers.length) {
                    event.preventDefault()
                    if (arrow === 'ArrowDown' && this.highlightIndex < (this.matchedUsers.length - 1)) {
                        this.highlightIndex++
                    } else if (arrow === 'ArrowUp' && this.highlightIndex !== -1) {
                        this.highlightIndex--
                    }
                }
            },
            handleInput () {
                this.inputValue = this.$refs.input.textContent.trim()
            },
            handleBlur () {
                if (!this.shouldUpdate) {
                    return true
                }
                this.focus = false
                if (this.inputValue) {
                    if (!this.exclude) {
                        if (!this.localValue.includes(this.inputValue)) {
                            this.localValue.splice(this.inputIndex, 0, this.inputValue)
                        }
                    } else {
                        const matchedUser = this.getMatchedUser(this.inputValue)
                        if (matchedUser) {
                            this.localValue.splice(this.inputIndex, 0, matchedUser['english_name'])
                        }
                    }
                }
            },
            handlePaste () {
                this.$nextTick(() => {
                    const values = [...new Set(this.inputValue.split(/,|;/))]
                    const pasteValue = []
                    values.forEach(name => {
                        const user = this.getMatchedUser(name)
                        if (user) {
                            pasteValue.push(user['english_name'])
                        }
                    })
                    this.localValue.splice(this.inputIndex, 0, ...pasteValue)
                    this.moveInput(pasteValue.length, { reset: true })
                })
            },
            handleCopy () {
                this.contextmenu = false
                this.$copyText(this.localValue.join(',')).then(() => {
                    this.$success(this.$t('复制成功'))
                }, () => {
                    this.$error(this.$t('复制失败'))
                })
            },
            handleContextmenu (event) {
                this.focus = false
                if (!this.localValue.length) {
                    return false
                }
                event.preventDefault()
                event.stopPropagation()
                const $contextmenu = this.$refs.contextmenu
                const $layout = event.currentTarget
                let $refrence = $layout
                let invisible = false
                while (!invisible && $refrence.nodeName !== 'BODY') {
                    $refrence = $refrence.parentElement
                    const overflow = window.getComputedStyle($refrence).getPropertyValue('overflow-y')
                    invisible = ['hidden', 'auto', 'scroll'].includes(overflow)
                }
                this.contextmenu = true
                this.$nextTick(() => {
                    const menuRect = $contextmenu.getBoundingClientRect()
                    const refrenceRect = $refrence.getBoundingClientRect()
                    const layoutRect = $layout.getBoundingClientRect()
                    let left = event.x - layoutRect.left
                    let top = event.y - layoutRect.top
                    if (event.x + menuRect.width + 5 < refrenceRect.right) {
                        left = left + 5
                    } else {
                        left = left - menuRect.width
                    }
                    if (event.y + menuRect > refrenceRect.bottom) {
                        top = top - menuRect.height
                    }
                    $contextmenu.style.left = left + 'px'
                    $contextmenu.style.top = top + 'px'
                })
            },
            getMatchedUser (nameToMatch) {
                const user = this.users.find(user => {
                    const enName = user['english_name']
                    const cnName = user['chinese_name']
                    const isMatch = [enName, cnName].some(name => name.toLowerCase() === nameToMatch.toLowerCase())
                    const isSelected = this.localValue.includes(enName)
                    return isMatch && !isSelected
                })
                return user
            },
            moveInput (step, option = {}) {
                this.$nextTick(() => {
                    this.inputIndex = this.inputIndex + step
                    const $refrenceTarget = this.$refs.selected ? this.$refs.selected[this.inputIndex] : null
                    this.$refs.container.insertBefore(this.$refs.input, $refrenceTarget)
                    this.setSelection(option)
                })
            },
            updateScroller () {
                this.$nextTick(() => {
                    const highlightIndex = this.highlightIndex
                    const $suggestionList = this.$refs.suggestionList
                    if (highlightIndex !== -1) {
                        const $suggestionItem = this.$refs.suggestionItem[highlightIndex]
                        const listClientHeight = $suggestionList.clientHeight
                        const listScrollTop = $suggestionList.scrollTop
                        // const listScrollHeight = $suggestionList.scrollHeight
                        const itemOffsetTop = $suggestionItem.offsetTop
                        const itemOffsetHeight = $suggestionItem.offsetHeight
                        if (itemOffsetTop >= listScrollTop && (itemOffsetTop + itemOffsetHeight) <= (listScrollTop + listClientHeight)) {
                            return false
                        } else if (itemOffsetTop <= listScrollTop) {
                            $suggestionList.scrollTop = itemOffsetTop
                        } else if ((itemOffsetTop + itemOffsetHeight) > (listScrollTop + listClientHeight)) {
                            $suggestionList.scrollTop = itemOffsetTop + itemOffsetHeight - listClientHeight
                        }
                    } else {
                        $suggestionList.scrollTop = 0
                    }
                })
            },
            updateSuggestionPosition () {
                this.updateTimer && clearTimeout(this.updateTimer)
                this.updateTimer = setTimeout(() => {
                    const $suggestionList = this.$refs.suggestionList
                    let $refrence = $suggestionList
                    let invisible = false
                    while (!invisible && $refrence.nodeName !== 'BODY') {
                        $refrence = $refrence.parentElement
                        const overflow = window.getComputedStyle($refrence).getPropertyValue('overflow-y')
                        invisible = ['hidden', 'auto', 'scroll'].includes(overflow)
                    }
                    const refrenceRect = $refrence.getBoundingClientRect()
                    const suggestionListRect = $suggestionList.getBoundingClientRect()
                    if (suggestionListRect.bottom + suggestionListRect.height < refrenceRect.bottom) {
                        this.suggestionListPostion = 'bottom'
                    } else {
                        this.suggestionListPostion = 'top'
                    }
                }, 0)
            },
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
    .form-objuser {
        height: 32px;
        font-size: 14px;
        cursor: text;
        .objuser-layout {
            position: relative;
            min-height: 100%;
            .objuser-loading {
                position: absolute;
                top: 8px;
                right: 8px;
                width: 16px;
                height: 16px;
                background-image: url("../../../assets/images/icon/loading.svg");
                z-index: 1;
            }
            .objuser-container {
                position: relative;
                min-width: 100%;
                min-height: 32px;
                padding: 3px 0;
                line-height: 1;
                border: 1px solid #c4c6cc;
                border-radius: 2px;
                background-color: #fff;
                white-space: nowrap;
                overflow: hidden;
                &.disabled {
                    cursor: not-allowed;
                    background-color: #fafbfd !important;
                    border-color: #dcdee5 !important
                }
                &.focus {
                    white-space: normal;
                    border-color: $cmdbBorderFocusColor;
                    z-index: 1;
                }
                &.ellipsis:after{
                    position: absolute;
                    bottom: 1px;
                    right: -1px;
                    height: 34px;
                    padding: 0 0 0 15px;
                    line-height: 34px;
                    font-size: 12px;
                    content: "";
                    border-right: 1px solid $cmdbBorderColor;
                    background: -webkit-linear-gradient(left, transparent, #fff 55%);
                    background: -o-linear-gradient(left, transparent, #fff 55%);
                    background: linear-gradient(to right, transparent, #fff 55%);
                }
                &.placeholder:after {
                    position: absolute;
                    left: 0;
                    top: 0;
                    height: 100%;
                    padding: 0 0 0 10px;
                    line-height: 30px;
                    content: attr(data-placeholder);
                    font-size: 12px;
                    color: #c3cdd7;
                }
            }
        }
    }
    .objuser-selected {
        display: inline-block;
        height: 20px;
        margin: 2px 3px;
        max-width: calc(100% - 4px);
        padding: 0 4px;
        line-height: 18px;
        vertical-align: top;
        border: 1px solid #d9d9d9;
        border-radius: 2px;
        cursor: default;
        @include ellipsis;
    }
    .objuser-input {
        display: inline-block;
        max-width: 100%;
        height: 20px;
        margin: 1px 0 0;
        padding: 0 4px;
        white-space: nowrap;
        line-height: 20px;
        vertical-align: top;
        outline: none;
        overflow: hidden;
    }
    .suggestion-list{
        position: absolute;
        left: 0;
        width: 100%;
        max-height: 162px;
        
        font-size: 14px;
        background: #fff;
        box-shadow: 0 0 1px 1px rgba(0, 0, 0, 0.1);
        border: 1px solid #c3cdd7;
        border-radius: 2px;
        z-index: 1;
        @include scrollbar-y;
        &.bottom {
            top: 100%;
            margin: 1px 0 0 0;
        }
        &.top {
            bottom: 100%;
            margin: 0 0 1px 0;
        }
        .suggestion-item {
            padding: 0 10px;
            height: 32px;
            line-height: 32px;
            cursor: pointer;
            @include ellipsis;
            &.highlight,
            &:hover{
                background-color: #f1f7ff;
            }
        }
    }
    .objuser-menu {
        position: absolute;
        top: 0;
        left: 0;
        padding: 5px 10px;
        background-color: #fff;
        box-shadow: 0 0 1px 1px rgba(0, 0, 0, 0.1);
        font-size: 14px;
        white-space: nowrap;
        z-index: 9999;
    }
</style>
