<template>
    <div class="cmdb-input">
        <bk-input :class="['cmdb-form-input', { 'has-icon': !!icon }]" type="text"
            v-model="localValue"
            :placeholder="placeholder"
            @enter="handleEnter">
        </bk-input>
        <i :class="[icon, 'input-icon']" v-if="icon" @click="handleIconClick"></i>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-input',
        props: {
            value: {
                type: String,
                default: ''
            },
            icon: {
                type: String,
                default: ''
            },
            placeholder: {
                type: String,
                default: ''
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
            handleEnter () {
                this.$emit('enter', this.localValue)
            },
            handleIconClick () {
                this.$emit('icon-click', this.localValue)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-input {
        position: relative;
        @include inlineBlock;
        .input-icon {
            position: absolute;
            font-size: 14px;
            right: 11px;
            top: 10px;
            cursor: pointer;
        }
    }
</style>
