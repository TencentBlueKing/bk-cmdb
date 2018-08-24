<template>
    <div class='bk-circle'>
        <svg :width="radius" :height="radius" viewBox="0 0 100 100" version="1.1">
            <circle 
                class="progress-background" 
                cx="50" 
                cy="50" 
                r="50" 
                fill="transparent" 
                :stroke-width='config.strokeWidth'
                :stroke='config.bgColor' />
            <circle 
                class="progress-bar" 
                cx="50" 
                cy="50" 
                r="50" 
                fill="transparent" 
                :class="'circle' + config.index"
                :stroke-width='config.strokeWidth' 
                :stroke='config.activeColor' 
                :stroke-dasharray="dashArray" 
                :stroke-dashoffset="dashOffset" />
          </svg>
        <div class="num" v-if="numShow" :style="numStyle">
            {{Math.round(percentFixed * 100)}}% 
        </div>
        <div class="title" v-if="title" :style="titleStyle">
            {{title}}
        </div>
    </div>
</template>
<script>
    export default {
        name: 'bk-round',
        props: {
            config: {
                type: Object,
                default () {
                    return {
                        strokeWidth: 5,
                        bgColor: 'gray',
                        activeColor: 'green',
                        index: 0
                    }
                }
            },
            percent: {
                type: Number,
                default: 0
            },
            title: {
                type: String
            },
            titleStyle: {
                type: Object,
                default () {
                    return {
                        fontSize: '16px'
                    }
                }
            },
            numShow: {
                type: Boolean,
                default: true
            },
            numStyle: {
                type: Object,
                default () {
                    return {
                        fontSize: '16px'
                    }
                }
            },
            radius: {
                type: String,
                default: '100px'
            }
        },
        computed: {
            dashOffset () {
                return this.percentFixed > 1 ? false : (1 - this.percentFixed) * this.dashArray
            },
            percentFixed () {
                return Number(this.percent.toFixed(2))
            }
        },
        data () {
            return {
                dashArray: Math.PI * 100
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/round-progress.scss';
</style>