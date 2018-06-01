<template>
    <span v-if="language === 'en'" class="form-validate-message" v-show="$parent.errors.has(name)">{{$parent.errors.first(name)}}</span>
    <span v-else class="form-validate-message" v-show="$parent.errors.has(name)">内容不合法</span>
</template>
<script>
    import {mapGetters} from 'vuex'
    export default {
        props: {
            value: {
                required: true
            },
            parentClass: {
                type: Boolean,
                default: true
            }
        },
        computed: {
            ...mapGetters([
                'language'
            ]),
            name () {
                return this.$attrs['name'] || this.$attrs['data-vv-name']
            }
        },
        watch: {
            value (val) {
                this.$emit('input', val)
            },
            '$parent.errors.items' () {
                if (this.parentClass) {
                    if (this.$parent.errors.has(this.name)) {
                        this.$el.parentElement.classList.add('form-validate-error')
                    } else {
                        this.$el.parentElement.classList.remove('form-validate-error')
                    }
                }
            }
        }
    }
</script>
<style lang="scss" scoped>
    input{
        display: none;
    }
    .form-validate-message{
        font-size: 12px;
        color: #ff5656;
    }
</style>