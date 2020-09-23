<template>
    <blueking-user-selector type="info"
        v-if="localValue.length"
        style="font-size: 12px;"
        :api="api"
        :value="localValue">
    </blueking-user-selector>
    <span v-else>--</span>
</template>

<script>
    import BluekingUserSelector from '@blueking/user-selector'
    export default {
        components: {
            BluekingUserSelector
        },
        props: {
            value: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                api: window.ESB.userManage
            }
        },
        computed: {
            localValue: {
                get () {
                    if (this.value) {
                        return this.value.split(',')
                    }
                    return []
                }
            }
        }
    }
</script>
