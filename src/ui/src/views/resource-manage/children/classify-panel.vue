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
  <div class="classify">
    <h4 class="classify-name" :title="classify['bk_classification_name']">
      <span class="classify-name-text">{{classify['bk_classification_name']}}</span>
    </h4>
    <div class="models-layout">
      <div class="models-link" v-for="(model, index) in models"
        :key="index"
        :title="model['bk_obj_name']"
        @click="redirect(model)">
        <i :class="['model-icon','icon', model['bk_obj_icon'], { 'nonpre-mode': !model['ispre'] }]"></i>
        <span class="model-name">{{model['bk_obj_name']}}</span>
        <i class="model-star bk-icon"
          :class="[isCollected(model) ? 'icon-star-shape' : 'icon-star']"
          @click.prevent.stop="toggleCustomNavigation(model)">
        </i>
        <div class="model-instance-count">
          <instance-count :obj-id="model.bk_obj_id" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import has from 'has'
  import { mapGetters } from 'vuex'
  import {
    MENU_RESOURCE_INSTANCE,
    MENU_RESOURCE_COLLECTION
  } from '@/dictionary/menu-symbol'
  import InstanceCount from './instance-count.vue'
  import { BUILTIN_MODELS, BUILTIN_MODEL_COLLECTION_KEYS, BUILTIN_MODEL_RESOURCE_MENUS } from '@/dictionary/model-constants.js'

  export default {
    components: {
      InstanceCount
    },
    props: {
      classify: {
        type: Object,
        required: true
      },
      collection: {
        type: Array,
        required: true,
        default: () => ([])
      }
    },
    data() {
      return {
        maxCustomNavigationCount: 8
      }
    },
    computed: {
      ...mapGetters('userCustom', ['usercustom']),
      collectedCount() {
        return this.collection.length
      },
      models() {
        return this.classify.bk_objects
      }
    },
    methods: {
      redirect(model) {
        if (has(BUILTIN_MODEL_RESOURCE_MENUS, model.bk_obj_id)) {
          this.$routerActions.redirect({
            name: BUILTIN_MODEL_RESOURCE_MENUS[model.bk_obj_id]
          })
        } else {
          this.$routerActions.redirect({
            name: MENU_RESOURCE_INSTANCE,
            params: {
              objId: model.bk_obj_id
            }
          })
        }
      },
      isCollected(model) {
        return this.collection.includes(model.bk_obj_id)
      },
      isBuiltinModel(model) {
        return Object.values(BUILTIN_MODELS).includes(model.bk_obj_id)
      },
      toggleCustomNavigation(model) {
        if (this.isBuiltinModel(model)) {
          this.toggleDefaultCollection(model)
        } else {
          let isAdd = false
          let newCollection
          const oldCollection = this.usercustom[MENU_RESOURCE_COLLECTION] || []
          if (oldCollection.includes(model.bk_obj_id)) {
            newCollection = oldCollection.filter(id => id !== model.bk_obj_id)
          } else {
            isAdd = true
            newCollection = [...oldCollection, model.bk_obj_id]
          }
          const promise = this.$store.dispatch('userCustom/saveUsercustom', {
            [MENU_RESOURCE_COLLECTION]: newCollection
          })
          promise.then(() => {
            this.$success(isAdd ? this.$t('添加导航成功') : this.$t('取消导航成功'))
          })
        }
      },
      async toggleDefaultCollection(model) {
        const isCollected = this.isCollected(model)
        try {
          const key =  BUILTIN_MODEL_COLLECTION_KEYS[model.bk_obj_id]
          await this.$store.dispatch('userCustom/saveUsercustom', {
            [key]: !isCollected
          })
          this.$success(isCollected ? this.$t('取消导航成功') : this.$t('添加导航成功'))
        } catch (e) {
          console.error(e)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .classify{
        margin: 0 0 20px 0;
        background-color: #fff;
        border: 1px solid #ebf0f5;
        box-shadow:0px 3px 6px 0px rgba(51,60,72,0.05);
    }
    .classify-name{
        padding: 13px 5px;
        margin: 0 20px;
        line-height: 20px;
        font-size: 0;
        color: $cmdbTextColor;
        border-bottom: 1px solid #ebf0f5;
        &-text {
            display: inline-block;
            padding: 0 2px 0 0;
            vertical-align: middle;
            max-width: calc(100% - 40px);
            font-size: 14px;
            @include ellipsis;
        }
        &-count {
            display: inline-block;
            width: 40px;
            vertical-align: middle;
            font-size: 14px;
        }
    }
    .models-layout{
        padding: 8px 0;
        .models-link{
            display: block;
            height: 38px;
            font-size: 0;
            position: relative;
            padding: 7px 25px;
            cursor: pointer;
            &:hover{
                background-color: #ecf3ff;
            }
            &:before{
                content: "";
                display: inline-block;
                height: 100%;
                vertical-align: middle;
            }
            &:hover .model-icon,
            &:hover .model-name{
                color: #3A84FF;
            }
            &:hover .model-star{
                display: inline-block;
            }
            .model-icon,
            .model-name{
                display: inline-block;
                vertical-align: middle;
            }
            .model-icon{
                font-size: 16px;
                color: #798AAD;
            }
            .nonpre-mode {
                color: #3A84FF !important;
            }
            .model-name{
                max-width: calc(100% - 100px);
                margin: 0 0 0 12px;
                font-size: 14px;
                line-height: 24px;
                color: $cmdbTextColor;
                @include ellipsis;
            }
            .model-instance-count {
                float: right;
                @include inlineBlock;
                width: 35px;
                font-size: 14px;
                height: 24px;
                line-height: 24px;
                color: #C4C6CC;
                text-align: right;
                display: flex;
                align-items: center;
                justify-content: flex-end;
            }
            .model-star{
                display: none;
                width: 24px;
                height: 24px;
                margin-left: 5px;
                line-height: 24px;
                text-align: center;
                font-size: 14px;
                cursor: pointer;
                vertical-align: middle;
                &.icon-star-shape{
                    color: #FFB400;
                    display: inline-block;
                }
            }
        }
    }
</style>
