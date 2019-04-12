<template>
    <div class="bk-page-count">
        <div class="bk-total-page">{{'共计' + totalPage + '页'}}</div>
        <bk-selector :placeholder="'页数'"
            :selected.sync="paginationIndex"
            :list="paginationListTmp"
            :setting-key="'id'"
            :display-key="'count'">
        </bk-selector>
    </div>
</template>
<script>
    import bkSelector from '../selector/selector.vue'
    export default {
        name: 'bk-pagination',
        components: {
            bkSelector
        },
        props: {
            paginationCount: {
                type: Number,
                default: 10,
                validator (value) {
                    return value >= 0
                }
            },
            totalPage: {
                type: Number,
                default: 5,
                validator (value) {
                    return value >= 0
                }
            },
            paginationList: {
                type: Array,
                default: () => [10, 20, 50, 100]
            }
        },
        data () {
            return {
                paginationIndex: this.paginationCount,
                paginationListTmp: []
            }
        },
        created () {
            this.initData()
        },
        watch: {
            paginationCount (value) {
                if (this.paginationList.includes(value)) {
                    this.paginationIndex = value
                } else {
                    this.paginationIndex = this.paginationList[0]
                }
            },
            paginationIndex (value) {
                this.$emit('update:paginationCount', value)
            }
        },
        methods: {
            initData () {
                this.paginationListTmp = this.paginationList.map(page => {
                    return {
                        id: page,
                        count: page
                    }
                })
            }
        }
    }
</script>

<style lang="scss">
    @import '../../bk-magic-ui/src/pagination.scss';
</style>
