<template>
    <div class="cmdb-form form-objuser"
        v-click-outside="handleClickOutside"
        @click="handleClick">
        <div class="objuser-layout">
            <div class="objuser-container"
                ref="container"
                :class="{disabled, focus, ellipsis}">
                <span class="objuser-selected"
                    v-for="(user, index) in localValue"
                    ref="selected"
                    :key="index"
                    @click.stop="handleSelectedClick($event, index)">
                    {{user}}
                </span>
                <span ref="input" class="objuser-input"
                    spellcheck="false"
                    contenteditable
                    v-show="focus"
                    @blur="handleBlur">
                </span>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-form-objuser',
        props: {
            value: {
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
            }
        },
        data () {
            return {
                localValue: ['admin', 'aaaaa', 'bbbbbbb', 'cccccccccc'],
                inputValue: '',
                focus: false,
                ellipsis: false
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
            }
        },
        mounted () {
            this.calcEllipsis()
        },
        methods: {
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
            handleClick () {
                if (this.disabled) {
                    return false
                }
                this.$refs.container.insertBefore(this.$refs.input, null)
                this.setSelection()
            },
            handleClickOutside () {
                this.focus = false
            },
            handleSelectedClick (event, index) {
                if (this.disabled) {
                    return false
                }
                let $refrenceTarget = event.target
                const offsetWidth = $refrenceTarget.offsetWidth
                const eventX = event.offsetX
                const $input = this.$refs.input
                if (eventX > (offsetWidth / 2)) {
                    $refrenceTarget = $refrenceTarget.nextElementSibling
                }
                if ($refrenceTarget === $input.nextElementSibling) {
                    this.setSelection()
                } else {
                    const $container = this.$refs.container
                    $container.insertBefore($input, $refrenceTarget)
                    this.setSelection(true)
                }
            },
            setSelection (reset = false) {
                if (reset) {
                    this.reset()
                }
                this.focus = true
                this.$nextTick(() => {
                    const $input = this.$refs.input
                    $input.focus()
                    if (window.getSelection) {
                        let range = window.getSelection()
                        range.selectAllChildren($input)
                        range.collapseToEnd()
                    } else if (document.selection) {
                        let range = document.selection.createRange()
                        range.moveToElementText($input)
                        range.collapse(false)
                        range.select()
                    }
                })
            },
            handleBlur () {
                
            },
            reset () {
                this.inputValue = ''
                this.$refs.input.innerHTML = ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-objuser {
        height: 36px;
        font-size: 14px;
        .objuser-layout {
            position: relative;
            height: 100%;
            .objuser-container {
                min-width: 100%;
                min-height: 100%;
                padding: 3px 0;
                border: 1px solid $cmdbBorderColor;
                border-radius: 2px;
                background-color: #fff;
                white-space: nowrap;
                overflow: hidden;
                &.disabled {
                    cursor: not-allowed;
                }
                &.focus {
                    white-space: normal;
                }
                &.ellipsis:after{
                    font-size: 12px;
                    content: ""; 
                    position: absolute; 
                    bottom: 1px; 
                    right: 1px; 
                    height: 34px;
                    line-height: 34px;
                    padding: 0 0 0 15px;
                    letter-spacing: 0;
                    background: -webkit-linear-gradient(left, transparent, #fff 55%);
                    background: -o-linear-gradient(left, transparent, #fff 55%);
                    background: linear-gradient(to right, transparent, #fff 55%);
                }
            }
        }
    }
    .objuser-selected {
        display: inline-block;
        height: 22px;
        margin: 3px;
        max-width: calc(100% - 4px);
        padding: 0 4px;
        line-height: 20px;
        vertical-align: top;
        border: 1px solid #d9d9d9;
        border-radius: 2px;
        @include ellipsis;
    }
    .objuser-input {
        display: inline-block;
        max-width: 100%;
        height: 22px;
        margin: 3px 0 0;
        padding: 0 4px;
        white-space: nowrap;
        line-height: 22px;
        vertical-align: top;
        outline: none;
        overflow: hidden;
    }
</style>