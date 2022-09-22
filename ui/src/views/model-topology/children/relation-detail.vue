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
  <div>
    <label class="form-label">
      <span class="label-text">
        {{$t('源模型')}}
      </span>
      <div class="cmdb-form-item">
        <bk-input type="text" class="cmdb-form-input" disabled
          :value="getModelName(relationInfo['bk_obj_id'])">
        </bk-input>
      </div>
    </label>
    <label class="form-label">
      <span class="label-text">
        {{$t('目标模型')}}
      </span>
      <div class="cmdb-form-item">
        <bk-input type="text" class="cmdb-form-input" disabled
          :value="getModelName(relationInfo['bk_asst_obj_id'])">
        </bk-input>
      </div>
    </label>
    <label class="form-label">
      <span class="label-text">
        {{$t('关联类型')}}
      </span>
      <div class="cmdb-form-item">
        <cmdb-selector style="width: 100%;"
          readonly
          :list="relationList"
          v-model="relationInfo.bk_asst_id"
        ></cmdb-selector>
      </div>
    </label>
    <label class="form-label">
      <span class="label-text">
        {{$t('源-目标约束')}}
      </span>
      <div class="cmdb-form-item">
        <cmdb-selector style="width: 100%;"
          readonly
          :list="mappingList"
          v-model="relationInfo.mapping"
        ></cmdb-selector>
      </div>
    </label>
    <label class="form-label">
      <span class="label-text">
        {{$t('关联描述')}}
      </span>
      <div class="cmdb-form-item" :class="{ 'is-error': errors.has('asstName') }">
        <bk-input type="text" class="cmdb-form-input"
          name="asstName"
          :readonly="relationInfo.ispre || !isEdit"
          v-model.trim="relationInfo['bk_obj_asst_name']"
          v-validate="'singlechar|length:256'">
        </bk-input>
        <p class="form-error">{{errors.first('asstName')}}</p>
      </div>
    </label>
    <div class="btn-group" v-if="isEdit && relationInfo.bk_asst_id !== 'bk_mainline'">
      <bk-button
        theme="primary"
        :loading="$loading('updateObjectAssociation')"
        :disabled="JSON.stringify(relationInfo) === relationInfoSnapshot"
        @click="saveRelation"
      >
        {{$t('保存')}}
      </bk-button>
      <bk-button
        theme="danger"
        outline
        :disabled="relationInfo.ispre || $loading('deleteObjectAssociation')"
        @click="deleteRelation"
      >
        {{$t('删除关联')}}
      </bk-button>
    </div>
  </div>
</template>

<script>
  import { mapActions, mapGetters } from 'vuex'
  export default {
    props: {
      objId: {
        type: String
      },
      asstId: {
        type: [String, Number],
        required: true
      },
      asstInfo: {
        type: Object
      },
      isEdit: {
        type: Boolean
      }
    },
    data() {
      return {
        relationList: [],
        mappingList: [{
          id: 'n:n',
          name: 'N-N'
        }, {
          id: '1:n',
          name: '1-N'
        }, {
          id: '1:1',
          name: '1-1'
        }],
        relationInfo: {
          ispre: false,
          id: 0,
          bk_obj_asst_id: '',
          bk_obj_asst_name: '',
          bk_obj_id: '',
          bk_asst_obj_id: '',
          bk_asst_id: '',
          mapping: '',
          on_delete: []
        },
        relationInfoSnapshot: ''
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['models'])
    },
    async created() {
      await this.initRelationList()
      this.initData()
    },
    methods: {
      ...mapActions('objectAssociation', [
        'searchObjectAssociation',
        'updateObjectAssociation',
        'deleteObjectAssociation',
        'searchAssociationType'
      ]),
      async initData() {
        const asstList = await this.searchObjectAssociation({
          params: {
            condition: {
              id: this.asstId
            }
          }
        })
        if (this.asstId !== '') {
          this.relationInfo = asstList.find(asst => asst.id === this.asstId)
          this.relationInfoSnapshot = JSON.stringify(this.relationInfo)
        }
      },
      async initRelationList() {
        const data = await this.searchAssociationType({ params: {} })
        const relationList = data.info.map(({ bk_asst_id: asstId, bk_asst_name: asstName }) => {
          if (asstName.length) {
            return {
              id: asstId,
              name: `${asstId}(${asstName})`
            }
          }
          return {
            id: asstId,
            name: asstId
          }
        })
        this.relationList = Object.freeze(relationList)
      },
      getModelName(objId) {
        const model = this.models.find(model => model.bk_obj_id === objId)
        if (model) {
          return model.bk_obj_name
        }
        return ''
      },
      async saveRelation() {
        await this.updateObjectAssociation({
          id: this.relationInfo.id,
          params: {
            bk_obj_asst_name: this.relationInfo.bk_obj_asst_name
          },
          config: {
            requestId: 'updateObjectAssociation'
          }
        })
        this.$emit('save', {
          type: 'update',
          params: {
            ...this.relationInfo
          }
        })
        this.$emit('cancel')
      },
      deleteRelation() {
        this.$bkInfo({
          title: this.$t('确定删除关联关系?'),
          confirmFn: async () => {
            await this.deleteObjectAssociation({
              id: this.relationInfo.id,
              config: {
                requestId: 'deleteObjectAssociation'
              }
            })
            this.$emit('save', {
              type: 'delete',
              params: {
                ...this.relationInfo,
                id: this.asstId
              }
            })
            this.$emit('cancel')
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .button-delete {
        background-color: #fff;
        color: #ff5656;
    }
</style>
