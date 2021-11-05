{#if txn == null}
    <p class="fg-dim">Select Txn</p>
{:else}
    <form class="txnform" autocomplete="off" on:submit|preventDefault={onSubmit}>
        <input class="txnform-desc bg-input fg-normal py-1 px-2" name="desc" id="desc" type="text" placeholder="Enter Description" bind:value={frm_desc}>
        <input class="txnform-date bg-input fg-normal py-1 px-2" name="date" id="date" type="date" bind:value={frm_date}>
        <select class="txnform-action py-1 px-2 bg-input fg-normal" id="action" name="action" placeholder="Deposit/Withdraw" bind:value={frm_action}>
            <option value="plus">{option_plus}</option>
            <option value="minus">{option_minus}</option>
        </select>
        <input class="txnform-amt bg-input fg-normal py-1 px-2" name="amt" id="amt" type="text" placeholder="Amount" bind:value={frm_amt} on:blur="{e => frm_amt = formatnum(frm_amt)}">
        <input class="txnform-refno bg-input fg-normal py-1 px-2" name="ref" id="ref" type="text" placeholder="Reference No" bind:value={frm_ref}>
        {#if mode == ""}
        <div class="txnform-btns flex flex-row justify-between">
            <div>
                {#if txn.txnid == 0}
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2">Submit</button>
                {:else}
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2">Update</button>
                {/if}
                <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault={onCancel}>Cancel</a>
            </div>
            <div>
                {#if txn.txnid != 0}
                <button class="mx-auto border border-normal py-1 px-2 bg-input" on:click|preventDefault={onDelete}>Delete</button>
                {/if}
            </div>
        </div>
        {:else if mode == "delete"}
        <div class="txnform-btns flex flex-row justify-left">
            <p class="self-center uppercase italic text-xs mr-4">Delete this transaction?</p>
            <div>
                <button class="mx-auto border border-normal py-1 px-2 bg-inputdel mr-2" on:click|preventDefault={onConfirmDelete}>Delete</button>
                <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault={onCancelDelete}>Cancel</a>
            </div>
        </div>
        {/if}
        {#if status != ""}
        <div class="txnform-status">
            <p class="uppercase italic text-xs">{status}</p>
        </div>
        {/if}
    </form>
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, del, formatnum} from "./helpers.js";
import * as data from "./data.js";

export let account = null;
export let txn = null;

let svcurl = "/api";
let ui = {};
let mode = "";
let status = "";

let frm_desc = txn.desc;
let frm_ref = txn.ref;

let frm_amt = "";
if (txn.amt != null) {
    frm_amt = formatnum(Math.abs(txn.amt).toString());
}

let option_plus;
let option_minus;
if (account.accounttype == 0) {
    option_plus = "Deposit";
    option_minus = "Withdraw";
} else {
    option_plus = "Buy";
    option_minus = "Sell";
}

let frm_action;
if (txn.amt >= 0) {
    frm_action = "plus"
} else {
    frm_action = "minus"
}

let frm_date;
if (txn.date) {
    frm_date = txn.date.substring(0,10);
} else {
    // If no date specified, use today's date.
    frm_date = new Date().toISOString().substring(0,10);
}

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    status = "processing";

    if (frm_ref == null) {
        frm_ref = "";
    }

    let t = {};
    t.txnid = txn.txnid;
    t.accountid = txn.accountid;
    if (t.accountid == 0) {
        t.accountid = account.accountid;
    }
    t.date = frm_date;
    t.ref = frm_ref;
    t.desc = frm_desc;

    if (frm_amt == "") {
        status = "please enter an amount";
        return;
    }

    t.amt = Math.abs(parseFloat(frm_amt));
    if (frm_action == "minus") {
        t.amt = -t.amt;
    }

    // If empty desc, put default action
    // (deposit/withdraw for accounts or buy/sell for stocks)
    if (t.desc.trim() == "") {
        if (t.amt >= 0) {
            t.desc = option_plus;
        } else {
            t.desc = option_minus;
        }
    }

    let sreq = `${svcurl}/txn`;
    let method = "PUT";
    if (t.txnid == 0) {
        method = "POST";
    }
    let err;
    [t, err] = await submit(sreq, method, t);
    if (err != null) {
        console.error(err);
        status = "server error submitting txn";
        return;
    }

    status = "";
    txn = t;
    dispatch("submit", t);
}

function onDelete(e) {
    mode = "delete";
}

function onCancelDelete(e) {
    mode = "";
}

async function onConfirmDelete(e) {
    status = "processing";

    let sreq = `${svcurl}/txn?id=${txn.txnid}`;
    let err = await del(sreq);
    if (err != null) {
        console.error(err);
        status = "server error deleting txn";
        return;
    }

    status = "";
    dispatch("submit");
}

function onCancel(e) {
    dispatch("cancel");
}

</script>
