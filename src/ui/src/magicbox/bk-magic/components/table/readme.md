<script>
    export default {
        data () {
            return {
                data: [{
                    'id': 1463,
                    'name': '张飞',
                    'age': 41,
                    'email': 'z.qgeyitbrzg@sdhra.by',
                    'url': 'nntp://zlqtkmfd.bf/latxycl'
                },
                {
                    'id': 1464,
                    'name': 'Mary Moore',
                    'age': 34,
                    'email': 't.igxjxd@vhpsh.uk',
                    'url': 'rlogin://pmeaerguys.sn/hjckcxjne'
                },
                {
                    'id': 1465,
                    'name': '曹操',
                    'age': 32,
                    'email': 'k.ngmmwcpjq@dfsxkywu.aw',
                    'url': 'ftp://lxqpcndr.sc/iqjzrn'
                },
                {
                    'id': 1466,
                    'name': 'Elizabeth Jones',
                    'age': 36,
                    'email': 'e.nscpdfwhc@bjdmp.tf',
                    'url': 'mailto://jfhmpdp.ma/xmwpidrl'
                },
                {
                    'id': 1467,
                    'name': 'Timothy White',
                    'age': 21,
                    'email': 't.fxcdx@wcjln.io',
                    'url': 'gopher://ugaat.gu/ybohv'
                },
                {
                    'id': 1468,
                    'name': '齐备',
                    'age': 30,
                    'email': 'g.rmena@ibuir.mc',
                    'url': 'ftp://qrwpjx.by/arbipfhwiv'
                }],
                tableConf1: {
                    tableFields: [{
                        name: 'id',
                        title: 'ID'
                    },
                    {
                        name: 'name',
                        title: '名称'
                    },
                    {
                        name: 'age',
                        title: '年龄'
                    },
                    {
                        name: 'email',
                        title: 'Email'
                    },
                    {
                        name: 'url',
                        title: '官网'
                    }]
                },
                tableConf2: {
                    tableFields: [{
                        name: 'id',
                        title: 'ID',
                        sortable: true
                    },
                    {
                        name: 'name',
                        title: '名称'
                    },
                    {
                        name: 'age',
                        title: '年龄',
                        sortable: true
                    },
                    {
                        name: 'email',
                        title: 'Email'
                    }]
                },
                tableConf3: {
                    curPage: 1,
                    totalPage: 3,
                    pageSize: 2,
                    tableFields: [{
                        name: 'id',
                        title: 'ID'
                    },
                    {
                        name: 'name',
                        title: '名称'
                    },
                    {
                        name: 'age',
                        title: '年龄'
                    },
                    {
                        name: 'email',
                        title: 'Email'
                    },
                    {
                        name: 'url',
                        title: '官网'
                    }]
                },
                tableConf4: {
                    curPage: 1,
                    totalPage: 3,
                    pageSize: 2,
                    tableFields: [{
                        name: 'id',
                        title: 'ID'
                    },
                    {
                        name: 'name',
                        title: '名称'
                    },
                    {
                        name: 'age',
                        title: '年龄'
                    },
                    {
                        name: 'email',
                        title: 'Email'
                    },
                    {
                        name: '__slot:actions',
                        title: '操作'
                    }]
                }
            }
        },
        methods: {
            editRow (data, index) {
                console.log(data)
                console.log(index)
                alert('edit')
            },
            deleteRow (data, index) {
                this.$refs.bkComplexTable.removeRow(data)
            },
            deleteSelectedRow () {
                this.$refs.bkComplexTable.removeSelectedRow()
            },
            renderBaseTableData (page) {
                this.$refs.bkBaseTable.renderCurPageData(page)
            },
            renderComplexTableData (page) {
                this.$refs.bkComplexTable.renderCurPageData(page)
            },
            removeRowCallback (data) {
                this.$refs.bkComplexTablePagination.totalPage = data.totalPage
                this.$refs.bkComplexTablePagination.curPage = data.curPage
            }
        }
    }
</script>

## Table 表格

用于展示数据，支持分页，可对数据进行排序、自定义操作

### 基础用法

基础的表格展示用法

:::demo 通过设置`data`来传入数据，`fields`来配置表头字段

```html
<template>
    <bk-table
        class-name="has-table-striped"
        :data="data"
        :fields="tableConf1.tableFields">
    </bk-table>
</template>
<script>
    export default {
        data () {
            return {
                data: [{
                    'id': 1463,
                    'name': 'Edward Brown',
                    'age': 41,
                    'email': 'z.qgeyitbrzg@sdhra.by',
                    'url': 'nntp://zlqtkmfd.bf/latxycl'
                }, 
                {
                    'id': 1464,
                    'name': 'Mary Moore',
                    'age': 34,
                    'email': 't.igxjxd@vhpsh.uk',
                    'url': 'rlogin://pmeaerguys.sn/hjckcxjne'
                }, 
                {
                    'id': 1465,
                    'name': 'Jessica Young',
                    'age': 32,
                    'email': 'k.ngmmwcpjq@dfsxkywu.aw',
                    'url': 'ftp://lxqpcndr.sc/iqjzrn'
                }, 
                {
                    'id': 1466,
                    'name': 'Elizabeth Jones',
                    'age': 36,
                    'email': 'e.nscpdfwhc@bjdmp.tf',
                    'url': 'mailto://jfhmpdp.ma/xmwpidrl'
                }, 
                {
                    'id': 1467,
                    'name': 'Timothy White',
                    'age': 21,
                    'email': 't.fxcdx@wcjln.io',
                    'url': 'gopher://ugaat.gu/ybohv'
                }],
                tableConf1: {
                    tableFields: [{
                        name: 'id',
                        title: 'ID'
                    }, 
                    {
                        name: 'name',
                        title: '名称'
                    }, 
                    {
                        name: 'age',
                        title: '年龄'
                    }, 
                    {
                        name: 'email',
                        title: 'Email'
                    }, 
                    {
                        name: 'url',
                        title: '官网'
                    }]
                }
            }
        }
    }
</script>
```
:::


### 带排序用法

基础的表格和分页组件结合来展示数据

:::demo 可以通过`fields`的`sortable`属性来显示排序

```html
<template>
    <bk-table
        ref="bkBaseTable"
        :data="data"
        :fields="tableConf2.tableFields">
    </bk-table>
</template>
<script>
    export default {
        data () {
            return {
                data: [{
                    'id': 1463,
                    'name': '张飞',
                    'age': 41,
                    'email': 'z.qgeyitbrzg@sdhra.by',
                    'url': 'nntp://zlqtkmfd.bf/latxycl'
                }, 
                {
                    'id': 1464,
                    'name': 'Mary Moore',
                    'age': 34,
                    'email': 't.igxjxd@vhpsh.uk',
                    'url': 'rlogin://pmeaerguys.sn/hjckcxjne'
                }, 
                {
                    'id': 1465,
                    'name': '曹操',
                    'age': 32,
                    'email': 'k.ngmmwcpjq@dfsxkywu.aw',
                    'url': 'ftp://lxqpcndr.sc/iqjzrn'
                }, 
                {
                    'id': 1466,
                    'name': 'Elizabeth Jones',
                    'age': 36,
                    'email': 'e.nscpdfwhc@bjdmp.tf',
                    'url': 'mailto://jfhmpdp.ma/xmwpidrl'
                }, 
                {
                    'id': 1467,
                    'name': 'Timothy White',
                    'age': 21,
                    'email': 't.fxcdx@wcjln.io',
                    'url': 'gopher://ugaat.gu/ybohv'
                }, 
                {
                    'id': 1468,
                    'name': '齐备',
                    'age': 30,
                    'email': 'g.rmena@ibuir.mc',
                    'url': 'ftp://qrwpjx.by/arbipfhwiv'
                }],
                tableConf2: {
                    tableFields: [{
                        name: 'id',
                        title: 'ID',
                        sortable: true
                    }, 
                    {
                        name: 'name',
                        title: '名称'
                    }, 
                    {
                        name: 'age',
                        title: '年龄',
                        sortable: true
                    }, 
                    {
                        name: 'email',
                        title: 'Email'
                    }]
                }
            }
        }
    }
</script>
```
:::


### 带分页用法

基础的表格和分页组件结合来展示数据

:::demo 和分页组件`bk-paging`结合使用

```html
<template>
    <bk-table
        ref="bkBaseTable"
        :data="data"
        :cur-page="tableConf3.curPage"
        :page-size="tableConf3.pageSize"
        :fields="tableConf3.tableFields">
    </bk-table>
    <div class="mt10">
        <bk-paging
          :cur-page.sync="tableConf3.curPage"
          @page-change="renderBaseTableData"
          :total-page="tableConf3.totalPage"
          :type="'compact'">
        </bk-paging>
    </div>
</template>
<script>
    export default {
        data () {
            return {
                data: [{
                    'id': 1463,
                    'name': '张飞',
                    'age': 41,
                    'email': 'z.qgeyitbrzg@sdhra.by',
                    'url': 'nntp://zlqtkmfd.bf/latxycl'
                }, 
                {
                    'id': 1464,
                    'name': 'Mary Moore',
                    'age': 34,
                    'email': 't.igxjxd@vhpsh.uk',
                    'url': 'rlogin://pmeaerguys.sn/hjckcxjne'
                }, 
                {
                    'id': 1465,
                    'name': '曹操',
                    'age': 32,
                    'email': 'k.ngmmwcpjq@dfsxkywu.aw',
                    'url': 'ftp://lxqpcndr.sc/iqjzrn'
                }, 
                {
                    'id': 1466,
                    'name': 'Elizabeth Jones',
                    'age': 36,
                    'email': 'e.nscpdfwhc@bjdmp.tf',
                    'url': 'mailto://jfhmpdp.ma/xmwpidrl'
                }, 
                {
                    'id': 1467,
                    'name': 'Timothy White',
                    'age': 21,
                    'email': 't.fxcdx@wcjln.io',
                    'url': 'gopher://ugaat.gu/ybohv'
                }, 
                {
                    'id': 1468,
                    'name': '齐备',
                    'age': 30,
                    'email': 'g.rmena@ibuir.mc',
                    'url': 'ftp://qrwpjx.by/arbipfhwiv'
                }],
                tableConf3: {
                    curPage: 1,
                    totalPage: 3,
                    pageSize: 2,
                    tableFields: [{
                        name: 'id',
                        title: 'ID'
                    }, 
                    {
                        name: 'name',
                        title: '名称'
                    }, 
                    {
                        name: 'age',
                        title: '年龄'
                    }, 
                    {
                        name: 'email',
                        title: 'Email'
                    }, 
                    {
                        name: 'url',
                        title: '官网'
                    }]
                }
            }
        },
        methods: {
            renderBaseTableData (page) {
                this.$refs.bkBaseTable.renderCurPageData(page)
            }
        }
    }
</script>
```
:::

### 复杂表格

可对数据进行自定义操作

:::demo 可以通过自定义模版实现扩展

```html
<template>
    <div class="bk-panel bk-demo">
        <div class="bk-panel-header" role="tab">
            <div class="bk-panel-info fl">
                <div class="panel-title">人员列表</div>
                <div class="panel-subtitle">提示：如果处理失败，会有....</div>
            </div>
            <div class="bk-panel-action fr">
                <div class="bk-form bk-inline-form bk-form-small">
                </div>
            </div>
        </div>
        <div class="bk-panel-body p0">
              <bk-table
                ref="bkComplexTable"
                :data="data"
                :page-size="tableConf4.pageSize"
                :cur-page="tableConf4.curPage"
                :fields="tableConf4.tableFields"
                @remove-row="removeRowCallback">
                <template slot="actions" scope="props">
                    <div class="table-button-container">
                        <button class="bk-button bk-primary bk-button-mini" @click="editRow(props.rowData, props.rowIndex)">编辑</button>
                        <button class="bk-button bk-danger bk-button-mini" @click="deleteRow(props.rowData, props.rowIndex)">删除</button>
                    </div>
                </template>
              </bk-table>
        </div>
        <div class="bk-panel-footer p10">
            <div class="fr">
                <bk-paging
                    ref="bkComplexTablePagination"
                    :cur-page.sync="tableConf4.curPage"
                    @page-change="renderComplexTableData"
                    :total-page="tableConf4.totalPage"
                    :type="'compact'">
                </bk-paging>
            </div>
        </div>
    </div>
</template>
<script>
    export default {
        data () {
            return {
                data: [{
                    'id': 1463,
                    'name': '张飞',
                    'age': 41,
                    'email': 'z.qgeyitbrzg@sdhra.by',
                    'url': 'nntp://zlqtkmfd.bf/latxycl'
                }, 
                {
                    'id': 1464,
                    'name': 'Mary Moore',
                    'age': 34,
                    'email': 't.igxjxd@vhpsh.uk',
                    'url': 'rlogin://pmeaerguys.sn/hjckcxjne'
                }, 
                {
                    'id': 1465,
                    'name': '曹操',
                    'age': 32,
                    'email': 'k.ngmmwcpjq@dfsxkywu.aw',
                    'url': 'ftp://lxqpcndr.sc/iqjzrn'
                }, 
                {
                    'id': 1466,
                    'name': 'Elizabeth Jones',
                    'age': 36,
                    'email': 'e.nscpdfwhc@bjdmp.tf',
                    'url': 'mailto://jfhmpdp.ma/xmwpidrl'
                }, 
                {
                    'id': 1467,
                    'name': 'Timothy White',
                    'age': 21,
                    'email': 't.fxcdx@wcjln.io',
                    'url': 'gopher://ugaat.gu/ybohv'
                }, 
                {
                    'id': 1468,
                    'name': '齐备',
                    'age': 30,
                    'email': 'g.rmena@ibuir.mc',
                    'url': 'ftp://qrwpjx.by/arbipfhwiv'
                }],
                tableConf4: {
                    curPage: 1,
                    totalPage: 3,
                    pageSize: 2,
                    tableFields: [{
                        name: 'id',
                        title: 'ID'
                    }, 
                    {
                        name: 'name',
                        title: '名称'
                    }, 
                    {
                        name: 'age',
                        title: '年龄'
                    }, 
                    {
                        name: 'email',
                        title: 'Email'
                    }, 
                    {
                        name: '__slot:actions',
                        title: '操作'
                    }]
                }
            }
        },
        methods: {
            editRow (data, index) {
                alert('edit')
            },
            deleteRow (data, index) {
                this.$refs.bkComplexTable.removeRow(data)
            },
            deleteSelectedRow () {
                this.$refs.bkComplexTable.removeSelectedRow()
            },
            renderComplexTableData (page) {
                this.$refs.bkComplexTable.renderCurPageData(page)
            },
            removeRowCallback (data) {
                console.log(data.rowData)
                this.$refs.bkComplexTablePagination.totalPage = data.totalPage
                this.$refs.bkComplexTablePagination.curPage = data.curPage
            }
        }
    }
</script>
```

:::

### 属性
| 参数      | 说明    | 类型      | 可选值       | 默认值   |
|---------- |-------- |---------- |-------------  |-------- |
|data   |表格数据 |Array    |——            | —— |
|cur-page   |显示第几页数据 |Number    |——            | 1 |
|page-size   |一页面展示多少条数据 |Number    |——            | 10 |
|fields  |配置表头展示  |Object    |——            |——  |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| remove-row | 当删除行时触发 | 删除的数据对象、总页数和当前页数 |

