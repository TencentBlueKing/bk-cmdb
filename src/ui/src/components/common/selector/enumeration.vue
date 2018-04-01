<template>
    <bk-select :selected.sync="localSelected" :disabled="disabled" @on-selected="handleSelected">
        <bk-select-option v-for="(option, index) in localOptions"
            :key="option.name"
            :label="option.name"
            :value="option.name">
        </bk-select-option>
    </bk-select>    
</template>
<script>
    export default {
        props: {
            disabled: {
                type: Boolean,
                default: false
            },
            selected: {
                required: true,
                validator: (selected) => {
                    return typeof selected === 'string' || !selected
                }
            },
            options: {
                type: [String, Array],
                required: true
            }
        },
        data () {
            return {
                localSelected: ''
            }
        },
        computed: {
            localOptions () {
                let localOptions = this.options
                if (!Array.isArray(localOptions)) {
                    try {
                        localOptions = JSON.parse(localOptions)
                    } catch (e) {
                        localOptions = []
                    }
                }
                return localOptions
            },
            defaultOption () {
                return this.localOptions.find(({is_default: isDefault}) => isDefault)
            }
        },
        watch: {
            selected (selected) {
                if (selected !== this.localSelected) {
                    this.setDefaultSelected()
                }
            },
            localSelected (localSelected) {
                this.$emit('update:selected', localSelected)
            },
            disabled (disabled) {
                this.setDefaultSelected()
            },
            localOptions (localOptions) {
                this.$nextTick(() => {
                    this.setDefaultSelected()
                })
            }
        },
        methods: {
            setDefaultSelected () {
                if (this.disabled) {
                    this.localSelected = this.selected ? this.selected : ''
                } else if (this.selected) {
                    this.localSelected = this.selected
                } else if (this.defaultOption) {
                    this.localSelected = this.defaultOption.name
                } else {
                    this.localSelected = ''
                }
            },
            handleSelected () {
                this.$emit('on-selected', ...arguments)
            }
        },
        created () {
            this.setDefaultSelected()
        }
    }
</script>
