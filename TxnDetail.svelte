<div class="txndetail webpanel">
{#if txn == null}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flexrow justify-between mb-2">
        <h1 class="h2">Transaction Details</h1>
    </div>

    <p class="p-1 border-b border-cell">{txn.date.substring(0, 10)}</p>
    <p class="p-1 border-b border-cell">{txn.desc}</p>
    {#if txn.amt > 0}
        <p class="p-1 border-b border-cell whitespace-nowrap fg-number-plus">{txn.fmtamt}</p>
    {:else}
        <p class="p-1 border-b border-cell whitespace-nowrap fg-number-minus">{txn.fmtamt}</p>
    {/if}
    {#if txn.ref != ""}
        <p class="p-1 border-b border-cell">{txn.ref}</p>
    {/if}
    {#if txn.memo != ""}
    <div class="memo p-1 border-b border-cell">
        {@html textToHtml(txn.memo)}
    </div>
    {/if}
{/if}
</div>

<script>
import {textToHtml} from "./helpers.js";
export let txn = null;

</script>

