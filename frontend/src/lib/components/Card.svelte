<script lang="ts">
    let { title, body, subtitle, icon }: { title: string, body: string, subtitle: string, icon: string } = $props()
</script>

<div class="card p-4 h-36 w-48">
    <div class="flex w-full justify-between mb-5">
        <p>{ title }</p> 
        <p class="text-3xl">{@html icon }</p>
    </div>
    <p class="text-2xl">{ body }</p>
    <p class="italic text-gray-400 inline-block overflow-hidden whitespace-nowrap">{ subtitle }</p>
</div>