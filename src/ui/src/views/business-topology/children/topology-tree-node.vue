<script>
  import { computed, defineComponent } from 'vue'
  import store from '@/store'
  import { t } from '@/i18n'
  import routerActions from '@/router/actions'
  import CmdbLoading from '@/components/loading/loading'
  import { MENU_BUSINESS_SET_TEMPLATE_DETAILS } from '@/dictionary/menu-symbol'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants'
  import { CONTAINER_OBJECTS } from '@/dictionary/container'
  import { getContainerObjectNames } from '@/service/container/common'

  export default defineComponent({
    components: {
      CmdbLoading
    },
    props: {
      node: {
        type: Object,
        default: () => ({})
      },
      data: {
        type: Object,
        default: () => ({})
      },
      isBlueKing: Boolean,
      editable: Boolean,
      nodeCountType: String
    },
    setup(props, { emit }) {
      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const getInternalNodeClass = (data) => {
        const nodeIconMap = {
          1: 'icon-cc-host-free-pool',
          2: 'icon-cc-host-breakdown',
          default: 'icon-cc-host-free-pool'
        }
        return nodeIconMap[data.default] || nodeIconMap.default
      }

      const getNodeIconTips = (node, data) => {
        const { bk_obj_id: objId, bk_obj_name: objName } = data

        // 容器节点
        if (data.is_container) {
          // 大类型为workload
          if (data.is_workload) {
            return `${getContainerObjectNames(CONTAINER_OBJECTS.WORKLOAD).FULL} (${getContainerObjectNames(objId).FULL})`
          }

          const tipsMap = {
            [CONTAINER_OBJECTS.CLUSTER]: getContainerObjectNames(objId).FULL,
            [CONTAINER_OBJECTS.FOLDER]: t('无Pod运行的节点'),
            [CONTAINER_OBJECTS.NAMESPACE]: getContainerObjectNames(objId).FULL
          }
          return tipsMap[objId]
        }

        // 传统节点
        const tipsMap = {
          [BUILTIN_MODELS.BUSINESS]: t('业务'),
          [BUILTIN_MODELS.SET]: `${t('集群')}${isTemplate(node) ? t('（模板创建）') : t('（自定义创建）')}`,
          [BUILTIN_MODELS.MODULE]: `${t('模块')}${isTemplate(node) ? t('（模板创建）') : t('（自定义创建）')}`,
        }
        return tipsMap[objId] || objName
      }

      const isTemplate = node => node.data.service_template_id || node.data.set_template_id

      const isShowCreate = (node, data) => {
        const isModule = data.bk_obj_id === 'module'
        const isIdleSet = data.is_idle_set
        const isContainer = data.is_container
        return !isModule && !isIdleSet && !isContainer
      }

      const getSetNodeTips = (node) => {
        const tips = document.createElement('div')
        const span = document.createElement('span')
        span.innerText = t('需在集群模板中新建')
        const link = document.createElement('a')
        link.innerText = t('立即跳转')
        link.href = 'javascript:void(0)'
        link.style.color = '#3a84ff'
        link.addEventListener('click', () => {
          routerActions.redirect({
            name: MENU_BUSINESS_SET_TEMPLATE_DETAILS,
            params: {
              templateId: node.data.set_template_id
            },
            history: true
          })
        })
        tips.appendChild(span)
        tips.appendChild(link)
        return tips
      }

      const getNodeCount = (data) => {
        const count = data[props.nodeCountType]
        if (typeof count === 'number') {
          return count
        }
        return 0
      }

      const handleSetNodeTipsToggle = (tips) => {
        const element = tips.reference.parentElement
        if (tips.state.isVisible) {
          element.classList.remove('hovering')
        } else {
          element.classList.add('hovering')
        }
        return true
      }

      const handleCreate = (node) => {
        emit('create', node)
      }

      return {
        bizId,
        getInternalNodeClass,
        getNodeIconTips,
        isTemplate,
        isShowCreate,
        getSetNodeTips,
        getNodeCount,
        handleSetNodeTipsToggle,
        handleCreate
      }
    }
  })
</script>

<template>
  <div :class="['topology-tree-node', { 'is-container': node.data.is_container, 'is-selected': node.selected }]">
    <div
      v-bk-tooltips.top.light="{ content: getNodeIconTips(node, data), disabled: data.default !== 0 }"
      :class="['node-icon', {
        'is-selected': node.selected,
        'is-template': isTemplate(node),
        'is-internal': data.default !== 0
      }]">
      <i v-if="data.is_folder" class="icon-cc-pod-folder"></i>
      <i v-else-if="data.default !== 0" :class="getInternalNodeClass(data)"></i>
      <span v-else>{{data.icon_text || data.bk_obj_name[0]}}</span>
    </div>

    <span class="node-name" :title="node.name">{{node.name}}</span>

    <div class="node-extra">
      <cmdb-auth v-if="isShowCreate(node, data)"
        class="node-create-trigger"
        :auth="{ type: $OPERATION.C_TOPO, relation: [bizId] }">
        <template #default="{ disabled }">
          <i v-if="isBlueKing && !editable"
            class="node-button disabled-node-button"
            v-bk-tooltips.top="{ content: $t('蓝鲸业务拓扑节点提示'), interactive: false }">
            {{$t('新建')}}
          </i>
          <i v-else-if="data.set_template_id"
            class="node-button disabled-node-button"
            v-bk-tooltips.top="{
              content: getSetNodeTips(node),
              interactive: true,
              onShow: handleSetNodeTipsToggle,
              onHide: handleSetNodeTipsToggle
            }">
            {{$t('新建')}}
          </i>
          <bk-button v-else class="node-button" v-test-id="'createNode'"
            theme="primary"
            :disabled="disabled"
            @click.stop="handleCreate(node)">
            {{$t('新建')}}
          </bk-button>
        </template>
      </cmdb-auth>

      <cmdb-loading :class="['node-count', { 'is-selected': node.selected }]"
        :loading="['pending', undefined].includes(data.status)">
        {{getNodeCount(data)}}
      </cmdb-loading>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.topology-tree-node {
  display: flex;

  &:hover {
    .node-extra {
      .node-create-trigger {
        display: inline-block;
        & ~ .node-count {
          display: none;
        }
      }
    }
  }

  &.is-container {
    .node-icon {
      border-radius: 4px;
    }
  }

  .node-icon {
    display: flex;
    flex: none;
    width: 20px;
    height: 20px;
    line-height: 20px;
    align-items: center;
    justify-content: center;
    margin: 8px 4px 8px 0;
    border-radius: 50%;
    background-color: #C4C6CC;
    font-size: 12px;
    font-style: normal;
    color: #FFF;
    &.is-template {
      background-color: #97aed6;
    }
    &.is-selected {
      background-color: #3A84FF;
      &.is-internal {
        color: #3A84FF;
      }
    }
    &.is-internal {
      font-size: 14px;
      color: #63656E;
      background-color: unset;
      &:hover {
        background-color: unset;
      }
    }

    &:hover {
      background-color: #3A84FF;
    }
  }

  .node-name {
    display: block;
    height: 36px;
    line-height: 36px;
    overflow: hidden;
    @include ellipsis;
  }

  .node-extra {
    margin-left: auto;
    display: flex;

    .node-create-trigger {
      display: none;
      font-size: 0;
      &.hovering {
        display: inline-block;
        & ~ .node-count {
          display: none;
        }
      }
      .node-button {
        height: 24px;
        padding: 0 6px;
        margin: 0 20px 0 4px;
        line-height: 22px;
        border-radius: 4px;
        font-size: 12px;
        min-width: auto;
        &.disabled-node-button {
          @include inlineBlock;
          line-height: 24px;
          font-style: normal;
          background-color: #dcdee5;
          color: #ffffff;
          outline: none;
          cursor: not-allowed;
        }
      }
    }

    .node-count {
      padding: 0 5px;
      margin: 9px 20px 9px 4px;
      height: 18px;
      line-height: 17px;
      border-radius: 2px;
      background-color: #f0f1f5;
      color: #979ba5;
      font-size: 12px;
      text-align: center;
      &.is-selected {
        background-color: #a2c5fd;
        color: #fff;
      }
      &.loading {
        background-color: transparent;
      }
    }
  }
}
</style>
