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
  <div class="step-pane">
    <div class="step-bar">
      <bk-steps :steps="steps" :cur-step="currentStepIndex"></bk-steps>
    </div>
    <div class="step-pane-body">
      <div class="step-container">
        <slot v-bind="{
          currentStep: currentStep,
          currentStepIndex: currentStepIndex
        }"></slot>
      </div>
    </div>
    <div class="action-bar">
      <bk-button @click="cancel" v-show="currentStepIndex < steps.length" class="cancel-button">{{t("取消")}}</bk-button>
      <bk-button v-show="currentStepIndex > 1 && currentStepIndex < steps.length" class="prev-step-button"
        @click="toPrevStep"
      >{{currentStep.prevButtonText || t("上一步")}}</bk-button
      >
      <span v-bk-tooltips="{
        content: currentStep.nextButtonDisabledTooltips,
        disabled: !currentStep.nextButtonDisabled
      }">
        <bk-button
          v-show="currentStepIndex <= steps.length && currentStep.nextButtonVisible"
          theme="primary"
          @click="toNextStep"
          :loading="isLoading"
          :disabled="currentStep.nextButtonDisabled"

          class="next-step-button"
        >{{ currentStep.nextButtonText || t("下一步") }}</bk-button
        >
      </span>
    </div>
  </div>
</template>

<script>
  import { defineComponent, ref, computed, toRef } from '@vue/composition-api'
  import { t } from '@/i18n'
  import has from 'has'
  import cloneDeep from 'lodash/cloneDeep'

  export default defineComponent({
    name: 'StepPane',
    props: {
      /**
       * 初始化时的默认当前步骤索引
       */
      defaultStepIndex: {
        type: Number,
        default: 1
      },
      /**
       * @property {Array.<Object>} steps 步骤条配置
       * @property {String} steps[].title 步骤标题，继承自 bk-step
       * @property {Number} steps[].icon 步骤标识，继承自 bk-step
       * @property {String} steps[].nextButtonText 下一步按钮文案，默认为「下一步」
       * @property {Boolean} steps[].nextButtonVisible 是否展示下一步按钮，默认 true
       * @property {Boolean|Function} steps[].nextButtonDisabled 是否置灰下一步按钮，默认 false
       * @property {Function} steps[].nextHandler 下一步按钮点击回调
       */
      steps: {
        type: Array,
        required: true
      },
    },
    setup(props, { emit }) {
      const { defaultStepIndex } = props
      const currentStepIndex = ref(defaultStepIndex)
      const theSteps = toRef(props, 'steps')
      const currentStep = computed(() => {
        const step = cloneDeep(theSteps.value[currentStepIndex.value - 1])

        if (has(step, 'nextButtonVisible')) {
          step.nextButtonVisible = Boolean(step.nextButtonVisible)
        } else {
          step.nextButtonVisible = true
        }

        if (has(step, 'nextButtonDisabled')) {
          if (typeof step.nextButtonDisabled === 'function') {
            step.nextButtonDisabled = step.nextButtonDisabled()
          } else {
            step.nextButtonDisabled = step.nextButtonDisabled
          }
        } else {
          step.nextButtonDisabled = false
        }

        return step
      })

      const isLoading = ref(false)

      /**
       * 进入下一步
       * 支持回调 done 和 stop，如果不传入回调，则默认进入下一步
       * done 表示完成并进入下一步
       * stop 表示停止 loading，并且不会进入下一步
       */
      const toNextStep = () => {
        if (currentStepIndex.value >= theSteps.value.length) return

        const stop = () => {
          isLoading.value = false
        }

        const done = () => {
          stop()
          currentStepIndex.value += 1
        }

        isLoading.value = true

        if (currentStep.value.nextHandler) {
          currentStep.value.nextHandler(done, stop)
        } else {
          done()
        }
      }

      const toPrevStep = () => {
        if (currentStepIndex.value <= 1) return
        currentStepIndex.value -= 1
      }

      const cancel = () => {
        emit('cancel')
      }

      return {
        t,
        currentStep,
        currentStepIndex,
        toNextStep,
        toPrevStep,
        cancel,
        isLoading
      }
    },
  })
</script>

<style lang="scss" scoped>
$stepBarHeight: 50px;
$actionBarHeight: 50px;

.step-pane {
  position: absolute;
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
  z-index: 999;
  background-color: #fff;
}

.step-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  height: $stepBarHeight;
  border-top: 1px solid #dcdee5;
  border-bottom: 1px solid #dcdee5;

  .bk-steps {
    margin: 0 195px;
  }
}

.step-pane-body {
  height: calc(100% - $actionBarHeight - $stepBarHeight);
}

.step-container {
  height: 100%;
}

.action-bar {
  position: absolute;
  bottom: 0;
  right: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: right;
  height: $actionBarHeight;
  box-sizing: border-box;
  background: #fff;
  border-top: 1px solid #e2e2e2;
  font-size: 14px;

  .next-step-button,
  .prev-step-button {
    width: 120px;
  }

  .bk-button {
    margin-right: 10px;

    &:last-child {
      margin-right: 24px;
    }
  }

  .cancel-button {
    width: 86px;
  }
}
</style>
