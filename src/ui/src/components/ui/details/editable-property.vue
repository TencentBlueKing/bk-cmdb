<script lang="ts">
  import { computed, defineComponent, PropType, ref, nextTick } from '@vue/composition-api'
  import PropertyFormElement from '@/components/ui/form/property-form-element.vue'

  interface IProperty {
    id: number,
    'bk_property_id': string,
    'bk_property_name': string,
    'bk_property_type': string,
    'bk_isapi': boolean,
    'bk_property_group': string
    'editable': boolean
  }

  export default defineComponent({
    components: {
      PropertyFormElement
    },
    props: {
      value: {
        type: [String, Number, Array, Boolean],
        default: ''
      },
      property: {
        type: Object as PropType<IProperty>,
        default: () => ({})
      },
      auth: Object,
      loading: Boolean,
      editState: {
        type: Object,
        default: () => ({})
      },
      actionActive: Boolean,
      mustRequired: {
        type: Boolean,
        default: null
      }
    },
    setup(props, { emit }) {
      const $propertyFormElement = ref(null)

      const isEditable = computed(() => props.property.editable && !props.property.bk_isapi)

      const isEditing = computed(() => props.property === props.editState.property)

      const setEditState = (property: IProperty) => {
        const value = (props.value === null || props.value === undefined) ? '' : props.value
        emit('update:editState', { value, property })
        nextTick(() => {
          const component = $propertyFormElement.value.$refs[`component-${property.bk_property_id}`]
          component?.focus?.()
        })
      }

      const confirmEvents = computed(() => {
        const { bk_property_type: type } = props.property

        let eventName = 'change'

        if (['singlechar', 'int', 'float'].includes(type)) {
          eventName = 'enter'
        }

        if (['list', 'enum'].includes(type)) {
          eventName = 'on-selected'
        }

        return { [eventName]: confirmEdit }
      })

      const confirmEdit = async () => {
        const valid = await $propertyFormElement.value?.$validator?.validate?.()

        if (!valid) {
          return
        }

        const changed = props.value !== props.editState.value

        emit('confirm', changed)
      }

      const handleClickOutSide = () => {
        if (isEditing.value) {
          confirmEdit()
        }
      }

      return {
        isEditable,
        isEditing,
        setEditState,
        handleClickOutSide,
        $propertyFormElement,
        confirmEvents
      }
    }
  })
</script>

<template>
  <div :class="['editable-property']">
    <!-- 详情态 -->
    <cmdb-property-value
      v-if="!isEditing"
      v-bk-overflow-tips
      :class="['property-value', { 'is-loading': loading }]"
      tag="div"
      :ref="`property-value-${property.bk_property_id}`"
      :value="value"
      :property="property">
    </cmdb-property-value>

    <!-- 非处理中（保存）状态 -->
    <template v-if="!loading">
      <!-- 显示编辑入口 -->
      <div class="property-actions" v-show="!isEditing">
        <cmdb-auth
          tag="i"
          class="icon-cc-edit-shape property-edit-button"
          :auth="auth"
          v-bk-tooltips="{
            disabled: isEditable,
            content: $t('系统限定不可修改'),
            placement: 'top'
          }"
          @click="setEditState(property)">
        </cmdb-auth>
        <slot name="more-action"></slot>
      </div>

      <!-- 编辑态，显示表单项 -->
      <div class="property-form" v-if="isEditing">
        <property-form-element ref="$propertyFormElement"
          @click.stop
          v-click-outside="{
            handler: handleClickOutSide
          }"
          :must-required="mustRequired"
          :property="property"
          :events="confirmEvents"
          v-model="editState.value">
        </property-form-element>
      </div>
    </template>
  </div>
</template>

<style lang="scss" scoped>
  .editable-property {
    display: flex;
    align-items: center;

    &:hover,
    &.action-active {
      .property-actions {
        visibility: visible;
      }
    }

    .property-value {
      font-size: 12px;
      color: #313238;

      &.is-loading {
        font-size: 0;
        &:before {
          content: "";
          display: inline-block;
          width: 16px;
          height: 16px;
          margin: 2px 0;
          background-image: url("@/assets/images/icon/loading.svg");
        }
      }
    }

    .property-actions {
      display: flex;
      align-items: center;
      visibility: hidden; // 避免位移
      margin-left: 12px;
    }

    .property-form {
      width: 100%;
    }

    .property-edit-button {
      cursor: pointer;
      font-size: 16px;

      & + .property-edit-button {
        margin-left: 4px;
      }

      &:hover {
        color: $primaryColor;
      }
    }
  }
</style>
