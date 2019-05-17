<template>
    <div>
        <!--div id="tester" style="width: 400px;height: 400px"></div-->
        <div style="width: 100%; height: 400px">
            <div id="myDiv" style="width: 50%;height: 400px;float:left;"></div>
            <div id="barDiv" style="width: 50%;height: 400px;float:right;"></div>
        </div>
        <div id="muDiv" style="width: 400px;width: 400px"></div>
    </div>
</template>

<script>
    import Plotly from 'plotly.js'

    export default {
        name: 'index',
        data () {
            return {
                businessList: []
            }
        },
        mounted () {
            // this.charts()
            this.pieCharts()
            this.barCharts()
            this.mutipChart()
        },
        methods: {
            charts () {
                const TESTER = document.getElementById('tester')
                const trace1 = {
                    x: [1, 2, 3, 4],
                    y: [10, 15, 13, 17],
                    type: 'scatter'
                }

                const trace2 = {
                    x: [1, 2, 3, 4],
                    y: [16, 5, 11, 9],
                    type: 'scatter'
                }

                const data = [trace1, trace2]

                Plotly.newPlot(TESTER, data)
            },
            pieCharts () {
                const myDiv = document.getElementById('myDiv')
                const data = [{
                    values: [19, 26, 55],
                    labels: ['Residential', 'Non-Residential', 'Utility'],
                    type: 'pie'
                }]

                const layout = {
                    height: 400,
                    width: 500,
                    title: '系统类型统计'
                }

                Plotly.newPlot(myDiv, data, layout)
            },
            barCharts () {
                const barDiv = document.getElementById('barDiv')
                const data = [
                    {
                        x: ['giraffes', 'orangutans', 'monkeys'],
                        y: [20, 14, 23],
                        name: '模型',
                        type: 'bar'
                    }
                ]
                const layout = { title: '业务下主机统计' }

                Plotly.newPlot(barDiv, data, layout)
            },
            mutipChart () {
                const muDiv = document.getElementById('muDiv')
                const trace1 = {
                    x: ['giraffes', 'orangutans', 'monkeys'],
                    y: [20, 14, 23],
                    name: 'SF Zoo',
                    type: 'bar'
                }

                const trace2 = {
                    x: ['giraffes', 'orangutans', 'monkeys'],
                    y: [12, 18, 29],
                    name: 'LA Zoo',
                    type: 'bar'
                }

                const data = [trace1, trace2]

                const layout = { barmode: 'stack', title: '模型统计' }

                Plotly.newPlot(muDiv, data, layout)
            }
        }
    }
</script>

<style scoped>

</style>
