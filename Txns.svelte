<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="min-width: 40rem;">
{#if account == null || account.txns == null}
    <p class="fg-dim">Select Account</p>
{:else}
    <h1 class="text-sm font-bold mb-2">{account.name}: Transactions</h1>
    {#each account.txns as t (t.txnid)}
    {#if ui.editid == t.txnid}
        <div class="p-2 border-b border-cell">
            <TxnForm txn={t} on:submit={txnform_submit} on:cancel={txnform_cancel} />
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
ui.editid = 0;

export function reset() {
    ui.selid = 0;
}

export function onEvent(e) {
}

function onseltxn(txn) {
    // If edit form is open, just cancel edit without selecting anything.
    if (ui.editid != 0) {
        ui.editid = 0;
        return;
    }

    ui.selid = txn.txnid;
    dispatch("select", txn);
}
function onedittxn(txn) {
    ui.editid = txn.txnid;
}

function txnform_submit(e) {
    ui.editid = 0;

    let updatedTxn = e.detail;
    updatedTxn.fmtamt = data.formattedAmt(updatedTxn.amt, account.currency);

    for (let i=0; i < account.txns.length; i++) {
        if (account.txns[i].txnid == updatedTxn.txnid) {
            account.txns[i] = updatedTxn;
        }
    }
    account = account;
}
function txnform_cancel(e) {
    ui.editid = 0;
}

</script>
