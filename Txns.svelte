<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="min-width: 40rem;">
{#if account == null || ui.txns == null}
    <p class="fg-dim">Select Account</p>
{:else}
    <div class="flex flex-row justify-between items-end mb-2">
        <h1 class="text-sm font-bold">Transactions for {account.name}</h1>
        <a class="text-xs pill" href="/" on:click|preventDefault={oncreate}>Create</a>
    </div>
    {#if ui.editid != 0}
    <!-- Don't show filter when Create form is visible. -->
        <div class="mb-2">
            <form autocomplete="off" on:submit|preventDefault="{e => {}}">
                <input class="block bg-input fg-normal py-1 px-2 w-full" name="filter" id="filter" type="text" placeholder="Filter" bind:value={ui.filter} on:keyup="{e => render(account)}">
            </form>
        </div>
    {/if}
    {#if ui.editid == 0}
        <div class="p-2 border-b border-cell">
            <TxnForm txn={ui.newtxn} on:submit={txnform_done} on:cancel={txnform_done} />
        </div>
    {/if}
    {#each ui.txns as t (t.txnid)}
        {#if ui.editid == t.txnid}
        <div class="p-2 border-b border-cell">
            <TxnForm txn={t} on:submit={txnform_done} on:cancel={txnform_done} />
        </div>
        {:else if ui.selid == t.txnid}
        <a href="/" on:click|preventDefault="{e => onedittxn(t)}">
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
        </a>
        {:else}
        <a href="/" on:click|preventDefault="{e => onseltxn(t)}">
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
        </a>
        {/if}
    {/each}
{/if}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";
import * as data from "./data.js";
import TxnForm from "./TxnForm.svelte";

export let account = null;

let svcurl = "/api";
let ui = {};
ui.selid = 0;
ui.editid = -1;
ui.newtxn = {
    txnid: 0,
    accountid: 0,
    date: "",
    ref: "",
    desc: "",
    amt: 0.0,
};

ui.txns = null;
ui.filter = "";

$: render(account);

function render(account) {
    console.log("Txns.svelte render()");
    if (account == null) {
        return;
    }
    processFilter(account, ui.filter);
}

function processFilter(account, sfilter) {
    if (account == null) {
        ui.txns = [];
        return;
    }
    sfilter = sfilter.trim();
    if (sfilter == "") {
        ui.txns = account.txns;
        return;
    }
    ui.txns = filterTxns(account.txns, sfilter);
}
function filterTxns(txns, sfilter) {
    sfilter = sfilter.toLowerCase();
    let tt = [];
    for (let i=0; i < txns.length; i++) {
        let t = txns[i];
        if (t.desc.toLowerCase().includes(sfilter)) {
            tt.push(t);
        }
    }
    return tt;
}

export function reset() {
    ui.selid = 0;
}

export function onEvent(e) {
}

function onseltxn(txn) {
    // If edit form is open, just cancel edit without selecting anything.
    if (ui.editid != -1) {
        ui.editid = -1;
        return;
    }

    ui.selid = txn.txnid;
    dispatch("select", txn);
}
function onedittxn(txn) {
    ui.editid = txn.txnid;
}
function oncreate(e) {
    if (account == null) {
        return;
    }
    ui.newtxn.accountid = account.accountid;
    ui.newtxn.date = new Date().toISOString(),
    ui.editid = 0;
}
function txnform_done(e) {
    ui.editid = -1;
}

</script>
