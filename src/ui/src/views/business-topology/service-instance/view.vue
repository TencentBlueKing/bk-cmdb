<template>
    <div class="view">
        <component :is="activeComponent"></component>
    </div>
</template>

<script>
    import ViewInstance from './instance/view'
    import ViewProcess from './process/view'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            ViewInstance,
            ViewProcess
        },
        data () {
            return {
                activeComponent: null,
                viewMap: Object.freeze({
                    'instance': ViewInstance.name,
                    'process': ViewProcess.name
                })
            }
        },
        created () {
            this.unwatchView = RouterQuery.watch('view', this.handleViewChange, { immediate: true })
            this.unwatchTab = RouterQuery.watch('tab', this.handleTabChange)
        },
        beforeDestroy () {
            this.unwatchView()
            this.unwatchTab()
        },
        methods: {
            handleViewChange (view = 'instance') {
                this.activeComponent = this.viewMap[view]
            },
            handleTabChange (tab) {
                if (tab !== 'serviceInstance') {
                    this.activeComponent = null
                } else {
                    const view = RouterQuery.get('view', 'instance')
                    RouterQuery.set({ view })
                    this.activeComponent = this.viewMap[view]
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .view {
        position: relative;
        padding: 15px 0;
    }
</style>
