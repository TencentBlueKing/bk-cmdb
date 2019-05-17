<template>
    <section :class="{'active': isActive, 'bk-tab2-pane': !show || !isActive}">
        <slot></slot>
    </section>
</template>
<script>
export default {
    name: 'bk-tabpanel',
    props: {
        name: {
            type: String,
            default: ''
        },
        title: {
            type: String,
            default: ''
        },
        show: {
            type: Boolean,
            default: true
        },
        tag: {
        }
    },
    data () {
        return {
            index: -1
        }
    },
    computed: {
        isActive () {
            return this.$parent.calcActiveName === this.name
        },
        panelParams () {
            return {
                title: this.title,
                name: this.name,
                show: this.show,
                tag: this.tag
            }
        }
    },
    watch: {
        panelParams: {
            deep: true,
            handler (val) {
                this.$parent.updateList(this.index, val)
            }
        }
    },
    mounted () {
        let panelParams = JSON.parse(JSON.stringify(this.panelParams))
        let pannelSlots = this.$slots
        let parentSlots = this.$parent.$slots
        for (let slot in pannelSlots) {
            if (slot === 'tag') {
                let vnode = pannelSlots['tag']
                panelParams.vnode = vnode
                panelParams.label = (h) => {
                    return h('div', {
                        class: ['bk-panel-label']
                    }, vnode)
                }
            }
        }
        this.index = this.$parent.addNavItem(panelParams)
    }
}
</script>
