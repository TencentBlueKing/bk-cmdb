<template>
    <bk-select
        v-if="type === 'selector'"
        v-model="localValue"
        v-bind="$attrs">
        <bk-option
            v-for="biz in businessList"
            :key="biz.bk_biz_id"
            :id="biz.bk_biz_id"
            :name="`[${biz.bk_biz_id}] ${biz.bk_biz_name}`">
        </bk-option>
    </bk-select>
    <span v-else>{{bizName}}</span>
</template>

<script>
    export default {
        props: {
            value: {
                type: [String, Number]
            },
            type: {
                type: String,
                default: 'selector',
                validator (type) {
                    return ['selector', 'info'].includes(type)
                }
            }
        },
        data () {
            return {
                businessList: []
            }
        },
        computed: {
            localValue: {
                get () {
                    return this.value
                },
                set (value) {
                    this.$emit('input', value)
                    this.$emit('change', value)
                }
            },
            bizName () {
                const biz = this.businessList.find(biz => biz.bk_biz_id === this.value)
                return biz ? biz.bk_biz_name : '--'
            }
        },
        created () {
            this.getFullAmountBusiness()
        },
        methods: {
            async getFullAmountBusiness () {
                try {
                    const data = await this.$http.get('biz/simplify?sort=bk_biz_id', {
                        requestId: 'auditBusinessSelector',
                        fromCache: true
                    })
                    this.businessList = Object.freeze(data.info || [])
                } catch (e) {
                    console.error(e)
                    this.businessList = []
                }
            }
        }
    }
</script>
