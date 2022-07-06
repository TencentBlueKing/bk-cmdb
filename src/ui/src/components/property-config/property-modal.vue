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
  <bk-dialog
    v-model="show"
    :draggable="false"
    :mask-close="false"
    :width="730"
    header-position="left"
    :title="$t('选择字段')"
    @value-change="handleVisibleChange"
    @confirm="handleConfirm">
    <bk-input
      class="search"
      type="text"
      :placeholder="$t('请输入字段名称搜索')"
      clearable
      right-icon="bk-icon icon-search"
      v-model.trim="searchName"
      @input="hanldeFilterProperty">
    </bk-input>
    <dl class="property-container">
      <div class="property-group" v-for="(group, groupIndex) in sortedGroups" :key="groupIndex">
        <dt class="group-title">{{group.bk_group_name}}</dt>
        <dd class="group-content">
          <ul class="property-list">
            <li
              class="property-item"
              v-for="property in groupedPropertyList[groupIndex]" :key="property.bk_property_id">
              <bk-checkbox
                :checked="isChecked(property)"
                :disabled="isDisabled(property)"
                :value="property.id"
                @change="handleChangeChecked(property, ...arguments)">
                <div
                  v-bk-tooltips.top-start="{ content: $t('该字段不支持配置'), disabled: !isDisabled(property) }">
                  {{property.bk_property_name}}
                </div>
              </bk-checkbox>
            </li>
          </ul>
        </dd>
      </div>
    </dl>
  </bk-dialog>
</template>

<script>
  export default {
    props: {
      visible: {
        type: Boolean,
        default: false
      },
      selectedList: {
        type: Array,
        default: () => ([])
      },
      sortedGroups: {
        type: Array,
        required: true,
        default: () => ([])
      },
      groupedProperties: {
        type: Array,
        required: true,
        default: () => ([])
      }
    },
    data() {
      return {
        show: this.visible,
        localSelected: [],
        searchName: '',
        propertyGroups: [],
        groupedPropertyList: []
      }
    },
    watch: {
      visible(val) {
        this.show = val
      },
      selectedList: {
        handler() {
          this.localSelected = this.selectedList.slice()
        },
        immediate: true
      },
      groupedProperties: {
        handler() {
          this.groupedPropertyList = this.groupedProperties.slice()
        },
        immediate: true
      }
    },
    methods: {
      isChecked(property) {
        return this.localSelected.some(target => target.id === property.id)
      },
      isDisabled(property) {
        return !property?.editable
      },
      handleChangeChecked(property, checked) {
        if (checked) {
          this.localSelected.push(property)
        } else {
          const index = this.localSelected.findIndex(target => target.id === property.id)
          index > -1 && this.localSelected.splice(index, 1)
        }
      },
      handleVisibleChange(val) {
        this.$emit('update:visible', val)
        if (!val) {
          this.groupedPropertyList = this.groupedProperties.slice()
          this.localSelected = this.selectedList.slice()
          this.searchName = ''
        }
      },
      handleConfirm() {
        this.$emit('update:selectedList', this.localSelected)
      },
      hanldeFilterProperty() {
        const keyword = this.searchName.toLowerCase()
        if (keyword) {
          this.groupedPropertyList = this.groupedProperties
            .map(properties => properties.filter(item => item.bk_property_name.toLowerCase().indexOf(keyword) > -1))
        } else {
          this.groupedPropertyList = this.groupedProperties.slice()
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
  .search {
      width: 280px;
      margin-bottom: 10px;
  }
  .property-container {
    height: 264px;
    @include scrollbar-y;
  }

  .property-group {
    margin-top: 14px;
    .group-title {
      position: relative;
      padding: 0 0 0 15px;
      line-height: 20px;
      font-size: 14px;
      font-weight: bold;
      color: #63656E;
      &:before {
        content: "";
        position: absolute;
        left: 0;
        top: 3px;
        width: 4px;
        height: 14px;
        background-color: #C4C6CC;
      }
    }
  }

  .property-list {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;

    .property-item {
      flex: 0 0 33.3333%;
      margin: 8px 0;
    }
  }
</style>
