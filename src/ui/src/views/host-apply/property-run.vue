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
  <div class="run-wrapper">
    <top-steps :current="3" v-if="hasSteps" />

    <div class="status-content">
      <div class="status-loading" v-if="$loading(Object.values(requestIds))">
        <div class="status-inner">
          <img class="status-icon" src="@/assets/images/icon/loading.svg">
          <p class="status-title" v-show="$loading(requestIds.run)">{{$t('正在应用中')}}</p>
          <p class="status-title" v-show="$loading(requestIds.task)">{{$t('正在查询应用状态')}}</p>
        </div>
      </div>

      <div class="status-result" v-else-if="requestErrors.run">
        <div class="result-icon">
          <span class="status-icon bk-icon icon-close icon-error"></span>
          <p class="result-title">
            <span>{{$t('应用异常')}}</span>
          </p>
          <p class="result-subtitle">
            <span>{{requestErrors.msg.run}}</span>
          </p>
          <div class="result-options">
            <bk-button @click="handleRetry">{{$t('重试')}}</bk-button>
          </div>
        </div>
      </div>

      <div class="status-result" v-else-if="taskStatus">
        <div class="result-icon">
          <span class="status-icon bk-icon icon-check-1 icon-success" v-show="taskSuccess"></span>
          <img class="status-icon" src="@/assets/images/icon/loading.svg" v-show="taskUndone">
          <span class="status-icon bk-icon icon-close icon-fail" v-show="taskFail"></span>
        </div>
        <p class="result-title">
          <span v-show="taskSuccess">{{$t('应用成功')}}</span>
          <span v-show="taskUndone">{{$t('正在应用中')}}</span>
          <span v-show="taskFail">{{$t('应用失败')}}</span>
        </p>
        <p class="result-subtitle">
          <span v-show="taskSuccess">{{$t('成功保存策略并应用到当前模块下主机')}}</span>
          <span v-show="taskUndone">{{$t('任务正在后台执行，请耐心等待')}}</span>
          <span v-show="taskFail">{{$t('本次任务运行失败，可点击下方按钮重试')}}</span>
        </p>
        <div class="result-options">
          <bk-button class="mr10" theme="primary" v-show="!taskUndone" @click="handleBack">{{$t('返回列表')}}</bk-button>
          <bk-button @click="handleRetry" v-show="taskFail">{{$t('重试')}}</bk-button>
          <bk-button @click="handleViewHost" v-show="taskSuccess">{{$t('查看主机')}}</bk-button>
        </div>
      </div>

      <bk-exception v-else class="exception-wrap-item exception-part" type="empty" scene="part">
        <span>{{$t('无应用任务')}}</span>
      </bk-exception>
    </div>

    <leave-confirm
      v-bind="leaveConfirmConfig"
      reverse
      :title="$t('是否退出应用')"
      :content="$t('应用未完成，退出后任务将在后台执行')"
      :ok-text="$t('退出')"
      :cancel-text="$t('取消')">
    </leave-confirm>
  </div>
</template>

<script>
  import { mapGetters, mapState, mapActions } from 'vuex'
  import topSteps from './children/top-steps.vue'
  import leaveConfirm from '@/components/ui/dialog/leave-confirm'
  import {
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_HOST_APPLY
  } from '@/dictionary/menu-symbol'
  import { TASK_STATUS, setTask, getTask, removeTask  } from './task-helper.js'
  import { CONFIG_MODE } from '@/services/service-template/index.js'

  export default {
    components: {
      topSteps,
      leaveConfirm,
    },
    data() {
      return {
        requestIds: {
          run: Symbol(),
          task: Symbol(),
        },
        requestErrors: {
          run: false,
          msg: {
            run: ''
          }
        },
        taskStatus: null,
        undoneStatus: [TASK_STATUS.NEW, TASK_STATUS.WAITING, TASK_STATUS.EXECUTING],
        pollTimer: null,
        runningTaskId: null,
        leaveConfirmConfig: {
          id: 'propertyRun',
          active: true
        }
      }
    },
    computed: {
      ...mapGetters(['userName']),
      ...mapGetters('objectBiz', ['bizId']),
      ...mapState('hostApply', ['propertyConfig']),
      mode() {
        return this.$route.params.mode
      },
      hasConfig() {
        return Object.keys(this.propertyConfig).length > 0
      },
      taskSuccess() {
        return this.taskStatus === TASK_STATUS.FINISHED
      },
      taskFail() {
        return this.taskStatus === TASK_STATUS.FAIL
      },
      taskUndone() {
        return this.undoneStatus.includes(this.taskStatus)
      },
      hasSteps() {
        return !['conflict-list'].includes(this.$route.query.from)
      },
      targetIdsKey() {
        const targetIdsKeys = {
          [CONFIG_MODE.MODULE]: 'bk_module_ids',
          [CONFIG_MODE.TEMPLATE]: 'service_template_ids'
        }
        return targetIdsKeys[this.mode]
      },
      requestConfigs() {
        return {
          [this.requestIds.run]: {
            [CONFIG_MODE.MODULE]: {
              action: 'runApply'
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'runTemplateApply'
            }
          },
          [this.requestIds.task]: {
            [CONFIG_MODE.MODULE]: {
              action: 'getApplyTaskStatus'
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'getTemplateApplyTaskStatus'
            }
          }
        }
      }
    },
    async created() {
      this.run()

      this.leaveConfirmConfig.active = false

      this.setBreadcrumbs()
    },
    destroyed() {
      this.clearTimer()
    },
    methods: {
      ...mapActions('hostApply', [
        'runApply',
        'runTemplateApply',
        'getApplyTaskStatus',
        'getTemplateApplyTaskStatus'
      ]),
      async run() {
        // 本地任务ID
        const localTaskId = getTask(`${this.userName}${this.bizId}`)

        // 存在配置数据，先执行应用
        if (this.hasConfig) {
          const [result, msg] = await this.saveAndApply()
          if (result === false) {
            this.requestErrors.run = true
            this.requestErrors.msg.run = msg
            return
          }

          // 存储新任务ID
          const { bk_biz_id: bizId, task_id: newTaskId } = result
          setTask(newTaskId, `${this.userName}${bizId}`)

          this.pollTaskStatus(newTaskId)
        } else if (localTaskId) {
          this.pollTaskStatus(localTaskId)
        }
      },
      setBreadcrumbs() {
        this.$store.commit('setTitle', this.$t('应用属性'))
      },
      async saveAndApply() {
        try {
          const requestConfig = this.requestConfigs[this.requestIds.run][this.mode]
          const res = await this[requestConfig.action]({
            params: {
              bk_biz_id: this.bizId,
              ...this.propertyConfig
            },
            config: {
              requestId: this.requestIds.run,
              globalError: false,
              transformData: false
            }
          })
          return Promise.resolve(res?.result ? [res.data] : [false, res.bk_error_msg])
        } catch (error) {
          console.error(error)
          return Promise.resolve(false)
        }
      },
      async pollTaskStatus(taskId) {
        // 记录当前运行中的任务，可用于重试
        this.runningTaskId = taskId

        const result = await this.getTaskStatus(taskId)

        // 请求失败或者任务处于未完成状态则开启轮询
        if (result === false || this.undoneStatus.includes(result?.status)) {
          this.clearTimer()
          this.pollTimer = setTimeout(() => this.pollTaskStatus(taskId), 5000)
        }

        const taskStatus = result?.status
        const undone = this.undoneStatus.includes(taskStatus)

        // 未完成退出时提示否则不提示
        this.leaveConfirmConfig.active = undone

        // 完成状态（失败或成功），移除本地任务ID
        if (!undone) {
          removeTask(`${this.userName}${this.bizId}`)
        }

        this.taskStatus = taskStatus
      },
      async getTaskStatus(taskId) {
        try {
          const requestConfig = this.requestConfigs[this.requestIds.task][this.mode]
          const result = await this[requestConfig.action]({
            params: {
              bk_biz_id: this.bizId,
              task_ids: [taskId]
            },
            config: {
              requestId: this.requestIds.task,
              globalError: false
            }
          })
          return Promise.resolve(result?.task_info?.find(({ task_id: id }) => id === taskId))
        } catch (error) {
          console.error(error)
          return Promise.resolve(false)
        }
      },
      handleBack() {
        const query = {}
        if (this.propertyConfig[this.targetIdsKey]?.length === 1) {
          // eslint-disable-next-line prefer-destructuring
          query.id = this.propertyConfig[this.targetIdsKey][0]
        }
        this.$store.commit('hostApply/clearRuleDraft')

        this.leaveConfirmConfig.active = false
        this.$nextTick(() => {
          this.$routerActions.redirect({
            name: MENU_BUSINESS_HOST_APPLY,
            params: {
              mode: this.mode
            },
            query
          })
        })
      },
      handleViewHost() {
        const query = {}
        if (this.propertyConfig.bk_module_ids?.length === 1) {
          query.node = `module-${this.propertyConfig.bk_module_ids[0]}`
        }

        this.leaveConfirmConfig.active = false
        this.$nextTick(() => {
          this.$routerActions.redirect({
            name: MENU_BUSINESS_HOST_AND_SERVICE,
            query,
            history: true
          })
        })
      },
      handleRetry() {
        this.run()
      },
      clearTimer() {
        this.pollTimer && clearTimeout(this.pollTimer)
      }
    }
  }
</script>

<style lang="scss" scoped>
  .status-content {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 350px;
    margin: 16px 24px;
    background: #fff;
    box-shadow: 0px 2px 4px 0px rgba(25, 25, 41, 0.05);
    text-align: center;

    .status-icon {
      display: block;
      width: 54px;
      height: 54px;
      margin: 0 auto;
      &:not(img) {
        line-height: 54px;
        border-radius: 50%;
        color: #FFF;
        font-size: 30px;
      }
      &.icon-check-1 {
        font-size: 50px;
      }
      &.icon-error {
        background-color: $dangerColor;
      }
      &.icon-success {
        background-color: #2dcb56;
      }
      &.icon-fail {
        background-color: #ff5656;
      }
      &.icon-abnormal {
        background-color: #ffb848;
      }
    }
  }

  .status-loading {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 40px 0;

    .status-title {
      margin: 30px 0 0;
      font-size: 24px;
      color: #313238;
    }
  }

  .status-result {
    .result-title {
      margin: 20px auto 0;
      line-height: 30px;
      font-size: 24px;
      color: #313238;
    }
    .result-subtitle {
      font-size: 14px;
      color: $textColor;
      margin: 12px auto 0;
      & + .result-stat {
        margin-top: 4px;
      }
    }
    .result-stat {
      margin: 28px auto 0;
      font-size: 14px;
      color: $textColor;
      .result-count {
        padding: 0 2px;
        font-weight: bold;
        &.fail {
          color: $dangerColor;
        }
      }
    }
    .result-options {
      font-size: 0;
      margin-top: 24px;

      .bk-button {
        min-width: 90px;
      }
    }
  }
</style>
