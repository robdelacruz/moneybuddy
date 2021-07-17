<div class:dim="{widgetstate == 'dim'}" class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="width: 40rem;">
{#if account == null || account.transactions == null}
    <p class="fg-dim">Select Account</p>
{:else}
    <h1 class="text-sm font-bold mb-2">{account.name}: Transactions</h1>
    {#each account.transactions as t, i}
    <a href="/" on:click|preventDefault="{e => onseltransaction(t, i)}">
        {#if ui.isel == i}
            <div class="flex flex-row flex-start p-1 border-b border-cell highlight">
                <p class="highlight cell-date">{t.date.substring(0, 10)}</p>
                <p class="highlight truncate cell-desc">{t.desc}</p>
                {#if t.amt >= 0}
                <p class="highlight text-right cell-amt">{t.fmtamt}</p>
                <p class="highlight text-right cell-amt"></p>
                {:else}
                <p class="highlight text-right cell-amt"></p>
                <p class="highlight text-right cell-amt">{t.fmtamt}</p>
                {/if}
            </div>
        {:else}
            <div class="flex flex-row flex-start p-1 border-b border-cell">
                <p class="fg-dim cell-date">{t.date.substring(0, 10)}</p>
                <p class="truncate cell-desc">{t.desc}</p>
                {#if t.amt >= 0}
                <p class="fg-dim text-right cell-amt">{t.fmtamt}</p>
                <p class="fg-dim text-right cell-amt"></p>
                {:else}
                <p class="fg-dim text-right cell-amt"></p>
                <p class="fg-dim text-right cell-amt">{t.fmtamt}</p>
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

export let widgetstate = "";
export let account = null;

let svcurl = "/api";
let ui = {};
ui.isel = -1;

function dispatchAction(action, entryid) {
    dispatch("action", {
        action: action,
        itemid: entryid,
    });
}

export function reset() {
    ui.isel = -1;
}

export function onEvent(e) {
    if (ui.account == null || ui.account.transactions.length == 0) {
        return;
    }

    if (e.key == "ArrowUp") {
        ui.isel--;
    } else if (e.key == "ArrowDown") {
        ui.isel++;
    }

    if (ui.isel < 0) {
        ui.isel = 0;
    }
    if (ui.isel > account.transactions.length-1) {
        ui.isel = account.transactions.length-1;
    }
}

function onseltransaction(transaction, i) {
    ui.isel = i;
    dispatch("select", transaction);
}

</script>
