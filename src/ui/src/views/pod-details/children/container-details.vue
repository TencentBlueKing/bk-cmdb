<script>
  import { computed, defineComponent, ref, watch, watchEffect } from 'vue'
  import store from '@/store'
  import routerActions from '@/router/actions'
  import RouterQuery from '@/router/query'
  import DetailsGeneralLayout from '@/components/ui/details/general-layout.vue'
  import TopoInstBaseInfo from '@/components/ui/details/topo-inst-base-info.vue'
  import ModelInstanceProperty from '@/components/model-instance/property.vue'
  import { CONTAINER_OBJECTS, CONTAINER_OBJECT_NAMES } from '@/dictionary/container'
  import { MENU_POD_DETAILS } from '@/dictionary/menu-symbol'
  import containerPropertyService from '@/service/container/property.js'
  import containerPropertGroupService from '@/service/container/property-group.js'
  import containerPodService from '@/service/container/pod.js'
  import containerConService from '@/service/container/container.js'

  export default defineComponent({
    components: {
      DetailsGeneralLayout,
      TopoInstBaseInfo,
      ModelInstanceProperty,
    },
    setup() {
      const requestIds = {
        property: Symbol(),
        list: Symbol()
      }

      const route = computed(() => RouterQuery.route)
      const active = ref(RouterQuery.get('tab', 'property'))

      const bizId = computed(() => store.getters['objectBiz/bizId'])
      const podId = computed(() => parseInt(route.value.params.podId, 10))
      const containerId = computed(() => parseInt(route.value.params.containerId, 10))

      const properties = ref([])
      const propertyGroups = ref([])

      const instance = ref({})
      const topoPaths = ref([])
      const model = ref({
        icon: 'icon-cc-container',
        name: CONTAINER_OBJECT_NAMES[CONTAINER_OBJECTS.CONTAINER].FULL
      })

      watchEffect(async () => {
        const podProperties = await containerPropertyService.getMany({
          objId: CONTAINER_OBJECTS.CONTAINER
        }, {
          requestId: requestIds.property,
          fromCache: true
        })
        properties.value = podProperties

        const objPropertyGroups = await containerPropertGroupService.getMany({ objId: CONTAINER_OBJECTS.CONTAINER })
        propertyGroups.value = objPropertyGroups

        // 实例
        const params = {
          id: containerId.value,
          podId: podId.value,
          bizId: bizId.value
        }
        instance.value = await containerConService.getOne(params)

        // 拓扑路径
        const { info } = await containerPodService.getPodPath({
          bk_biz_id: bizId.value,
          ids: [podId.value]
        })
        topoPaths.value = info
      })

      const topologyList = computed(() => topoPaths.value.map(topo => ({
        ...topo,
        path: `${topo.biz_name} / ${topo.cluster_name} / ${topo.namespace} / ${topo.workload_name}`
      })))

      watch(active, (active) => {
        RouterQuery.set('tab', active)
      })

      const handlePathClick = () => {
        routerActions.redirect({
          name: MENU_POD_DETAILS,
          params: {
            bizId: bizId.value,
            podId: podId.value
          },
          history: true
        })
      }

      return {
        active,
        properties,
        propertyGroups,
        instance,
        model,
        topologyList,
        handlePathClick
      }
    }
  })
</script>

<template>
  <details-general-layout>
    <template #top>
      <topo-inst-base-info
        :model="model"
        :inst="instance"
        :topology-list="topologyList"
        @path-click="handlePathClick">
      </topo-inst-base-info>
    </template>
    <template #main>
      <bk-tab class="details-tab"
        type="unborder-card"
        :active.sync="active">
        <bk-tab-panel name="property" :label="$t('Container属性')">
          <model-instance-property
            :properties="properties"
            :property-groups="propertyGroups"
            :inst="instance"
            :readonly="true">
          </model-instance-property>
        </bk-tab-panel>
      </bk-tab>
    </template>
  </details-general-layout>
</template>
