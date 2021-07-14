<div class="p-2 border border-gray-200">
{#if event != null}
    <p class="">key: {event.key}</p>
    <p class="">code: {event.code}</p>
{:else}
    <p class="">(no event)</p>
{/if}
</div>

<script>
let event = null;

export function setEvent(e) {
    event = e;
}
</script>
