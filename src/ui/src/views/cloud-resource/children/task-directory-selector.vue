<template>
    <bk-select
        v-if="display === 'selector'"
        v-model="selected"
        ref="selector"
        searchable
        size="small"
        font-size="small"
        :loading="loading"
        :popover-width="260"
        :clearable="false"
        :placeholder="$t('请选择xx', { name: $t('资源目录') })"
        @toggle="handleSelectToggle">
        <bk-option v-for="directory in directories"
            :key="directory.bk_module_id"
            :id="directory.bk_module_id"
            :disabled="!directory.authorized"
            :name="directory.bk_module_name">
            <div v-cursor="getCursorData(directory)">
                {{directory.bk_module_name}}
            </div>
        </bk-option>
        <bk-option class="create-option"
            v-if="createMode"
            ref="createOptionComponent"
            :id="createDirectoryId"
            :name="$t('新增目录')"
            :disabled="true"
            @click.native.stop>
            <bk-input ref="input"
                size="small"
                font-size="small"
                :placeholder="$t('请输入目录名称，回车结束')"
                v-model.trim="newDirectory"
                @enter="handleConfirmCreate">
            </bk-input>
        </bk-option>
        <cmdb-auth tag="a" href="javascript:void(0)" class="extension-link" slot="extension"
            :auth="{ type: $OPERATION.C_RESOURCE_DIRECTORY }"
            :onclick="hideSelectorPanel"
            @click="handleCreateDirectory">
            <i class="bk-icon icon-plus-circle"></i>
            {{$t('新增目录')}}
        </cmdb-auth>
    </bk-select>
    <span v-else>{{getInfo()}}</span>
</template>

<script>
    import symbols from '../common/symbol'
    import AuthProxy from '@/components/ui/auth/auth-queue'
    export default {
        name: 'task-directory-selector',
        props: {
            value: {
                type: [String, Number],
                default: ''
            },
            display: {
                type: String,
                default: 'selector'
            }
        },
        data () {
            return {
                loading: true,
                createMode: false,
                createDirectoryId: 'createDirectoryId',
                newDirectory: '',
                directories: [],
                request: {
                    findMany: symbols.get('directory'),
                    create: Symbol('create')
                }
            }
        },
        computed: {
            selected: {
                get () {
                    return this.value
                },
                set (value, oldValue) {
                    this.$emit('input', value)
                    this.$emit('change', value, oldValue)
                }
            },
            useIAM () {
                return window.CMDB_CONFIG.site.authscheme === 'iam'
            }
        },
        created () {
            this.getDirectories()
        },
        methods: {
            getCursorData (directory) {
                return {
                    active: this.useIAM ? !directory.authorized : false,
                    auth: { type: this.$OPERATION.C_RESOURCE_HOST, relation: [directory.bk_module_id] },
                    onclick: this.hideSelectorPanel
                }
            },
            hideSelectorPanel () {
                this.$refs.selector.close()
            },
            async getDirectories () {
                try {
                    this.loading = true
                    const { info } = await this.$store.dispatch('resource/directory/findMany', {
                        params: {
                            sort: 'bk_module_name'
                        },
                        config: {
                            requestId: this.request.findMany,
                            fromCache: true
                        }
                    })
                    if (this.display === 'selector' && this.useIAM) {
                        await this.injectAuth(info)
                    }
                    info.sort((dirA, dirB) => {
                        const aAuth = dirA.authorized ? 1 : 0
                        const bAuth = dirB.authorized ? 1 : 0
                        return bAuth - aAuth
                    })
                    this.directories = info
                    if (!this.selected && info.length) {
                        this.selected = info[0].bk_module_id
                    }
                } catch (e) {
                    this.directories = []
                    console.error(e)
                } finally {
                    this.loading = false
                }
            },
            injectAuth (directories) {
                return new Promise(resolve => {
                    const fakeComponent = {
                        auth: directories.map(directory => ({ type: this.$OPERATION.C_RESOURCE_HOST, relation: [directory.bk_module_id] })),
                        updateAuth: results => {
                            directories.forEach(directory => {
                                const result = results.find(result => result.parent_layers[0].resource_id === directory.bk_module_id)
                                this.$set(directory, 'authorized', result ? result.is_pass : false)
                            })
                            resolve()
                        }
                    }
                    AuthProxy.add({
                        component: fakeComponent,
                        data: fakeComponent.auth
                    })
                })
            },
            handleSelectToggle (isVisible) {
                if (!isVisible) {
                    this.toggleCreate(false)
                }
            },
            handleCreateDirectory () {
                this.toggleCreate(true)
            },
            toggleCreate (isCreateMode) {
                this.createMode = isCreateMode
                if (isCreateMode) {
                    this.$nextTick(() => {
                        this.$refs.input.focus()
                    })
                } else {
                    this.newDirectory = ''
                    this.$refs.selector.close()
                }
            },
            async handleConfirmCreate () {
                if (!this.newDirectory.length) {
                    return false
                }
                try {
                    const data = await this.$store.dispatch('resource/directory/create', {
                        params: {
                            bk_module_name: this.newDirectory
                        },
                        config: {
                            requestId: this.request.create
                        }
                    })
                    this.directories.push({
                        bk_module_id: data.created.id,
                        bk_module_name: this.newDirectory
                    })
                    this.selected = data.created.id
                    this.toggleCreate(false)
                    this.injectAuth(this.directories)
                } catch (e) {
                    console.error(e)
                }
            },
            getInfo () {
                const directory = this.directories.find(directory => directory.bk_module_id === this.value)
                if (directory) {
                    return directory.bk_module_name
                }
                return null
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-option.is-disabled {
        cursor: pointer;
        color: $textColor;
    }
    .extension-link.disabled {
        color: $textDisabledColor;
    }
</style>
