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
  <cmdb-sticky-layout class="export">
    <bk-steps
      class="export-steps"
      :steps="steps"
      :cur-step="currentStep"
      v-show="showSteps"></bk-steps>
    <keep-alive>
      <component :is="stepComponent"></component>
    </keep-alive>
    <div :class="['options', { 'is-sticky': sticky }]" slot="footer" slot-scope="{ sticky }">
      <template v-if="currentStep < stepsLen">
        <bk-button class="mr10" theme="primary"
          :disabled="nextStepDisabled"
          @click="nextStep">
          {{$t('下一步')}}
        </bk-button>
      </template>
      <template v-if="currentStep > 1 && currentStep <= stepsLen">
        <bk-button class="mr10" theme="default" @click="previousStep">{{$t('上一步')}}</bk-button>
      </template>
      <template v-if="currentStep === stepsLen">
        <bk-button
          class="mr10"
          theme="primary"
          :disabled="exportDisabled"
          @click="startTask">{{$t('开始导出')}}</bk-button>
      </template>
      <bk-button theme="default" @click="close" v-if="currentStep <= stepsLen">{{$t('取消')}}</bk-button>
    </div>
  </cmdb-sticky-layout>
</template>

<script>
  import exportProperty from './export-property'
  import exportRelation from './export-relation'
  import exportStatus from './export-status'
  import useState from './state'
  import useTask from './task'
  import { computed } from 'vue'
  export default {
    name: 'export-file',
    components: {
      [exportProperty.name]: exportProperty,
      [exportRelation.name]: exportRelation,
      [exportStatus.name]: exportStatus
    },
    setup() {
      const [{
        step: currentStep,
        fields,
        presetFields,
        exportRelation: allowExportRelation,
        relations,
        steps
      }, { setState }] = useState()

      const nextStep = () => setState({ step: currentStep.value + 1 })
      const previousStep = () => setState({ step: currentStep.value - 1 })
      const close = () => setState({ visible: false })

      const stepComponent = computed(() => {
        const component = [exportProperty.name, exportRelation.name]
        const status = exportStatus.name
        const map = component.slice(0, stepsLen.value)
        map.push(status)
        return map[currentStep.value - 1]
      })
      const nextStepDisabled = computed(() => fields.value.length <= presetFields.value.length)
      const exportDisabled = computed(() => {
        if (!allowExportRelation.value) {
          return false
        }
        return Object.keys(relations.value).length === 0
      })
      const stepsLen = computed(() => steps.value.length)
      const showSteps = computed(() => currentStep.value < stepsLen.value + 1 && stepsLen.value > 1)

      const [, { start }] = useTask()
      const startTask = () => {
        nextStep()
        start()
      }

      return {
        nextStepDisabled,
        exportDisabled,
        currentStep,
        nextStep,
        previousStep,
        stepComponent,
        startTask,
        close,
        steps,
        showSteps,
        stepsLen
      }
    }
  }
</script>

<style lang="scss" scoped>
  .export {
    height: 100%;
    padding: 20px 24px 0;
    @include scrollbar-y;
    .export-steps {
      width: 350px;
      margin: 0 auto;
    }
    .options {
      height: 50px;
      width: calc(100% + 56px);
      margin: 20px 0 0 -28px;
      padding: 0 28px;
      display: flex;
      align-items: center;
      background-color: #fff;
      &.is-sticky {
        border-top: 1px solid $borderColor;
      }
    }
  }
</style>
