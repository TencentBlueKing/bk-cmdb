<script>
  import { computed, defineComponent } from 'vue'
  import store from '@/store'

  export default defineComponent({
    setup() {
      const selectedNode = computed(() => store.getters['businessHost/selectedNode'])
      const isFolder = computed(() =>  selectedNode.value?.data?.is_folder)
      const nodeId = computed(() =>  (isFolder.value ? '--' : selectedNode.value?.data?.bk_inst_id ?? '--'))
      const nodeName = computed(() => selectedNode.value?.data?.bk_inst_name ?? '--')

      return {
        nodeId,
        nodeName
      }
    }
  })
</script>

<template>
  <div class="default-node-info">
    <div class="info-item">
      <label class="name">{{$t('ID')}}:</label>
      <span class="value">{{nodeId}}</span>
    </div>
    <div class="info-item">
      <label class="name">{{$t('节点名称')}}</label>
      <span class="value">{{nodeName}}</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.default-node-info {
  padding: 20px 0 20px 36px;
  display: flex;
  .info-item {
    flex: auto;
    max-width: 400px;
    font-size: 14px;
    .name {
      color: #63656e;
    }
    .value {
      color: #313238;
    }
  }
}
</style>
