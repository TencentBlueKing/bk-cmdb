<script>
  import { computed, defineComponent, ref, watch, watchEffect } from 'vue'
  import store from '@/store'
  import routerActions from '@/router/actions'
  import RouterQuery from '@/router/query'
  import DetailsGeneralLayout from '@/components/ui/details/general-layout.vue'
  import TopoInstBaseInfo from '@/components/ui/details/topo-inst-base-info.vue'
  import ModelInstanceProperty from '@/components/model-instance/property.vue'
  import ContainerList from './children/container-list.vue'
  import { CONTAINER_OBJECTS, CONTAINER_OBJECT_NAMES } from '@/dictionary/container'
  import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
  import containerPropertyService from '@/service/container/property.js'
  import containerPropertGroupService from '@/service/container/property-group.js'
  import containerPodService from '@/service/container/pod.js'

  export default defineComponent({
    components: {
      DetailsGeneralLayout,
      TopoInstBaseInfo,
      ModelInstanceProperty,
      ContainerList
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

      const properties = ref([])
      const propertyGroups = ref([])

      const instance = ref({})
      const topoPaths = ref([])
      const model = ref({
        icon: 'icon-cc-pod',
        name: CONTAINER_OBJECT_NAMES[CONTAINER_OBJECTS.POD].FULL
      })

      watchEffect(async () => {
        const podProperties = await containerPropertyService.getMany({
          objId: CONTAINER_OBJECTS.POD
        }, {
          requestId: requestIds.property,
          fromCache: true
        })
        properties.value = podProperties

        const objPropertyGroups = await containerPropertGroupService.getMany({ objId: CONTAINER_OBJECTS.POD })
        propertyGroups.value = objPropertyGroups

        // 实例
        const params = { id: podId.value, bizId: bizId.value }
        instance.value = await containerPodService.getOne(params)

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

      const handlePathClick = (path) => {
        routerActions.redirect({
          name: MENU_BUSINESS_HOST_AND_SERVICE,
          query: {
            node: `${path.kind}-${path.bk_workload_id}`,
            topo_path: `cluster-${path.bk_cluster_id},namespace-${path.bk_namespace_id},${path.kind}-${path.bk_workload_id}`
          },
          params: {
            bizId: bizId.value
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
        <bk-tab-panel name="property" :label="$t('Pod属性')">
          <model-instance-property
            :properties="properties"
            :property-groups="propertyGroups"
            :inst="instance"
            :readonly="true">
          </model-instance-property>
        </bk-tab-panel>
        <bk-tab-panel name="containers" label="Container(s)">
          <container-list v-if="active === 'containers'"></container-list>
        </bk-tab-panel>
      </bk-tab>
    </template>
  </details-general-layout>
</template>
