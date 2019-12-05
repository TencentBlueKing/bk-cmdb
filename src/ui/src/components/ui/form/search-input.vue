<template>
    <div class="cmdb-search-input">
        <div class="search-input-wrapper">
            <textarea ref="textarea"
                v-model="localValue"
                :rows="rows"
                :placeholder="placeholder || $t('请输入关键词')"
                :disabled="disabled"
                :maxlength="maxlength"
                @focus="handleFocus"
                @blur="handleBlur"
                @input="setValue"
                @keydown.enter="handleEnter"
                @keydown.delete="handleDelete">
            </textarea>
            <i class="bk-icon icon-close"
                :class="{
                    'is-show': focus
                }"
                v-if="localValue.length"
                @click="handleClear">
            </i>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-search-input',
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
            maxlength: {
                type: Number,
                default: 2000
            }
        },
        data () {
            return {
                localValue: this.value,
                rows: 1,
                timer: null,
                focus: false
            }
        },
        watch: {
            value (value) {
                this.setLocalValue()
            }
        },
        created () {
            if (this.focus) {
                this.setRows()
            }
        },
        methods: {
            setLocalValue () {
                if (this.localValue !== this.value) {
                    this.localValue = this.value
                    this.$emit('on-change', this.localValue)
                }
            },
            setValue () {
                this.$emit('input', this.localValue)
            },
            handleClear () {
                this.timer && clearTimeout(this.timer)
                this.localValue = ''
                this.rows = 1
                this.$refs.textarea.focus()
                this.setValue()
                this.$emit('clear')
            },
            setRows () {
                const rows = this.localValue.split('\n').length
                this.rows = Math.min(5, Math.max(rows, 1))
            },
            handleFocus () {
                this.setRows()
                this.focus = true
            },
            handleBlur () {
                this.focus = false
                this.timer = setTimeout(() => {
                    this.rows = 1
                    if (this.$refs.textarea) {
                        this.$refs.textarea.scrollTop = 0
                    }
                }, 200)
            },
            handleEnter () {
                this.rows = Math.min(this.rows + 1, 5)
                this.$emit('enter', this.localValue)
            },
            handleDelete () {
                this.$nextTick(() => {
                    this.setRows()
                })
            }
        }
    }
</script>
<style lang="scss" scoped>
    .cmdb-search-input {
        position: relative;
        .search-input-wrapper {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            line-height: 22px;
            z-index: 100;
            &:hover {
                .icon-close {
                    display: block;
                }
            }
            textarea {
                display: block;
                width: 100%;
                padding: 4px 20px 4px 10px;
                border: 1px solid #c3cdd7;
                resize: none;
                font-size: 14px;
                @include scrollbar-y;
            }
            .icon-close {
                display: none;
                position: absolute;
                top: 50%;
                right: 4px;
                width: 28px;
                height: 28px;
                line-height: 28px;
                text-align: center;
                transform: translate3d(0, -50%, 0) scale(.5);
                font-size: 12px;
                border-radius: 50%;
                background-color: #C4C6CC;
                color: #fff;
                cursor: pointer;
                &.is-show {
                    display: block;
                }
                &:hover {
                    background-color: #979BA5;
                }
            }
        }
    }
</style>
