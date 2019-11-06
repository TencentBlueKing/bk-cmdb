<template>
    <div class="host-search-layout">
        <div :class="['search-bar', { 'has-scroll': hasScroll && !showEllipsis }]">
            <bk-input class="search-input"
                ref="searchInput"
                type="textarea"
                :placeholder="$t('请输入IP，多个值换行分隔')"
                :rows="rows"
                v-model="searchContent"
                @focus="handleFocus"
                @blur="handleBlur"
                @keydown="handleKeydown">
            </bk-input>
            <i class="bk-icon search-btn icon-search" @click="handleSearch"></i>
            <span v-if="showEllipsis" class="search-text" @click="handleSearchInput">{{searchText}}</span>
        </div>
    </div>
</template>

<script>
    import { MENU_RESOURCE_HOST } from '@/dictionary/menu-symbol'
    export default {
        data () {
            return {
                rows: 1,
                searchText: '',
                searchContent: '',
                textarea: '',
                hasScroll: false,
                showEllipsis: false,
                textareaDom: null
            }
        },
        watch: {
            searchContent () {
                this.handleObserverScroll()
            }
        },
        mounted () {
            const textarea = this.$refs.searchInput && this.$refs.searchInput.$refs.textarea
            this.textareaDom = textarea
        },
        methods: {
            setRows () {
                const rows = this.searchContent.split('\n').length || 1
                this.rows = Math.min(10, rows)
            },
            handleKeydown (value, keyEvent) {
                if (['Enter', 'NumpadEnter'].includes(keyEvent.code)) {
                    this.$nextTick(() => {
                        this.rows = Math.min(this.rows + 1, 10)
                    })
                } else if (keyEvent.code === 'Backspace') {
                    this.$nextTick(() => {
                        this.setRows()
                    })
                }
            },
            handleObserverScroll () {
                this.$nextTick(() => {
                    if (this.textareaDom) {
                        this.hasScroll = this.textareaDom.scrollHeight > this.textareaDom.offsetHeight
                        this.$nextTick(() => {
                            this.setRows()
                        })
                    }
                })
            },
            handleFocus () {
                this.$emit('focus-status', true)
                this.setRows()
            },
            handleBlur () {
                this.$emit('focus-status', false)
                const data = this.searchContent.split('\n').filter(text => text)
                if (data.length) {
                    this.showEllipsis = true
                    this.searchText = data.join(',')
                }
                this.$nextTick(() => {
                    this.rows = 1
                    this.textareaDom && (this.textareaDom.scrollTop = 0)
                })
            },
            handleSearchInput () {
                this.showEllipsis = false
                this.textareaDom && this.textareaDom.focus()
            },
            async handleSearch () {
                const searchList = this.searchContent.split('\n').filter(ip => ip)
                if (searchList.length) {
                    const validateQueue = []
                    searchList.forEach(text => {
                        validateQueue.push(this.$validator.verify(text, 'ip'))
                    })
                    const results = await Promise.all(validateQueue)
                    const isPassValidate = results.every(res => res.valid)
                    if (!isPassValidate) {
                        this.$warn('请输入完整IP进行搜索，多个IP用换行分割')
                        return
                    } else if (results.length > 500) {
                        this.$warn('目前最多支持搜索500个IP')
                        return
                    }
                    this.$store.commit('hosts/setFilterIP', {
                        text: searchList.join('\n'),
                        exact: true
                    })
                    this.$store.commit('hosts/setIsHostSearch', true)
                    this.$router.push({ name: MENU_RESOURCE_HOST })
                } else {
                    this.searchContent = ''
                    this.textareaDom && this.textareaDom.focus()
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-search-layout {
        width: 100%;
        max-width: 640px;
        margin: 0 auto;
    }
    .search-bar {
        position: relative;
        &.has-scroll {
            .search-btn {
                margin-right: 17px !important;
            }
        }
    }
    .search-input {
        /deep/ {
            .bk-textarea-wrapper {
                border: 0;
            }
            .bk-form-textarea {
                min-height: 42px;
                line-height: 30px;
                font-size: 14px;
                border: 1px solid #C4C6CC;
                padding: 5px 50px 5px 16px;
            }
        }
    }
    .search-btn {
        position: absolute;
        right: 0;
        top: 0;
        width: 50px;
        height: 42px;
        line-height: 42px;
        color: #C3CDD7;
        font-size: 18px;
        text-align: center;
        z-index: 5;
        cursor: pointer;
        &.icon-close-circle-shape:hover {
            color: #979BA5;
        }
    }
    .search-text {
        position: absolute;
        right: 0;
        top: 0;
        width: 100%;
        height: 42px;
        line-height: 30px;
        font-size: 14px;
        color: #63656E;
        border: 1px solid #C4C6CC;
        background-color: #FFFFFF;
        padding: 5px 50px 5px 16px;
        z-index: 1;
        cursor: text;
        @include ellipsis;
    }
</style>
