<template>
    <div class="bk-tag-selector" @click="foucusInputer($event)">
		<div :class="['bk-tag-input', {'active': isEdit, 'disabled': disabled}]">
			<ul class="tag-list" ref="tagList">
	            <li class="key-node" v-for="(tag, index) in localTagList" :key="index" @click.stop="selectTag($event, tag)">
	                <span class="tag">{{tag[displayKey]}}</span>
	                <a href="javascript:void(0)" class="remove-key" @click.stop="removeTag($event, tag, index)" v-if="!disabled && hasDeleteIcon">
	                    <i class="bk-icon icon-close"></i>
	                </a>
	            </li>
	            <li ref="staffInput" id="staffInput">
	                <input
	                    type="text"
	                    class="input"
                        ref="input"
	                    v-model="curInputValue"
	                    v-if="!disabled"
						@input="input"
                        @focus="focusInput"
                        @paste="paste"
	                    @blur="blurHandler"
	                    @keydown="keyupHandler">
	            </li>
	        </ul>
	        <p class="placeholder" v-show="!isEdit && !localTagList.length && !curInputValue.length">{{placeholder}}</p>
		</div>
		<transition name="optionList">
			<div class="bk-selector-list" v-show="showList && renderList.length">
	            <ul ref="selectorList" :style="{'max-height': `${contentMaxHeight}px`}" class="outside-ul">
					<li v-for="(data, index) in renderList"
	                    class="bk-selector-list-item"
						:class="activeClass(index)"
	                    :key="index"
                        @click="setValTab(data, 'select')">
                        <Render :node="data" :displayKey="displayKey" :tpl='tpl' />
	                </li>
	            </ul>
	        </div>
		</transition>
    </div>
</template>

<script>
    import Render from './render'

    export default {
        name: 'bk-tag-input',
        components: { Render },
        props: {
            placeholder: {
                type: String,
                default: '请输入并按Enter结束'
            },
            value: {
                type: Array,
                default () {
                    return []
                }
            },
            // tags: {
            //     type: Array,
            //     default () {
            //         return []
            //     }
            // },
            disabled: {
                type: Boolean,
                default: false
            },
            hasDeleteIcon: {
                type: Boolean,
                default: false
            },
            separator: {
                type: String,
                default: ''
            },
            maxData: {
                type: Number,
                default: -1
            },
            maxResult: {
                type: Number,
                default: 5
            },
            isBlurTrigger: {
                type: Boolean,
                default: true
            },
            saveKey: {
                type: String,
                default: 'id'
            },
            displayKey: {
                type: String,
                default: 'name'
            },
            searchKey: {
                type: String,
                default: 'name'
            },
            list: {
                type: Array,
                default: []
            },
            contentMaxHeight: {
                type: Number,
                default: 300
            },
            allowCreate: {
                type: Boolean,
                default: false
            },
            tpl: Function,
            pasteFn: Function
        },
        data () {
            return {
                curInputValue: '',
                cacheVal: '',
                timer: 0,
                focusList: this.allowCreate ? -1 : 0, // 列表选中项
                isEdit: false,
                showList: false,
                isCanRemoveTag: false,
                tagList: [],
                localTagList: [],
                renderList: [],
                initData: []
            }
        },
        created () {
            this.getData()
        },
        watch: {
            curInputValue (newVal, oldVal) {
                if (newVal !== '' && this.renderList.length) {
                    this.showList = true
                } else {
                    setTimeout(() => { this.showList = false }, 100)
                }
            },
            showList (val) {
                if (val) {
                    this.$nextTick(() => {
                        this.$refs.selectorList.scrollTop = 0
                    })
                }
            },
            list (val) {
                if (val) {
                    this.getData()
                }
            }
        },
        methods: {
            // 获取字符长度，汉字两个字节
            getCharLength (str) {
                let len = str.length
                let bitLen = 0

                for (let i = 0; i < len; i++) {
                    if ((str.charCodeAt(i) & 0xff00) !== 0) {
                        bitLen++
                    }
                    bitLen++
                }

                return bitLen
            },
            // 过滤数据
            filterData (val) {
                this.renderList = [...this.initData.filter(item => item[this.searchKey].indexOf(val) > -1)]
                if (this.renderList.length > this.maxResult) {
                    this.renderList = [...this.renderList.slice(0, this.maxResult)]
                }
            },
            // 更新样式
            activeClass (i) {
                return {
                    'bk-selector-selected': i === this.focusList
                }
            },
            // 获取元素位置
            getSiteInfo () {
                let res = {
                    index: 0,
                    temp: []
                }
                let nodes = this.$refs.tagList.childNodes

                for (let i = 0; i < nodes.length; i++) {
                    let node = nodes[i]

                    if (!(node.nodeType === 3 && !/\S/.test(node.nodeValue))) {
                        res.temp.push(node)
                    }
                }

                Object.keys(res.temp).forEach(key => {
                    if (res.temp[key].id === 'staffInput') res.index = key
                })

                return res
            },
            getData () {
                this.initData = [...this.list]
                
                if (this.value.length) {
                    this.value.map(tag => {
                        this.initData.filter(val => {
                            if (tag === val[this.saveKey]) {
                                this.localTagList.push(val)
                                this.tagList.push(val[this.saveKey])
                            } else if (this.allowCreate && !tag.includes(tag)) {
                                let temp = {}
                                
                                temp[this.saveKey] = tag
                                temp[this.displayKey] = tag
                                this.localTagList.push(temp)
                                this.tagList.push(tag)
                            }
                        })
                    })

                    this.value.forEach(tag => {
                        this.initData = this.initData.filter(val => !tag.includes(val[this.saveKey]))
                    })
                }
            },
            selectTag (event, tag) {
                if (this.disabled) return

                let domLen = event.target.parentNode.offsetWidth
                let resSite = (event.target.parentNode.offsetWidth - 20) / 2

                if (event.offsetX > resSite) {
                    this.insertAfter(this.$refs.staffInput, event.target.parentNode)
                } else {
                    this.$refs.tagList.insertBefore(this.$refs.staffInput, event.target.parentNode)
                }

                this.foucusInputer(event)
                this.$refs.input.focus()
                this.$refs.input.style.width = 12 + 'px'
            },
            input (event) {
                if (this.maxData === -1 || this.maxData > this.tagList.length) {
                    let { value } = event.target
                    let charLen = this.getCharLength(value)

                    this.cacheVal = value
                    if (charLen) {
                        this.isCanRemoveTag = false
                        this.filterData(value)
                        this.$refs.input.style.width = (charLen * 12) + 'px'
                    } else {
                        this.isCanRemoveTag = true
                    }
                } else {
                    this.blurHandler()
                    this.curInputValue = ''
                    this.showList = false
                }
                
                this.isEdit = true
                // 重置下拉菜单选中信息
                this.focusList = this.allowCreate ? -1 : 0
            },
            focusInput () {
                this.isCanRemoveTag = true
            },
            paste (event) {
                event.preventDefault()

                let value = event.clipboardData.getData('text')
                let valArr = this.pasteFn ? this.pasteFn(value) : this.defaultPasteFn(value)
                let tags = []

                valArr.map(val => tags.push(val[this.saveKey]))
                
                if (tags.length) {
                    let nodes = this.$refs.tagList.childNodes
                    let result = this.getSiteInfo(nodes)
                    let localTags = []
                    let localInitDara = []
                    
                    this.initData.map(data => {
                        localInitDara.push(data[this.saveKey])
                    })
                    tags = tags.filter(tag => { return tag && tag.trim() && !this.tagList.includes(tag) && localInitDara.includes(tag) })

                    if (this.maxData !== -1) {
                        if (this.tagList.length < this.maxData) {
                            let differ = this.maxData - this.tagList.length
                            if (tags.length > differ) {
                                tags = [...tags.slice(0, (differ))]
                            }
                        } else {
                            tags = []
                        }
                    }
                    
                    tags.map(tag => {
                        let temp = {}
                        temp[this.saveKey] = tag
                        temp[this.displayKey] = tag
                        localTags.push(temp)
                    })
                    
                    if (tags.length) {
                        this.tagList = [...this.tagList.slice(0, result.index), ...tags, ...this.tagList.slice(result.index, this.tagList.length)]
                        this.localTagList = [...this.localTagList.slice(0, result.index), ...localTags, ...this.localTagList.slice(result.index, this.localTagList.length)]

                        let site = nodes[parseInt(result.index) + 1]
                        this.insertAfter(this.$refs.staffInput, site)
                        this.$refs.input.focus()
                        this.$refs.input.style.width = 12 + 'px'
                        tags.map(tag => {
                            this.initData = this.initData.filter(val => !tag.includes(val[this.saveKey]))
                        })
              
                        this.handlerChange('select')
                    }
                }
            },
            defaultPasteFn (val) {
                let target = []
                let textArr = val.split(';')
                
                textArr.map(item => {
                    if (item.match(/^[a-zA-Z][a-zA-Z_]+/g)) {
                        let finalItem = item.match(/^[a-zA-Z][a-zA-Z_]+/g).join('')
                        let temp = {}
                        temp[this.saveKey] = finalItem
                        temp[this.displayKey] = finalItem
                        target.push(temp)
                    }
                })
                return target
            },
            updateScrollTop () {
                // 获取下拉列表容器的位置信息，用于判断上下键选中的元素是否在可视区域，若不在则需滚动至可视区域
                const panelObj = this.$el.querySelector('.bk-selector-list .outside-ul')
                let panelInfo = {
                    height: panelObj.clientHeight,
                    yAxios: panelObj.getBoundingClientRect().y
                }

                this.$nextTick(() => {
                    const activeObj = this.$el.querySelector('.bk-selector-list .bk-selector-selected')
                    let activeInfo = {
                        height: activeObj.clientHeight,
                        yAxios: activeObj.getBoundingClientRect().y
                    }

                    // 选中元素顶部坐标大于容器顶部坐标时，则该元素有部分或者全部区域不在可视区域，执行滚动
                    if (activeInfo.yAxios < panelInfo.yAxios) {
                        let currentScTop = panelObj.scrollTop
                        panelObj.scrollTop = currentScTop - (panelInfo.yAxios - activeInfo.yAxios)
                    }

                    let distanceToBottom = activeInfo.yAxios + activeInfo.height - panelInfo.yAxios

                    // 选中元素底部坐标大于容器顶部坐标，且超出容器的实际高度，则该元素有部分或者全部区域不在可视区域，执行滚动
                    if (distanceToBottom > panelInfo.height) {
                        let currentScTop = panelObj.scrollTop
                        panelObj.scrollTop = currentScTop + distanceToBottom - panelInfo.height
                    }
                })
            },
            keyupHandler (event) {
                let target
                let val = event.target.value
                let valLen = this.getCharLength(val)
                let result = this.getSiteInfo()
                let nodes = this.$refs.tagList.childNodes

                switch (event.code) {
                    case 'ArrowUp':
                        // 上
                        event.preventDefault()
                        this.focusList--
                        this.focusList = this.focusList < 0 ? -1 : this.focusList
                        if (this.focusList === -1) {
                            this.focusList = this.renderList.length - 1
                        }
                        this.updateScrollTop()
                        break
                    case 'ArrowDown':
                        // 下
                        event.preventDefault()
                        this.focusList++
                        this.focusList = this.focusList > this.renderList.length - 1 ? this.renderList.length : this.focusList
                        if (this.focusList === this.renderList.length) {
                            this.focusList = 0
                        }
                        this.updateScrollTop()
                        break
                    case 'ArrowLeft':
                        this.isEdit = true
                        if (!valLen) {
                            if (parseInt(result.index) > 1) {
                                let leftsite = nodes[parseInt(result.index) - 2]
                                this.insertAfter(this.$refs.staffInput, leftsite)
                                this.$refs.input.value = ''
                                this.$refs.input.style.width = 12 + 'px'
                            } else {
                                let nodes = this.$refs.tagList.childNodes
                                this.$refs.tagList.insertBefore(this.$refs.staffInput, nodes[0])
                            }
                            this.$refs.input.focus()
                        }
                        break
                    case 'ArrowRight':
                        this.isEdit = true
                        if (!valLen) {
                            let rightsite = nodes[parseInt(result.index) + 1]
                            this.insertAfter(this.$refs.staffInput, rightsite)
                            this.$refs.input.focus()
                        }
                        break
                    case 'Enter':
                    case 'NumpadEnter':
                        if ((!this.allowCreate && this.showList) || (this.allowCreate && this.focusList >= 0 && this.showList)) {
                            this.setValTab(this.renderList[this.focusList], 'select')
                            this.showList = false
                        } else if (this.allowCreate) {
                            let tag = this.curInputValue
                            this.setValTab(tag, 'custom')
                        }
                        this.cacheVal = ''
                        break
                    case 'Backspace':
                        if (parseInt(result.index) !== 0 && !this.curInputValue.length) {
                            target = this.localTagList[result.index - 1]
                            this.backspaceHandler(result.index, target)
                        }
                        break
                    default:
                        break
                }
            },
            // 选中标签
            setValTab (item, type) {
                let nodes = this.$refs.tagList.childNodes
                let result = this.getSiteInfo(nodes)
                let isSelected = false
                let tags = []
                let newVal
                
                if (type === 'custom') {
                    if (this.separator) {
                        let localTags = []

                        tags = item.split(this.separator)
                        tags = tags.filter(tag => { return tag && tag.trim() && !this.tagList.includes(tag) })
                        tags = [...new Set(tags)]
                        tags.map(tag => {
                            let temp = {}
                            temp[this.saveKey] = tag
                            temp[this.displayKey] = tag
                            localTags.push(temp)
                        })
                        
                        if (tags.length) {
                            this.tagList = [...this.tagList.slice(0, result.index), ...tags, ...this.tagList.slice(result.index, this.tagList.length)]
                            this.localTagList = [...this.localTagList.slice(0, result.index), ...localTags, ...this.localTagList.slice(result.index, this.localTagList.length)]
                            isSelected = true
                        }
                    } else {
                        if (typeof item === 'object') {
                            newVal = item[this.saveKey]
                            if (newVal && !this.tagList.includes(newVal)) {
                                newVal = newVal.replace(/\s+/g, '')
                                
                                if (newVal.length) {
                                    this.tagList = [...this.tagList.slice(0, result.index), newVal, ...this.tagList.slice(result.index, this.tagList.length)]
                                    this.localTagList = [...this.localTagList.slice(0, result.index), item, ...this.localTagList.slice(result.index, this.localTagList.length)]
                                    isSelected = true
                                }
                            }
                        } else {
                            let localItem = {}
                            newVal = item.trim()
                            localItem[this.saveKey] = newVal
                            localItem[this.displayKey] = newVal

                            if (newVal.length && !this.tagList.includes(newVal)) {
                                this.tagList = [...this.tagList.slice(0, result.index), newVal, ...this.tagList.slice(result.index, this.tagList.length)]
                                this.localTagList = [...this.localTagList.slice(0, result.index), localItem, ...this.localTagList.slice(result.index, this.localTagList.length)]
                                isSelected = true
                            }
                        }
                    }
                } else {
                    newVal = item[this.saveKey]
                    this.tagList = [...this.tagList.slice(0, result.index), newVal, ...this.tagList.slice(result.index, this.tagList.length)]
                    this.localTagList = [...this.localTagList.slice(0, result.index), item, ...this.localTagList.slice(result.index, this.localTagList.length)]
                    isSelected = true
                }

                if (isSelected) {
                    let site = nodes[parseInt(result.index) + 1]
                    this.insertAfter(this.$refs.staffInput, site)
                    this.$refs.input.focus()
                    this.$refs.input.style.width = 12 + 'px'
                    if (this.allowCreate && this.separator) {
                        tags.map(tag => {
                            this.initData = this.initData.filter(val => !tag.includes(val[this.saveKey]))
                        })
                    } else {
                        this.initData = this.initData.filter(val => !newVal.includes(val[this.saveKey]))
                    }
                }
               
                this.handlerChange('select')
                this.clearInput()
            },
            // 输入清除
            backspaceHandler (index, target) {
                // 如果清空输入
                if (!this.curInputValue) {
                    if (this.isCanRemoveTag) {
                        this.tagList = [...this.tagList.slice(0, index - 1), ...this.tagList.slice(index, this.tagList.length)]
                        this.localTagList = [...this.localTagList.slice(0, index - 1), ...this.localTagList.slice(index, this.localTagList.length)]

                        let nodes = this.$refs.tagList.childNodes
                        let result = this.getSiteInfo(nodes)
                        let key = parseInt(result.index) === 1 ? 1 : parseInt(result.index) - 2
                        let site = nodes[key]

                        this.insertAfter(this.$refs.staffInput, site)
                        this.$refs.input.focus()
                        let isExistInit = this.list.some(item => {
                            return item === target[this.saveKey]
                        })
                        if ((this.allowCreate && isExistInit) || !this.allowCreate) {
                            this.initData.push(target)
                        }

                        this.$refs.input.style.width = 12 + 'px'
                        this.handlerChange('remove')
                    }
                    this.isCanRemoveTag = true
                }
            },
            // 删除标签
            removeTag (event, data, index) {
                this.tagList.splice(index, 1)
                this.localTagList.splice(index, 1)

                let tags = []
                let isExistInit = this.list.some(item => {
                    return item === data[this.saveKey]
                })

                if ((this.allowCreate && isExistInit) || !this.allowCreate) {
                    this.initData.push(data)
                }

                this.$refs.input.style.width = 12 + 'px'
                this.resetInput()
                this.handlerChange('remove')
            },
            handlerChange (type) {
                this.$emit('input', this.tagList)
                this.$emit('change', this.tagList)
                this.$emit(type)
                this.$emit('update:tags', this.tagList)
            },
            // 清空输入框
            clearInput () {
                this.curInputValue = ''
            },
            blurHandler () {
                this.resetInput()
                this.timer = setTimeout(() => {
                    this.clearInput()
                    this.isEdit = false
                }, 300)
            },
            // 输入框获取焦点时触发
            foucusInputer (event) {
                if (this.disabled) return
                
                if (event.target.className === 'bk-tag-input active' || event.target.className === 'tag-list') {
                    setTimeout(() => {
                        this.curInputValue = this.cacheVal
                    }, 100)
                } else {
                    this.cacheVal = ''
                }

                clearTimeout(this.timer)
                this.isEdit = true
                this.$nextTick(() => {
                    this.$el.querySelector('.input').focus()
                })
            },
            // 改变元素位置
            insertAfter (newElement, targetElement) {
                let parent = targetElement.parentNode

                if (parent.lastChild === targetElement) {
                    parent.appendChild(newElement)
                } else {
                    parent.insertBefore(newElement, targetElement.nextSibling)
                }
            },
            // 重置input框位置
            resetInput () {
                let nodes = this.$refs.tagList.childNodes
                let result = this.getSiteInfo(nodes)

                if (result.index !== result.temp.length) {
                    this.clearInput()
                    let site = nodes[nodes.length - 1]
                    
                    this.insertAfter(this.$refs.staffInput, site)
                }
            }
        }
    }
</script>

<style lang="scss">
    @import '../../bk-magic-ui/src/tag-input.scss'
</style>
