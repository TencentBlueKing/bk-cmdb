<template>
    <section class="view-instance">
        <instance-options class="instance-options"></instance-options>
        <instance-list class="instance-list"></instance-list>
    </section>
</template>

<script>
    import InstanceOptions from './options'
    import InstanceList from './list'
    export default {
        name: 'view-instance',
        components: {
            InstanceOptions,
            InstanceList
        },
        data () {
            return {}
        }
    }
</script>
