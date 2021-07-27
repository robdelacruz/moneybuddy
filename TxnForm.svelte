{#if txn == null}
    <p class="fg-dim">Select Txn</p>
{:else}
    <form class="">
        <div class="flex flex-row mb-2">
            <input class="bg-input fg-normal py-1 px-2 mr-2" name="date" id="date" type="date" placeholder="" value={txn.date.substring(0,10)}>
            <input class="bg-input fg-normal py-1 px-2 flex-grow" name="desc" id="desc" type="text" placeholder="Enter Description" value={txn.desc}>
        </div>
        <div class="flex flex-row mb-2">
            {#if txn.amt >= 0}
            <select class="py-1 px-2 bg-input fg-normal mr-2" id="amttype" name="amttype" placeholder="Deposit/Withdraw" value="deposit">
                <option value="deposit">Deposit</option>
                <option value="withdraw">Withdraw</option>
            </select>
            {:else}
            <select class="py-1 px-2 bg-input fg-normal mr-2" id="amttype" name="amttype" placeholder="Deposit/Withdraw" value="withdraw">
                <option value="deposit">Deposit</option>
                <option value="withdraw">Withdraw</option>
            </select>
            {/if}
            <input class="block bg-input fg-normal py-1 px-2 cell-amt mr-2" name="amt" id="amt" type="number" placeholder="0.00" value={Math.abs(txn.amt)}>
            <input class="block bg-input fg-normal py-1 px-2 flex-grow" name="ref" id="ref" type="text" placeholder="ref no" value={txn.ref}>
        </div>
        <div class="flex flex-row justify-start">
            <div>
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2" on:click|preventDefault={onUpdate}>Update</button>
                <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault={onCancel}>Cancel</a>
            </div>
        </div>
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

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

function onUpdate(e) {
    dispatch("update", txn);
}

function onCancel(e) {
    dispatch("cancel");
}

</script>
