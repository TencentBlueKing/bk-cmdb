<template>
    <div class="relevance-topo-wrapper">
        <div id="topo" class="topo"></div>
    </div>
</template>

<script>
    export default {
        props: {
            isShow: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.init()
                }
            }
        },
        methods: {
            init () {

            }
        }
    }
</script>

<style lang="scss" scoped>
    .relevance-topo-wrapper {
        height: calc(100% - 64px);
    }
</style>
