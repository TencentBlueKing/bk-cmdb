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
  <bk-dialog class="bk-dialog-no-padding edit-label-dialog"
    v-model="isShow"
    :width="580"
    :mask-close="false"
    :esc-close="false"
    @after-leave="handleHidden">
    <div slot="header">
      {{$t('编辑标签')}}
    </div>
    <label-dialog-content ref="labelComp"
      :default-labels="labels">
    </label-dialog-content>
    <div class="edit-label-dialog-footer" slot="footer">
      <bk-button theme="primary" :loading="$loading(Object.values(request))" @click.stop="handleSubmit">
        {{$t('确定')}}
      </bk-button>
      <bk-button theme="default" class="ml5" @click.stop="close">{{$t('取消')}}</bk-button>
    </div>
  </bk-dialog>
</template>

<script>
  import LabelDialogContent from './label-dialog-content'
  import { mapGetters } from 'vuex'
  export default {
    components: {
      LabelDialogContent
    },
    props: {
      serviceInstance: Object,
      updateCallback: Function
    },
    data() {
      return {
        isShow: false,
        request: {
          update: Symbol('update')
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      labels() {
        if (!this.serviceInstance.labels) {
          return []
        }
        return Object.entries(this.serviceInstance.labels).map(([key, value], index) => ({ id: index, key, value }))
      }
    },
    methods: {
      show() {
        this.isShow = true
      },
      close() {
        this.isShow = false
      },
      async handleSubmit() {
        try {
          const { labelComp } = this.$refs
          const validateResult = await labelComp.$validator.validateAll()
          if (!validateResult) {
            return false
          }
          const list = labelComp.submitList

          const labelSet = {}
          list.forEach((label) => {
            labelSet[label.key] = label.value
          })
          await this.$store.dispatch('instanceLabel/updateInstanceLabel', {
            params: {
              bk_biz_id: this.bizId,
              instance_ids: [this.serviceInstance.id],
              labels: labelSet
            },
            config: {
              requestId: this.request.update
            }
          })
          this.updateCallback && this.updateCallback(list)
          this.$success(this.$t('保存成功'))
          this.isShow = false
        } catch (error) {
          console.error(error)
        }
      },
      handleHidden() {
        this.$emit('close')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .edit-label-dialog {
        /deep/ .bk-dialog-header {
            text-align: left !important;
            font-size: 24px;
            color: #444444;
            margin-top: -15px;
        }
        .edit-label-dialog-footer {
            .bk-button {
                min-width: 76px;
            }
        }
    }
</style>
