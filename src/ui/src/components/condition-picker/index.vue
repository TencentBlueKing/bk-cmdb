<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <bk-popover ref="popover" :tippy-options="{
    delay: [300, 0],
    hideOnClick: true,
    interactive: true,
    placement: 'top',
    animateFill: false,
    sticky: true,
    theme: 'light',
    boundary: 'window',
    trigger: 'click',
    zIndex: 9999,
    onHidden: () => {
      this.confirm()
    }
  }">
    <bk-button class="form-condition-button" :style="{ marginTop: selected.length ? '10px' : 0 }"
      :icon="icon"
      :text="true"
      :disabled="disabled"
      @click="isShow = true">
      {{text}}
    </bk-button>
    <property-selector
      slot="content"
      v-if="isShow"
      ref="addConditionComp"
      :selected="selected"
      :disabled-property-map="disabledProperties"
      :models="models"
      :property-map="propertyMap"
      :height="height">
    </property-selector>
  </bk-popover>
</template>

<script>
  import { mapGetters } from 'vuex'
  import FilterStore from '@/components/filters/store'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import PropertySelector from '@/components/condition-picker/property-selector.vue'

  export default {
    components: {
      PropertySelector
    },
    props: {
      objId: String, // type = 2需要
      icon: {
        type: String,
        default: ''
      },
      disabled: {
        type: Boolean,
        default: false
      },
      text: String,
      type: {
        type: Number,
        default: 1 // 1动态分组添加条件  2.资源实例高级筛选添加条件 3. 主机高级筛选添加条件
      },
      selected: {
        type: Array,
        default: () => ([])
      },
      propertyMap: {
        type: [Object, Array],
        default: () => ({})
      },
      handler: Function
    },
    inject: {
      dynamicGroupForm: {
        value: 'dynamicGroupForm',
        default: null
      }
    },
    data() {
      return {
        height: 490,
        isShow: false,
        addConditionComp: null,
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
      appHeight() {
        return this.$store.state.appHeight
      },
      groups() {
        const sequence = ['host', 'module', 'set', 'node', 'biz']
        return Object.keys(this.propertyMap).map((modelId) => {
          const model = this.getModelById(modelId) || {}
          return {
            id: modelId,
            name: model.bk_obj_name,
            children: this.propertyMap[modelId]
          }
        })
          .sort((groupA, groupB) => sequence.indexOf(groupA.id) - sequence.indexOf(groupB.id))
      },
      target() {
        return this.dynamicGroupForm?.formData?.bk_obj_id
      },
      dynamicGroupModels() {
        if (this.target === 'host') {
          return this.dynamicGroupForm?.availableModels
        }
        return this.dynamicGroupForm?.availableModels.filter(model => model.bk_obj_id === this.target)
      },
      models() {
        if (this.type === 1) return this.dynamicGroupModels
        return this.groups.map(group => ({
          id: group.id,
          bk_obj_name: group.name,
          bk_obj_id: group.id
        }))
      },
      disabledProperties() {
        if (this.type === 1) return this.dynamicGroupForm?.disabledPropertyMap
        const disabledPropertyMap = {}
        this.groups.forEach((group) => {
          disabledPropertyMap[group.id] = group.children
            .filter(item => item.bk_property_type === PROPERTY_TYPES.INNER_TABLE)
            .map(item => item.bk_property_id)
        })
        return disabledPropertyMap
      }
    },
    watch: {
      isShow(val) {
        if (val) {
          const { bottom = 0 } = this.$refs?.popover?.$el?.getClientRects()?.[0]
          const dis = this.appHeight - bottom
          if (dis > 370 && dis < 500) {
            this.height = dis - 10
          } else {
            this.height = 490
          }
        }
      }
    },
    methods: {
      confirm() {
        const selected = this.$refs?.addConditionComp?.localSelected ?? this.selected
        this.isShow = false
        if (this.type !== 3) return this.handler([...selected])
        FilterStore.updateSelected(selected)
        FilterStore.updateUserBehavior(selected)
      },
    }
  }
</script>
<style  lang="scss" scoped>
.form-condition-button {
  :deep(>div) {
    display: flex;
    align-items: center;
    .bk-icon {
      top: 0;
    }
  }
}
</style>
