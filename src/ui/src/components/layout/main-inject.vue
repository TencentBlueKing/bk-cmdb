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
            const $mainLayout = document.querySelector('.main-layout')
            if (this.injectType === 'append') {
                $mainLayout.append(this.$el)
            } else {
                $mainLayout.insertBefore(this.$el, $mainLayout.firstChild)
            }
        },
        beforeDestroy () {
            this.$el.remove()
        }
    }
</script>