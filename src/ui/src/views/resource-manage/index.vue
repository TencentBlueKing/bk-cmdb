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
  <div class="classify-layout clearfix">
    <div class="classify-filter">
      <bk-input class="filter-input"
        clearable
        :placeholder="$t('请输入xx', { name: $t('关键字') })"
        right-icon="icon-search"
        v-model.trim="filter">
      </bk-input>
    </div>
    <div v-show="!isEmpty">
      <div class="classify-waterfall fl"
        v-for="col in classifyColumns.length"
        :key="col">
        <cmdb-classify-panel
          v-for="classify in classifyColumns[col - 1]"
          :key="classify['bk_classification_id']"
          :classify="classify"
          :collection="collection">
        </cmdb-classify-panel>
      </div>
    </div>
    <no-search-results v-if="isEmpty && !globalLoading" :text="$t('搜不到相关资源')" />
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import debounce from 'lodash.debounce'
  import noSearchResults from '@/views/status/no-search-results.vue'
  import cmdbClassifyPanel from './children/classify-panel'
  import useInstanceCount from './children/use-instance-count.js'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'

  export default {
    components: {
      cmdbClassifyPanel,
      noSearchResults
    },
    data() {
      return {
        filter: '',
        debounceFilter: null,
        matchedModels: null
      }
    },
    computed: {
      ...mapGetters(['globalLoading']),
      ...mapGetters('objectModelClassify', ['classifications', 'models']),
      ...mapGetters('userCustom', { collection: 'resourceCollection' }),
      filteredClassifications() {
        const result = []
        this.classifications.forEach((classification) => {
          const models = classification.bk_objects.filter((model) => {
            const isInvisible = model.bk_ishidden
            const isPaused = model.bk_ispaused
            const isMatched = this.matchedModels ? this.matchedModels.includes(model.bk_obj_id) : true

            // 集群/模块暂不允许查看实例
            const isModuleOrSet = [BUILTIN_MODELS.MODULE, BUILTIN_MODELS.SET].includes(model.bk_obj_id)

            return !isInvisible && !isPaused && isMatched && !isModuleOrSet
          })
          if (models.length) {
            result.push({
              ...classification,
              bk_objects: models
            })
          }
        })
        return result
      },
      modelIds() {
        return this.filteredClassifications.map(item => item.bk_objects.map(obj => obj.bk_obj_id))
      },
      classifyColumns() {
        const colHeight = [0, 0, 0, 0]
        const classifyColumns = [[], [], [], []]
        this.filteredClassifications.forEach((classify) => {
          const minColHeight = Math.min(...colHeight)
          const rowIndex = colHeight.indexOf(minColHeight)
          classifyColumns[rowIndex].push(classify)
          colHeight[rowIndex] = colHeight[rowIndex] + this.calcWaterfallHeight(classify)
        })
        return classifyColumns
      },
      isEmpty() {
        return this.classifyColumns.every(column => !column.length)
      }
    },
    watch: {
      filter() {
        this.debounceFilter()
      }
    },
    created() {
      this.debounceFilter = debounce(this.filterModel, 300)
      const { fetchData: getInstanceCount } = useInstanceCount({ modelIds: this.modelIds }, this)
      getInstanceCount()
    },
    methods: {
      filterModel() {
        if (this.filter) {
          const models = this.models.filter(model => model.bk_obj_name.indexOf(this.filter) > -1)
          this.matchedModels = models.map(model => model.bk_obj_id)
        } else {
          this.matchedModels = null
        }
      },
      calcWaterfallHeight(classify) {
        // 46px 分类高度
        // 16px 模型列表padding
        // 36 模型高度
        return 46 + 16 + (classify.bk_objects.length * 36)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .classify-layout{
        padding: 15px 20px 20px;
    }
    .classify-filter {
        padding: 0 20px 20px 0;
        .filter-input {
            width: 240px;
        }
    }
    .classify-waterfall{
        width: calc((100% - 80px) / 4);
        margin: 0 0 0 20px;
        &:first-child{
            margin: 0;
        }
    }
</style>
