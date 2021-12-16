{#if currency == null}
    <p class="fg-dim">Select Currency</p>
{:else}
    <form class="" autocomplete="off" on:submit|preventDefault={onSubmit}>
        <div class="mb-2">
            <input class="input w-full" name="name" id="name" type="text" placeholder="Currency Name" bind:value={frm_name}>
        </div>
        <div class="mb-2">
            <input class="input w-full" name="usdrate" id="usdrate" type="number" placeholder="USD Rate" step="any" min="0.0" bind:value={frm_usdrate}>
        </div>
        {#if mode == ""}
        <div class="flexrow justify-between">
            <div>
                {#if currency.currencyid == 0}
                <button class="btn bg-inputok mr-2">Create</button>
                {:else}
                <button class="btn bg-inputok mr-2">Update</button>
                {/if}
                <a href="/" class="action" on:click|preventDefault={onCancel}>Cancel</a>
            </div>
            <div>
                {#if currency.currencyid != 0}
                <button class="btn bg-input" on:click|preventDefault={onDelete}>Delete</button>
                {/if}
            </div>
        </div>
        {:else if mode == "delete"}
        <div class="">
            <p class="prompt mb-2">Delete this currency?</p>
            <div>
                <button class="btn bg-inputdel mr-2" on:click|preventDefault={onConfirmDelete}>Delete</button>
                <a href="/" class="action" on:click|preventDefault={onCancelDelete}>Cancel</a>
            </div>
        </div>
        {/if}
        {#if status != ""}
        <div class="">
            <p class="prompt">{status}</p>
        </div>
        {/if}
    </form>
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, del} from "./helpers.js";
import * as data from "./data.js";

export let currency = null;
export let userid = 0;

console.log(currency);

let svcurl = "/api";
let mode = "";
let status = "";
let frm_name = currency.name;
let frm_usdrate = currency.usdrate;

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    status = "processing";

    if (frm_name == "") {
        status = "enter the currency name (Ex. USD)";
        return;
    }
    if (frm_usdrate == null || frm_usdrate == 0) {
        status = "enter the currency rate in terms of USD";
        return;
    }

    let c = {};
    c.currencyid = currency.currencyid;
    c.name = frm_name;
    c.usdrate = frm_usdrate;
    c.userid = userid;

    let sreq = `${svcurl}/currency`;
    let method = "PUT";
    if (c.currencyid == 0) {
        method = "POST";
    }
    let err;
    [c, err] = await submit(sreq, method, c);
    if (err != null) {
        console.error(err);
        status = "server error submitting currency";
        return;
    }

    status = "";
    currency = c;
    dispatch("submit", c);
}

function onDelete(e) {
    mode = "delete";
}

function onCancelDelete(e) {
    mode = "";
}

async function onConfirmDelete(e) {
    status = "processing";

    let sreq = `${svcurl}/currency?id=${currency.currencyid}`;
    let err = await del(sreq);
    if (err != null) {
        console.error(err);
        status = "server error deleting currency";
        return;
    }

    status = "";
    dispatch("submit");
}


function onCancel(e) {
    dispatch("cancel");
}

</script>
