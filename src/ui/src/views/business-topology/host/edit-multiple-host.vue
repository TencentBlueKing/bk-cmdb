<template>
    <bk-sideslider
        v-transfer-dom
        :is-show.sync="slider.show"
        :title="slider.title"
        :width="800"
        :before-close="handleSliderBeforeClose">
        <component :is="slider.component" slot="content"
            v-if="slider.show"
            ref="multipleForm"
            v-bind="slider.props"
            @on-submit="handleMultipleSave"
            @on-cancel="handleSliderBeforeClose">
        </component>
    </bk-sideslider>
</template>

<script>
    import editMultipleHost from '@/mixins/edit-multiple-host'
    export default {
        mixins: [editMultipleHost]
    }
</script>
