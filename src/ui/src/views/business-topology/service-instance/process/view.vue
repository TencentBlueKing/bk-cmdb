<template>
    <section class="view-process">
        <process-options></process-options>
        <process-list class="process-list"></process-list>
    </section>
</template>

<script>
    import ProcessOptions from './options'
    import ProcessList from './list'
    export default {
        name: 'view-process',
        components: {
            ProcessOptions,
            ProcessList
        },
        data () {
            return {}
        }
    }
</script>
