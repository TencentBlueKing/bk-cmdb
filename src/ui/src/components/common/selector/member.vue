<template>
    <div  v-click-outside="memberInputReset"
        :data-placeholder="$t(`Common['${placeholder}']`)"
        :class="['member-selector', {'active': focus, 'placeholder': !focus && !localSelected.length, 'disabled': disabled}]">
        <div class="member-wrapper">
            <div ref="memberContainer" :class="['member-container', {'active': focus , 'ellipsis': showEllipsis}]" @click="handleSelectorClick(localSelected.length)">
                <span ref="memberSelected" class="member-selected" 
                    v-for="(selectedMember, index) in localSelected"
                    :title="selectedMember"
                    @click.stop="handleSelectorClick(index, $event)">
                    <span class="member-input-emitter before" @click.stop="locateInputPosition(index)"></span>
                    {{selectedMember}}
                    <span class="member-input-emitter after" @click.stop="locateInputPosition(index + 1)"></span>
                </span>
                <span ref="memberInput"  contenteditable class="member-input" spellcheck="false"
                    v-show="focus"
                    :data-shadow-member="shadowMember"
                    @click.stop
                    @keyup="setMemberInputText($event)"
                    @keydown="handleKeydown($event)">
                </span>
            </div>
            <transition name="toggle-slide">
                <ul ref="memberList" class="member-list" v-show="focus && filterMembers.length">
                    <li ref="memberItem" v-for="(member, index) in filterMembers"
                        :class="['member-item', {'highlight': selectedIndex === index}]" 
                        :title="getLable(member)"
                        @click.stop="handleMemberItemClick(member)">
                        {{getLable(member)}}
                    </li>
                </ul>
            </transition>
        </div>
    </div>
</template>
<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            placeholder: {
                type: String,
                default: '请选择'
            },
            exclude: {
                type: Boolean,
                default: false
            },
            selected: {
                validator (selected) {
                    return typeof selected === 'string' || Array.isArray(selected) || selected === undefined
                }
            },
            disabled: {
                type: Boolean,
                default: false
            },
            multiple: {
                type: Boolean,
                default: true
            },
            visible: {
                type: Boolean,
                default: true
            }
        },
        data () {
            return {
                focus: false,
                localSelected: [],
                memberInputText: '',
                inputIndex: null,
                selectedIndex: null,
                filterMembers: [],
                showEllipsis: false
            }
        },
        computed: {
            ...mapGetters({
                'members': 'memberList',
                'memberLoading': 'memberLoading'
            }),
            /* 跟随在光标后的提示人员 */
            shadowMember () {
                if (this.selectedIndex !== null) {
                    let highlightMemberName = this.filterMembers[this.selectedIndex]['english_name']
                    return highlightMemberName.startsWith(this.memberInputText) ? highlightMemberName.substring(this.memberInputText.length) : ''
                }
                return ''
            }
        },
        watch: {
            selected (selected) {
                this.setLocalSelected()
            },
            localSelected (localSelected) {
                let localSelectedStr = localSelected.join(',')
                if (typeof this.selected === 'string' && localSelectedStr !== this.selected) {
                    this.$emit('update:selected', localSelectedStr)
                } else if (Array.isArray(this.selected) && this.selected.join(',') !== this.localSelectedStr) {
                    this.$emit('update:selected', localSelected)
                } else if (this.selected === undefined) {
                    this.$emit('update:selected', localSelectedStr)
                }
                this.setFilterMember() // 更新人员列表，已选择的不显示在人员列表中
                this.calcEllipsis() // 计算是否溢出
                this.calcEmmiter()  // 计算已选人员tag前后光标定位元素的宽度
            },
            memberInputText (memberInputText) {
                this.setFilterMember()
                this.$refs.memberList.scrollTop = 0
            },
            focus (focus) {
                this.calcEmmiter() // 获得焦点时，选择器会展开, 需要重新计算已选人员tag前后光标定位元素的宽度
                if (focus) {
                    this.updateMemberListPosition()
                } else {
                    this.selectedIndex = null
                    this.$refs.memberList.scrollTop = 0
                }
            },
            selectedIndex (selectedIndex) {
                // 根据当前上下移动选择的索引计算列表滚动
                let memberListItemHeight = 32
                let scrollCount = 4
                if (selectedIndex !== null && selectedIndex > scrollCount) {
                    this.$refs.memberList.scrollTop = (selectedIndex - scrollCount) * memberListItemHeight
                } else {
                    this.$refs.memberList.scrollTop = 0
                }
                this.updateMemberListPosition()
            },
            members (members) {
                this.setFilterMember()
            },
            active (active) {
                this.calcEllipsis()
                this.calcEmmiter()
            }
        },
        created () {
            if (!this.members.length && !this.memberLoading) {
                this.getMemberList()
            }
            this.setLocalSelected()
            this.setFilterMember()
        },
        methods: {
            ...mapActions(['getMemberList']),
            /* 设置本地存储数据 */
            setLocalSelected () {
                let selected = this.selected
                let localSelected = [...this.localSelected]
                if (typeof selected === 'string' && localSelected.join(',') !== selected) {
                    localSelected = !selected ? [] : selected.split(',')
                } else if (Array.isArray(selected) && selected.join(',') !== localSelected.join(',')) {
                    localSelected = [...selected]
                } else if (selected === undefined && localSelected.length) {
                    localSelected = []
                }
                if (!this.exclude) {
                    localSelected = localSelected.filter(selected => this.members.some(({english_name: englishName}) => englishName === selected))
                }
                this.localSelected = localSelected
            },
            /* 根据当前输入筛选人员列表 */
            setFilterMember () {
                let filterVal = this.memberInputText.toLowerCase()
                this.selectedIndex = null
                this.filterMembers = this.members.filter(member => {
                    let enInclude = member['english_name'].toLowerCase().indexOf(filterVal) !== -1
                    let cnInclude = member['chinese_name'].toLowerCase().indexOf(filterVal) !== -1
                    let isSelected = this.localSelected.includes(member['english_name'])
                    return (enInclude || cnInclude) && !isSelected
                })
            },
            /* 计算是否显示省略符号 */
            calcEllipsis () {
                this.$nextTick(() => {
                    let $memberSelected = this.$refs.memberSelected
                    let selectedMargin = 8
                    let containerPadding = 8
                    let memberSelectedWidth = 0
                    if ($memberSelected && $memberSelected.length) {
                        $memberSelected.forEach($selected => {
                            memberSelectedWidth = memberSelectedWidth + $selected.offsetWidth
                        })
                        memberSelectedWidth = memberSelectedWidth + $memberSelected.length * selectedMargin
                        this.showEllipsis = memberSelectedWidth > (this.$refs.memberContainer.offsetWidth - containerPadding)
                    } else {
                        this.showEllipsis = false
                    }
                })
            },
            /* 计算每个已选人员后面的输入定位元素的宽度 */
            calcEmmiter () {
                this.$nextTick(() => {
                    let $memberSelected = this.$refs.memberSelected
                    if ($memberSelected) {
                        let rows = {}
                        let selectedMarginTop = 8
                        let selectedHeight = 16
                        let memberInputDefaultWidth = 2
                        let rowOffsetTop = selectedMarginTop + selectedHeight
                        // 将已选人员的DOM按行换分
                        $memberSelected.forEach($selected => {
                            let row = ($selected.offsetTop - selectedMarginTop) / rowOffsetTop + 1
                            if (rows.hasOwnProperty(row)) {
                                rows[row].push($selected)
                            } else {
                                rows[row] = [$selected]
                            }
                        })
                        // 将每行的最后一个已选人员的光标定位宽度设置为该行的剩余宽度
                        for (let row in rows) {
                            let $rowSelected = rows[row]
                            let $lastRowSelected = $rowSelected.pop()
                            let afterEmitterWidth = this.$refs.memberContainer.offsetWidth - ($lastRowSelected.offsetLeft + $lastRowSelected.offsetWidth) - memberInputDefaultWidth
                            $lastRowSelected.querySelector('.member-input-emitter.after').style.width = `${afterEmitterWidth}px`
                            $rowSelected.forEach($selected => {
                                $selected.querySelectorAll('.member-input-emitter').forEach($emitter => {
                                    $emitter.style.width = `${selectedMarginTop}px`
                                })
                            })
                        }
                    }
                })
            },
            updateMemberListPosition () {
                this.$nextTick(() => {
                    this.$refs.memberList.style.top = `${this.$refs.memberContainer.offsetHeight}px`
                })
            },
            /* 激活人员选择器，定位光标位置，如果是点击的已选人员，未超过一半，设置在其前，否则在其后 */
            handleSelectorClick (index, event) {
                if (!this.disabled) {
                    if (event) {
                        let offsetWidth = event.target.offsetWidth
                        let eventX = event.offsetX
                        index = eventX > (offsetWidth / 2) ? index + 1 : index
                    }
                    this.locateInputPosition(index)
                }
            },
            /* 移动模拟输入元素的DOM位置并获取焦点 */
            locateInputPosition (index) {
                if (!this.disabled) {
                    if (!this.multiple) {
                        index = 0
                        this.inputIndex = 0
                        this.localSelected = []
                    }
                    let $memberInput = this.$refs.memberInput
                    let $memberSelected = this.$refs.memberSelected
                    if (index !== this.inputIndex) {
                        let $refrenceElement = $memberSelected && $memberSelected[index] ? $memberSelected[index] : null
                        this.inputIndex = $memberSelected && $memberSelected.length ? index : null
                        this.$refs.memberContainer.insertBefore($memberInput, $refrenceElement)
                        this.memberInputFocus()
                    } else {
                        this.setCaretPosition()
                    }
                }
            },
            /* 同样的位置，将光标定位至最后 */
            setCaretPosition () {
                let $memberInput = this.$refs.memberInput
                this.focus = true
                this.$nextTick(() => {
                    $memberInput.focus()
                    if (window.getSelection) {
                        let range = window.getSelection()
                        range.selectAllChildren($memberInput)
                        range.collapseToEnd()
                    } else if (document.selection) {
                        let range = document.selection.createRange()
                        range.moveToElementText($memberInput)
                        range.collapse(false)
                        range.select()
                    }
                })
            },
            /* 处于输入状态时对不同按键的响应 */
            handleKeydown (event) {
                let eventKey = event.key
                let keyFunc = {
                    'Enter': this.handleConfirm,
                    'Backspace': this.handleBackspace,
                    'Delete': this.handleBackspace,
                    'ArrowLeft': this.handleArrow,
                    'ArrowRight': this.handleArrow,
                    'ArrowDown': this.handleArrow,
                    'ArrowUp': this.handleArrow
                }
                if (keyFunc.hasOwnProperty(eventKey)) {
                    keyFunc[eventKey](event)
                }
            },
            /* 按下回车确认输入，添加人员 */
            handleConfirm (event) {
                event.preventDefault()
                let memberInputText = this.memberInputText
                if (memberInputText.length) {
                    let member = this.exclude ? {'english_name': memberInputText} : this.filterMembers.find(({english_name: englishName}) => englishName === memberInputText)
                    if (this.exclude) {
                        this.localSelected.splice(this.inputIndex, 0, memberInputText)
                    } else {
                        const member = this.selectedIndex !== null ? this.filterMembers[this.selectedIndex] : this.filterMembers.find(({english_name: englishName}) => englishName === memberInputText)
                        if (member) {
                            this.localSelected.splice(this.inputIndex, 0, member['english_name'])
                        }
                    }
                }
                this.selectedIndex = null
                this.inputIndex = null
                this.memberInputReset()
            },
            /* 删除已选人员 */
            handleBackspace () {
                if (!this.memberInputText && this.inputIndex > 0) {
                    let memberInput = this.$refs.memberInput
                    this.localSelected.splice(--this.inputIndex, 1)
                    this.$refs.memberContainer.insertBefore(memberInput, memberInput.previousSibling)
                    this.memberInputFocus()
                }
            },
            /* 按下箭头 */
            handleArrow (event) {
                let $memberInput = this.$refs.memberInput
                let $memberSelected = this.$refs.memberSelected
                let arrow = event.key
                if (!this.memberInputText && arrow === 'ArrowLeft' && this.inputIndex > 0) { // 无输入时，向左移动光标位置
                    this.inputIndex--
                    this.$refs.memberContainer.insertBefore($memberInput, $memberSelected[this.inputIndex])
                    this.memberInputFocus()
                } else if (!this.memberInputText && arrow === 'ArrowRight' && this.inputIndex < $memberSelected.length) { // 无输入时，向右移动光标位置
                    this.inputIndex++
                    this.$refs.memberContainer.insertBefore($memberInput, $memberSelected[this.inputIndex])
                    this.memberInputFocus()
                } else if (arrow === 'ArrowDown' && this.filterMembers.length) { // 向下选择列表人员
                    event.preventDefault()
                    if (this.selectedIndex === null) {
                        this.selectedIndex = 0
                    } else if (this.selectedIndex < this.filterMembers.length - 1) {
                        this.selectedIndex++
                    }
                } else if (arrow === 'ArrowUp' && this.filterMembers.length) { // 向上选择列表人员
                    event.preventDefault()
                    if (this.selectedIndex > 0 && this.selectedIndex !== null) {
                        this.selectedIndex--
                    } else {
                        this.selectedIndex = null
                    }
                }
            },
            /* 光标获取焦点 */
            memberInputFocus () {
                this.memberInputReset()
                this.focus = true
                this.$nextTick(() => {
                    let $memberInput = this.$refs.memberInput
                    $memberInput.focus()
                })
            },
            /* 清空输入内容 */
            memberInputReset () {
                this.focus = false
                this.$refs.memberInput.innerHTML = ''
                this.$nextTick(() => {
                    this.memberInputText = ''
                })
            },
            /* 从过滤出来的人员列表中选中指定人员 */
            handleMemberItemClick (member) {
                if (!this.localSelected.some(selected => selected === member['english_name'])) {
                    this.localSelected.push(member['english_name'])
                }
                this.memberInputReset()
            },
            /* 输入时设置本地存储文本，用于筛选列表 */
            setMemberInputText (event) {
                this.memberInputText = this.$refs.memberInput.textContent.trim()
            },
            /* 获取人员列表展示的文本 */
            getLable (member) {
                if (!member['chinese_name']) {
                    return member['english_name']
                }
                return `${member['english_name']}(${member['chinese_name']})`
            }
        }
    }
</script>
<style lang="scss" scoped>
    .member-selector{
        width: 100%;
        height: 36px;
        cursor: text;
        position: relative;
        z-index: 1;
        &.active{
            z-index: 999;
        }
        &.disabled{
            cursor: not-allowed;
            background-color: #fafafa;
        }
        .member-wrapper{
            height: 100%;
        }
        .member-container{
            width: 100%;
            min-height: 100%;
            padding: 0 4px 8px 4px;
            border: 1px solid #c3cdd7;
            font-size: 0;
            border-radius: 2px;
            white-space: nowrap;
            overflow: hidden;
            background-color: #fff;
            max-height: 114px;
            line-height: 24px;
            @include scrollbar;
            &.active{
                overflow-y: auto;
                overflow-x: hidden;
                white-space: normal;
                position: absolute;
                top: 0;
                left: 0;
                border-color: $borderFocusColor;
            }
            &.active.ellipsis:after{
                display: none;
            }
            &.ellipsis:after{
                font-size: 12px;
                content: "···"; 
                position: absolute; 
                bottom: 1px; 
                right: 1px; 
                height: 34px;
                line-height: 34px;
                padding: 0 2px 0 10px;
                letter-spacing: 0;
                background: -webkit-linear-gradient(left, transparent, #fff 55%);
                background: -o-linear-gradient(left, transparent, #fff 55%);
                background: linear-gradient(to right, transparent, #fff 55%);
            }
        }
        &.placeholder:after{
            content: attr(data-placeholder);
            position: absolute;
            left: 4px;
            top: 0;
            height: 34px;
            line-height: 34px;
            font-size: 12px;
            color: #c3cdd7;
            pointer-events: none;
        }
    }
    .member-selected{
        display: inline-block;
        vertical-align: top;
        height: 18px;
        line-height: 16px;
        padding: 0 4px;
        margin: 8px 4px 0 4px;
        font-size: 12px;
        background-color: #fafafa;
        border: 1px solid #d9d9d9;
        border-radius: 2px;
        position: relative;
        .member-input-emitter {
            position: absolute;
            width: 8px;
            height: 34px;
            top: -9px;
            z-index: 1;
            &.before{
                right: 100%;
            }
            &.after{
                left: 100%;
            }
        }
    }
    .member-input{
        display: inline-block;
        vertical-align: top;
        font-size: 12px;
        min-width: 2px;
        height: 18px;
        line-height: 18px;
        margin: 8px 0 0 0;
        outline: 0;
        position: relative;
        z-index: 2;
        &:after{
            display: none;
            content: attr(data-shadow-member);
            color: #c3cdd7;
        }
    }
    .member-list{
        position: absolute;
        top: 100%;
        left: 0;
        width: 100%;
        font-size: 14px;
        max-height: 162px;
        overflow: auto;
        background: #fff;
        box-shadow: 0 0 1px 1px rgba(0, 0, 0, 0.1);
        border: 1px solid #c3cdd7;
        border-radius: 2px;
        @include scrollbar;
        z-index: 1000;
        .member-item{
            padding: 0 10px;
            height: 32px;
            line-height: 32px;
            cursor: pointer;
            @include ellipsis;
            &.highlight,
            &:hover{
                background-color: #f5f5f5;
            }
        }
    }
</style>