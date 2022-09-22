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
  <section class="export-content">
    <div class="header">
      <div class="subtitle-wrapper">
        <i18n class="subtitle" tag="h2" path="分批下载副标题">
          <template #count><strong class="count">{{count}}</strong></template>
          <template #limit><span>{{limit}}</span></template>
        </i18n>
        <span class="process-counter">
          <span class="finished">{{finishedTask.length}}</span>
          <i>&nbsp;/&nbsp;&nbsp;</i>
          <span class="total">{{all.length}}</span>
        </span>
      </div>
    </div>
    <ul class="list" ref="list">
      <li class="list-item"
        v-for="(task, index) in all"
        ref="listItem"
        :key="index">
        <span :class="['state', task.state]">
          <i :class="['state-icon', iconMapping[task.state]]" v-if="task.state !== 'waiting'"></i>
          {{textMapdding[task.state]}}
        </span>
        <span class="info">
          <span class="info-name">{{`${task.name}.xlsx`}}</span>
          <span class="info-error"
            v-if="task === current && current.state === 'error'">
            {{message}}
          </span>
        </span>
      </li>
    </ul>
  </section>
</template>

<script>
  import useState from './state'
  import useTask from './task'
  import { computed } from '@vue/composition-api'
  export default {
    setup() {
      const [{ count, limit }] = useState()
      const [taskState] = useTask()
      const finishedTask = computed(() => taskState.all.value.filter(task => task.state === 'finished'))
      return {
        count,
        limit,
        finishedTask,
        ...taskState
      }
    }
  }
</script>

<style lang="scss" scoped>
    .export-content {
      .header {
        .subtitle-wrapper {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-top: 32px;
          line-height: 20px;
          .subtitle {
              flex: 1;
              font-size: 14px;
              font-weight: normal;
              color: $textColor;
              .count {
                  font-weight: bold;
              }
          }
          .process-counter {
              font-size: 14px;
              color: $textColor;
              font-weight: bold;
              padding-right: 8px;
              display: inline-flex;
              .finished {
                  color: $successColor;
              }
          }
        }
      }
    }
    .list {
        position: relative;
        margin-top: 12px;
        .list-item {
          display: flex;
          align-items: center;
          margin-bottom: 10px;
          height: 52px;
          border: 1px solid #dcdee5;
          border-radius: 2px;
          box-shadow: 0px 2px 4px 0px rgba(0,0,0,0.1);
        }
    }
    .state {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 60px;
        height: 100%;
        font-size: 12px;
        flex-direction: column;
        &.pending {
            background-color: #e1ecff;
        }
        &.finished {
            background-color: #e4faf0;
            .state-icon {
                color: $successColor;
            }
        }
        &.waiting {
            background-color: #f0f1f5;
        }
        &.error {
            background-color: #fedddc;
            .state-icon {
                color: #ea3636;
            }
        }
        .state-icon {
            font-size: 20px;
            &.loading {
                display: inline-block;
                width: 16px;
                height: 16px;
                background-color: transparent;
                background-image: url("../../assets/images/icon/loading.svg");
                background-position: center center;
                background-size: 16px;
                background-repeat: no-repeat;
            }
        }
    }
    .info {
        flex: 1;
        display: flex;
        justify-content: flex-start;
        flex-direction: column;
        padding: 0 17px;
        .info-name {
            font-size: 14px;
            font-weight: bold;
            color: $textColor;
        }
        .info-error {
            font-size: 12px;
            color: $dangerColor;
        }
    }
</style>
