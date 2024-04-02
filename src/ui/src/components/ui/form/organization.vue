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
  <div :class="[
    'cmdb-organization',
    size,
    {
      'is-focus': focused,
      'is-disabled': disabled,
      'is-unselected': unselected
    }
  ]">
    <bk-select
      :class="['selector', { 'active': isActive }]"
      ref="select"
      :popover-width="400"
      :scroll-height="400"
      :popover-options="{
        sticky: true,
        flipOnUpdate: true
      }"
      :searchable="true"
      :multiple="localMultiple"
      font-size="normal"
      v-model="checked"
      :remote-method="handleSearch"
      :display-tag="true"
      :disabled="disabled"
      :tag-fixed-height="false"
      :show-empty="false"
      :placeholder="placeholder"
      :loading="$loading([searchRequestId])"
      @tab-remove="handleSelectTabRemove"
      @clear="handleClear"
      @toggle="handleToggle">
      <bk-big-tree
        ref="tree"
        class="big-tree"
        v-bind="treeProps"
        :use-default-empty="true"
        @check-change="handleTreeCheckChange">
      </bk-big-tree>
    </bk-select>
  </div>
</template>

<script setup>
  import { computed, ref } from 'vue'
  import { useStore } from '@/store'
  import debounce from 'lodash.debounce'
  import { isEmptyPropertyValue } from '@/utils/tools'

  const props = defineProps({
    value: {
      type: [Array, String, Number],
      default: () => []
    },
    disabled: {
      type: Boolean,
      default: false
    },
    readonly: Boolean,
    multiple: {
      type: Boolean,
      default: false
    },
    clearable: Boolean,
    size: String,
    placeholder: {
      type: String,
      default: ''
    },
    zIndex: {
      type: Number,
      default: 2500
    },
    formatter: Function
  })

  const emit = defineEmits(['on-checked', 'input', 'toggle'])

  defineExpose({
    focus: () => select?.value?.show?.()
  })

  const focused = ref(false)
  const unselected = ref(false)

  const tree = ref(null)
  const select = ref(null)

  const store = useStore()
  const searchRequestId = Symbol('orgSearch')

  const initValue = props.value
  const localMultiple = computed(() => {
    if (Array.isArray(initValue) && initValue.length > 1 && !props.multiple) {
      return true
    }
    return props.multiple
  })

  const treeProps = computed(() => ({
    showCheckbox: localMultiple.value,
    checkOnClick: !localMultiple.value,
    checkStrictly: false,
    showLinkLine: false,
    enableTitleTip: true,
    expandOnClick: localMultiple.value,
    selectable: !localMultiple.value,
    lazyMethod,
    lazyDisabled
  }))

  const checked = computed({
    get() {
      if (localMultiple.value) {
        if (this.value && !Array.isArray(this.value)) {
          return [this.value]
        }
        return this.value || []
      }
      if (Array.isArray(this.value)) {
        return this.value[0] || ''
      }
      return this.value || ''
    },
    set(value) {
      emit('on-checked', value)
      emit('change', value)
      emit('input', value)
    }
  })

  const loadTree = async () => {
    const { data: topData } = await getLazyData()

    if (!isEmptyPropertyValue(checked.value)) {
      const defaultChecked = Array.isArray(checked.value) ? checked.value : [checked.value]
      const checkedRes = await getSearchData({
        lookup_field: 'id',
        exact_lookups: defaultChecked.join(','),
        with_ancestors: true
      })

      // 可能为空，节点数据存在才获取相关联数据
      const chcekedData = checkedRes.results || []
      if (chcekedData.length) {
        // 已选中节点的树形数据
        const chekcedTreeData = getTreeSearchData(chcekedData)
        // 已选中节点的完整树形数据（含兄弟节点）
        const fullCheckedTreeData = await getCheckedFullTreeData(chekcedTreeData)
        // 将匹配的树分支替换以合并
        fullCheckedTreeData.forEach((checkedNode) => {
          const matchedIndex = topData.findIndex(top => top.id === checkedNode.id)
          if (matchedIndex !== -1) {
            topData[matchedIndex] = checkedNode
          }
        })

        // 设置树数据，选中节点数据已被完整包含
        setTreeData(topData)

        // 将选中节点全部展开
        defaultChecked.forEach(id => tree.value.setExpanded(id))
        // 设置为选中状态
        tree.value.setChecked(defaultChecked)
      } else {
        setTreeData(topData)
      }
    } else {
      setTreeData(topData)
    }
  }

  const setTreeData = (data) => {
    tree?.value?.setData?.(data)

    // hack将默认tree.setData注册给select的选项中的name替换为full_name，子级被选中时要显示完全名称
    const replaceOptionName = (nodes) => {
      nodes.forEach((node) => {
        select.value.optionsMap[node.id].name = node.full_name.split('/').join(' / ')
        if (node.children) {
          replaceOptionName(node.children)
        }
      })
    }
    replaceOptionName(data)
  }

  // 懒加载数据，结果中不含关联数据只是一层，用于点击时展示下一层
  const getLazyData = async (parentId) => {
    try {
      const params = {
        lookup_field: 'level',
        exact_lookups: 0
      }
      const config = {
        fromCache: !parentId,
        requestId: `get_org_department_${!parentId ? '0' : parentId}`
      }
      if (parentId) {
        params.lookup_field = 'parent'
        params.exact_lookups = parentId
      }
      const res = await store.dispatch('organization/getDepartment', { params, ...config })
      const data = res.results || []
      return { data }
    } catch (e) {
      console.error(e)
    }
  }

  // 搜索数据，结果中会带上相关联的层级数据，用于搜索和加载默认选中数据，这两种场景都需要回溯展示
  const getSearchData = async params => store.dispatch('organization/getDepartment', {
    params,
    requestId: searchRequestId
  })

  const resetTree = () => {
    checked.value = []
    loadTree()
  }

  const lazyMethod = async (node) => {
    const results = await getLazyData(node.id)
    const { data } = results

    // 默认tree.setData会执行select的registerOption，此处直接注册在树的lazyMethod方法需要手动执行
    data.forEach((node) => {
      select.value.registerOption({
        id: node.id,
        name: node.full_name.split('/').join(' / '),
        disabled: false,
        unmatched: false,
        isHighlight: false
      })
    })

    return results
  }

  const lazyDisabled = node => !node.data.has_children

  const getCheckedFullTreeData = async (chekcedTreeData) => {
    // 获取所有节点id
    const ids = []
    const getId = (nodes) => {
      nodes.forEach((node) => {
        ids.push(node.id)
        if (node.children) {
          getId(node.children)
        }
      })
    }
    getId(chekcedTreeData)

    // 获取所有节点的子节点
    const childNodeRes = await getSearchData({
      lookup_field: 'parent',
      exact_lookups: ids.join(','),
      with_ancestors: false
    })
    const childNodeList = childNodeRes.results || []

    // 将子节点补齐到对应的目标节点
    const appendChild = (nodes) => {
      nodes.forEach((node) => {
        childNodeList.forEach((child) => {
          if (child.parent === node.id) {
            if (node.children) {
              const childIds = node.children.map(item => item.id)
              if (childIds.indexOf(child.id) === -1) {
                node.children.push(child)
              }
            } else {
              node.children = [child]
            }
          }
        })

        if (node.children) {
          appendChild(node.children)
        }
      })
    }
    appendChild(chekcedTreeData)

    return chekcedTreeData
  }

  const getTreeSearchData = (data) => {
    // 将偏平的数据组装成树形结构
    const treeData = []
    data.forEach((item) => {
      const ancestorLength = item.ancestors.length
      const curNode = {
        id: item.id,
        name: item.name,
        level: ancestorLength,
        full_name: item.full_name.split('/').join(' / ')
      }
      const ids = [curNode.id]
      const treeNode = {}
      for (let i = ancestorLength - 1; i >= 0; i--) {
        const node = item.ancestors[i]
        ids.push(node.id)
        node.level = i
        node.children = [item.ancestors[i + 1] ? item.ancestors[i + 1] : curNode]
        node.full_name = item.full_name.split('/', i + 1).join(' / ')
      }

      treeNode.ids = ids.reverse()
      if (item.ancestors[0]) {
        // eslint-disable-next-line prefer-destructuring
        treeNode.map = item.ancestors[0]
      } else {
        treeNode.map = curNode
      }

      treeData.push(treeNode)
    })

    // 合并与去重
    for (let i = 0; i < treeData.length; i++) {
      const node = treeData[i]
      const path = node.ids.join('-')
      for (let j = i + 1; j < treeData.length; j++) {
        const nodeNext = treeData[j]
        let k = nodeNext.ids.length
        while (k) {
          const pathNext = nodeNext.ids.slice(0, k).join('-')
          // 路径比较，将被比较对象除重复部分外的数据合并至比较对象
          if (path.indexOf(pathNext) !== -1) {
            const nextRest = data[j].ancestors.slice(k - 1)
            const appendToNode = data[i].ancestors[k - 1]
            if (appendToNode && nextRest.length) {
              // 合并时去重
              const exists = appendToNode.children.map(item => item.id)
              nextRest[0].children.forEach((item) => {
                if (exists.indexOf(item.id) === -1) {
                  appendToNode.children.push(item)
                }
              })
            }
            nodeNext.remove = true
            break
          } else if (pathNext.indexOf(path) !== -1) {
            // 如果路径被反向包含则可直接删除
            node.remove = true
            break
          }
          k-- // eslint-disable-line no-plusplus
        }
      }
    }

    // 得到最终用于树的数据
    const finalTreeData = treeData.filter(item => !item.remove).map(item => item.map)
    return finalTreeData
  }

  const setTreeSearchData = (data) => {
    const finalTreeData = getTreeSearchData(data)
    setTreeData(finalTreeData)
    finalTreeData.forEach((node) => {
      tree.value.setExpanded(node.id)
    })
  }

  const handleSelectTabRemove = (options) => {
    if (!tree.value?.getNodeById(options.id)) {
      return
    }

    tree.value?.setChecked(options.id, { emitEvent: true, checked: false })
  }

  const handleTreeCheckChange = (ids, node) => {
    if (localMultiple.value) {
      checked.value = [...ids]
    } else {
      checked.value = [node.id]
      select.value?.close()
    }
  }

  const handleSearch = debounce(async (value) => {
    const keyword = value.trim()
    tree?.value?.filter?.(keyword)
    try {
      if (keyword.length) {
        const res = await getSearchData({
          lookup_field: 'name',
          fuzzy_lookups: keyword,
          with_ancestors: true
        })
        const data = res.results || []
        setTreeSearchData(data)
      } else if (!value.length) {
        loadTree()
      }
    } catch (e) {
      console.error(e)
    }
  }, 160)

  const handleClear = () => {
    resetTree()
  }
  const isActive = ref(false)

  const handleToggle = (active) => {
    isActive.value = active
    emit('toggle', active)
  }

  loadTree()
</script>

<script>
  export default {
    name: 'cmdb-form-organization'
  }
</script>

<style lang="scss" scoped>
.cmdb-organization {
  position: relative;
  width: 100%;
  height: 32px;
  .selector {
    width: 100%;
    &.active {
      position: absolute;
      z-index: 2;
    }
  }
}
:deep(.bk-big-tree-empty) {
  position: static;
}
</style>
