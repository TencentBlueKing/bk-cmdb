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
  <div class="footer">
    <p class="contact" v-html="contact"></p>
    <p class="copyright">
      {{copyright}} {{verison}}
    </p>
  </div>
</template>

<script>
  import moment from 'moment'
  import xss from 'xss'

  export default {
    name: 'TheFooter',
    props: {
      previewContact: {
        type: String,
        default: ''
      },
      previewCopyright: {
        type: String,
        default: ''
      },
    },
    computed: {
      setting() {
        return this.$store.state.globalConfig.config
      },
      contact() {
        if (this.previewContact) return this.parseMarkdownLink(this.previewContact)
        return this.parseMarkdownLink(this.setting.footer.contact)
      },
      copyright() {
        if (this.previewCopyright) return this.parseTimeVars(this.previewCopyright)
        return this.parseTimeVars(this.setting.footer.copyright)
      },
      verison() {
        return this.$Site.buildVersion
      }
    },
    methods: {
      parseMarkdownLink(markdown) {
        return xss(markdown, { stripIgnoreTagBody: true })
          ?.replace(/\[([^\]]+)\]\(([^)]+)\)/ig, '<a target="_blank" class="contact-link" href="$2">$1</a>')
      },
      /**
       * 转换时间变量
       * @param {string} content 用户输入的 copyright 信息
       */
      parseTimeVars(content) {
        const currentYear = moment().format('YYYY')
        const currentMonth = moment().format('MM')
        const currentDay = moment().format('DD')

        return content
          .replace(/\{\{current_year\}\}/ig, currentYear)
          .replace(/\{\{current_month\}\}/ig, currentMonth)
          .replace(/\{\{current_day\}\}/ig, currentDay)
      },
    },
  }
</script>

<style lang="scss" scoped>
.footer {
    position: absolute;
    left: 25px;
    right: 25px;
    bottom: 0;
    padding-top: 8px;
    height: 52px;
    font-size: 12px;
    text-align: center;
    color: $textColor;
    border-top: 1px solid #DCDEE5;
    background-color: #F5F6FA;
    z-index: 2;
}

.copyright {
  line-height: 24px;
}

::v-deep .contact-link {
  color: $primaryColor;
}
</style>
