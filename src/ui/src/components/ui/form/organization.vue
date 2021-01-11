<template>
    <div class="cmdb-organization-select"
        :class="{
            'is-focus': focused,
            'is-disabled': disabled,
            'is-readonly': readonly,
            'is-unselected': unselected
        }"
        :data-placeholder="placeholder">
        <i class="select-loading" v-if="$loading([searchRequestId]) && searchValue === undefined"></i>
        <i class="select-angle bk-icon icon-angle-down"></i>
        <i class="select-clear bk-icon icon-close"
            v-if="clearable && !unselected && !disabled && !readonly"
            @click.prevent.stop="handleClear">
        </i>
        <bk-popover class="select-dropdown"
            ref="selectDropdown"
            trigger="click"
            placement="bottom-start"
            theme="light select-dropdown"
            animation="slide-toggle"
            :disabled="disabled"
            :z-index="zIndex"
            :arrow="false"
            :offset="-1"
            :distance="12"
            :on-show="handleDropdownShow"
            :on-hide="handleDropdownHide">
            <div class="select-name"
                :title="displayName">
                {{displayName}}
            </div>
            <div slot="content" :style="{ width: popoverWidth + 'px' }" class="select-dropdown-content">
                <div class="search-bar">
                    <bk-input
                        :placeholder="$t('搜索')"
                        ext-cls="search-input"
                        right-icon="bk-icon icon-search"
                        clearable
                        v-model.trim="searchValue"
                        @input="handleSearch">
                    </bk-input>
                </div>
                <bk-big-tree class="org-tree"
                    ref="tree"
                    v-bkloading="{ isLoading: $loading([searchRequestId]) }"
                    v-bind="treeProps"
                    @check-change="handleCheckChange">
                    <div class="tree-node" slot-scope="{ node, data: nodeData }"
                        :class="{ 'is-selected': node.selected }"
                        :title="nodeData.name">
                        <div class="node-name">{{nodeData.name}}</div>
                    </div>
                </bk-big-tree>
            </div>
        </bk-popover>
    </div>
</template>

<script>
    import debounce from 'lodash.debounce'
    export default {
        name: 'cmdb-form-organization',
        props: {
            value: {
                type: [Array, String],
                default: []
            },
            disabled: {
                type: Boolean,
                default: false
            },
            readonly: Boolean,
            multiple: Boolean,
            clearable: Boolean,
            placeholder: {
                type: String,
                default: ''
            },
            zIndex: {
                type: Number,
                default: 2500
            },
            formatter: Function
        },
        data () {
            return {
                focused: false,
                checked: this.value || [],
                popoverWidth: 0,
                searchRequestId: Symbol('orgSearch'),
                searchValue: undefined,
                displayName: '',
                viewName: '',
                treeProps: {
                    showCheckbox: true,
                    checkOnClick: false,
                    checkStrictly: false,
                    lazyMethod: this.lazyMethod,
                    lazyDisabled: this.lazyDisabled
                }
            }
        },
        computed: {
            unselected () {
                return !this.checked.length
            }
        },
        watch: {
            value (value) {
                this.checked = value
            },
            checked (checked) {
                this.$emit('input', checked)
                this.$emit('on-checked', checked)
                this.setDisplayName()
            },
            focused (focused) {
                this.$emit('toggle', focused)
            }
        },
        created () {
            this.initTree()
        },
        methods: {
            async initTree () {
                await this.loadTree()
                this.setDisplayName()
            },
            async loadTree () {
                const { data: topData } = await this.getLazyData()
                const defaultChecked = this.checked
                const tree = this.$refs.tree

                if (defaultChecked.length) {
                    const checkedRes = await this.getSearchData({
                        lookup_field: 'id',
                        exact_lookups: defaultChecked.join(','),
                        with_ancestors: true
                    })

                    // 可能为空，节点数据存在才获取相关联数据
                    const chcekedData = checkedRes.results || []
                    if (chcekedData.length) {
                        // 已选中节点的树形数据
                        const chekcedTreeData = this.getTreeSearchData(chcekedData)
                        // 已选中节点的完整树形数据（含兄弟节点）
                        const fullCheckedTreeData = await this.getCheckedFullTreeData(chekcedTreeData)
                        // 将匹配的树分支替换以合并
                        fullCheckedTreeData.forEach(checkedNode => {
                            const matchedIndex = topData.findIndex(top => top.id === checkedNode.id)
                            if (matchedIndex !== -1) {
                                topData[matchedIndex] = checkedNode
                            }
                        })

                        // 设置树数据，选中节点数据已被完整包含
                        tree.setData(topData)

                        // 将选中节点全部展开
                        defaultChecked.forEach(id => tree.setExpanded(id))
                        // 设置为选中状态
                        tree.setChecked(defaultChecked)
                    } else {
                        tree.setData(topData)
                    }
                } else {
                    tree.setData(topData)
                }
            },
            async setViewData () {
                if (!this.checked.length) {
                    this.viewName = '--'
                    return
                }

                const res = await this.getSearchData({
                    lookup_field: 'id',
                    exact_lookups: this.checked.join(',')
                })
                const names = (res.results || []).map(item => item.full_name)
                this.viewName = this.formatName(names)
            },
            async getLazyData (parentId) {
                try {
                    const params = {
                        lookup_field: 'level',
                        exact_lookups: 0
                    }
                    const config = {
                        fromCache: !parentId,
                        requestId: `get_org_department_${!parentId ? '0' : parentId}`
                    }
                    if (parentId) {
                        params.lookup_field = 'parent'
                        params.exact_lookups = parentId
                    }
                    const res = await this.$store.dispatch('organization/getDepartment', { params, ...config })
                    const data = res.results || []
                    return { data }
                } catch (e) {
                    console.error(e)
                }
            },
            getSearchData (params) {
                return this.$store.dispatch('organization/getDepartment', {
                    params,
                    requestId: this.searchRequestId
                })
            },
            resetTree () {
                this.checked = []
                this.loadTree()
            },
            lazyMethod (node) {
                return this.getLazyData(node.id)
            },
            lazyDisabled (node) {
                return !node.data.has_children
            },
            setDisplayName () {
                const tree = this.$refs.tree
                const nodes = this.checked.map(id => tree.getNodeById(id)).filter(node => !!node)
                const displayNames = nodes.map(node => node.data.full_name)
                this.displayName = this.formatName(displayNames)
            },
            formatName (names) {
                let name = ''
                if (this.formatter) {
                    name = this.formatter(names)
                } else {
                    name = names.join('; ')
                }
                return name
            },
            async getCheckedFullTreeData (chekcedTreeData) {
                // 获取所有节点id
                const ids = []
                const getId = (nodes) => {
                    nodes.forEach(node => {
                        ids.push(node.id)
                        if (node.children) {
                            getId(node.children)
                        }
                    })
                }
                getId(chekcedTreeData)

                // 获取所有节点的子节点
                const childNodeRes = await this.getSearchData({
                    lookup_field: 'parent',
                    exact_lookups: ids.join(','),
                    with_ancestors: false
                })
                const childNodeList = childNodeRes.results || []

                // 将子节点补齐到对应的目标节点
                const appendChild = (nodes) => {
                    nodes.forEach(node => {
                        childNodeList.forEach(child => {
                            if (child.parent === node.id) {
                                if (node.children) {
                                    const childIds = node.children.map(item => item.id)
                                    if (childIds.indexOf(child.id) === -1) {
                                        node.children.push(child)
                                    }
                                } else {
                                    node.children = [child]
                                }
                            }
                        })

                        if (node.children) {
                            appendChild(node.children)
                        }
                    })
                }
                appendChild(chekcedTreeData)

                return chekcedTreeData
            },
            getTreeSearchData (data) {
                // 将偏平的数据组装成树形结构
                const treeData = []
                data.forEach(item => {
                    const ancestorLength = item.ancestors.length
                    const curNode = {
                        id: item.id,
                        name: item.name,
                        level: ancestorLength,
                        full_name: item.full_name
                    }
                    const ids = [curNode.id]
                    const treeNode = {}
                    for (let i = ancestorLength - 1; i >= 0; i--) {
                        const node = item.ancestors[i]
                        ids.push(node.id)
                        node.level = i
                        node.children = [item.ancestors[i + 1] ? item.ancestors[i + 1] : curNode]
                        node.full_name = item.full_name.split('/', i + 1).join('/')
                    }

                    treeNode.ids = ids.reverse()
                    if (item.ancestors[0]) {
                        treeNode.map = item.ancestors[0]
                    } else {
                        treeNode.map = curNode
                    }

                    treeData.push(treeNode)
                })

                // 合并与去重
                for (let i = 0; i < treeData.length; i++) {
                    const node = treeData[i]
                    const path = node.ids.join('-')
                    for (let j = i + 1; j < treeData.length; j++) {
                        const nodeNext = treeData[j]
                        let k = nodeNext.ids.length
                        while (k) {
                            const pathNext = nodeNext.ids.slice(0, k).join('-')
                            // 路径比较，将被比较对象除重复部分外的数据合并至比较对象
                            if (path.indexOf(pathNext) !== -1) {
                                const nextRest = data[j].ancestors.slice(k - 1)
                                const appendToNode = data[i].ancestors[k - 1]
                                if (appendToNode && nextRest.length) {
                                    // 合并时去重
                                    const exists = appendToNode.children.map(item => item.id)
                                    nextRest[0].children.forEach(item => {
                                        if (exists.indexOf(item.id) === -1) {
                                            appendToNode.children.push(item)
                                        }
                                    })
                                }
                                nodeNext.remove = true
                                break
                            } else if (pathNext.indexOf(path) !== -1) {
                                // 如果路径被反向包含则可直接删除
                                node.remove = true
                                break
                            }
                            k--
                        }
                    }
                }

                // 得到最终用于树的数据
                const finalTreeData = treeData.filter(item => !item.remove).map(item => item.map)
                return finalTreeData
            },
            setTreeSearchData (data) {
                const tree = this.$refs.tree
                const finalTreeData = this.getTreeSearchData(data)
                tree.setData(finalTreeData)
                finalTreeData.forEach(node => {
                    tree.setExpanded(node.id)
                })
            },
            handleSearch: debounce(async function (value) {
                const keyword = value.trim()
                try {
                    if (keyword.length) {
                        const res = await this.getSearchData({
                            lookup_field: 'name',
                            fuzzy_lookups: keyword,
                            with_ancestors: true
                        })
                        const data = res.results || []
                        this.setTreeSearchData(data)
                    } else if (!value.length) {
                        this.loadTree()
                    }
                } catch (e) {
                    console.error(e)
                }
            }, 160),
            handleClear () {
                this.resetTree()
            },
            handleCheckChange (ids, node) {
                if (this.multiple) {
                    this.checked = ids
                } else {
                    const tree = this.$refs.tree
                    tree.removeChecked({ emitEvent: false })
                    tree.setChecked(node.id, { emitEvent: false })
                    this.checked = [node.id]
                    this.$refs.selectDropdown.instance.hide()
                }
            },
            handleDropdownShow () {
                this.popoverWidth = this.$el.offsetWidth
                this.focused = true
            },
            handleDropdownHide () {
                this.focused = false
            },
            focus () {
                this.$refs.selectDropdown.instance.show()
            }
        }
    }
</script>

<style lang="scss">
    .tippy-tooltip {
        &.select-dropdown-theme {
            padding: 0;
            box-shadow: 0 3px 9px 0 rgba(0, 0, 0, .1);
        }
    }

    .search-input {
        .bk-form-input {
            &:focus {
                border-color: #c4c6cc !important;
            }
        }
    }
</style>
<style lang="scss" scoped>
    .cmdb-organization-select {
        position: relative;
        width: 100%;
        border: 1px solid #c4c6cc;
        background-color: #fff;
        border-radius: 2px;
        line-height: 30px;
        color: #63656e;
        cursor: pointer;
        font-size: 12px;

        &.is-focus {
            border-color: #3a84ff;
            box-shadow:0px 0px 4px rgba(58, 132, 255, 0.4);
            .select-angle {
                transform: rotate(-180deg);
            }
        }
        &.is-disabled {
            background-color: #fafbfd;
            border-color: #dcdee5;
            color: #c4c6cc;
            cursor: not-allowed;
        }
        &.is-readonly,
        &.is-loading {
            background-color: #fafbfd;
            border-color: #dcdee5;
            cursor: default;
        }

        &.is-unselected::before {
            position: absolute;
            height: 100%;
            content: attr(data-placeholder);
            left: 10px;
            top: 0;
            color: #c3cdd7;
            pointer-events: none;
        }

        &:hover {
            .select-clear {
                display: block;
            }
        }

        .select-angle {
            position: absolute;
            right: 2px;
            top: 4px;
            font-size: 22px;
            color: #979ba5;
            transition: transform .3s cubic-bezier(0.4, 0, 0.2, 1);
            pointer-events: none;
        }

        .select-clear {
            display: none;
            position: absolute;
            right: 6px;
            top: 8px;
            width: 14px;
            height: 14px;
            line-height: 14px;
            background-color: #c4c6cc;
            border-radius: 50%;
            text-align: center;
            font-size: 14px;
            color: #fff;
            z-index: 100;
            &:before {
                display: block;
                transform: scale(.7);
            }
            &:hover {
                background-color: #979ba5;
            }
        }

        .select-loading {
            position: absolute;
            top: 8px;
            left: 8px;
            width: 16px;
            height: 16px;
            background-image: url("../../../assets/images/icon/loading.svg");
            z-index: 1;
        }

        .select-dropdown {
            display: block;

            .select-name {
                height: 30px;
                padding: 0 36px 0 10px;
                @include ellipsis;
            }
        }
    }

    .select-dropdown-content {
        border: 1px solid #dcdee5;
        border-radius: 2px;
        line-height: 32px;
        background: #fff;
        color: #63656e;
        overflow: hidden;

        .search-bar {
            padding: 10px;
        }

        .org-tree {
            height: 220px !important;

            .tree-node {
                .node-name {
                    @include ellipsis;
                }
            }
        }
    }

    /deep/.bk-tooltip {
        > .bk-tooltip-ref {
            display: block;
        }
    }
</style>
