<template>
    <div>
        <div @mouseover="show" @mouseout="hide" :class="['bk-toolbox',{'is-show': isShow}]">
            <slot name="toolbox-trigger"></slot>
            <slot name="toolbox-content"></slot>
        </div>
    </div>
</template>
<script>
export default {
    name: 'bk-toolbox',
    data () {
        return {
            isShow: false
        }
    },
    methods: {
        show () {
            this.isShow = true
            this.$emit('show')
        },
        hide () {
            this.isShow = false
            this.$emit('hide')
        }
    }
}
</script>
