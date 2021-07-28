{#if txn == null}
    <p class="fg-dim">Select Txn</p>
{:else}
    <form class="" on:submit|preventDefault={onSubmit}>
        <div class="flex flex-row mb-2">
            <input class="bg-input fg-normal py-1 px-2 mr-1 flex-grow" name="desc" id="desc" type="text" placeholder="Enter Description" bind:value={txn.desc}>
            {#if txn.amt >= 0}
            <select class="py-1 px-2 bg-input fg-normal mr-1 flex-shrink" id="amttype" name="amttype" placeholder="Deposit/Withdraw" value="deposit">
                <option value="deposit">Deposit</option>
                <option value="withdraw">Withdraw</option>
            </select>
            {:else}
            <select class="py-1 px-2 bg-input fg-normal mr-1 flex-shrink" id="amttype" name="amttype" placeholder="Deposit/Withdraw" value="withdraw">
                <option value="deposit">Deposit</option>
                <option value="withdraw">Withdraw</option>
            </select>
            {/if}
            <input class="block bg-input fg-normal py-1 px-2 cell-amt" name="amt" id="amt" type="number" placeholder="0.00" value={Math.abs(txn.amt)}>
        </div>
        <div class="flex flex-row mb-2">
            <input class="bg-input fg-normal py-1 px-2 cell-date mr-1" name="date" id="date" type="date" placeholder="" value={txn.date.substring(0,10)}>
            <input class="block bg-input fg-normal py-1 px-2 cell-amt mr-1" name="ref" id="ref" type="text" placeholder="ref no" bind:value={txn.ref}>
            <input class="block bg-input fg-normal py-1 px-2 flex-grow" name="memo" id="memo" type="text" placeholder="memo" value="">
        </div>
        <div class="flex flex-row justify-start">
            <div>
                {#if txn.txnid == 0}
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2">Create</button>
                {:else}
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2" on:click|preventDefault={onSubmit}>Update</button>
                {/if}
                <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault={onCancel}>Cancel</a>
            </div>
        </div>
        {#if ui.submitstatus != ""}
        <div class="">
            <p class="uppercase italic text-xs">{ui.submitstatus}</p>
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

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    ui.status = "processing";

    let sreq = `${svcurl}/txn`;
    let method = "PUT";
    if (txn.txnid == 0) {
        method = "POST";
    }
    let [savedtxn, err] = await submit(sreq, method, txn);
    console.log(savedtxn);
    if (err != null) {
        console.error(err);
        ui.submitstatus = "server error submitting txn";
        return;
    }

    ui.status = "";
    txn = savedtxn;
    dispatch("submit", savedtxn);
}

function onCancel(e) {
    dispatch("cancel");
}

</script>
