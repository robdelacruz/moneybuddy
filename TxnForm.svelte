{#if ui.txn == null}
    <p class="fg-dim">Select Txn</p>
{:else}
    <form class="" on:submit|preventDefault={onSubmit}>
        <div class="flex flex-row mb-2">
            <input class="bg-input fg-normal py-1 px-2 mr-1 flex-grow" name="desc" id="desc" type="text" placeholder="Enter Description" bind:value={ui.txn.desc}>
            <select class="py-1 px-2 bg-input fg-normal mr-1 flex-shrink" id="action" name="action" placeholder="Deposit/Withdraw" bind:value={ui.txn.action}>
                <option value="deposit">Deposit</option>
                <option value="withdraw">Withdraw</option>
            </select>
            <input class="block bg-input fg-normal py-1 px-2 cell-amt" name="amt" id="amt" type="number" placeholder="0.00" step="any" bind:value={ui.txn.absamt}>
        </div>
        <div class="flex flex-row mb-2">
            <input class="bg-input fg-normal py-1 px-2 cell-date mr-1" name="date" id="date" type="date" bind:value={ui.txn.isodate}>
            <input class="block bg-input fg-normal py-1 px-2 cell-amt mr-1" name="ref" id="ref" type="text" placeholder="ref no" bind:value={ui.txn.ref}>
            <input class="block bg-input fg-normal py-1 px-2 flex-grow" name="memo" id="memo" type="text" placeholder="memo" bind:value={ui.txn.memo}>
        </div>
        <div class="flex flex-row justify-start">
            <div>
                {#if ui.txn.txnid == 0}
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

ui.txn = {};
ui.txn.txnid = txn.txnid;
ui.txn.accountid = txn.accountid;
ui.txn.isodate = new Date().toISOString().substring(0,10);
ui.txn.date = txn.date;
ui.txn.ref = txn.ref;
ui.txn.desc = txn.desc;
ui.txn.amt = txn.amt;

onMount(function() {
    if (ui.txn == null) {
        return;
    }

    if (ui.txn.amt >= 0) {
        ui.txn.action = "deposit";
    } else {
        ui.txn.action = "withdraw";
    }
    ui.txn.absamt = Math.abs(ui.txn.amt);
    ui.txn.isodate = ui.txn.date.substring(0,10);
    ui.txn.memo = "";
});

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    ui.status = "processing";

    ui.txn.amt = Math.abs(ui.txn.absamt);
    if (ui.txn.action == "withdraw") {
        ui.txn.amt = -txn.amt;
    }
    ui.txn.date = ui.txn.isodate;

    let sreq = `${svcurl}/txn`;
    let method = "PUT";
    if (ui.txn.txnid == 0) {
        method = "POST";
    }
    let [savedtxn, err] = await submit(sreq, method, ui.txn);
    if (err != null) {
        console.error(err);
        ui.submitstatus = "server error submitting txn";
        return;
    }

    ui.status = "";
    ui.txn = savedtxn;
    dispatch("submit", ui.txn);
}

function onCancel(e) {
    dispatch("cancel");
}

</script>
