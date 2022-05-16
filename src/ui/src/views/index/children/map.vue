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
  <div class="map">
    <img src="../../../assets/images/map.svg" :style="imgStyles">
  </div>
</template>

<script>
  import {
    addResizeListener,
    removeResizeListener
  } from '@/utils/resize-events.js'
  export default {
    name: 'cmdb-index-map',
    data() {
      return {
        ratio: {
          height: 404 / 857,
          top: 181 / 767
        },
        imgStyles: {
          width: '857px',
          height: '404px',
          visibility: 'hidden'
        },
        resizeHandler: null
      }
    },
    mounted() {
      this.initResizeEvent()
    },
    beforeDestroy() {
      removeResizeListener(this.$parent.$el, this.resizeHandler)
    },
    methods: {
      initResizeEvent() {
        this.resizeHandler = () => {
          const parentRect = this.$parent.$el.getBoundingClientRect()
          const imgStyles = {
            width: `${Math.floor(parentRect.width * 0.66)}px`,
            height: `${Math.floor(parentRect.width * 0.66 * this.ratio.height)}px`,
            left: `${Math.floor((parentRect.width * 0.17) + parentRect.left)}px`
          }
          this.imgStyles = imgStyles
        }
        addResizeListener(this.$parent.$el, this.resizeHandler)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .map {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        z-index: -1;
        pointer-events: none;
        overflow: visible;
        svg, img {
            position: absolute;
            top: 181px;
        }
        img {
            opacity: 0.5745;
        }
    }
</style>
