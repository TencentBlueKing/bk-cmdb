<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="custom-layout">
    <div class="input-wrapper">
      <textarea class="ip-input"
        :placeholder="$t('请输入IP，换行分隔')"
        :class="{ 'has-error': invalidList.length }"
        v-model="value"
        @focus="handleFocus">
            </textarea>
      <p class="ip-error" v-show="invalidList.length">{{$t('IP不正确或主机不存在')}}</p>
      <bk-button class="ip-confirm" outline theme="primary" @click="handleConfirm">{{$t('添加至列表')}}</bk-button>
    </div>
    <div class="table-wrapper" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
      <host-table :list="hostList" :selected="selected" @select-change="handleHostSelectChange" />
    </div>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import HostTable from './host-table.vue'
  import isIP from 'validator/es/lib/isIP'

  export default {
    components: {
      HostTable
    },
    props: {
      selected: {
        type: Array,
        default: () => ([])
      }
    },
    data() {
      return {
        value: '',
        validList: [],
        invalidList: [],
        hostList: [],
        request: {
          host: Symbol('host')
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('businessHost', ['getDefaultSearchCondition'])
    },
    activated() {
      this.value = ''
      this.validList = []
      this.invalidList = []
    },
    methods: {
      handleFocus() {
        this.invalidList = []
      },
      async handleConfirm() {
        try {
          this.validateList()
          if (this.validList.length) {
            const results = await Promise.all(this.getReqs())
            const list = []
            results.forEach(({ info }) => {
              list.push(...info)
            })
            const unexistList = this.validList.filter((ip) => {
              const exist = list.some(target => target.host.bk_host_innerip === ip
                || target.host.bk_host_innerip_v6 === ip)
              return !exist
            })
            const newHostList = list.filter(({ host }) => !this.hostList
              .some(target => target.host.bk_host_id === host.bk_host_id))
            this.hostList.push(...newHostList)
            this.invalidList.push(...unexistList)
          }
          this.value = this.invalidList.join('\n')
        } catch (e) {
          console.error(e)
        }
      },
      validateList() {
        const list = [...new Set(this.value.split('\n').map(ip => ip.trim())
          .filter(ip => ip.length))]
        const invalidList = []
        const validList = []
        list.forEach((text) => {
          if (isIP(text)) {
            validList.push(text)
          } else {
            invalidList.push(text)
          }
        })
        this.validList = validList
        this.invalidList = invalidList
      },
      getReqs() {
        const IPList = []
        const IPv6List = []
        this.validList.forEach((ip) => {
          if (isIP(ip, 4)) {
            IPList.push(ip)
          }
          if (isIP(ip, 6)) {
            IPv6List.push(ip)
          }
        })

        const defaultParams = () => ({
          bk_biz_id: this.bizId,
          condition: this.getDefaultSearchCondition()
        })

        const reqs = []

        if (IPList.length) {
          const params = defaultParams()
          params.ip = { data: IPList, exact: 1, flag: 'bk_host_innerip' }
          reqs.push(this.$store.dispatch('hostSearch/searchHost', {
            params,
            config: {
              requestId: this.request.host
            }
          }))
        }

        if (IPv6List.length) {
          const params = defaultParams()
          const hostCond = params.condition.find(item => item.bk_obj_id === 'host')
          hostCond.condition = [{
            field: 'bk_host_innerip_v6',
            operator: '$in',
            value: IPv6List
          }]
          reqs.push(this.$store.dispatch('hostSearch/searchHost', {
            params,
            config: {
              requestId: this.request.host
            }
          }))
        }


        return reqs
      },
      handleHostSelectChange(data) {
        this.$emit('select-change', data)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .custom-layout {
        position: relative;
        display: flex;
        height: 100%;
        padding-top: 24px;

        .input-wrapper {
            position: relative;
            width: 280px;
        }
        .table-wrapper {
            flex: auto;
            margin-left: 20px;
        }
    }
    .ip-input {
        display: block;
        width: 100%;
        height: calc(100% - 60px);
        padding: 5px 10px;
        font-size: 12px;
        line-height: 20px;
        background-color: #FFF;
        border-radius:2px;
        border:1px solid #C4C6CC;
        cursor: text;
        outline: 0;
        resize: none;
        @include scrollbar;
        &.has-error {
            color: $dangerColor;
            text-decoration: underline;
        }
    }
    .ip-error {
        position: absolute;
        bottom: 45px;
        left: 0;
        line-height: 16px;
        font-size: 12px;
        color: $dangerColor;
    }
    .ip-confirm {
        display: block;
        width: 100%;
        margin: 15px 0;
    }
</style>
