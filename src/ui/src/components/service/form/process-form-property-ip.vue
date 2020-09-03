<template>
    <cmdb-input-select
        name="ip"
        style="width: 100%"
        :placeholder="$t('请选择或输入IP')"
        :options="IPList"
        v-bind="$attrs"
        v-model="localValue">
    </cmdb-input-select>
</template>

<script>
    export default {
        props: {
            value: {
                type: String,
                default: ''
            }
        },
        inject: ['form'],
        data () {
            return {
                IPList: []
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
            requestId () {
                return `getInstanceIpByHost_${this.form.hostId}`
            }
        },
        created () {
            this.getBindIPList()
        },
        beforeDestroy () {
            this.$http.cancel(this.requestId)
        },
        methods: {
            async getBindIPList () {
                try {
                    const { options } = await this.$store.dispatch('serviceInstance/getInstanceIpByHost', {
                        hostId: this.form.hostId,
                        config: {
                            requestId: this.requestId,
                            fromCache: true
                        }
                    })
                    this.IPList = options.map(ip => ({ id: ip, name: ip }))
                } catch (error) {
                    this.IPList = []
                    console.error(error)
                }
            }
        }
    }
</script>
