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
  <div class="instance-association" :class="{ 'full-screen': fullScreen }">
    <bk-button class="exit" size="small" theme="default"
      v-show="fullScreen"
      @click="toggleFullScreen(false)">
      <i class="icon-cc-resize-small"></i>
      {{$t('退出')}}
    </bk-button>
    <div class="container" ref="container"></div>
    <a v-show="false" ref="tooltip" href="javascript:void(0)" @click="handleShowDetails">
      {{$t('详情信息')}}
    </a>
  </div>
</template>

<script>
  import useInstanceAssociation from '@/hooks/instance/association'
  import useAllAssocation from '@/hooks/association/all'
  import Graphics from './graphics'
  import { computed, ref } from '@vue/composition-api'
  import tippy from 'bk-magic-vue/lib/utils/tippy'
  export default {
    name: 'cmdb-instance-association',
    props: {
      bizId: Number,
      objId: {
        type: String,
        required: true
      },
      instId: {
        type: Number,
        required: true
      },
      instName: String
    },
    setup() {
      const [{ map: instanceMap, pending: instancePending }, findInstanceAssociation] = useInstanceAssociation()
      const [{ associations, pending: associationPending }, findAssociation] = useAllAssocation()
      const pending = computed(() => (instancePending.value || associationPending.value))
      return {
        pending,
        instanceMap,
        associations,
        findInstanceAssociation,
        findAssociation,
        fullScreen: ref(false)
      }
    },
    async mounted() {
      await this.findAssociation()
      this.graphics = new Graphics(this.$refs.container)
      this.graphics.on('click', 'node', this.handleNodeClick)
      this.graphics.on('mouseover', 'node', this.handleNodeMouseover)
      this.setAssociation({
        bk_obj_id: this.objId,
        bk_inst_id: this.instId,
        bk_inst_name: this.instName
      }, true)
    },
    methods: {
      getInstances(source) {
        return this.instanceMap[source.bk_obj_id][source.bk_inst_id] || []
      },
      getDirection(instance) {
        const define = this.associations.find(association => association.bk_asst_id === instance.bk_asst_id)
        if (!define) {
          return {
            label: instance.bk_asst_id,
            arrow: 'none'
          }
        }
        // 根据作为源还是目标，设置线条名称及是否翻转箭头指向
        return {
          label: instance.target ? define.src_des : define.dest_des,
          arrow: !instance.target && define.direction === 'src_to_dest' ? 'dest_to_src' : define.direction }
      },
      createNodes(instances, parentNode) {
        return instances.map((instance) => {
          const model = this.$store.getters['objectModelClassify/getModelById'](instance.bk_obj_id)
          return {
            data: {
              id: `${instance.bk_obj_id}_${instance.bk_inst_id}`,
              name: instance.bk_inst_name,
              icon: model ? model.bk_obj_icon : 'icon-cc-defalut',
              parentId: parentNode && parentNode.data.id,
              loaded: !parentNode,
              instance
            },
            selected: !parentNode,
            group: 'nodes'
          }
        })
      },
      createEdges(instances, parentNode) {
        return instances.map((instance) => {
          const parentId = parentNode.data.id
          const targetId = `${instance.bk_obj_id}_${instance.bk_inst_id}`
          const direction = this.getDirection(instance)
          return {
            data: {
              id: instance.id,
              source: parentId,
              target: targetId,
              label: direction.label,
              direction: direction.arrow
            },
            group: 'edges'
          }
        })
      },
      async setAssociation(source, isTop = false) {
        await this.findInstanceAssociation(source)
        const { data, root } = this.getInstances(source)
        const [rootNode] = this.createNodes([root], null)
        const nodes = this.createNodes(data, rootNode)
        const edges = this.createEdges(data, rootNode)
        const elements = this.graphics.add([rootNode, ...nodes, ...edges])
        if (isTop) {
          const [topElement] = elements.toArray()
          this.graphics.setCurrent(topElement)
        }
      },
      async handleNodeClick(event) {
        const node = event.target
        const data = event.target.data()
        if (data.loaded || data.instance.deleted) {
          return
        }
        await this.setAssociation({
          bk_obj_id: data.instance.bk_obj_id,
          bk_inst_id: data.instance.bk_inst_id,
          bk_inst_name: data.instance.bk_inst_name
        })
        node.data('loaded', true)
      },
      /**
       * 悬浮node时，延时添加显示详情的tooltip, 并销毁上一个节点的tooltip
       */
      handleNodeMouseover(event) {
        this.tooltip && this.tooltip.instance.destroy()
        const contentRef = this.$refs.tooltip
        contentRef.style.display = 'inline-block'
        this.tooltip = {
          instance: tippy(event.target.popperRef(), {
            content: contentRef,
            placement: 'right',
            sticky: true,
            hideOnClick: true,
            showOnInit: true,
            interactive: true,
            animateFill: false,
            arrow: true,
            theme: 'light'
          }),
          data: event.target.data('instance')
        }
      },
      async handleShowDetails() {
        const showInstanceDetails = await import('@/components/instance/details')
        const { data, instance } = this.tooltip
        const model = this.$store.getters['objectModelClassify/getModelById'](data.bk_obj_id)
        showInstanceDetails.default({
          bk_biz_id: this.bizId,
          bk_obj_id: data.bk_obj_id,
          bk_inst_id: data.bk_inst_id,
          title: `${model.bk_obj_name}-${data.bk_inst_name}`
        })
        instance.hide()
      },
      toggleFullScreen(fullScreen) {
        this.$store.commit('setLayoutStatus', { mainFullScreen: fullScreen })
        this.fullScreen = fullScreen
      }
    }
  }
</script>

<style lang="scss" scoped>
  .instance-association {
    width: 100%;
    height: 100%;
    background-color: #f9f9f9;
    &.full-screen {
      position: fixed;
      top: 0;
      right: 0;
      bottom: 0;
      left: 0;
      .exit {
        position: absolute;
        top: 20px;
        right: 20px;
        z-index: 9999;
      }
    }
    .container {
      width: 100%;
      height: 100%;
    }
  }
</style>
