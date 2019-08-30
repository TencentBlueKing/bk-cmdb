<template>
    <div class="map">
        <img src="../../../assets/images/map.svg"
            ref="img"
            :style="imgStyles">
        <transition-group name="map-highlight-fade">
            <svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
                viewBox="0 0 857 404"
                v-for="(coordinate, index) in coordinates"
                v-show="randomIndex.includes(index)"
                :style="imgStyles"
                :key="coordinate.id">
                <g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
                    <g transform="translate(-255, -182)">
                        <g transform="translate(255, 182)">
                            <circle :cx="coordinate.x" :cy="coordinate.y" r="1.5" stroke-width="1" fill="#D3D7DE"></circle>
                            <circle :cx="coordinate.x" :cy="coordinate.y" r="0" stroke="#D3D7DE" stroke-width="1.5" fill-opacity="0">
                                <animate attributeName="r" id="ani1" begin="0" from="0" :to="getRandomRadius(coordinate)" dur="4.8s" repeatCount="indefinite"></animate>
                                <animate attributeName="opacity" begin="0" from="0.8" to="0" dur="4.8s" repeatCount="indefinite"></animate>
                            </circle>
                            <circle :cx="coordinate.x" :cy="coordinate.y" r="0" stroke="#D3D7DE" stroke-width="1.5" fill-opacity="0">
                                <animate attributeName="r" begin="ani1.begin + 1.6s" from="0" :to="getRandomRadius(coordinate)" dur="4.8s" repeatCount="indefinite"></animate>
                                <animate attributeName="opacity" begin="ani1.begin + 1.6s" from="0.8" to="0" dur="4.8s" repeatCount="indefinite"></animate>
                            </circle>
                            <circle :cx="coordinate.x" :cy="coordinate.y" r="0" stroke="#D3D7DE" stroke-width="1.5" fill-opacity="0">
                                <animate attributeName="r" begin="ani1.begin + 3.2s" from="0" :to="getRandomRadius(coordinate)" dur="4.8s" repeatCount="indefinite"></animate>
                                <animate attributeName="opacity" begin="ani1.begin + 3.2s" from="0.8" to="0" dur="4.8s" repeatCount="indefinite"></animate>
                            </circle>
                        </g>
                    </g>
                </g>
            </svg>
        </transition-group>
    </div>
</template>

<script>
    import {
        addResizeListener,
        removeResizeListener
    } from '@/utils/resize-events.js'
    export default {
        name: 'cmdb-index-map',
        data () {
            return {
                timer: null,
                randomIndex: [],
                randomRadius: {},
                ratio: {
                    height: 404 / 857,
                    top: 181 / 767
                },
                imgStyles: {
                    width: '857px',
                    height: '404px',
                    visibility: 'hidden'
                },
                resizeHandler: null,
                coordinates: [{
                    id: 'node-1',
                    x: 307.0145,
                    y: 17.4003
                }, {
                    id: 'node-2',
                    x: 474.1875,
                    y: 24.1063
                }, {
                    id: 'node-3',
                    x: 179.9595,
                    y: 45.8113
                }, {
                    id: 'node-4',
                    x: 306.5745,
                    y: 74.7513
                }, {
                    id: 'node-5',
                    x: 23.8035,
                    y: 96.4573
                }, {
                    id: 'node-6',
                    x: 125.0935,
                    y: 125.3963
                }, {
                    id: 'node-7',
                    x: 516.3915,
                    y: 81.9853
                }, {
                    id: 'node-8',
                    x: 626.1225,
                    y: 53.0463
                }, {
                    id: 'node-9',
                    x: 706.3135,
                    y: 89.2213
                }, {
                    id: 'node-10',
                    x: 765.4005,
                    y: 60.2813
                }, {
                    id: 'node-11',
                    x: 782.2825,
                    y: 103.6923
                }, {
                    id: 'node-12',
                    x: 444.6425,
                    y: 161.5713
                }, {
                    id: 'node-13',
                    x: 575.4795,
                    y: 154.3383
                }, {
                    id: 'node-14',
                    x: 693.6515,
                    y: 197.7483
                }, {
                    id: 'node-15',
                    x: 205.2825,
                    y: 219.4543
                }, {
                    id: 'node-16',
                    x: 277.0315,
                    y: 284.5703
                }, {
                    id: 'node-17',
                    x: 217.9435,
                    y: 385.8603
                }, {
                    id: 'node-18',
                    x: 465.7455,
                    y: 270.0983
                }, {
                    id: 'node-19',
                    x: 702.0915,
                    y: 299.0403
                }, {
                    id: 'node-20',
                    x: 765.4005,
                    y: 364.1573
                }]
            }
        },
        computed: {
            randomCoordinates () {
                const uniqueIndex = [...new Set(this.randomIndex)]
                return uniqueIndex.map(index => this.coordinates[index])
            }
        },
        mounted () {
            this.initResizeEvent()
            this.randomHighlight()
        },
        beforeDestroy () {
            clearTimeout(this.timer)
            removeResizeListener(this.$parent.$el, this.resizeHandler)
        },
        methods: {
            initResizeEvent () {
                this.resizeHandler = () => {
                    const parentRect = this.$parent.$el.getBoundingClientRect()
                    const imgStyles = {
                        width: Math.floor(parentRect.width * 0.66) + 'px',
                        height: Math.floor(parentRect.width * 0.66 * this.ratio.height) + 'px',
                        left: Math.floor(parentRect.width * 0.17 + parentRect.left) + 'px'
                    }
                    this.imgStyles = imgStyles
                }
                addResizeListener(this.$parent.$el, this.resizeHandler)
            },
            randomHighlight () {
                const index = Math.floor(Math.random() * this.coordinates.length)
                this.randomIndex.unshift(index)
                this.randomIndex.splice(5)
                this.timer = setTimeout(() => {
                    this.randomHighlight()
                }, 3000)
            },
            getRandomRadius (coordinate) {
                const radius = this.randomRadius[coordinate.id]
                if (!radius) {
                    this.randomRadius[coordinate.id] = Math.floor(Math.random() * 30) + 10
                }
                return this.randomRadius[coordinate.id]
            }
        }
    }
</script>

<style lang="scss" scoped>
    .map {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        z-index: -1;
        pointer-events: none;
        overflow: visible;
        svg, img {
            position: absolute;
            top: 181px;
        }
        img {
            opacity: 0.5745;
        }
    }
    .map-highlight-fade-enter-active,
    .map-highlight-fade-leave-active {
        transition: opacity linear 1s;
    }

    .map-highlight-fade-enter,
    .map-highlight-fade-leave-active {
        opacity: 0;
    }
</style>
