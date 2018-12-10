<template>
    <bk-selector :list="list" :selected.sync="localSelected"></bk-selector>
</template>

<script>
    export default {
        props: ['value', 'type'],
        data () {
            return {
                operatorMap: {
                    'common': [{
                        id: '$eq',
                        name: this.$t('Common[\'等于\']')
                    }, {
                        id: '$ne',
                        name: this.$t('Common[\'不等于\']')
                    }],
                    'char': [{
                        id: '$regex',
                        name: this.$t('Common[\'包含\']')
                    }, {
                        id: '$eq',
                        name: this.$t('Common[\'等于\']')
                    }, {
                        id: '$ne',
                        name: this.$t('Common[\'不等于\']')
                    }],
                    'name': [{
                        id: '$in',
                        name: 'IN'
                    }, {
                        id: '$eq',
                        name: this.$t('Common[\'等于\']')
                    }, {
                        id: '$ne',
                        name: this.$t('Common[\'不等于\']')
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