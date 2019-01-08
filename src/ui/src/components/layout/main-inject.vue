<template>
    <div class="cmdb-main-inject">
        <slot></slot>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-main-inject',
        props: {
            injectType: {
                type: String,
                default: 'append',
                validator (val) {
                    return ['append', 'prepend'].includes(val)
                }
            }
        },
        mounted () {
            const $scroller = document.querySelector('.main-layout')
            if (this.injectType === 'append') {
                $scroller.append(this.$el)
            } else {
                this.$el.insertBefore($scroller.firstChild)
            }
        },
        beforeDestroy () {
            this.$el.remove()
        }
    }
</script>