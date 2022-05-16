/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
        this.slider.props.loading = true
        await this.$store.dispatch('hostUpdate/updateHost', {
          params: {
            ...changedValues,
            bk_host_id: this.selection.map(row => row.host.bk_host_id).join(',')
          }
        })
        this.slider.show = false
        this.slider.props.loading = false
        RouterQuery.set({
          _t: Date.now(),
          page: 1
        })
        this.$success(this.$t('修改成功'))
      } catch (e) {
        this.slider.props.loading = false
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
