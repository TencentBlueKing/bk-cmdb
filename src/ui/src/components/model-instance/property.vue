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
  <div class="property">
    <div v-for="(group, index) in $sortedGroups" :key="index" class="group">
      <h2 class="group-name">{{ group.bk_group_name }}</h2>
      <ul class="property-list">
        <li
          v-for="property in $groupedProperties[index]"
          :id="`property-item-${property.id}`"
          :key="property.id"
          :class="['property-item', property.bk_property_type]">
          <div v-bk-overflow-tips class="property-name">
            {{ property.bk_property_name }}
          </div>
          <template
            v-if="property.bk_property_type !== PROPERTY_TYPES.INNER_TABLE">
            <div
              v-if="property !== editState.property"
              :class="[
                'property-value',
                { 'is-loading': loadingState.includes(property) },
              ]">
              <cmdb-property-value
                :ref="`property-value-${property.bk_property_id}`"
                tag="div"
                :is-show-overflow-tips="isShowOverflowTips(property)"
                :value="instState[property.bk_property_id]"
                :property="property">
              </cmdb-property-value>
            </div>
            <template v-if="!loadingState.includes(property)">
              <template v-if="!readonly && !isPropertyEditable(property)">
                <i
                  v-bk-tooltips="{
                    content: $t('系统限定不可修改'),
                    placement: 'top',
                    onShow: () => {
                      setFocus(`#property-item-${property.id}`, true)
                    },
                    onHide: () => {
                      setFocus(`#property-item-${property.id}`, false)
                    },
                  }"
                  class="is-related property-edit icon-cc-edit-shape">
                </i>
              </template>
              <template v-else-if="!readonly">
                <cmdb-auth
                  v-show="property !== editState.property"
                  style="margin: 8px 0 0 8px; font-size: 0"
                  :auth="authData">
                  <bk-button
                    slot-scope="{ disabled }"
                    text
                    theme="primary"
                    class="property-edit-btn"
                    :disabled="disabled"
                    @click="setEditState(property)">
                    <i class="property-edit icon-cc-edit-shape"></i>
                  </bk-button>
                </cmdb-auth>
                <div
                  v-if="property === editState.property"
                  class="property-form">
                  <div :class="['form-component', property.bk_property_type]">
                    <component
                      :is="`cmdb-form-${property.bk_property_type}`"
                      v-bind="$tools.getValidateEvents(property)"
                      :ref="`component-${property.bk_property_id}`"
                      v-model.trim="editState.value"
                      v-validate="$tools.getValidateRules(property)"
                      :class="[
                        property.bk_property_type,
                        { error: errors.has(property.bk_property_id) },
                      ]"
                      :unit="property.unit"
                      :options="property.option || []"
                      :data-vv-name="property.bk_property_id"
                      :data-vv-as="property.bk_property_name"
                      :placeholder="$tools.getPropertyPlaceholder(property)"
                      :auto-check="false"
                      :multiple="property.ismultiple">
                    </component>
                  </div>
                  <i
                    class="form-confirm bk-icon icon-check-1"
                    @click="confirm"></i>
                  <i
                    class="form-cancel bk-icon icon-close"
                    @click="exitForm"></i>
                  <span
                    v-if="errors.has(property.bk_property_id)"
                    class="form-error">
                    {{ errors.first(property.bk_property_id) }}
                  </span>
                </div>
              </template>
              <template
                v-if="
                  !$tools.isEmptyPropertyValue(
                    instState[property.bk_property_id]
                  ) && property !== editState.property
                ">
                <div class="copy-box">
                  <i
                    class="property-copy icon-cc-copy"
                    @click="handleCopy(property.bk_property_id)">
                  </i>
                  <transition name="fade">
                    <span
                      v-if="showCopyTips === property.bk_property_id"
                      class="copy-tips"
                      :style="{
                        width: $i18n.locale === 'en' ? '100px' : '70px',
                      }">
                      {{ $t('复制成功') }}
                    </span>
                  </transition>
                </div>
              </template>
            </template>
          </template>
          <template v-else>
            <cmdb-form-innertable
              :ref="`component-${property.bk_property_id}`"
              class="form-component-innertable"
              :mode="'update'"
              :disabled="!isPropertyEditable(property)"
              :disabled-tips="$t('系统限定不可修改')"
              :property="property"
              :obj-id="objId"
              :instance-id="instanceId"
              :auth="authData" />
          </template>
        </li>
      </ul>
    </div>
    <slot name="append"></slot>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import isEqual from 'lodash/isEqual'

import {
  BUILTIN_MODELS,
  BUILTIN_MODEL_PROPERTY_KEYS,
  BUILTIN_MODEL_RESOURCE_TYPES,
} from '@/dictionary/model-constants.js'
import { PROPERTY_TYPES } from '@/dictionary/property-constants'
import businessSetService from '@/service/business-set/index.js'
import projectService from '@/service/project/index.js'
import formMixins from '@/mixins/form'

import authMixin from './mixin-auth'

export default {
  filters: {
    filterShowText(value, unit) {
      return value === '--' ? '--' : value + unit
    },
  },
  mixins: [formMixins, authMixin],
  props: {
    inst: {
      type: Object,
      required: true,
    },
    resourceType: {
      type: String,
      default: '',
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    objId: {
      type: String,
    },
  },
  data() {
    return {
      PROPERTY_TYPES,
      instState: this.inst,
      editState: {
        property: null,
        value: null,
      },
      loadingState: [],
      showCopyTips: false,
    }
  },
  computed: {
    ...mapGetters('objectModelClassify', ['models', 'getModelById']),
    authData() {
      const auths = {
        [BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS]]:
          this.INST_AUTH.U_BUSINESS,
        [BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS_SET]]:
          this.INST_AUTH.U_BUSINESS_SET,
        [BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.PROJECT]]:
          this.INST_AUTH.U_PROJECT,
      }
      return auths[this.resourceType] || this.INST_AUTH.U_INST
    },
    instanceId() {
      return this.instState?.[
        BUILTIN_MODEL_PROPERTY_KEYS[this.objId]?.ID || 'bk_inst_id'
      ]
    },
  },
  watch: {
    inst(val) {
      this.instState = this.$tools.getInstFormValues(
        this.properties,
        val,
        false
      )
    },
  },
  methods: {
    ...mapActions('objectCommonInst', ['updateInst']),
    ...mapActions('objectBiz', ['updateBusiness']),

    setFocus(id, focus) {
      const item = this.$el.querySelector(id)
      focus ? item.classList.add('focus') : item.classList.remove('focus')
    },
    getPlaceholder(property) {
      const placeholderTxt = ['enum', 'list', 'organization'].includes(
        property.bk_property_type
      )
        ? '请选择xx'
        : '请输入xx'
      return this.$t(placeholderTxt, { name: property.bk_property_name })
    },
    isPropertyEditable(property) {
      return property.editable && !property.bk_isapi
    },
    isShowOverflowTips(property) {
      const complexTypes = [PROPERTY_TYPES.MAP, PROPERTY_TYPES.ENUMQUOTE]
      return !complexTypes.includes(property.bk_property_type)
    },
    setEditState(property) {
      const value = this.instState[property.bk_property_id]
      this.editState.value = value === null || value === undefined ? '' : value
      this.editState.property = property
      this.$nextTick(() => {
        const component = this.$refs[`component-${property.bk_property_id}`]
        component[0] && component[0].focus && component[0].focus()
      })
    },
    async confirm() {
      const { property, value } = this.editState
      try {
        const isValid = await this.$validator.validateAll()
        if (!isValid) {
          return false
        }
        this.exitForm()
        const oldValue = this.instState[property.bk_property_id]
        if (isEqual(oldValue, value)) return

        this.loadingState.push(property)
        const values = {
          [property.bk_property_id]: this.$tools.formatValue(value, property),
        }

        if (
          this.resourceType ===
          BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS]
        ) {
          await this.updateBusiness({
            bizId: this.instState.bk_biz_id,
            params: values,
          })
        } else if (
          this.resourceType ===
          BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS_SET]
        ) {
          const MODEL_ID_KEY =
            BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID
          await businessSetService.update({
            bk_biz_set_ids: [this.instState[MODEL_ID_KEY]],
            data: {
              bk_biz_set_attr: { ...values },
            },
          })
        } else if (
          this.resourceType ===
          BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.PROJECT]
        ) {
          const params = {
            ids: [this.instState.id],
            data: values,
          }
          await projectService.update(params)
        } else {
          await this.updateInst({
            objId: this.instState.bk_obj_id,
            instId: this.instState.bk_inst_id,
            params: values,
          })
        }

        this.$success(this.$t('修改成功'))

        this.instState = { ...this.instState, ...values }

        this.loadingState = this.loadingState.filter(
          exist => exist !== property
        )

        this.$emit('after-update')
      } catch (e) {
        console.error(e)
        this.loadingState = this.loadingState.filter(
          exist => exist !== property
        )
      }
    },
    exitForm() {
      this.editState.property = null
      this.editState.value = null
    },
    handleCopy(propertyId) {
      const [component] = this.$refs[`property-value-${propertyId}`]
      const copyText = component?.getCopyValue() ?? ''
      this.$copyText(copyText).then(
        () => {
          this.showCopyTips = propertyId
          const timer = setTimeout(() => {
            this.showCopyTips = false
            clearTimeout(timer)
          }, 200)
        },
        () => {
          this.$error(this.$t('复制失败'))
        }
      )
    },
  },
}
</script>

<style lang="scss" scoped>
.property {
  height: 100%;
  overflow: auto;

  @include scrollbar-y;
}

.group {
  margin: 22px 0 0;

  .group-name {
    line-height: 21px;
    font-size: 16px;
    font-weight: normal;
    color: #333948;

    &::before {
      content: '';
      display: inline-block;
      vertical-align: -2px;
      width: 4px;
      height: 14px;
      margin-right: 9px;
      background-color: $cmdbBorderColor;
    }
  }
}

.property-list {
  width: 1208px;
  margin: 25px 0 0;
  color: #63656e;
  display: flex;
  flex-wrap: wrap;

  .property-item {
    flex: 0 0 50%;
    max-width: 50%;
    padding-bottom: 8px;
    display: flex;

    &:hover,
    &.focus {
      .property-edit {
        opacity: 1;
      }

      .property-copy {
        display: inline-block;
      }
    }

    &.innertable {
      flex: 0 0 100%;
      max-width: 100%;

      .property-name {
        flex: none;
      }

      .form-component-innertable {
        flex: none;
        width: calc(100% - 260px);
        margin-top: 6px;
      }
    }

    .property-name {
      position: relative;
      width: 260px;
      line-height: 32px;
      padding: 0 16px 0 36px;
      font-size: 14px;
      color: #63656e;
      text-align: right;

      @include ellipsis;

      &::after {
        position: absolute;
        right: 2px;
        content: '：';
      }
    }

    .property-value {
      margin: 6px 0 0 4px;
      max-width: 286px;
      font-size: 14px;
      color: #313237;
      overflow: hidden;
      text-overflow: ellipsis;
      word-break: break-all;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;

      &.is-loading {
        font-size: 0;

        &::before {
          content: '';
          display: inline-block;
          width: 16px;
          height: 16px;
          margin: 2px 0;
          background-image: url('../../assets/images/icon/loading.svg');
        }
      }
    }

    .property-edit-btn {
      height: auto;
      font-size: 0;
    }

    .property-edit {
      font-size: 16px;
      opacity: 0;

      &.is-related {
        display: inline-block;
        vertical-align: middle;
        width: 16px;
        height: 16px;
        margin: 8px 0 0 8px;
        line-height: 1;
      }

      &:hover {
        opacity: 0.8;
      }
    }

    .property-copy {
      margin: 8px 0 0 8px;
      color: #3c96ff;
      cursor: pointer;
      display: none;
      font-size: 16px;
    }

    .copy-box {
      position: relative;
      font-size: 0;

      .copy-tips {
        position: absolute;
        top: -22px;
        left: -18px;
        min-width: 70px;
        height: 26px;
        line-height: 26px;
        font-size: 12px;
        color: #fff;
        text-align: center;
        background-color: #9f9f9f;
        border-radius: 2px;
      }

      .fade-enter-active,
      .fade-leave-active {
        transition: all 0.5s;
      }

      .fade-enter {
        top: -14px;
        opacity: 0;
      }

      .fade-leave-to {
        top: -28px;
        opacity: 0;
      }
    }
  }
}

.property-form {
  font-size: 0;
  position: relative;

  .bk-icon {
    display: inline-block;
    vertical-align: middle;
    width: 32px;
    height: 32px;
    margin: 0 0 0 6px;
    border-radius: 2px;
    border: 1px solid #c4c6cc;
    line-height: 30px;
    font-size: 20px;
    text-align: center;
    cursor: pointer;

    &.form-confirm {
      color: #0082ff;
      font-size: 20px;

      &::before {
        display: inline-block;
      }
    }

    &.form-cancel {
      color: #979ba5;
      font-size: 20px;

      &::before {
        display: inline-block;
      }
    }

    &:hover {
      font-weight: bold;
    }
  }

  .form-error {
    position: absolute;
    top: 100%;
    left: 0;
    font-size: 12px;
    line-height: 1;
    color: $cmdbDangerColor;
  }

  .form-component {
    display: inline-block;
    vertical-align: middle;
    height: 32px;
    width: 260px;
    margin: 0 4px 0 0;

    &.bool {
      width: 42px;
      height: 24px;
    }

    &.longchar {
      height: auto;
      vertical-align: top;
    }
  }
}
</style>
