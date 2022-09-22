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

export default {
  data() {
    return {
      active: false
    }
  },
  watch: {
    active(active) {
      this.$emit('active-change', active)
      this.hackEnterEvent()
    }
  },
  methods: {
    handleToggle(active) {
      this.active = active
    },
    hackEnterEvent() {
      if (this.active) {
        window.addEventListener('keyup', this.handleEnter, true)
        this.$el.style.position = 'relative'
        // eslint-disable-next-line no-underscore-dangle
        this.$el.style.zIndex = window.__bk_zIndex_manager.nextZIndex()
      } else {
        this.$el.style.position = ''
        this.$el.style.zIndex = ''
        window.removeEventListener('keyup', this.handleEnter, true)
      }
    },
    handleEnter(event) {
      if (event.key.toLowerCase() !== 'enter') return
      this.$emit('enter', event)
    }
  }
}
