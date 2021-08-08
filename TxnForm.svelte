{#if txn == null}
    <p class="fg-dim">Select Txn</p>
{:else}
    <form class="" on:submit|preventDefault={onSubmit}>
        <div class="flex flex-row mb-2">
            <input class="bg-input fg-normal py-1 px-2 mr-1 flex-grow" name="desc" id="desc" type="text" placeholder="Enter Description" bind:value={ui.frm.desc}>
            <select class="py-1 px-2 bg-input fg-normal mr-1 flex-shrink" id="action" name="action" placeholder="Deposit/Withdraw" bind:value={ui.frm.action}>
                <option value="deposit">Deposit</option>
                <option value="withdraw">Withdraw</option>
            </select>
            <input class="block bg-input fg-normal py-1 px-2 cell-amt" name="amt" id="amt" type="number" placeholder="Amount" step="any" min="0.0" bind:value={ui.frm.amt}>
        </div>
        <div class="flex flex-row mb-2">
            <input class="bg-input fg-normal py-1 px-2 cell-date mr-1" name="date" id="date" type="date" bind:value={ui.frm.date}>
            <input class="block bg-input fg-normal py-1 px-2 cell-amt mr-1" name="ref" id="ref" type="text" placeholder="ref no" bind:value={ui.frm.ref}>
            <input class="block bg-input fg-normal py-1 px-2 flex-grow" name="memo" id="memo" type="text" placeholder="memo" bind:value={ui.frm.memo}>
        </div>
        <div class="flex flex-row justify-start">
            <div>
                {#if txn.txnid == 0}
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2">Create</button>
                {:else}
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2">Update</button>
                {/if}
                <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault={onCancel}>Cancel</a>
            </div>
        </div>
        {#if ui.status != ""}
        <div class="">
            <p class="uppercase italic text-xs">{ui.status}</p>
        </div>
        {/if}
    </form>
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";
import * as data from "./data.js";

export let txn = null;

let svcurl = "/api";
let ui = {};
ui.status = "";

$: init();

function init() {
    if (txn == null) {
        return;
    }

    ui.frm = {};
    ui.frm.desc = txn.desc;
    ui.frm.ref = txn.ref;
    ui.frm.memo = "";

    if (txn.amt >= 0) {
        ui.frm.action = "deposit";
    } else {
        ui.frm.action = "withdraw";
    }
    ui.frm.amt = Math.abs(txn.amt);

    if (txn.date) {
        ui.frm.date = txn.date.substring(0,10);
    } else {
        // If no date specified, use today's date.
        ui.frm.date = new Date().toISOString().substring(0,10);
    }
}

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    ui.status = "processing";

    let t = {};
    t.txnid = txn.txnid;
    t.accountid = txn.accountid;
    t.date = ui.frm.date;
    t.ref = ui.frm.ref;
    t.desc = ui.frm.desc;

    t.amt = Math.abs(ui.frm.amt);
    if (ui.frm.action == "withdraw") {
        t.amt = -t.amt;
    }

    // If empty desc, just put "deposit" or "withdraw" depending on action.
    if (t.desc.trim() == "") {
        t.desc = ui.frm.action;
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
        ui.status = "server error submitting txn";
        return;
    }

    ui.status = "";
    txn = t;
    dispatch("submit", t);
}

function onCancel(e) {
    dispatch("cancel");
}

</script>
