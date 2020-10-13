<template>
    <span class="search-input-wrapper" v-if="multiple">
        <bk-input class="search-input" type="number" v-model="start"></bk-input>
        <span class="search-input-grep">-</span>
        <bk-input class="search-input" type="number" v-model="end"></bk-input>
    </span>
    <bk-input class="search-input" type="number" v-model="localValue" v-else></bk-input>
</template>

<script>
    export default {
        name: 'cmdb-search-float',
        props: {
            value: {
                type: [Number, String, Array],
                default: ''
            }
        },
        computed: {
            multiple () {
                return Array.isArray(this.value)
            },
            localValue: {
                get () {
                    return this.value
                },
                set (value) {
                    let newValue
                    if (Array.isArray(value)) {
                        newValue = value.map(number => number.toString().length ? Number(number) : number)
                    } else {
                        newValue = value.toString().length ? Number(value) : value
                    }
                    this.$emit('input', newValue)
                    this.$emit('change', newValue)
                }
            },
            start: {
                get () {
                    const [start] = this.value
                    return start
                },
                set (start) {
                    const [, end = ''] = this.value
                    this.localValue = [start, end]
                }
            },
            end: {
                get () {
                    const [, end] = this.value
                    return end
                },
                set (end) {
                    const [start = ''] = this.value
                    this.localValue = [start, end]
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .search-input-wrapper {
        display: inline-flex;
        .search-input-grep {
            flex: 20px 0 0;
            text-align: center;
        }
        .search-input {
            flex: 1;
        }
    }
</style>
