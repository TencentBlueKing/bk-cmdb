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
  <bk-sideslider class="filter-form-sideslider"
    v-transfer-dom
    :is-show.sync="isShow"
    :width="400"
    :show-mask="false"
    :transfer="true"
    :quick-close="false"
    @hidden="handleClosed">
    <div class="filter-form-header" slot="header">
      {{$t('高级筛选')}}
      <template v-if="collectable && collection">
        {{`(${collection.name})`}}
      </template>
      <i class="bk-icon icon-close" @click="close"></i>
    </div>
    <cmdb-sticky-layout class="filter-layout" slot="content">
      <bk-form class="filter-form" form-type="vertical">
        <bk-form-item class="filter-ip" label="IP">
          <bk-input type="textarea"
            ref="ip"
            :rows="4"
            :placeholder="$t('主机搜索提示语')"
            data-vv-name="ip"
            data-vv-validate-on="blur"
            v-validate="'ipSearchMaxCloud|ipSearchMaxCount'"
            v-focus
            v-model.trim="IPCondition.text"
            @focus="errors.remove('ip')">
          </bk-input>
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
        </bk-form-item>
        <bk-form-item class="filter-item"
          v-for="property in selected"
          :key="property.id"
          :class="`filter-item-${property.bk_property_type}`">
          <label class="item-label">
            {{property.bk_property_name}}
            <span class="item-label-suffix">({{getLabelSuffix(property)}})</span>
          </label>
          <div class="item-content-wrapper">
            <operator-selector class="item-operator"
              v-if="!withoutOperator.includes(property.bk_property_type)"
              :property="property"
              v-model="condition[property.id].operator"
              @change="handleOperatorChange(property, ...arguments)">
            </operator-selector>
            <component class="item-value"
              :is="getComponentType(property)"
              :placeholder="getPlaceholder(property)"
              :ref="`component-${property.id}`"
              v-bind="getBindProps(property)"
              v-model.trim="condition[property.id].value"
              v-bk-tooltips.top="{
                disabled: !property.placeholder,
                theme: 'light',
                trigger: 'click',
                content: property.placeholder
              }"
              @active-change="handleComponentActiveChange(property, ...arguments)">
            </component>
          </div>
          <i class="item-remove bk-icon icon-close" @click="handleRemove(property)"></i>
        </bk-form-item>
        <bk-form-item>
          <bk-button class="filter-add-button ml10" type="primary" text @click="handleSelectProperty">
            {{$t('添加其他条件')}}
          </bk-button>
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
            class="option-search mr10"
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
  import PropertySelector from './property-selector'
  import FilterStore from './store'
  import OperatorSelector from './operator-selector'
  import { mapGetters } from 'vuex'
  import Utils from './utils'
  import { isContainerObject } from '@/service/container/common'

  export default {
    components: {
      OperatorSelector
    },
    directives: {
      focus: {
        inserted: (el) => {
          const input = el.querySelector('textarea')
          setTimeout(() => {
            input.focus()
          }, 0)
        }
      }
    },
    data() {
      return {
        isShow: false,
        withoutOperator: ['date', 'time', 'bool', 'service-template'],
        IPCondition: Utils.getDefaultIP(),
        condition: {},
        selected: [],
        collectionForm: {
          name: '',
          error: ''
        }
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
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
        handler() {
          const newCondition = this.$tools.clone(FilterStore.condition)
          Object.keys(newCondition).forEach((id) => {
            if (has(this.condition, id)) {
              newCondition[id] = this.condition[id]
            }
          })
          this.condition = newCondition
          const filterCondition = ['bk_host_innerip_v6', 'bk_host_outerip_v6']
          this.selected = [...this.storageSelected].filter(item => !filterCondition.includes(item.bk_property_id))
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
    methods: {
      getLabelSuffix(property) {
        const model = this.getModelById(property.bk_obj_id)
        return model ? model.bk_obj_name : model.bk_obj_id
      },
      getComponentType(property) {
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId,
          bk_property_type: propertyType
        } = property
        const normal = `cmdb-search-${propertyType}`
        // 业务名在包含与非包含操作符时使用输入联想组件
        if (modelId === 'biz' && propertyId === 'bk_biz_name' && this.condition[property.id].operator !== '$regex') {
          return `cmdb-search-${modelId}`
        }

        // 资源-主机下无业务
        if (!FilterStore.bizId) {
          return normal
        }

        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'

        // 在业务视图并且非模糊查询的情况，模块与集群名称使用专属的输入联想组件，否则使用与propertyType匹配的相应组件
        if ((isSetName || isModuleName) && this.condition[property.id].operator !== '$regex') {
          return `cmdb-search-${modelId}`
        }

        return normal
      },
      getBindProps(property) {
        const props = Utils.getBindProps(property)
        if (!FilterStore.bizId) {
          return props
        }
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId,
          bk_property_type: propertyType
        } = property

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
      handleSelectProperty() {
        PropertySelector.show()
      },
      handleSearch() {
        // tag-input组件在blur时写入数据有200ms的延迟，此处等待更长时间，避免无法写入
        this.searchTimer && clearTimeout(this.searchTimer)
        this.searchTimer = setTimeout(() => {
          FilterStore.resetPage(true)
          FilterStore.updateSelected(this.selected) // 此处会额外触发一次watch
          FilterStore.setCondition({
            condition: this.$tools.clone(this.condition),
            IP: this.$tools.clone(this.IPCondition)
          })
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
        this.IPCondition.text = ''
        Object.keys(this.condition).forEach((id) => {
          const property = this.selected.find(property => property.id.toString() === id.toString())
          const propertyCondititon = this.condition[id]
          const defaultValue = Utils.getOperatorSideEffect(property, propertyCondititon.operator, null)
          propertyCondititon.value = defaultValue
        })
        this.errors.clear()
      },
      focusCollectionName() {
        this.$refs.collectionName.$refs.input.focus()
      },
      clearCollectionName() {
        this.collectionForm.name = ''
        this.closeCollectionForm.error = ''
      },
      handleClosed() {
        this.$emit('closed')
      },
      open() {
        this.isShow = true
      },
      close() {
        this.isShow = false
      },
      focusIP() {
        this.$refs?.ip?.$el.querySelector('textarea').focus()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .filter-form-sideslider {
        pointer-events: none;
        /deep/ {
            .bk-sideslider-wrapper {
                pointer-events: initial;
            }
            .bk-sideslider-closer {
                display: none;
            }
            .bk-sideslider-title {
                border-bottom: none;
            }
        }
    }
    .filter-form-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-left: -30px;
        .icon-close {
            font-size: 32px;
            margin-right: 6px;
            margin-top: -14px;
            cursor: pointer;
            opacity: .75;
            &:hover {
                opacity: 1;
            }
        }
    }
    .filter-layout {
        height: 100%;
        @include scrollbar-y;
    }
    .filter-form {
        padding: 0 10px;
    }
    .filter-ip {
        padding: 0 10px 10px;
        .filter-ip-error {
            line-height: initial;
            font-size: 12px;
            color: $dangerColor;
        }
        :deep(.bk-form-textarea) {
          resize: vertical;
        }
    }
    .filter-item {
        padding: 2px 10px 10px;
        margin-top: 5px !important;
        &:hover {
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
            align-items: center;
        }
        .item-operator {
            flex: 110px 0 0;
            margin-right: 8px;
            & ~ .item-value {
                max-width: calc(100% - 118px);
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
        padding: 10px 20px;
        &.is-sticky {
            border-top: 1px solid $borderColor;
            background-color: #fff;
        }
        .option-collect,
        .option-collect-wrapper {
            & ~ .option-reset {
                margin-left: auto;
            }
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
