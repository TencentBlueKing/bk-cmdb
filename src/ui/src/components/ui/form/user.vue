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
  <div class="cmdb-form-user">
    <div class="prepend" v-if="$slots.prepend">
      <slot name="prepend" />
    </div>
    <blueking-user-selector class="cmdb-form-objuser"
      ref="userSelector"
      display-list-tips
      v-bind="props"
      v-model="localValue"
      :class="{ 'has-fast-select': fastSelect }"
      :empty-text="$t('无匹配人员')"
      @focus="$emit('focus')"
      @blur="$emit('blur')">
    </blueking-user-selector>
  </div>
</template>

<script>
  import BluekingUserSelector from '@blueking/user-selector'
  import { mapGetters } from 'vuex'
  import Vue from 'vue'
  export default {
    name: 'cmdb-form-objuser',
    components: {
      BluekingUserSelector
    },
    props: {
      value: {
        type: String,
        default: ''
      },
      fastSelect: Boolean
    },
    computed: {
      ...mapGetters(['userName']),
      api() {
        return window.ESB.userManage
      },
      localValue: {
        get() {
          return (this.value && this.value.length) ? this.value.split(',') : []
        },
        set(val) {
          this.$emit('input', val.toString())
          this.$emit('change', val.toString, this.value)
        }
      },
      props() {
        const props = { ...this.$attrs }
        if (this.api) {
          try {
            const url = new URL(this.api)
            props.api = `${window.API_HOST}proxy/get/usermanage${url.pathname}`
          } catch (e) {
            console.error(e)
          }
        } else {
          props.fuzzySearchMethod = this.fuzzySearchMethod
          props.exactSearchMethod = this.exactSearchMethod
          props.pasteValidator = this.pasteValidator
        }
        return props
      }
    },
    mounted() {
      this.setupFastSelect()
    },
    methods: {
      setupFastSelect() {
        if (!this.fastSelect) return
        const FastSelect = new Vue({
          // eslint-disable-next-line no-unused-vars
          render: h => (
            <span class="fast-select"
                on-click={ this.handleFastSelect }>
                { this.$i18n.locale === 'en' ? 'me' : '我' }
            </span>
          )
        })
        FastSelect.$mount()
        // eslint-disable-next-line no-underscore-dangle
        FastSelect.$el.setAttribute([this.$options._scopeId], true)
        const { container } = this.$refs.userSelector.$refs
        container.parentElement.append(FastSelect.$el)
      },
      focus() {
        this.$refs.userSelector.focus()
      },
      async fuzzySearchMethod(keyword) {
        const users = await this.$http.get(`${window.API_HOST}user/list`, {
          params: {
            fuzzy_lookups: keyword
          },
          config: {
            cancelPrevious: true
          }
        })
        return {
          next: false,
          results: users.map(user => ({
            username: user.english_name,
            display_name: user.chinese_name
          }))
        }
      },
      exactSearchMethod(usernames) {
        const isBatch = Array.isArray(usernames)
        return Promise.resolve(isBatch ? usernames.map(username => ({ username })) : { username: usernames })
      },
      pasteValidator(usernames) {
        return Promise.resolve(usernames)
      },
      handleFastSelect(event) {
        event.stopPropagation()
        const value = [...this.localValue]
        const exist = value.includes(this.userName)
        if (exist) return
        if (this.$refs.userSelector.multiple) {
          value.push(this.userName)
        } else {
          value.splice(0, value.length, this.userName)
        }
        this.localValue = value
      }
    }
  }
</script>

<style lang="scss" scoped>
    .cmdb-form-user {
      display: flex;

      .prepend {
        margin-right: -1px;
      }
    }
    .cmdb-form-objuser {
        width: 100%;
        &.has-fast-select {
            /deep/ .user-selector-container {
                padding-right: 20px;
            }
        }

        &[size="small"] {
          height: 26px !important;

          /deep/ .user-selector-container:not(.focus) {
              height: 26px !important;

              &.placeholder:after {
                line-height: 24px;
              }

              .user-selector-input {
                margin-top: 2px;
                height: 20px;
                line-height: 20px;
              }

              .user-selector-selected,
              .user-selector-overflow-tag {
                margin: 2px 0 2px 6px;
                line-height: 20px;
              }
          }
        }
    }
    .fast-select {
        position: absolute;
        top: 50%;
        right: 4px;
        font-size: 12px;
        line-height: 16px;
        margin-top: -8px;
        color: $textColor;
        z-index: 2;
        cursor: pointer;
    }
</style>
<style lang="scss">
  .tippy-box[data-theme="light user-selector-popover"] {
    bottom: -4px
  }
</style>
