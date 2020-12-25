<template>
    <div class="cagetory-input" v-click-outside="handleCancel">
        <bk-input type="text"
            :ref="inputRef"
            :placeholder="placeholder"
            v-model="localValue"
            @enter="handleConfirm">
        </bk-input>
        <div class="operation">
            <span class="text-primary btn-confirm"
                @click.stop="handleConfirm">{{$t('确定')}}
            </span>
            <span class="text-primary" @click="handleCancel">{{$t('取消')}}</span>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            value: {
                type: String,
                default: ''
            },
            placeholder: {
                type: String,
                default: ''
            },
            inputRef: {
                type: String,
                default: ''
            },
            editId: {
                type: Number,
                default: 0
            }
        },
        data () {
            return {
                localValue: this.value
            }
        },
        watch: {
            value (value) {
                this.localValue = value
            },
            localValue (localValue) {
                this.$emit('input', localValue)
            }
        },
        methods: {
            handleConfirm () {
                this.$emit('on-confirm', this.localValue, this.editId)
            },
            handleCancel () {
                this.$emit('on-cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cagetory-input {
        @include space-between;
        width: 100%;
        font-weight: normal;
        .bk-form-control {
            flex: 1;
            font-size: 0;
            height: 32px;
            line-height: 32px;
            margin-right: 10px;
            /deep/ .bk-form-input {
                font-size: 14px;
                line-height: 32px;
            }
        }
        .text-primary {
            display: inline-block;
            line-height: normal;
            font-size: 12px;
            &.btn-confirm {
                position: relative;
                margin-right: 6px;
                &::after {
                    content: '';
                    position: absolute;
                    top: 3px;
                    right: -6px;
                    display: inline-block;
                    width: 1px;
                    height: 14px;
                    background-color: #dcdee5;
                }
            }
        }
    }
</style>
