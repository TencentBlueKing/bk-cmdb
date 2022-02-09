/**
 * 主机详情只读状态 mixin
 */
export const readonlyMixin = {
  computed: {
    readonly() {
      return this.$route.meta.readonly
    }
  },
}
