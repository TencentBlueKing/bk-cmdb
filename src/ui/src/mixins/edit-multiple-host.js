import RouterQuery from '@/router/query'
import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
export default {
  props: {
    properties: {
      type: Array,
      required: true
    },
    selection: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      slider: {
        show: false,
        title: '',
        component: null,
        props: {}
      },
      request: {
        propertyGroups: Symbol('propertyGroups')
      }
    }
  },
  computed: {
    isGlobalView() {
      // eslint-disable-next-line prefer-destructuring
      const topRoute = this.$route.matched[0]
      return topRoute ? topRoute.name !== MENU_BUSINESS : true
    },
    saveAuth() {
      return this.selection.map(({ host, biz, module }) => {
        const isBizHost = biz[0].default === 0
        if (isBizHost) {
          return {
            type: this.$OPERATION.U_HOST,
            relation: [biz[0].bk_biz_id, host.bk_host_id]
          }
        }
        return {
          type: this.$OPERATION.U_RESOURCE_HOST,
          relation: [module[0].bk_module_id, host.bk_host_id]
        }
      })
    }
  },
  methods: {
    async handleMultipleEdit() {
      try {
        this.slider.show = true
        this.slider.title = this.$t('主机属性')
        const groups = await this.getPropertyGroups()
        setTimeout(() => {
          this.slider.component = 'cmdb-form-multiple'
          this.slider.props.propertyGroups = groups
          this.slider.props.properties = this.properties
          this.slider.props.saveAuth = this.saveAuth
        }, 200)
      } catch (e) {
        console.error(e)
      }
    },
    async handleMultipleSave(changedValues) {
      try {
        await this.$store.dispatch('hostUpdate/updateHost', {
          params: {
            ...changedValues,
            bk_host_id: this.selection.map(row => row.host.bk_host_id).join(',')
          }
        })
        this.slider.show = false
        RouterQuery.set({
          _t: Date.now(),
          page: 1
        })
        this.$success(this.$t('修改成功'))
      } catch (e) {
        console.error(e)
      }
    },
    handleSliderBeforeClose() {
      const $form = this.$refs.multipleForm
      if (Object.keys($form.changedValues).length) {
        return new Promise((resolve) => {
          this.$bkInfo({
            title: this.$t('确认退出'),
            subTitle: this.$t('退出会导致未保存信息丢失'),
            extCls: 'bk-dialog-sub-header-center',
            confirmFn: () => {
              this.slider.show = false
              this.slider.component = null
              resolve(true)
            },
            cancelFn: () => {
              resolve(false)
            }
          })
        })
      }
      this.slider.component = null
      this.slider.show = false
    },
    getPropertyGroups() {
      return this.$store.dispatch('objectModelFieldGroup/searchGroup', {
        objId: 'host',
        params: this.isGlobalView ? {} : { bk_biz_id: parseInt(this.$route.params.bizId, 10) },
        config: {
          requestId: this.request.propertyGroups
        }
      })
    }
  }
}
