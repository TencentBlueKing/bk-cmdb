<template>
  <bk-select
    v-if="displayType === 'selector'"
    multiple
    searchable
    display-tag
    v-bind="$attrs"
    v-model="localValue"
    @clear="() => $emit('clear')"
    @toggle="handleToggle">
    <bk-option
      v-for="template in list"
      :key="template.id"
      :id="template.id"
      :name="template.name">
    </bk-option>
  </bk-select>
  <span v-else>
    <slot name="info-prepend"></slot>
    {{info}}
  </span>
</template>

<script>
  import activeMixin from './mixins/active'
  import { mapGetters } from 'vuex'
  const requestId = Symbol('serviceTemplate')
  export default {
    name: 'cmdb-search-service-template',
    mixins: [activeMixin],
    props: {
      value: {
        type: [Array, String],
        default: () => ([])
      },
      displayType: {
        type: String,
        default: 'selector',
        validator(type) {
          return ['selector', 'info'].includes(type)
        }
      }
    },
    data() {
      return {
        list: [],
        requestId
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      localValue: {
        get() {
          const { value } = this
          if (Array.isArray(value)) {
            return value
          }
          return value.split(',').map(id => parseInt(id, 10))
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      },
      info() {
        const info = []
        this.localValue.forEach((id) => {
          const data = this.list.find(data => data.id === id)
          data && info.push(data.name)
        })
        return info.join(' | ') || '--'
      }
    },
    created() {
      this.getServiceTemplate()
    },
    methods: {
      async getServiceTemplate() {
        try {
          const { info } = await this.$store.dispatch('serviceTemplate/searchServiceTemplate', {
            params: {
              bk_biz_id: this.bizId
            },
            config: {
              requestId: this.requestId,
              fromCache: true
            }
          })
          this.list = this.$tools.localSort(info, 'name')
        } catch (error) {
          console.error(error)
          this.list = []
        }
      }
    }
  }
</script>
