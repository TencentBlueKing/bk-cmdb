<template>
  <cmdb-form-objuser
    v-model="userValue"
    :exclude="false"
    :multiple="true"
    :placeholder="$t('请输入xx', { name: $t('账号') })">
    <bk-select class="user-option" slot="prepend"
      :clearable="false"
      v-model="selectValue">
      <bk-option id="in" :name="$t('包含')"></bk-option>
      <bk-option id="not_in" :name="$t('不包含')"></bk-option>
    </bk-select>
  </cmdb-form-objuser>
</template>

<script>
  export default {
    props: {
      value: {
        type: Array,
        default: []
      }
    },
    computed: {
      userValue: {
        get() {
          return this.value[1]
        },
        set(values) {
          this.$emit('input', [this.selectValue, values])
        }
      },
      selectValue: {
        get() {
          return this.value[0]
        },
        set(value) {
          this.$emit('input', [value, this.userValue])
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
  .user-option {
    width: 90px;
    border-color: #c4c6cc;
    box-shadow: none;
  }
</style>
