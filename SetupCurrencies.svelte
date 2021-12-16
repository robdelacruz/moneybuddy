{#if root == null || root.currencies.length == 0}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flexrow justify-between mb-2">
        <h2 class="h2">Setup Currencies</h2>
        <a class="pill text-xs" href="/" on:click|preventDefault={oncreate}>New</a>
    </div>
    {#if editid == 0}
        <div class="p-2 border-b border-cell mb-2">
            <CurrencyForm currency={newcurrency} userid={userid} on:submit={currencyform_done} on:cancel={currencyform_done} />
        </div>
    {/if}
    {#each root.currencies as c (c.currencyid)}
        <a class="tblrow" class:sel="{selid == c.currencyid}" href="/" on:click|preventDefault="{e => onclickrow(c, c.currencyid)}">
            <p class="cell-desc">{c.name}</p>
        </a>
        {#if editid == c.currencyid}
        <div class="p-2 border-b border-cell">
            <CurrencyForm currency={c} userid={userid} on:submit={currencyform_done} on:cancel={currencyform_done} />
        </div>
        {/if}
    {/each}
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";
import * as data from "./data.js";
import CurrencyForm from "./CurrencyForm.svelte";

export let root = null;
export let userid = 0;

let selid = 0;
let editid = -1;
let newcurrency = {
    currencyid: 0,
    name: "",
    usdrate: null,
    userid: userid,
};

export function reset() {
    selid = 0;
    editid = -1;
}

function onclickrow(item, itemid) {
    // If row already selected, edit it.
    if (selid == itemid && selid != editid) {
        editid = itemid;
        return;
    }

    // row not selected

    // If edit form is open, just cancel edit without selecting anything.
    if (editid != -1) {
        editid = -1;
        return;
    }

    selid = itemid;
}
function oncreate(c) {
    editid = 0;
    selid = 0;
}
function currencyform_done(e) {
    editid = -1;
}

</script>


