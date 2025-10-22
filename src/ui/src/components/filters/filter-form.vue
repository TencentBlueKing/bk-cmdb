<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <bk-sideslider class="filter-form-sideslider"
    v-transfer-dom
    :is-show.sync="isShow"
    :width="400"
    :show-mask="false"
    :transfer="true"
    :before-close="handleSliderBeforeClose"
    @hidden="handleHidden">
    <div class="filter-form-header" slot="header">
      {{$t('高级筛选')}}
      <template v-if="collectable && collection">
        {{`(${collection.name})`}}
      </template>
    </div>
    <cmdb-sticky-layout class="filter-layout" slot="content" ref="propertyList" v-scroll="{
      targetClass: 'last-item',
      orientation: 'bottom',
      distance: 63
    }">
      <bk-form class="filter-form" form-type="vertical">
        <bk-form-item class="filter-ip filter-item">
          <label class="item-label">
            IP
          </label>
          <editable-block
            ref="ipEditableBlock"
            class="ip-editable-block"
            :enter-search="false"
            :is-exact="IPCondition.exact"
            :placeholder="editBlockPlaceholder"
            v-model="IPCondition.text">
          </editable-block>
          <input type="hidden"
            ref="ip"
            name="ip"
            data-vv-validate-on="change"
            data-vv-name="ip"
            v-validate="'ipSearchMaxCloud|ipSearchMaxCount'"
            v-model="IPCondition.text" />
          <p class="filter-ip-error" v-if="errors.has('ip')">
            {{errors.first('ip')}}
          </p>
          <div>
            <bk-checkbox class="mr20" v-model="IPCondition.inner" @change="handleIPOptionChange('outer', ...arguments)">
              {{$t('内网IP')}}
            </bk-checkbox>
            <bk-checkbox class="mr20" v-model="IPCondition.outer" @change="handleIPOptionChange('inner', ...arguments)">
              {{$t('外网IP')}}
            </bk-checkbox>
            <bk-checkbox v-bk-tooltips.top="$t('ipv6暂不支持模糊搜索')" v-model="IPCondition.exact">{{$t('精确')}}</bk-checkbox>
          </div>
          <div class="filter-operate">
            <condition-picker
              ref="conditionPicker"
              :text="$t(conditionText)"
              :icon="icon"
              :selected="selected"
              :property-map="propertyMap"
              :type="3">
            </condition-picker>
            <bk-popconfirm
              :content="$t('确定清空筛选条件')"
              width="280"
              trigger="click"
              :confirm-text="$t('确定')"
              :cancel-text="$t('取消')"
              @confirm="handleClearCondition">
              <bk-button :text="true" class="mr10" theme="primary"
                :disabled="!selected.length">
                {{$t('清空条件')}}
              </bk-button>
            </bk-popconfirm>
          </div>

        </bk-form-item>
        <bk-form-item class="filter-item"
          v-for="(property, index) in selected"
          :key="property.id"
          :class="[`filter-item-${property.bk_property_type}`, {
            'last-item': index === selected.length - 1 && scrollToBottom
          }]">
          <label class="item-label">
            {{property.bk_property_name}}
            <span class="item-label-suffix">({{getLabelSuffix(property)}})</span>
          </label>
          <div class="item-content-wrapper">
            <operator-selector class="item-operator"
              v-if="!withoutOperator.includes(property.bk_property_type)"
              :property="property"
              :custom-type-map="customOperatorTypeMap"
              :symbol-map="operatorSymbolMap"
              :desc-map="operatorDescMap"
              v-model="condition[property.id].operator"
              @change="handleOperatorChange(property, ...arguments)">
            </operator-selector>
            <component class="item-value r0"
              :is="getComponentType(property)"
              :placeholder="getPlaceholder(property)"
              :property="property"
              :is-paste-split="getPasteSplit(property.bk_property_id)"
              :ref="`component-${property.id}`"
              v-bind="getBindProps(property)"
              v-model.trim="condition[property.id].value"
              v-bk-tooltips.top="{
                disabled: !property.placeholder,
                theme: 'light',
                trigger: 'click',
                content: property.placeholder
              }"
              @active-change="handleComponentActiveChange(property, ...arguments)"
              @change="handleChange"
              @inputchange="hanleInputChange"
              @click.native="() => handleClick(`component-${property.id}`)"
              :popover-options="{
                duration: 0,
                onShown: handleShow,
                onHidden: handlePopoverHidden
              }">
            </component>
          </div>
          <i class="item-remove bk-icon icon-close" @click="handleRemove(property)"></i>
        </bk-form-item>
      </bk-form>
      <div class="filter-options"
        slot="footer"
        slot-scope="{ sticky }"
        :class="{ 'is-sticky': sticky }">
        <span v-bk-tooltips="{
          disabled: !searchDisabled,
          content: $t('条件无效，Node条件属性与其他条件属性不能同时设置')
        }">
          <bk-button
            class="option-search mr10 search-btn"
            theme="primary"
            :disabled="errors.any() || searchDisabled"
            @click="handleSearch">
            {{$t('查询')}}
          </bk-button>
        </span>
        <template v-if="collectable">
          <span class="option-collect-wrapper" v-if="collection"
            v-bk-tooltips="{
              disabled: allowCollect,
              content: $t('请先填写筛选条件')
            }">
            <bk-button class="option-collect" theme="default"
              :disabled="!allowCollect"
              @click="handleUpdateCollection">
              {{$t('更新条件')}}
            </bk-button>
          </span>
          <bk-popover v-else
            class="option-collect"
            ref="collectionPopover"
            placement="top-end"
            theme="light"
            trigger="manual"
            :width="280"
            :z-index="99999"
            :tippy-options="{
              interactive: true,
              hideOnClick: false,
              onShown: focusCollectionName,
              onHidden: clearCollectionName
            }"
            v-bk-tooltips="{
              disabled: allowCollect,
              content: $t('请先填写筛选条件')
            }">
            <bk-button theme="default" :disabled="!allowCollect || isMixCondition" @click="handleCreateCollection">
              {{$t('收藏此条件')}}
            </bk-button>
            <section class="collection-form" slot="content">
              <label class="collection-title">{{$t('收藏此条件')}}</label>
              <bk-input class="collection-name"
                ref="collectionName"
                :placeholder="$t('请填写名称')"
                :data-vv-as="$t('名称')"
                data-vv-name="collectionName"
                v-model="collectionForm.name"
                v-validate="'required|length:256'"
                @focus="handleCollectionFormFocus"
                @enter="handleSaveCollection">
              </bk-input>
              <p class="collection-error" v-if="collectionForm.error || errors.has('collectionName')">
                {{collectionForm.error || errors.first('collectionName')}}
              </p>
              <div class="collection-options">
                <bk-button class="mr10"
                  theme="primary"
                  size="small"
                  :disabled="!collectionForm.name.length"
                  :loading="$loading('createCollection')"
                  @click="handleSaveCollection">
                  {{$t('确定')}}
                </bk-button>
                <bk-button theme="default" size="small" @click="closeCollectionForm">{{$t('取消')}}</bk-button>
              </div>
            </section>
          </bk-popover>
        </template>
        <bk-button class="option-reset" theme="default" @click="handleReset">{{$t('清空')}}</bk-button>
      </div>
    </cmdb-sticky-layout>
  </bk-sideslider>
</template>

<script>
  import has from 'has'
  import FilterStore from './store'
  import OperatorSelector from './operator-selector'
  import { mapGetters } from 'vuex'
  import Utils from './utils'
  import { isContainerObject } from '@/service/container/common'
  import ConditionPicker from '@/components/condition-picker'
  import { setCursorPosition, getConditionSelect, updatePropertySelect, isPasteSplit } from '@/utils/util'
  import useSideslider from '@/hooks/use-sideslider'
  import isEqual from 'lodash/isEqual'
  import EditableBlock from '@/components/editable-block/index.vue'
  import { QUERY_OPERATOR, QUERY_OPERATOR_HOST_SYMBOL, QUERY_OPERATOR_HOST_DESC } from '@/utils/query-builder-operator'
  import { POSITIVE_INTEGER } from '@/dictionary/property-constants'

  export default {
    components: {
      OperatorSelector,
      ConditionPicker,
      EditableBlock
    },
    props: {
      type: {
        type: String,
        default: '' // index - 系统首页高级筛选
      },
      searchAction: {
        type: Function,
        default: () => {}
      },
      icon: {
        type: String,
        default: ''
      },
      conditionText: {
        type: String,
        default: '添加其他条件'
      }
    },
    data() {
      const { IN, NIN, LIKE, CONTAINS, EQ, NE, GTE, LTE, RANGE } = QUERY_OPERATOR
      return {
        scrollToBottom: false,
        isShow: false,
        withoutOperator: ['date', 'time', 'bool', 'service-template'],
        IPCondition: Utils.getDefaultIP(),
        originIPCondition: { ...FilterStore.IP },
        condition: {},
        originCondition: {},
        selected: [],
        collectionForm: {
          name: '',
          error: ''
        },
        customOperatorTypeMap: {
          float: [EQ, NE, GTE, LTE, RANGE, IN],
          int: [EQ, NE, GTE, LTE, RANGE, IN],
          longchar: [IN, NIN, CONTAINS, LIKE],
          singlechar: [IN, NIN, CONTAINS, LIKE],
          array: [IN, NIN, CONTAINS, LIKE],
          object: [IN, NIN, CONTAINS, LIKE]
        },
        operatorSymbolMap: QUERY_OPERATOR_HOST_SYMBOL,
        operatorDescMap: QUERY_OPERATOR_HOST_DESC
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
      editBlockPlaceholder() {
        const { exact } = this.IPCondition
        const placeholder = exact ? '主机搜索提示语' : '主机模糊搜索提示语'
        return this.$t(placeholder)
      },
      propertyMap() {
        let modelPropertyMap = { ...FilterStore.modelPropertyMap }
        const ignoreHostProperties = ['bk_host_innerip', 'bk_host_outerip', '__bk_host_topology__', 'bk_host_innerip_v6', 'bk_host_outerip_v6']
        modelPropertyMap.host = modelPropertyMap.host
          ?.filter(property => !ignoreHostProperties.includes(property.bk_property_id))

        // 暂时不支持node对象map类型的字段
        modelPropertyMap.node = modelPropertyMap.node
          ?.filter(property => !['map'].includes(property.bk_property_type))

        const getPropertyMapExcludeBy = (exclude = []) => {
          const excludes = !Array.isArray(exclude) ? [exclude] : exclude
          const propertyMap = []
          for (const [key, value] of Object.entries(modelPropertyMap)) {
            if (!excludes.includes(key)) {
              propertyMap[key] = value
            }
          }
          return propertyMap
        }

        // 资源-主机视图
        if (!FilterStore.bizId) {
          // 非已分配
          if (!FilterStore.isResourceAssigned) {
            return getPropertyMapExcludeBy('node')
          }
          return modelPropertyMap
        }

        // 当前处于业务节点，使用除业务外全量的字段(包括node)
        if (FilterStore.isBizNode) {
          return getPropertyMapExcludeBy('biz')
        }

        // 容器拓扑
        if (FilterStore.isContainerTopo) {
          return {
            host: modelPropertyMap.host || [],
            node: modelPropertyMap.node || [],
          }
        }

        // 业务拓扑主机，不需要业务和Node模型字段
        modelPropertyMap = {
          host: modelPropertyMap.host || [],
          module: modelPropertyMap.module || [],
          set: modelPropertyMap.set || []
        }
        return modelPropertyMap
      },
      storageSelected() {
        return FilterStore.selected
      },
      storageIPCondition() {
        return FilterStore.IP
      },
      collectable() {
        return FilterStore.collectable
      },
      collection() {
        return FilterStore.activeCollection
      },
      allowCollect() {
        const hasIP = !!this.IPCondition.text.trim().length
        const hasCondition = Object.keys(this.condition).some((id) => {
          const { value } = this.condition[id]
          return !Utils.isEmptyCondition(value)
        })
        return hasIP || hasCondition
      },
      isMixCondition() {
        const hasNodeField = Utils.hasNodeField(this.selected, this.condition)
        const hasNormalTopoField = Utils.hasNormalTopoField(this.selected, this.condition)
        return hasNormalTopoField && hasNodeField
      },
      searchDisabled() {
        return this.isMixCondition
      }
    },
    watch: {
      storageSelected: {
        immediate: true,
        handler(val) {
          const filterCondition = ['bk_host_innerip_v6', 'bk_host_outerip_v6']
          const { addSelect, deleteSelect } = getConditionSelect(val, this.selected)

          this.scrollToBottom = this.hasAddSelected(val, this.selected, addSelect)
          this.condition = this.setCondition(this.condition)
          updatePropertySelect(this.selected, this.handleRemove, addSelect, deleteSelect, 'push', filterCondition)
        }
      },
      storageIPCondition: {
        immediate: true,
        handler() {
          this.IPCondition = {
            ...this.storageIPCondition
          }
        }
      }
    },
    created() {
      setTimeout(() => {
        this.focusIP()
      }, 0)
      this.originCondition = this.setCondition(this.originCondition)
      const subTitle = this.type === 'index' ? '离开将会导致表单填写的内容丢失' : '离开将会导致未保存信息丢失'
      const { beforeClose, setChanged } = useSideslider('', { subTitle })
      this.beforeClose = beforeClose
      this.setChanged = setChanged
    },
    methods: {
      getPasteSplit(id) {
        return isPasteSplit(id)
      },
      hasAddSelected(val, oldVal, addSelect) {
        return val[0] && oldVal[0] && addSelect.length > 0
      },
      handleClearCondition() {
        this.clearCondition()
        this.selected = []
        FilterStore.updateSelected([...this.selected])
        FilterStore.updateUserBehavior(this.selected)
      },
      handleClick(e) {
        const parent = this.$refs[e][0].$el
        this.target = parent.getElementsByClassName('bk-select-tag-container')[0]
          || parent.getElementsByClassName('bk-tag-input')[0]

        if (~this.target?.className.indexOf('is-focus')) {
          // select专属
          return
        }
        this.calcPosition('click')
      },
      handleChange() {
        this.calcPosition()
      },
      hanleInputChange() {
        this.calcPosition()
      },
      handleShow() {
        this.calcPosition()
      },
      handlePopoverHidden() {
        this.$refs.propertyList.$el.classList.remove('over-height')
      },
      calcPosition(type = 'change') {
        if (type === 'click') this.$refs.propertyList.$el.classList.remove('over-height')
        if (!this.target) return

        this.$nextTick(() => {
          const limit = document.querySelector('.sticky-footer').getClientRects()[0].top
          const { bottom } = this.target.getClientRects()[0]
          if (bottom > Math.ceil(limit)) {
            this.$refs.propertyList.$el.classList.add('over-height')
          }
        })
      },
      setCondition(nowCondition) {
        const newCondition = this.$tools.clone(FilterStore.condition)
        Object.keys(nowCondition).forEach((id) => {
          if (has(nowCondition, id)) {
            newCondition[id] = nowCondition[id]
          }
        })
        return newCondition
      },
      getLabelSuffix(property) {
        const model = this.getModelById(property.bk_obj_id)
        return model ? model.bk_obj_name : model.bk_obj_id
      },
      getComponentType(property) {
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId,
          bk_property_type: propertyType,
          id
        } = property
        const {
          operator
        } = this.condition[id]
        const normal = `cmdb-search-${propertyType}`

        // 业务名在包含与非包含操作符时使用输入联想组件
        if (modelId === 'biz' && propertyId === 'bk_biz_name'
          && ![QUERY_OPERATOR.CONTAINS, QUERY_OPERATOR.LIKE].includes(this.condition[property.id].operator)) {
          return `cmdb-search-${modelId}`
        }

        // 数字类型int 和 float支持in操作符
        if (Utils.numberUseIn(property, operator)) {
          return 'cmdb-search-singlechar'
        }

        // 资源-主机下无业务
        if (!FilterStore.bizId) {
          return normal
        }

        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'

        // 在业务视图并且非模糊查询的情况，模块与集群名称使用专属的输入联想组件，否则使用与propertyType匹配的相应组件
        if ((isSetName || isModuleName)
          && ![QUERY_OPERATOR.CONTAINS, QUERY_OPERATOR.LIKE].includes(this.condition[property.id].operator)) {
          return `cmdb-search-${modelId}`
        }
        return normal
      },
      getBindProps(property) {
        const props = Utils.getBindProps(property)
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId,
          bk_property_type: propertyType,
          id
        } = property
        const {
          operator
        } = this.condition[id]

        if (POSITIVE_INTEGER.includes(propertyId)) {
          if (!props.options) props.options = {}
          props.options.min = 1
        }
        // 数字类型int 和 float支持in操作符
        if (Utils.numberUseIn(property, operator)) {
          props.onlyNumber = true
          props.fuzzy = false
        }
        if (!FilterStore.bizId) {
          return props
        }

        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'
        if (isSetName || isModuleName) {
          return Object.assign(props, { bizId: FilterStore.bizId })
        }

        // 容器对象标签属性，需要注入标签kv数据作为选项
        if (isContainerObject(modelId) && propertyType === 'map') {
          return Object.assign(props, { options: FilterStore.containerPropertyMapValue?.[modelId]?.[propertyId] })
        }

        return props
      },
      getPlaceholder(property) {
        return Utils.getPlaceholder(property)
      },
      handleIPOptionChange(negativeType, value) {
        if (!(value || this.IPCondition[negativeType])) {
          this.IPCondition[negativeType] = true
        }
      },
      handleOperatorChange(property, operator) {
        const { value } = this.condition[property.id]
        const effectValue = Utils.getOperatorSideEffect(property, operator, value)
        this.condition[property.id].value = effectValue
      },
      // 人员选择器参考定位空间不足，备选面板左移了，此处将其通过offset配置移到最右边
      handleComponentActiveChange(property, active) {
        if (!active) {
          return false
        }
        const { id, bk_property_type: type } = property
        if (type !== 'objuser') {
          return false
        }
        const [component] = this.$refs[`component-${id}`]
        try {
          this.$nextTick(() => {
            const reference = component.$el.querySelector('.user-selector-input')
            // eslint-disable-next-line no-underscore-dangle
            reference._tippy.setProps({
              offset: [240, 5]
            })
          })
        } catch (error) {
          console.error(error)
        }
      },
      async handleRemove(property) {
        const index = this.selected.indexOf(property)
        index > -1 && this.selected.splice(index, 1)
        if (this.collection) return
        await this.$nextTick()
        FilterStore.updateSelected([...this.selected])
        FilterStore.updateUserBehavior(this.selected)
      },
      handleSearch() {
        // tag-input组件在blur时写入数据有200ms的延迟，此处等待更长时间，避免无法写入
        this.searchTimer && clearTimeout(this.searchTimer)
        this.searchTimer = setTimeout(() => {
          const condition = {
            condition: this.$tools.clone(this.condition),
            IP: this.$tools.clone(this.IPCondition)
          }
          if (this.type === 'index') {
            return this.searchAction(condition)
          }

          FilterStore.resetPage(true)
          FilterStore.updateSelected(this.selected) // 此处会额外触发一次watch
          FilterStore.setCondition(condition)
          this.close()
        }, 300)
      },
      handleCreateCollection() {
        const { instance } = this.$refs.collectionPopover
        this.errors.clear()
        instance.show()
      },
      closeCollectionForm() {
        const { collectionPopover } = this.$refs
        if (!collectionPopover) {
          return false
        }
        const { instance } = this.$refs.collectionPopover
        instance.hide()
      },
      handleCollectionFormFocus() {
        this.collectionForm.error = null
      },
      async handleSaveCollection() {
        try {
          const isValid = await this.$validator.validate('collectionName')
          if (!isValid) {
            return false
          }
          await FilterStore.createCollection({
            bk_biz_id: FilterStore.bizId,
            name: this.collectionForm.name,
            info: JSON.stringify(this.IPCondition),
            query_params: this.getCollectionQueryParams()
          })
          this.$success(this.$t('收藏成功'))
          this.closeCollectionForm()
        } catch (error) {
          this.collectionForm.error = error.bk_error_msg
          console.error(error)
        }
      },
      getCollectionQueryParams() {
        const params = this.selected.map(property => ({
          bk_obj_id: property.bk_obj_id,
          field: property.bk_property_id,
          operator: this.condition[property.id].operator,
          value: this.condition[property.id].value
        }))
        return JSON.stringify(params)
      },
      async handleUpdateCollection() {
        try {
          await FilterStore.updateCollection({
            ...this.collection,
            info: JSON.stringify(this.IPCondition),
            query_params: this.getCollectionQueryParams()
          })
          this.$success(this.$t('更新收藏成功'))
        } catch (error) {
          console.error(error)
        }
      },
      handleReset() {
        this.$refs.ipEditableBlock.clear()
        this.clearCondition()
        this.errors.clear()
      },
      clearCondition() {
        Object.keys(this.condition).forEach((id) => {
          const property = this.selected.find(property => property.id.toString() === id.toString())
          const propertyCondititon = this.condition[id]
          const defaultValue = Utils.getOperatorSideEffect(property, propertyCondititon.operator, '')
          propertyCondititon.value = defaultValue
        })
      },
      focusCollectionName() {
        this.$refs.collectionName.$refs.input.focus()
      },
      clearCollectionName() {
        this.collectionForm.name = ''
        this.closeCollectionForm.error = ''
      },
      handleSliderBeforeClose() {
        const changedIPCondtion = !isEqual(this.IPCondition, this.originIPCondition)
        const changedCondition =  !isEqual(this.condition, this.originCondition)
        const { isShow } = this.$refs.conditionPicker

        if (isShow) return
        if (changedIPCondtion || changedCondition) {
          this.setChanged(true)
          return this.beforeClose(() => {
            this.close()
          })
        }
        this.close()
      },
      handleHidden() {
        this.$emit('closed')
      },
      open() {
        this.isShow = true
      },
      close() {
        this.isShow = false
      },
      focusIP() {
        const ele = this.$refs.ipEditableBlock
        if (ele) {
          ele?.focus()
          setCursorPosition(ele?.$refs?.searchInput, ele?.searchContent?.length)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .over-height {
      .g-expand {
        bottom: 0;
      }
    }
    .ip-editable-block {
      :deep(.search-input) {
        min-height: 82px;
        font-size: 12px;
        line-height: 24px;
      }

      :deep(.search-close) {
        font-size: 14px;
        &:hover {
            color: #979ba5;
        }
      }
    }
    .filter-form-sideslider {
        /deep/ {
            .bk-sideslider-wrapper {
                pointer-events: initial;
            }
        }
    }
    .filter-form-header {
        cursor: pointer;
        @include ellipsis;
    }
    .filter-layout {
        height: 100%;
        @include scrollbar-y;
    }
    .filter-form {
        padding: 0 14px;
    }
    .filter-ip {
        padding: 7px 10px 0px !important;
        position: sticky;
        top: 0;
        z-index: 9999;
        background: white;

        .filter-ip-error {
            line-height: initial;
            font-size: 12px;
            color: $dangerColor;
        }
        :deep(.bk-form-textarea) {
          resize: vertical;
        }
        .filter-operate {
          @include space-between;
        }
    }
    .filter-item {
        padding: 2px 10px 10px;
        &:not(.filter-ip):hover {
            background: #f5f6fa;
            .item-remove {
                opacity: 1;
            }
        }
        .item-label {
            display: block;
            font-size: 14px;
            font-weight: 400;
            line-height: 24px;
            @include ellipsis;
            .item-label-suffix {
                font-size: 12px;
                color: #979ba5;
            }
        }
        .item-content-wrapper {
            display: flex;
            align-items: flex-start;
            min-height: 32px;
        }
        .item-operator {
            flex: 128px 0 0;
            margin-right: 8px;
            & ~ .item-value {
                max-width: calc(100% - 136px);
            }
        }
        .item-value {
            flex: 1;
        }
        .item-remove {
            position: absolute;
            width: 24px;
            height: 24px;
            display: flex;
            justify-content: center;
            align-items: center;
            right: -10px;
            top: 3px;
            font-size: 20px;
            opacity: 0;
            cursor: pointer;
            color: $textColor;
            &:hover {
                color: $dangerColor;
            }
        }
    }
    .filter-options {
        display: flex;
        align-items: center;
        padding: 10px 24px;
        &.is-sticky {
            border-top: 1px solid $borderColor;
            background-color: #fff;
        }
        .option-reset {
            margin-left: auto;
        }
    }
    .collection-form {
        .collection-title {
            display: block;
            font-size: 13px;
            color: #63656E;
            line-height:17px;
        }
        .collection-name {
            margin-top: 13px;
        }
        .collection-error {
            color: $dangerColor;
            position: absolute;
        }
        .collection-options {
            display: flex;
            padding: 20px 0 10px;
            justify-content: flex-end;
        }
    }
</style>
