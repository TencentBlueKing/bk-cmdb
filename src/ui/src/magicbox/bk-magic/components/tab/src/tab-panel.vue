<template>
    <section
        :class="{
            'active': isActive,
            'bk-tab2-pane': !show || !isActive
        }">
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
        }
    },
    mounted () {
        this.index = this.$parent.addNavItem({
            title: this.title,
            name: this.name,
            show: this.show
        })
    },
    updated () {
        this.$parent.updateList(this.index, {
            title: this.title,
            name: this.name,
            show: this.show
        })
    }
}
</script>
