<script>
  import { defineComponent } from '@vue/composition-api'

  export default defineComponent({
    props: {
      auth: Object,
      authorized: Boolean,
      tag: {
        type: String,
        default: 'span'
      }
    },
    render(h) {
      return h(this.tag, {
        directives: [
          {
            name: 'cursor',
            value: {
              auth: this.auth,
              active: !this.authorized
            }
          }
        ],
        staticClass: 'auth-mask'
      }, this.$slots.default)
    }
  })
</script>
