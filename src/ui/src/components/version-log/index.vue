<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <bk-version-detail
    :class="{ 'version-detail': versionList.length === 0 }"
    :current-version="currentVersion"
    :finished="true"
    :show.sync="isShowChangeLogs"
    :version-list="versionList"
    :version-detail="versionDetail"
    :get-version-detail="handleGetVersionDetail">
    <template slot-scope="content">
      <div v-if="content.detail">
        <h2 style="margin-top:5px">{{content.detail}} {{$t('版本日志')}} </h2>
        <div class="markdowm-container" v-html="currentContent"></div>
      </div>
      <div class="exception-wrap" v-if="versionList.length === 0">
        <bk-exception class="exception-wrap-item" type="empty">
          <span style="font-size: 12px;">{{$t('暂无版本日志')}}</span>
        </bk-exception>
      </div>
    </template>
  </bk-version-detail>
</template>

  <script>
  import {  mapActions } from 'vuex'
  import { marked } from 'marked'
  import xss from 'xss'
  export default {
    name: 'detail',
    props: {
      currentVersion: {
        type: String,
        default: ''
      },
      show: {
        type: Boolean,
        default: false
      },
      versionList: {
        type: Array,
        default: () => []
      }
    },
    data() {
      return {
        versionDetail: '',
        currentContent: ''
      }
    },
    computed: {
      isShowChangeLogs: {
        get() {
          return this.show
        },
        set(newVal) {
          this.$emit('update:show', newVal)
        }
      }
    },
    methods: {
      ...mapActions('versionLog', [
        'getLogDetail'
      ]),
      async handleGetVersionDetail(version) {
        if (!version) return
        const params = { version: version.title }
        try {
          const logDetail = await this.getLogDetail(params)
          this.versionDetail =  version.title
          this.currentContent = xss(marked(logDetail), { stripIgnoreTagBody: true })
        } catch (e) {
          console.error(e)
        }
      }
    }
  }
  </script>

<style lang="scss" scoped>
.exception-wrap {
  display: flex;
  flex-wrap: wrap;
}
.exception-wrap .exception-wrap-item {
  margin: 10px;
  height: 420px;
  padding-top: 22px;
}
::v-deep .markdowm-container {
  margin-top: 14px;
  font-size: 14px;
  color: #63656e;
  h1,
  h2,
  h3,
  h4,
  h5 {
    height: auto;
    margin: 10px 0;
    font: normal 14px/1.5 "Helvetica Neue", Helvetica, Arial, "Lantinghei SC",
      "Hiragino Sans GB", "Microsoft Yahei", sans-serif;
    font-weight: bold;
    color: #34383e;
  }
  h1 {
    font-size: 30px;
  }
  h2 {
    font-size: 24px;
  }
  h3 {
    font-size: 18px;
  }
  h4 {
    font-size: 16px;
  }
  h5 {
    font-size: 14px;
  }
  em {
    font-style: italic;
  }
  div,
  p,
  font,
  span,
  li {
    line-height: 1.3;
  }
  p {
    margin: 0 0 1em;
  }
  table,
  table p {
    margin: 0;
  }
  ul,
  ol {
    padding: 0;
    margin: 0 0 1em 2em;
    text-indent: 0;
  }
  ul {
    padding: 0;
    margin: 10px 0 15px 15px;
    list-style-type: none;
  }
  ol {
    padding: 0;
    margin: 10px 0 10px 25px;
  }
  ol > li {
    line-height: 1.8;
    white-space: normal;
  }
  ul > li {
    padding-left: 15px !important;
    line-height: 1.8;
    white-space: normal;
    &::before {
      display: inline-block;
      width: 6px;
      height: 6px;
      margin-right: 9px;
      margin-left: -15px;
      background: #63656e;
      border-radius: 50%;
      content: "";
    }
  }
  li > ul {
    margin-bottom: 10px;
  }
  li ol {
    padding-left: 20px !important;
  }
  ul ul,
  ul ol,
  ol ol,
  ol ul {
    margin-left: 20px;
  }
  ul.list-type-1 > li {
    padding-left: 0 !important;
    margin-left: 15px !important;
    list-style: circle !important;
    background: none !important;
  }
  ul.list-type-2 > li {
    padding-left: 0 !important;
    margin-left: 15px !important;
    list-style: square !important;
    background: none !important;
  }
  ol.list-type-1 > li {
    list-style: lower-greek !important;
  }
  ol.list-type-2 > li {
    list-style: upper-roman !important;
  }
  ol.list-type-3 > li {
    list-style: cjk-ideographic !important;
  }
  pre,
  code {
    width: 95%;
    padding: 0 3px 2px;
    font-family: Monaco, Menlo, Consolas, "Courier New", monospace;
    font-size: 14px;
    color: #333;
    border-radius: 3px;
  }
  code {
    padding: 2px 4px;
    font-family: Consolas, monospace, tahoma, Arial;
    color: #d14;
    border: 1px solid #e1e1e8;
  }
  pre {
    display: block;
    padding: 9.5px;
    margin: 0 0 10px;
    font-family: Consolas, monospace, tahoma, Arial;
    font-size: 13px;
    word-break: break-all;
    word-wrap: break-word;
    white-space: pre-wrap;
    background-color: #f6f6f6;
    border: 1px solid #ddd;
    border: 1px solid rgb(0 0 0 / 15%);
    border-radius: 2px;
  }
  pre code {
    padding: 0;
    white-space: pre-wrap;
    border: 0;
  }
  blockquote {
    padding: 0 0 0 14px;
    margin: 0 0 20px;
    border-left: 5px solid #dfdfdf;
  }
  blockquote p {
    margin-bottom: 0;
    font-size: 14px;
    font-weight: 300;
    line-height: 25px;
  }
  blockquote small {
    display: block;
    line-height: 20px;
    color: #999;
  }
  blockquote small::before {
    content: "\2014 \00A0";
  }
  blockquote::before,
  blockquote::after {
    content: "";
  }
}
.version-detail {
  ::v-deep .bk-version-left {
    display: none;
  }
}
</style>

