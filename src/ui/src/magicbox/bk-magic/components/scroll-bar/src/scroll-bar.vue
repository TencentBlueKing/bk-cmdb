<template>
    <div class="bk-scrollbar"
        :class="{
            'vertical-bar': vertical,
            'horizontal-bar': horizontal}">
        <div class="bk-scrollbar-slider"></div>
    </div>
</template>

<script>
    export default {
        name: 'bk-scrollbar',
        data () {
            return {
                vertical: true,
                horizontal: false,
                direction: 'vertical'
            }
        }
    }
</script>
