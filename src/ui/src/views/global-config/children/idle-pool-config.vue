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
  <div class="idle-machine-pool-config" v-bkloading="{ isLoading: globalConfig.loading }">
    <bk-form ref="idleFormRef" :model="idleForm" :rules="idleFormRules">
      <!-- 内置空闲机集群 -->
      <bk-form-item :label="$t('集群')" property="set" required :icon-offset="iconOffset">
        <ModuleBuilder
          :removeable="false"
          module-id="set"
          module-id-disabled
          @before-confirm="beforeClusterConfirm"
          :module-name-placeholder="globalConfig.config.validationRules
            .businessTopoInstNames.message"
          @cancel="clearError"
          :module-name.sync="idleForm.set"></ModuleBuilder>
      </bk-form-item>

      <!-- 内置空闲机模块 -->
      <bk-form-item
        v-for="(idleModuleName, idleModuleKey, idleModuleIndex) in idleForm.buildInModules"
        :key="idleModuleKey"
        required
        :label="idleModuleIndex === 0 ? $t('模块') : ''"
        :property="idleModuleKey"
        :icon-offset="iconOffset">
        <ModuleBuilder
          indent-line
          module-id-disabled
          :module-name-placeholder="globalConfig.config.validationRules.businessTopoInstNames.message"
          @before-confirm="beforeModuleConfirm({
            moduleKey: idleModuleKey,
            moduleName: idleModuleName
          }, $event)"
          @cancel="clearError"
          :node-indent="nodeIndent"
          :removeable="false"
          :module-id="idleModuleKey"
          :module-name.sync="idleForm.buildInModules[idleModuleKey]">
        </ModuleBuilder>
      </bk-form-item>

      <!-- 新增模块 -->
      <bk-form-item
        :icon-offset="iconOffset"
        v-for="(userModule, userModuleIndex) in idleForm.userModules"
        :key="userModule.ruleKey"
        :property="userModule.moduleKey"
        :rules="userModule.rules">
        <ModuleBuilder
          indent-line
          :module-id.sync="userModule.moduleKey"
          :module-id-disabled="!userModule.isNew"
          :module-name.sync="userModule.moduleName"
          :node-indent="nodeIndent"
          :state.sync="userModule.state"
          :module-id-placeholder="$t('模块 ID，英文/数字')"
          :module-name-placeholder="globalConfig.config.validationRules.businessTopoInstNames.message"
          @before-confirm="beforeUserModuleConfirm(userModule, $event)"
          @cancel="handleUserModuleCancel(userModuleIndex, $event)"
          @remove="removeUserModule(userModule.moduleKey, userModule.moduleName)">
        </ModuleBuilder>
      </bk-form-item>

      <bk-form-item ext-cls="action-form-item">
        <bk-button text @click="addUserModule"
          :disabled="addBtnDisabled">{{$t('添加模块')}}</bk-button>
      </bk-form-item>
    </bk-form>
  </div>
</template>

<script>
  import { ref, computed, reactive, defineComponent, onMounted } from 'vue'
  import ModuleBuilder from './module-builder.vue'
  import store from '@/store'
  import { bkInfoBox, bkMessage } from 'bk-magic-vue'
  import { t } from '@/i18n'
  import to from 'await-to-js'
  import cloneDeep from 'lodash/cloneDeep'
  import EventBus from '@/utils/bus'
  import { Validator } from 'vee-validate'

  export default defineComponent({
    name: 'idle-machine-pool-config',
    components: {
      ModuleBuilder,
    },
    setup() {
      const globalConfig = computed(() => store.state.globalConfig)
      const nodeIndent = 40
      const iconOffset = 38
      const defaultIdleForm = {
        set: '', // 内置集群
        buildInModules: {}, // 内置模块
        userModules: [] // 用户自定义模块
      }
      const idleForm = reactive(cloneDeep(defaultIdleForm))
      const veeVlidate = new Validator()

      const businessTopoInstNameRule = value => ({
        required: true,
        message: globalConfig.value.config.validationRules.businessTopoInstNames.message,
        validator: async () => {
          const { valid } = await veeVlidate.verify(value(), 'businessTopoInstNames')
          return valid
        },
        trigger: 'blur'
      })

      const idleFormRules = reactive({
        set: [
          {
            required: true,
            message: t('请输入集群名称'),
            trigger: 'blur'
          },
          businessTopoInstNameRule(() => idleForm.set)
        ]
      })
      const idleFormRef = ref(null)
      const addBtnDisabled = computed(() => idleForm.userModules.some(userModule => userModule.state === 'editting'))

      /**
       * 生成内置模块验证规则
       * @param {string} moduleKey 模块 ID
       */
      const generateBuildInModuleRules = moduleKey => [
        {
          required: true,
          message: t('请输入模块 ID'),
          validator: () => idleForm.buildInModules[moduleKey] !== '',
          trigger: 'blur'
        },
        businessTopoInstNameRule(() => idleForm.buildInModules[moduleKey])
      ]
      /**
       * 生成扩展模块的验证规则
       * @param {Symbol} ruleKey 对应模块的标记，通过这个标记来找到规则所在模块
       */
      const generateUserModuleRules = ruleKey => ([
        {
          required: true,
          validator: async () => {
            const newModule = idleForm.userModules.find(userModule => userModule.ruleKey === ruleKey)
            return newModule.moduleKey?.trim() && newModule.moduleName?.trim()
          },
          message: () => {
            const newModule = idleForm.userModules.find(userModule => userModule.ruleKey === ruleKey)
            if (!newModule?.moduleKey?.trim()) {
              return t('请输入模块 ID')
            }
            if (!newModule?.moduleName?.trim()) {
              return t('请输入名称')
            }
          },
          trigger: 'blur'
        },
        businessTopoInstNameRule(() => idleForm.userModules
          .find(userModule => userModule.ruleKey === ruleKey).moduleName)
      ])

      const initForm = () => {
        const { idlePool, set } = cloneDeep(globalConfig.value.config)

        // 用户自定义模块处理
        const userModules = idlePool?.userModules?.map(({ moduleKey, moduleName }) => {
          const ruleKey = Symbol('ruleKey')
          return {
            ruleKey,
            moduleKey,
            moduleName,
            rules: generateUserModuleRules(ruleKey)
          }
        })
        delete idlePool.userModules
        Object.assign(idleForm, cloneDeep(defaultIdleForm), {
          set,
          buildInModules: idlePool,
          userModules: cloneDeep(userModules)
        })

        // 加入内置模块规则
        const buildInRules = {}
        Object.keys(idleForm.buildInModules).forEach((buildInModuleKey) => {
          buildInRules[buildInModuleKey] = generateBuildInModuleRules(buildInModuleKey)
        })
        Object.assign(idleFormRules, buildInRules)

        idleFormRef.value.clearError()
      }

      const clearError = () => idleFormRef.value.clearError()

      onMounted(() => {
        initForm()
        EventBus.$on('globalConfig/fetched', initForm)
      })

      const updateSet = done => idleFormRef.value.validate().then(() => {
        bkInfoBox({
          type: 'warning',
          title: t('确认修改集群？'),
          subTitle: t('执行修改，将在所有业务中立即生效，请谨慎操作'),
          okText: t('确认修改'),
          cancelText: t('取消'),
          confirmLoading: true,
          confirmFn: async () => {
            const [err] = await to(store.dispatch('globalConfig/updateIdleSet', {
              setKey: 'set',
              setName: idleForm.set,
            }))
            if (err) return false
            done()
            initForm()
            bkMessage({
              theme: 'success',
              message: t('修改成功')
            })
            return true
          }
        })
      })
        .catch((err) => {
          console.log(err)
        })

      const updateModule = ({ moduleKey, moduleName }, done) => {
        idleFormRef.value.validate()
          .then(() => {
            bkInfoBox({
              type: 'warning',
              title: t('确认修改模块？'),
              subTitle: t('执行修改，将在所有业务中立即生效，请谨慎操作'),
              okText: t('确认修改'),
              cancelText: t('取消'),
              confirmLoading: true,
              confirmFn: async () => {
                const [err] = await to(store.dispatch('globalConfig/updateIdleModule', {
                  moduleKey,
                  moduleName,
                }))
                if (err) return false
                done()
                initForm()
                bkMessage({
                  theme: 'success',
                  message: t('修改成功')
                })
                return true
              }
            })
          })
          .catch((err) => {
            console.log(err)
          })
      }

      const createModule = ({ moduleKey, moduleName }, done) => {
        idleFormRef.value.validate().then(() => {
          bkInfoBox({
            type: 'warning',
            title: t('确认新增模块？'),
            subTitle: t('执行新增，将在所有业务中立即生效，请谨慎操作'),
            okText: t('确认新增'),
            cancelText: t('取消'),
            confirmLoading: true,
            confirmFn: async () => {
              const [err] = await to(store.dispatch('globalConfig/createIdleModule', {
                moduleKey,
                moduleName,
              }))
              if (err) {
                return false
              }
              done()
              initForm()
              bkMessage({
                theme: 'success',
                message: t('新增成功')
              })
              return true
            }
          })
        })
      }

      // 只有用户自定义的模块可以删除
      const removeUserModule = (moduelId, moduleName) => {
        bkInfoBox({
          type: 'warning',
          title: t('确认删除模块？'),
          subTitle: t('执行删除，将在所有业务中立即生效，请谨慎操作'),
          okText: t('确认删除'),
          cancelText: t('取消'),
          confirmLoading: true,
          confirmFn: async () => {
            const [err] = await to(store.dispatch('globalConfig/deleteIdleModule', {
              moduleKey: moduelId,
              moduleName,
            }))
            if (err) return false
            initForm()
            bkMessage({
              theme: 'success',
              message: t('删除成功')
            })
            return true
          }
        })
      }

      // 增加新增模块的输入控件
      const addUserModule = () => {
        const ruleKey = Symbol('ruleKey')
        idleForm.userModules.push({
          ruleKey,
          isNew: true, // 标识模块是还没同步到远程的模块
          state: 'editting',
          moduleKey: '',
          moduleName: '',
          rules: generateUserModuleRules(ruleKey)
        })
      }

      const beforeClusterConfirm = (done) => {
        updateSet(done)
      }

      const beforeModuleConfirm = ({ moduleKey, moduleName }, done) => {
        updateModule({ moduleKey, moduleName }, done)
      }

      // 新增模块取消新增时，需要删除对应的输入控件
      const handleUserModuleCancel = (userModuleIndex, userModule) => {
        if (!userModule.moduleKey?.trim() && !userModule.moduleName?.trim()) {
          idleForm.userModules.splice(userModuleIndex, 1)
        }
        clearError()
      }

      // 新增模块包含已新增的模块和即将新增的模块，已新增的模块直接更新即可
      const beforeUserModuleConfirm = (userModule, done) => {
        if (userModule.isNew) {
          createModule({
            moduleKey: userModule.moduleKey,
            moduleName: userModule.moduleName
          }, done)
        } else {
          updateModule({
            moduleKey: userModule.moduleKey,
            moduleName: userModule.moduleName
          }, done)
        }
      }

      return {
        nodeIndent,
        iconOffset,
        globalConfig,
        addBtnDisabled,
        idleForm,
        idleFormRules,
        idleFormRef,
        addUserModule,
        removeUserModule,
        clearError,
        beforeClusterConfirm,
        beforeModuleConfirm,
        beforeUserModuleConfirm,
        handleUserModuleCancel
      }
    },
  })
</script>

<style lang="scss" scoped>
.idle-machine-pool-config{
  width: 700px;
}

::v-deep .action-form-item {
  margin-top: 16px;

  .bk-form-content{
    line-height: normal;
    height: auto;
  }

  .bk-button-text{
    margin-left: 40px;
  }
}
</style>
