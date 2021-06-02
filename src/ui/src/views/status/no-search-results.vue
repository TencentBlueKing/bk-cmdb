<template functional>
  <div class="no-search-results">
    <img src="../../assets/images/full-text-search.png" alt="no-search-results">
    <p>{{props.text}}</p>
  </div>
</template>

<style lang="scss" scoped>
  .no-search-results {
    width: 90%;
    margin: 0 auto;
    padding: 32px 0;
    text-align: center;
    color: #63656E;
    font-size: 16px;
    img {
      width: 104px;
    }
  }
</style>
