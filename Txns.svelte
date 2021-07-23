<div class:dim="{widgetstate == 'dim'}" class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="min-width: 40rem;">
{#if account == null || account.txns == null}
    <p class="fg-dim">Select Account</p>
{:else}
    <h1 class="text-sm font-bold mb-2">{account.name}: Transactions</h1>
    {#each account.txns as t, i}
    <a href="/" on:click|preventDefault="{e => onseltxn(t, i)}">
    {#if ui.selid == t.txnid}
            <div class="flex flex-row flex-start p-1 border-b border-cell highlight">
                <p class="highlight cell-date">{t.date.substring(0, 10)}</p>
                <p class="highlight truncate cell-desc">{t.desc}</p>
                {#if t.amt >= 0}
                <p class="highlight text-right cell-amt mr-1">{t.fmtamt}</p>
                <p class="highlight text-right cell-amt mr-1"></p>
                {:else}
                <p class="highlight text-right cell-amt mr-1"></p>
                <p class="highlight text-right cell-amt mr-1">{t.fmtamt}</p>
                {/if}
                <a class="cell-icon pl-1 py-1" href="/" on:click|preventDefault="{e => onseledittxn(t)}">
                    <svg class="w-2 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M12.3 3.7l4 4L4 20H0v-4L12.3 3.7zm1.4-1.4L16 0l4 4-2.3 2.3-4-4z"/></svg>
                </a>
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
                <p class="cell-icon pl-1 py-1"></p>
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

export let widgetstate = "";
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
