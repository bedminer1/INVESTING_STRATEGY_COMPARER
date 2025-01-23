<script lang="ts">
    import { onMount } from 'svelte';
    import { Chart, LineController, LineElement, PointElement, LinearScale, Title, CategoryScale, Tooltip, Legend } from 'chart.js';

    // Register necessary Chart.js components
    Chart.register(LineController, LineElement, PointElement, LinearScale, Title, CategoryScale, Tooltip, Legend);

    let { stats, label, xAxisLabels }: { stats: DataSet[], label: string, xAxisLabels: string[]} = $props()

    let chart: Chart | null = null
    Chart.defaults.color = 'rgb(250,255,255)'
    Chart.defaults.font.size = 14
    let chartCanvas: HTMLCanvasElement


    onMount(() => {
        if (chart) { // update
            chart.destroy();
        }

        chart = new Chart(chartCanvas, {
            type: 'line',
            data: {
                labels: xAxisLabels,
                datasets: stats.map(stat => ({
                    label: stat.label + " " + label,
                    data: stat.data,
                    borderColor: stat.borderColor || 'rgba(75, 192, 192, 1)',
                    backgroundColor: stat.backgroundColor || 'rgba(75, 192, 192, 0.2)',
                    fill: true,
                }))
            },
            options: {
                responsive: true,
                scales: {
                    y: {
                        title: {
                            display: true,
                            text: "Value (USD)",
                        }
                    },
                    x: {
                        title: {
                            display: true,
                            text: 'Months'
                        }
                    }
                },
            }
        });

        return () => {
            chart?.destroy();
        };
    });
</script>

<canvas class="w-full" bind:this={chartCanvas}></canvas>