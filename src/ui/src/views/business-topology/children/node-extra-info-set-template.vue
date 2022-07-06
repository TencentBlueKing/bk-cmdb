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
  <div class="info-item clearfix" v-if="template"
    :title="`${$t('集群模板')} : ${template.name}`">
    <span class="name fl">{{$t('集群模板')}}</span>
    <div class="value fl">
      <template v-if="instance.set_template_id">
        <div class="template-value set-template fl" @click="linkToTemlate">
          <span class="text link">{{template.name}}</span>
          <i class="icon-cc-share"></i>
        </div>
        <cmdb-auth v-if="!isBizSet" :auth="{ type: $OPERATION.U_TOPO, relation: [bizId] }">
          <bk-button slot-scope="{ disabled }"
            :class="['sync-set-btn', 'ml5', { 'has-change': hasChange }]"
            :disabled="!hasChange || disabled"
            @click="handleSyncSetTemplate">
            {{$t('同步集群')}}
          </bk-button>
        </cmdb-auth>
      </template>
      <span class="text" v-else>{{template.name}}</span>
    </div>
  </div>
</template>

<script>
  import { MENU_BUSINESS_SET_TOPOLOGY, MENU_BUSINESS_SET_TEMPLATE_DETAILS } from '@/dictionary/menu-symbol'
  export default {
    name: 'set-template-info',
    props: {
      instance: {
        type: Object,
        required: true
      }
    },
    data() {
      return {
        template: null,
        hasChange: false
      }
    },
    computed: {
      isBizSet() {
        return this.$route.name === MENU_BUSINESS_SET_TOPOLOGY
      },
      bizId() {
        const { objectBiz, bizSet } = this.$store.state
        return this.isBizSet ? bizSet.bizId : objectBiz.bizId
      },
      selectedNode() {
        return this.$store.state.businessHost.selectedNode
      }
    },
    watch: {
      instance: {
        immediate: true,
        handler(instance) {
          instance && this.setTemplate()
        }
      }
    },
    methods: {
      async setTemplate() {
        if (this.instance.set_template_id) {
          try {
            const [template, hasChange] = await Promise.all([
              this.getTemplate(),
              this.getTemplateDiff()
            ])
            this.template = template
            this.hasChange = hasChange
          } catch (error) {
            this.template = null
            this.hasChange = false
            console.error(error)
          }
        } else {
          this.template = {
            name: this.$t('无')
          }
          this.hasChange = false
        }
      },
      getTemplate() {
        return this.$store.dispatch('setTemplate/getSingleSetTemplateInfo', {
          bizId: this.bizId,
          setTemplateId: this.instance.set_template_id,
          config: {
            requestId: 'getSingleSetTemplateInfo'
          }
        })
      },
      async linkToTemlate() {
        try {
          const data = await this.$store.dispatch('setTemplate/getSingleSetTemplateInfo', {
            setTemplateId: this.instance.set_template_id,
            bizId: this.bizId,
            config: {
              globalError: false
            }
          })
          if (!data) {
            return this.$error(this.$t('跳转失败，集群模板已经被删除'))
          }
        } catch (error) {
          console.error(error)
          this.$error(error.message)
        }
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SET_TEMPLATE_DETAILS,
          params: {
            templateId: this.instance.set_template_id
          },
          history: true
        })
      },
      handleSyncSetTemplate() {
        this.$store.commit('setFeatures/setSyncIdMap', {
          id: `${this.bizId}_${this.instance.set_template_id}`,
          instancesId: [this.instance.bk_set_id]
        })
        this.$routerActions.redirect({
          name: 'setSync',
          params: {
            setTemplateId: this.instance.set_template_id,
            moduleId: this.selectedNode.data.bk_inst_id
          },
          history: true
        })
      },
      async getTemplateDiff() {
        try {
          const data = await this.$store.dispatch('setSync/diffTemplateAndInstances', {
            bizId: this.bizId,
            setTemplateId: this.instance.set_template_id,
            params: {
              bk_set_id: this.instance.bk_set_id
            },
            config: {
              requestId: 'diffTemplateAndInstances'
            }
          })
          const diff = data.difference ? (data.difference || {}).module_diffs : []
          const len = diff.filter(_module => _module.diff_type !== 'unchanged').length
          return !!len
        } catch (e) {
          console.error(e)
          return false
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .info-item {
        width: 50%;
        max-width: 400px;
        line-height: 26px;
        margin-bottom: 12px;
        .name {
            position: relative;
            padding-right: 16px;
            &::after {
                content: ":";
                position: absolute;
                right: 10px;
            }
        }
        .value {
            width: calc(80% - 10px);
            padding-right: 10px;
            .text {
                @include inlineBlock;
                @include ellipsis;
                max-width: calc(100% - 16px);
                font-size: 14px;
            }
            .template-value {
                width: 100%;
                font-size: 0;
                color: #3a84ff;
                cursor: pointer;
                &.set-template {
                    width: auto;
                    max-width: calc(100% - 75px);
                }
            }
            .icon-cc-share {
                @include inlineBlock;
                font-size: 12px;
                margin-left: 4px;
            }
        }
    }
    .sync-set-btn {
        position: relative;
        height: 26px;
        line-height: 24px;
        padding: 0 10px;
        font-size: 12px;
        margin-top: -2px;
        &.has-change::before {
            content: '';
            position: absolute;
            top: -4px;
            right: -4px;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background-color: #EA3636;
        }
    }
</style>
