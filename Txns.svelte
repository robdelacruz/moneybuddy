<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="min-width: 40rem;">
{#if account == null}
    <p class="fg-dim">Select Account</p>
{:else}
    <div class="flex flex-row justify-between items-end mb-2">
        <h1 class="text-sm font-bold">Transactions for {account.name}</h1>
        <a class="text-xs pill" href="/" on:click|preventDefault={oncreate}>Create</a>
    </div>
    {#if editid != 0}
    <!-- Don't show filter when Create form is visible. -->
        <div class="mb-2">
            <form autocomplete="off" on:submit|preventDefault="{e => {}}">
                <input class="block bg-input fg-normal py-1 px-2 w-full" name="filter" id="txnfilter" type="text" placeholder="Filter" bind:value={frm_filter}>
            </form>
        </div>
    {/if}
    {#if editid == 0}
        <div class="p-2 border-b border-cell">
            <TxnForm txn={newtxn} on:submit={txnform_done} on:cancel={txnform_done} />
        </div>
    {/if}
    {#each displaytxns as t (t.txnid)}
        {#if editid == t.txnid}
        <div class="p-2 border-b border-cell">
            <TxnForm txn={t} on:submit={txnform_done} on:cancel={txnform_done} />
        </div>
        {:else if selid == t.txnid}
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
let selid = 0;
let editid = -1;
let newtxn = {
    txnid: 0,
    accountid: 0,
    date: "",
    ref: "",
    desc: "",
    amt: 0.0,
};

let frm_filter = "";
let displaytxns = [];

// account + frm_filter --> displaytxns

$: displaytxns = filterTxns(account, frm_filter);

function filterTxns(account, sfilter) {
    console.log("Txns.svelte filterTxns()");
    if (account == null) {
        return [];
    }
    if (account.txns == null) {
        return [];
    }
    sfilter = sfilter.trim().toLowerCase();
    if (sfilter == "") {
        return account.txns;
    }
    let tt = [];
    for (let i=0; i < account.txns.length; i++) {
        let t = account.txns[i];
        if (t.desc.toLowerCase().includes(sfilter)) {
            tt.push(t);
        }
    }
    return tt;
}

export function reset() {
    selid = 0;
}

export function onEvent(e) {
}

function onseltxn(txn) {
    // If edit form is open, just cancel edit without selecting anything.
    if (editid != -1) {
        editid = -1;
        return;
    }

    selid = txn.txnid;
    dispatch("select", txn);
}
function onedittxn(txn) {
    editid = txn.txnid;
}
function oncreate(e) {
    if (account == null) {
        return;
    }
    newtxn.accountid = account.accountid;
    newtxn.date = new Date().toISOString(),
    editid = 0;
}
function txnform_done(e) {
    editid = -1;
}

</script>
