<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="min-width: 40rem;">
{#if account == null || account.txns == null}
    <p class="fg-dim">Select Account</p>
{:else}
    <h1 class="text-sm font-bold mb-2">{account.name}: Transactions</h1>
    {#each account.txns as t, i}
    <a href="/" on:click|preventDefault="{e => onseltxn(t, i)}">
    {#if ui.selid == t.txnid}
            <div class="flex flex-row flex-start p-1 border-b border-cell highlight-1">
                <p class="cell-date">{t.date.substring(0, 10)}</p>
                <p class="truncate cell-desc">{t.desc}</p>
                {#if t.amt >= 0}
                <p class="text-right cell-amt mr-1">{t.fmtamt}</p>
                <p class="text-right cell-amt mr-1"></p>
                {:else}
                <p class="text-right cell-amt mr-1"></p>
                <p class="text-right cell-amt mr-1">{t.fmtamt}</p>
                {/if}
            </div>
        {:else}
            <div class="flex flex-row flex-start p-1 border-b border-cell">
                <p class="fg-dim cell-date">{t.date.substring(0, 10)}</p>
                <p class="truncate cell-desc">{t.desc}</p>
                {#if t.amt >= 0}
                <p class="fg-dim text-right cell-amt mr-1">{t.fmtamt}</p>
                <p class="fg-dim text-right cell-amt mr-1"></p>
                {:else}
                <p class="fg-dim text-right cell-amt mr-1"></p>
                <p class="fg-dim text-right cell-amt mr-1">{t.fmtamt}</p>
                {/if}
            </div>
        {/if}
    </a>
    {/each}
{/if}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";

export let account = null;

let svcurl = "/api";
let ui = {};
ui.selid = 0;

export function reset() {
    ui.selid = 0;
}

export function onEvent(e) {
}

function onseltxn(txn) {
    ui.selid = txn.txnid;
    dispatch("select", txn);
}

function onseledittxn(txn) {
    console.log("onseledittxn");
}

</script>
