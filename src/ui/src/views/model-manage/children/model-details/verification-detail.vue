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
  <div class="model-slider-content verification-detail-wrapper">
    <bk-form form-type="vertical">
      <bk-form-item :label="$t('规则类型')" :desc="$t('唯一校验规则类型提示')" desc-type="icon">
        <bk-radio-group v-model="ruleType">
          <bk-radio class="mr30" value="single" :disabled="readonly">{{$t('单独唯一')}}</bk-radio>
          <bk-radio value="union" :disabled="readonly">{{$t('联合唯一')}}</bk-radio>
        </bk-radio-group>
      </bk-form-item>
      <bk-form-item :label="$t('校验规则')" required>
        <bk-select class="rules-selector" v-model="rules"
          searchable
          :readonly="readonly"
          :multiple="isUnion">
          <bk-option v-for="property in displayProperties"
            :key="property.id"
            :id="property.id"
            :name="property.bk_property_name">
          </bk-option>
        </bk-select>
        <p class="rules-error" v-if="validateResult">
          <template v-if="validateResult.type === SET_TYPE.superset">
            <i18n path="唯一校验子集提示">
              <template #name>
                <span class="rules-error-name"
                  v-for="(rule, index) in validateResult.rules"
                  :key="rule.id"
                  @click="handleErrorRuleClick(rule)">
                  {{getValidateRuleName(rule)}}
                  <template v-if="index !== (validateResult.rules.length - 1)">、</template>
                </span>
              </template>
            </i18n>
          </template>
          <template v-else>{{$t('唯一校验超集提示')}}</template>
        </p>
      </bk-form-item>
      <bk-form-item class="mt20" v-if="!readonly">
        <bk-button theme="primary" class="mr10"
          :disabled="submitDisabled"
          :loading="$loading([request.create, request.update])"
          @click="handleSubmit">
          {{this.id ? $t('保存') : $t('提交')}}
        </bk-button>
        <bk-button theme="default" @click="cancel">{{$t('取消')}}</bk-button>
      </bk-form-item>
    </bk-form>
  </div>
</template>

<script>
  export default {
    props: {
      ruleList: {
        type: Array,
        default: () => ([])
      },
      attributeList: {
        type: Array,
        default: () => ([])
      },
      id: {
        type: Number,
        default: null
      },
      readonly: Boolean
    },
    data() {
      return {
        ruleType: 'single',
        rules: '',
        singleRuleTypes: Object.freeze(['singlechar', 'int', 'float']),
        unionRuleTypes: Object.freeze(['singlechar', 'int', 'float', 'enum', 'date', 'list']),
        validateResult: null,
        SET_TYPE: Object.freeze({
          subset: Symbol('subset'),
          superset: Symbol('superset')
        }),
        request: Object.freeze({
          create: Symbol('create'),
          update: Symbol('update')
        })
      }
    },
    computed: {
      singleRuleProperties() {
        return this.attributeList.filter(property => this.singleRuleTypes.includes(property.bk_property_type))
      },
      unionRuleProperties() {
        return this.attributeList.filter(property => this.unionRuleTypes.includes(property.bk_property_type))
      },
      displayProperties() {
        return this.isUnion ? this.unionRuleProperties : this.singleRuleProperties
      },
      isUnion() {
        return this.ruleType === 'union'
      },
      submitDisabled() {
        if (!this.isChanged()) {
          return true
        }
        if (this.isUnion) {
          return this.rules.length === 0
        }
        return this.rules === ''
      },
      submitRules() {
        if (this.isUnion) {
          return this.rules
        }
        return this.rules ? [this.rules] : []
      },
      existRules() {
        return this.ruleList
          .filter(item => item.id !== this.id)
          .map(item => ({ id: item.id, rules: item.keys.map(key => key.key_id) }))
      },
      selectedName() {
        const name = this.submitRules.map((id) => {
          const property = this.attributeList.find(property => property.id === id)
          return property ? property.bk_property_name : id
        })
        return name.join('+')
      }
    },
    watch: {
      isUnion() {
        this.resetRules()
      },
      rules() {
        this.resetValidateResult()
      },
      id() {
        this.setupInitialRules()
      }
    },
    created() {
      this.setupInitialRules()
    },
    methods: {
      /**
       * 根据是否有id设置初始数据
       */
      setupInitialRules() {
        if (!this.id) return
        const record = this.ruleList.find(record => record.id === this.id)
        const rules = record.keys.map(rule => rule.key_id)
        const isUnion = rules.length > 1
        this.ruleType = isUnion ? 'union' : 'single'
        this.validateResult = null
        setTimeout(() => { // 等待watch中调用resetRules后再赋值, resetRules中已使用nextTick，此处使用setTimeout进一步延迟执行
          this.rules = isUnion ? rules : rules[0]
        }, 0)
      },
      async resetRules() {
        await this.$nextTick() // 等待bk-select内部切换完成单选多选等关联数据后再赋值
        this.rules = this.isUnion ? [] : ''
      },
      resetValidateResult() {
        this.validateResult = null
      },
      getValidateRuleName(rule) {
        const name = rule.rules.map((id) => {
          const property = this.attributeList.find(property => property.id === id)
          return property ? property.bk_property_name : id
        })
        return name.join('+')
      },
      /**
       * 判断当前规则是否是某些规则的子集，并返回已存在的超集
       */
      validateSubset() {
        const rules = this.submitRules
        const superset = this.existRules.filter(existRule => rules.every(id => existRule.rules.includes(id)))
        return superset.length ? superset : false
      },
      /**
       * 判断当前规则是否是某些规则的超集，并返回已存在的子集
       */
      validateSuperset() {
        const rules = this.submitRules
        const subset = this.existRules.filter(existRule => existRule.rules.every(id => rules.includes(id)))
        return subset.length ? subset : false
      },
      validateRules() {
        if (!this.submitRules.length) {
          return false
        }
        const subset = this.validateSubset()
        if (subset) {
          this.validateResult = {
            type: this.SET_TYPE.subset,
            rules: subset
          }
          return false
        }
        const superset = this.validateSuperset()
        if (superset) {
          this.validateResult = {
            type: this.SET_TYPE.superset,
            rules: superset
          }
          return false
        }
        return true
      },
      handleSubmit() {
        const valid = this.validateRules()
        if (!valid) return
        this.id ? this.updateRule() : this.createRule()
      },
      async createRule() {
        try {
          await this.$store.dispatch('objectUnique/createObjectUniqueConstraints', {
            objId: this.$route.params.modelId,
            params: {
              must_check: false, // 该参数对应UI已从界面移除
              keys: this.submitRules.map(id => ({ key_kind: 'property', key_id: id }))
            },
            config: {
              requestId: this.request.create
            }
          })
          this.$success('创建成功')
          this.$emit('save')
        } catch (error) {
          console.error(error)
        }
      },
      async updateRule() {
        try {
          const oldRule = this.ruleList.find(oldRule => oldRule.id === this.id)
          await this.$store.dispatch('objectUnique/updateObjectUniqueConstraints', {
            objId: this.$route.params.modelId,
            id: this.id,
            params: {
              must_check: oldRule.must_check,
              keys: this.submitRules.map(id => ({ key_kind: 'property', key_id: id }))
            },
            config: {
              requestId: this.request.update
            }
          })
          this.$success('修改成功')
          this.$emit('save')
        } catch (error) {
          console.error(error)
        }
      },
      handleErrorRuleClick(rule) {
        this.$emit('change-mode', rule)
        this.$success(this.$t('已切换至校验规则', { name: this.getValidateRuleName(rule) }))
      },
      /**
       * 外部调用
       * 用于判断当前是否进行了操作，用于侧滑关闭时进行提示
       */
      isChanged() {
        if (this.id) {
          const initialRecord = this.ruleList.find(rule => rule.id === this.id)
          const initialRules = initialRecord.keys.map(rule => rule.key_id)
          return initialRules.toString() !== this.submitRules.toString()
        }
        return !!this.submitRules.length
      },
      cancel() {
        this.$emit('cancel')
      }
    }
  }
</script>

<style lang="scss" scoped>
  .rules-selector {
    display: block;
  }
  .rules-error {
    font-size: 12px;
    color: $dangerColor;
    .rules-error-name {
      color: $primaryColor;
      cursor: pointer;
    }
  }
</style>

