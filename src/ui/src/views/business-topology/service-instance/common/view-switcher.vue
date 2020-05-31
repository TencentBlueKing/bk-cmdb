<template>
    <div class="bk-button-group">
        <bk-button v-for="item in items"
            :key="item.id"
            :class="{ 'is-selected': active === item.id }"
            @click="handleSwitch(item.id)">
            {{item.name}}
        </bk-button>
    </div>
</template>

<script>
    import RouterQuery from '@/router/query'
    export default {
        data () {
            return {
                active: RouterQuery.get('view', 'instance'),
                items: [{
                    id: 'instance',
                    name: this.$t('实例')
                }, {
                    id: 'process',
                    name: this.$t('进程')
                }]
            }
        },
        methods: {
            handleSwitch (active) {
                RouterQuery.set({ 'view': active })
            }
        }
    }
</script>
