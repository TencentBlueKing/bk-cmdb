<template>
  <div class="form-label">
    <span class="label-text">{{$t('字段设置')}}</span>
    <label class="cmdb-form-checkbox cmdb-checkbox-small" v-if="isEditableShow">
      <input type="checkbox" tabindex="-1" v-model="localValue.editable" :disabled="isReadOnly || ispre">
      <span class="cmdb-checkbox-text">
        {{$t('可编辑')}}
      </span>
      <i class="bk-cc-icon icon-cc-tips disabled-tips"
        v-if="modelId === 'host'"
        v-bk-tooltips="$t('主机属性设置为不可编辑状态后提示')"></i>
    </label>
    <label class="cmdb-form-checkbox cmdb-checkbox-small" v-if="isRequiredShow && !isMainLineModel">
      <input type="checkbox" tabindex="-1" v-model="localValue.isrequired" :disabled="isReadOnly || ispre">
      <span class="cmdb-checkbox-text">{{$t('必填')}}</span>
    </label>
  </div>
</template>

<script>
  export default {
    props: {
      isReadOnly: {
        type: Boolean,
        default: false
      },
      type: {
        type: String,
        required: true
      },
      editable: {
        type: Boolean,
        default: true
      },
      isrequired: {
        type: Boolean,
        default: false
      },
      isMainLineModel: {
        type: Boolean,
        default: false
      },
      ispre: Boolean
    },
    data() {
      return {
        editableMap: [
          'singlechar',
          'int',
          'float',
          'enum',
          'date',
          'time',
          'longchar',
          'objuser',
          'timezone',
          'bool',
          'list',
          'organization'
        ],
        isrequiredMap: [
          'singlechar',
          'int',
          'float',
          'date',
          'time',
          'longchar',
          'objuser',
          'timezone',
          'list',
          'organization'
        ],
        localValue: {
          editable: this.editable,
          isrequired: this.isrequired
        }
      }
    },
    inject: ['customObjId'], // 来源于自定义字段编辑
    computed: {
      isEditableShow() {
        return this.editableMap.indexOf(this.type) !== -1
      },
      isRequiredShow() {
        return this.isrequiredMap.indexOf(this.type) !== -1
      },
      modelId() {
        return this.$route.params.modelId ?? this.customObjId
      }
    },
    watch: {
      editable(editable) {
        this.localValue.editable = editable
      },
      isrequired(isrequired) {
        this.localValue.isrequired = isrequired
      },
      'localValue.editable'(editable) {
        this.$emit('update:editable', editable)
      },
      'localValue.isrequired'(isrequired) {
        if (!isrequired && this.isOnlyShow) {
          this.localValue.isonly = false
        }
        this.$emit('update:isrequired', isrequired)
      }
    }
  }
</script>

<style lang="scss" scoped>
  .disabled-tips {
    font-size: 12px;
    margin-left: 6px;
  }
</style>
