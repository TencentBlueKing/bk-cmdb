<template>
    <bk-select
        v-model="localSelected"
        :clearable="false"
        :popover-min-width="75"
        :disabled="disabled"
        :popover-options="{
            boundary: 'window'
        }">
        <bk-option
            v-for="(option, index) in list"
            :key="index"
            :id="option.id"
            :name="option.name">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        // eslint-disable-next-line
        props: ['value', 'type', 'disabled'],
        data () {
            return {
                operatorMap: {
                    'common': [{
                        id: '$eq',
                        name: this.$t('等于')
                    }, {
                        id: '$ne',
                        name: this.$t('不等于')
                    }],
                    'char': [{
                        id: '$multilike',
                        name: this.$t('包含')
                    }, {
                        id: '$in',
                        name: this.$t('等于')
                    }, {
                        id: '$ne',
                        name: this.$t('不等于')
                    }],
                    'name': [{
                        id: '$multilike',
                        name: this.$t('包含')
                    }, {
                        id: '$in',
                        name: this.$t('等于')
                    }, {
                        id: '$nin',
                        name: this.$t('不等于')
                    }],
                    'enum': [{
                        id: '$in',
                        name: this.$t('等于')
                    }, {
                        id: '$nin',
                        name: this.$t('不等于')
                    }]
                },
                localSelected: ''
            }
        },
        computed: {
            list () {
                let type = this.type
                if (!this.operatorMap.hasOwnProperty(type)) {
                    type = 'common'
                }
                return this.operatorMap[type]
            }
        },
        watch: {
            value (value) {
                this.setLocalSelected()
            },
            localSelected (localSelected) {
                this.$emit('input', localSelected)
                this.$emit('on-selected', localSelected)
            }
        },
        created () {
            this.setLocalSelected()
        },
        methods: {
            setLocalSelected () {
                if (this.list.some(item => item.id === this.value)) {
                    this.localSelected = this.value
                } else {
                    this.localSelected = this.list[0].id
                }
            }
        }
    }
</script>
