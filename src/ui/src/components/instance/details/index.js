import useProperty from '@/hooks/model/property'
import useGroup from '@/hooks/model/group'
import usePending from '@/hooks/utils/pending'
import useInstance from '@/hooks/instance/instance'
import { reactive, computed, toRef } from '@vue/composition-api'
import Vue from 'vue'
import store from '@/store'
import i18n from '@/i18n'
const state = reactive({
  bk_obj_id: null,
  bk_inst_id: null,
  bk_biz_id: null,
  visible: false,
  title: null
})
const visible = toRef(state, 'visible')
const title = toRef(state, 'title')
const modelOptions = computed(() => ({ bk_obj_id: state.bk_obj_id, bk_biz_id: state.bk_biz_id }))
const instanceOptions = computed(() => ({ bk_inst_id: state.bk_inst_id, ...modelOptions.value }))
const [{ properties, pending: propertyPending }] = useProperty(modelOptions)
const [{ groups, pending: groupPneding }] = useGroup(modelOptions)
const [{ instance, pending: instancePending }] = useInstance(instanceOptions)
const pending = usePending([propertyPending, groupPneding, instancePending], true)

const createDetails = () => {
  const Component = Vue.extend({
    mounted() {
      document.body.appendChild(this.$el)
    },
    render() {
      const directives = [{ name: 'bkloading', value: { isLoading: pending.value } }]
      const close = () => {
        visible.value = false
        setTimeout(() => this.$destroy(), 200)
      }
      return (
        <bk-sideslider
          is-show={ visible.value }
          width={ 700 }
          title={ title.value }
          { ...{ on: { 'update:isShow': close } } }>
          <cmdb-details slot="content"
            { ...{ directives } }
            show-options={ false }
            inst={ instance.value }
            properties={ properties.value }
            property-groups={ groups.value }>
          </cmdb-details>
        </bk-sideslider>
      )
    }
  })
  return new Component({ store, i18n })
}
export default function (options) {
  Object.assign(state, options)
  const details = createDetails()
  details.$mount()
  visible.value = true
}
