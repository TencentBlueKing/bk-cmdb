<template functional>
  <div class="top-steps">
    <bk-steps
      class="steps"
      :steps="[
        { title: '属性设置', icon: 1 },
        { title: '信息确认', icon: 2 },
        { title: '属性应用', icon: 3 }
      ]"
      :cur-step="props.current">
    </bk-steps>
  </div>
</template>

<style lang="scss" scoped>
  .top-steps {
    display: flex;
    height: 50px;
    align-items: center;
    justify-content: center;
    padding: 0 20px;
    background-color: #fff;
    border-bottom: 1px solid $borderColor;
    .steps {
      width: 65%;
    }
  }
</style>
