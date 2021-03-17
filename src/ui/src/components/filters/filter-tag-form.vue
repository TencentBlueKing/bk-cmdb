<template>
  <bk-form form-type="vertical">
    <bk-form-item>
      <label class="form-label">
        {{property.bk_property_name}}
        <span class="form-label-suffix">({{labelSuffix}})</span>
      </label>
      <div class="form-wrapper">
        <operator-selector class="form-operator"
          v-if="!withoutOperator.includes(property.bk_property_type)"
          :type="property.bk_property_type"
          v-model="operator"
          @change="handleOperatorChange"
          @toggle="handleActiveChange">
        </operator-selector>
        <component class="form-value"
          :is="getComponentType()"
          :placeholder="getPlaceholder()"
          v-bind="getBindProps()"
          v-model.trim="value"
          @active-change="handleActiveChange">
        </component>
      </div>
      <div class="form-options">
        <bk-button class="mr10" text @click="handleConfirm">{{$t('确定')}}</bk-button>
        <bk-button class="mr10" text @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </bk-form-item>
  </bk-form>
</template>

<script>
  import OperatorSelector from './operator-selector'
  import FilterStore from './store'
  import Utils from './utils'
  import { mapGetters } from 'vuex'
  export default {
    components: {
      OperatorSelector
    },
    props: {
      property: {
        type: Object,
        required: true
      }
    },
    data() {
      return {
        withoutOperator: ['date', 'time', 'bool', 'service-template'],
        localOperator: null,
        localValue: null,
        active: false
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
      labelSuffix() {
        const model = this.getModelById(this.property.bk_obj_id)
        return model ? model.bk_obj_name : this.property.bk_obj_id
      },
      operator: {
        get() {
          return this.localOperator || FilterStore.condition[this.property.id].operator
        },
        set(operator) {
          this.localOperator = operator
        }
      },
      value: {
        get() {
          if (this.localValue === null) {
            return FilterStore.condition[this.property.id].value
          }
          return this.localValue
        },
        set(value) {
          this.localValue = value
        }
      }
    },
    methods: {
      getPlaceholder() {
        return Utils.getPlaceholder(this.property)
      },
      getComponentType() {
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId,
          bk_property_type: propertyType
        } = this.property
        const normal = `cmdb-search-${propertyType}`
        if (!FilterStore.bizId) {
          return normal
        }
        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'
        if (isSetName || isModuleName) {
          return `cmdb-search-${modelId}`
        }
        return normal
      },
      getBindProps() {
        const props = Utils.getBindProps(this.property)
        if (!FilterStore.bizId) {
          return props
        }
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId
        } = this.property
        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'
        if (isSetName || isModuleName) {
          return Object.assign(props, { bizId: FilterStore.bizId })
        }
        return props
      },
      resetCondition() {
        this.operator = null
        this.value = null
      },
      handleOperatorChange(operator) {
        this.value = Utils.getOperatorSideEffect(this.property, operator, this.value)
      },
      // 当失去焦点时，激活状态的改变做一个延时，避免点击表单外部时直接隐藏了表单对应的tooltips
      handleActiveChange(active) {
        this.timer && clearTimeout(this.timer)
        if (active) {
          this.active = active
        } else {
          this.timer = setTimeout(() => {
            this.active = active
          }, 100)
        }
      },
      handleConfirm() {
        FilterStore.updateCondition(this.property, this.operator, this.value)
        this.$emit('confirm')
      },
      handleCancel() {
        this.$emit('cancel')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-label {
        display: block;
        font-size: 14px;
        font-weight: 400;
        line-height: 32px;
        @include ellipsis;
        .form-label-suffix {
            font-size: 12px;
            color: #979ba5;
        }
    }
    .form-wrapper {
        width: 380px;
        display: flex;
        .form-operator {
            flex: 110px 0 0;
            margin-right: 8px;
            align-self: baseline;
            & ~ .form-value {
                max-width: calc(100% - 120px);
            }
        }
        .form-value {
            flex: 1;
        }
    }
    .form-options {
        display: flex;
        height: 32px;
        align-items: center;
        justify-content: flex-end;
    }
</style>
