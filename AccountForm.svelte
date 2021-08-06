{#if ui.account == null}
    <p class="fg-dim">Select Account</p>
{:else}
    <form class="" on:submit|preventDefault={onSubmit}>
        <div class="mb-2">
            <input class="block bg-input fg-normal py-1 px-2 w-full" name="accountname" id="accountname" type="text" placeholder="Enter Account Name" bind:value={ui.account.name}>
        </div>
        <div class="flex flex-row mb-2">
            <div class="mr-2 w-1/2">
                <select class="py-1 px-2 bg-input fg-normal w-full" id="accounttype" name="accounttype" placeholder="Account Type" bind:value={ui.account.accounttype}>
                    <option value={0}>Bank Account</option>
                    <option value={1}>Stock</option>
                </select>
            </div>
            <div class="w-1/2">
                <select class="py-1 px-2 bg-input fg-normal w-full" id="currency" name="currency" placeholder="Currency" bind:value={ui.account.currencyid}>
                    {#each currencies as currency}
                    <option value={currency.currencyid}>{currency.currency}</option>
                    {/each}
                </select>
            </div>
        </div>
        <div class="flex flex-row justify-start">
            <div>
                {#if ui.account.accountid == 0}
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

export let account = null;
export let currencies = [];

let svcurl = "/api";
let ui = {};
ui.status = "";

ui.account = {};
ui.account.accountid = account.accountid;
ui.account.code = account.code;
ui.account.name = account.name;
ui.account.accounttype = account.accounttype;
ui.account.currencyid = account.currencyid;

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    ui.status = "processing";

    let sreq = `${svcurl}/account`;
    let method = "PUT";
    if (ui.account.accountid == 0) {
        method = "POST";
    }
    let [savedaccount, err] = await submit(sreq, method, ui.account);
    if (err != null) {
        console.error(err);
        ui.status = "server error submitting account";
        return;
    }

    ui.status = "";
    ui.account = savedaccount;
    dispatch("submit", ui.account);
}

function onCancel(e) {
    dispatch("cancel");
}

</script>
