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
  <div class="simple-node-info">
    <bk-alert type="info" class="alert-tips" :title="$t('空Pod节点提示语')" closable></bk-alert>
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
  </div>
</template>

<style lang="scss" scoped>
.simple-node-info {
  .alert-tips {
    margin-top: 24px;
  }
}
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
