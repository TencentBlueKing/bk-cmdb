<script>
  import { defineComponent, ref } from 'vue'

  export default defineComponent({
    props: {
      model: {
        type: Object,
        default: () => ({})
      },
      inst: {
        type: Object,
        default: () => ({})
      },
      topologyList: {
        type: Array,
        default: () => ([])
      }
    },
    setup(props, { emit }) {
      const topHeight = ref(80)

      const handlePathClick = (path) => {
        emit('path-click', path)
      }

      return {
        topHeight,
        handlePathClick
      }
    }
  })
</script>

<template>
  <div class="model-base-info">
    <div class="basic">
      <i :class="['model-icon', model.icon]"></i>
      <span class="inst-name">{{inst.name}}</span>
      <span class="model-name">{{model.name}}</span>
    </div>
    <div class="topology">
      <div class="topology-label">
        {{$t('所属拓扑')}}:
      </div>
      <ul class="topology-list">
        <li :class="['topology-item']"
          v-for="(item, index) in topologyList"
          :key="index">
          <span class="topology-path" v-bk-overflow-tips @click="handlePathClick(item)">{{item.path}}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.model-base-info {
  .basic {
    display: flex;
    align-items: center;

    .model-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      color: #3A84FF;
      width: 38px;
      height: 38px;
      background: #FFFFFF;
      border: 1px solid #DDE4EB;
      border-radius: 50%;
      font-size: 16px;
    }

    .inst-name {
      font-weight: 700;
      color: #313238;
      margin-left: 8px;
    }

    .model-name {
      font-size: 12px;
      color: #14A568;
      padding: 0 10px;
      height: 22px;
      line-height: 22px;
      background: #E4FAF0;
      border-radius: 2px;
      margin-left: 6px;
    }
  }

  .topology {
    display: flex;
    font-size: 12px;
    margin: 4px 0 0 46px;

    .topology-label {
      font-weight: 700;
    }

    .topology-list {
      margin-left: 4px;
      .topology-path {
        cursor: pointer;
        &:hover {
          color: $primaryColor;
        }
      }
    }
  }
}
</style>
