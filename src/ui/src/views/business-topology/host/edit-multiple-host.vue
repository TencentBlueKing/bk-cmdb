<template>
  <bk-sideslider
    v-transfer-dom
    :is-show.sync="slider.show"
    :title="slider.title"
    :width="800"
    :before-close="handleSliderBeforeClose">
    <div class="form-content" slot="content" style="height: 100%"
      v-bkloading="{ isLoading: $loading(Object.values(request)) }">
      <component :is="slider.component"
        ref="multipleForm"
        v-bind="slider.props"
        @on-submit="handleMultipleSave"
        @on-cancel="handleSliderBeforeClose">
      </component>
    </div>
  </bk-sideslider>
</template>

<script>
  import editMultipleHost from '@/mixins/edit-multiple-host'
  export default {
    mixins: [editMultipleHost]
  }
</script>
