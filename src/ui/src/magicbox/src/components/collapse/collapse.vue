<template>
    <div class="bk-collapse">
        <slot></slot>
    </div>
</template>

<script>
export default {
    name: 'bkCollapse',
    props: {
        accordion: {
            type: Boolean,
            default: false
        },
        value: {
            type: [Array, String]
        }
    },
    watch: {
        value (val) {
            this.currentValue = val
        },
        currentValue () {
            this.setActive()
        }
    },
    data () {
        return {
            currentValue: this.value
        }
    },
    
    methods: {
        /**
         * @description 设置默认展开
         */
        setActive () {
            const activeKey = this.getActiveKey()
            this.$children.forEach((child, index) => {
                const name = child.name || index.toString()
                child.isActive = activeKey.indexOf(name) > -1
                child.index = index
            })
        },
        /**
         * @description 获取处于激活态的
         */
        getActiveKey () {
            let activeKey = this.currentValue || []
            const accordion = this.accordion
            if (!Array.isArray(activeKey)) {
                activeKey = [activeKey]
            }
            if (accordion && activeKey.length > 1) {
                activeKey = [activeKey[0]]
            }
            for (let i = 0; i < activeKey.length; i++) {
                activeKey[i] = activeKey[i].toString()
            }
            return activeKey
        },
        toggle (data) {
            const name = data.name.toString()
            let newActiveKey = []
            if (this.accordion) {
                if (!data.isActive) {
                    newActiveKey.push(name)
                }
            } else {
                let activeKey = this.getActiveKey()
                const nameIndex = activeKey.indexOf(name)
                if (data.isActive) {
                    if (nameIndex > -1) {
                        activeKey.splice(nameIndex, 1)
                    }
                } else {
                    if (nameIndex < 0) {
                        activeKey.push(name)
                    }
                }
                newActiveKey = activeKey
            }
            this.currentValue = newActiveKey
            this.$emit('input', newActiveKey)
            this.$emit('item-click', newActiveKey)
        }
    },
    mounted () {
        this.setActive()
    }
}
</script>
