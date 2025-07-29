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

<script setup>
  import { computed, ref, watch } from 'vue'
  import { useHttp, jsonp } from '@/api'
  import { t } from '@/i18n'
  import { useStore } from '@/store'
  import fieldTemplateService from '@/service/field-template'
  import queryBuilderOperator, { QUERY_OPERATOR } from '@/utils/query-builder-operator'
  import { BUILTIN_MODELS, UNCATEGORIZED_GROUP_ID } from '@/dictionary/model-constants'
  import { escapeRegexChar } from '@/utils/util'

  const props = defineProps({
    defaultFilter: {
      type: Array,
      default: () => []
    }
  })

  const store = useStore()
  const http = useHttp()

  const emit = defineEmits(['search'])

  const searchSelectComp = ref(null)
  const isTyeing = ref(false)
  const isFocus = ref(false)
  const filterMenus = [
    {
      id: 'templateName',
      name: t('模板名称'),
      multiable: true,
      remote: true
    },
    {
      id: 'modelName',
      name: t('模型名称'),
      multiable: true,
      remote: true
    },
    {
      id: 'modifier',
      name: t('更新人'),
      multiable: true,
      remote: true
    }
  ]
  const filter = ref([])

  watch(() => props.defaultFilter, (defaultFilter) => {
    filter.value = defaultFilter.map(item => ({
      id: item.id,
      name: filterMenus.find(menu => menu.id === item.id)?.name,
      values: (item.value?.split(',') || []).map(val => ({ name: val }))
    }))
  }, { immediate: true })

  const displayFilterMenus = computed(() => {
    if (filter.value?.length) {
      return filterMenus.filter(menu => !filter.value.some(item => item.id === menu.id))
    }
    return filterMenus.slice()
  })

  watch(displayFilterMenus, () => {
    if (!isFocus.value) {
      return
    }
    // fix菜单项丢失
    searchSelectComp.value?.showMenu()
  })

  const fetchOptions = async (val, menu) => {
    const fetchs = {
      templateName: fetchTemplate,
      modelName: fetchModel,
      modifier: fetchMember
    }

    return fetchs[menu.id](val, menu)
  }

  const fetchTemplate = async (val, menu) => {
    const params = {
      template_filter: {
        field: 'name',
        operator: queryBuilderOperator(QUERY_OPERATOR.LIKE),
        value: escapeRegexChar(val)
      },
      fields: ['id', 'name'],
      page: {
        start: 0,
        limit: 100,
        sort: 'id'
      }
    }

    if (!isTyeing.value || !val?.length || val === `${menu.name}：`) {
      Reflect.deleteProperty(params, 'template_filter')
    }

    const { list } = await fieldTemplateService.find(params, {
      cancelPrevious: true,
      globalPermission: false
    })
    return list
  }

  const fetchModel = async (val, menu) => {
    const classifications = await store.dispatch('objectModelClassify/searchClassificationsObjects', {
      config: { fromCache: true }
    })

    const excludeModelIds = (store.getters['objectMainLineModule/mainLineModels'] || [])
      .filter(model => model.bk_obj_id !== BUILTIN_MODELS.HOST)
      .map(model => model.bk_obj_id)

    const list = []
    classifications.forEach((classification) => {
      list.push({
        ...classification,
        bk_objects: classification.bk_objects
          .filter(model => !model.bk_ispaused && !model.bk_ishidden && !excludeModelIds.includes(model.bk_obj_id))
      })
    })

    const modelList = list
      .filter(item => item.bk_objects.length > 0)
      .sort((a, b) => (b.bk_classification_id === UNCATEGORIZED_GROUP_ID ? -1 : 0))
      .reduce((acc, cur) => acc.concat(cur.bk_objects), [])
      .map(model => ({
        id: model.id,
        name: model.bk_obj_name,
        bk_obj_id: model.bk_obj_id
      }))

    if (!isTyeing.value || !val?.length || val === `${menu.name}：`) {
      return modelList
    }

    const reg = new RegExp(val, 'i')
    return modelList.filter(item => reg.test(item.name))
  }

  const fetchMember = async (val, menu) => {
    let query = val
    if (!isTyeing.value || !query?.length || query === `${menu.name}：`) {
      query = 'a'
    }

    let result = []
    if (window.ESB.userManage) {
      const params = {
        app_code: 'bk-magicbox',
        page: 1,
        page_size: 100,
        fuzzy_lookups: query
      }
      const api = window.ESB.userManage
      const response = await jsonp(api, params)
      if (response.code !== 0) {
        console.error(response?.message)
        return []
      }
      result = (response?.data?.results || []).map(item => ({
        id: item.id,
        username: item.username,
        // name: `${item.username}${item.display_name ? `(${item.display_name})` : ''}`
        name: item.username
      }))
    } else {
      const data = await http.get(`${window.API_HOST}user/list`, {
        params: {
          fuzzy_lookups: val
        },
        config: {
          cancelPrevious: true
        }
      })
      result = (data || []).map(user => ({
        id: user.english_name,
        username: user.english_name,
        name: user.chinese_name
      }))
    }

    return result
  }

  const handleSearchSelectChange = (list) => {
    const ids = filterMenus.map(item => item.id)
    list.forEach((item) => {
      const nameItem = list.find(searchItem => searchItem.id === 'templateName')
      if (!ids.includes(item.id)) {
        // 存在多个就合并
        if (nameItem) {
          nameItem.values.push({ name: item.name })
        } else {
          item.id = 'templateName'
          item.values = [{ name: item.name }]
        }
      }
    })

    // 合并完之后，被合并的项还会遗留，在这里去掉
    const newList = list.filter(item => ids.includes(item.id))
    emit('search', newList)
  }
  const handleInputChange = () => {
    isTyeing.value = true
  }
  const handleInputFocus = () => {
    isFocus.value = true
  }
  const handleInputClickOutside = () => {
    isFocus.value = false
  }
</script>

<template>
  <bk-search-select
    ref="searchSelectComp"
    class="search-select"
    :clearable="false"
    :placeholder="$t('请输入模板名称/模型/更新人')"
    :filter="true"
    :show-condition="false"
    :show-popover-tag-change="false"
    :data="displayFilterMenus"
    v-model="filter"
    :remote-method="fetchOptions"
    @input-change="handleInputChange"
    @input-focus="handleInputFocus"
    @input-click-outside="handleInputClickOutside"
    @change="handleSearchSelectChange">
  </bk-search-select>
</template>

<style lang="scss" scoped>
.search-select {
  width: 480px;
  background: #fff;
}
</style>
<style lang="scss">
.bk-search-list {
  max-width: 360px;
  .bk-search-list-menu-item {
    width: 100%;
    .item-name {
      width: 100%;
      @include ellipsis;
    }
  }
}
</style>
