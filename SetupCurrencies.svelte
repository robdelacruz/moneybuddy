{#if root == null || root.currencies.length == 0}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flex flex-row justify-between items-end mb-2">
        <h2 class="text-sm font-bold fg-h1 bg-normal">Setup Currencies</h2>
        <a class="text-xs pill" href="/" on:click|preventDefault={oncreate}>New</a>
    </div>
    {#if editid == 0}
        <div class="p-2 border-b border-cell mb-2">
            <CurrencyForm currency={newcurrency} userid={userid} on:submit={currencyform_done} on:cancel={currencyform_done} />
        </div>
    {/if}
    {#each root.currencies as c (c.currencyid)}
        {#if selid == c.currencyid && selid != editid}
        <a class="flex flex-row justify-start p-1 border-b border-cell highlight" href="/" on:click|preventDefault="{e => onedit(c)}">
            <p class="flex-grow truncate mr-2">{c.name}</p>
        </a>
        {:else}
        <a class="flex flex-row justify-start p-1 border-b border-cell" href="/" on:click|preventDefault="{e => onselect(c)}">
            <p class="flex-grow truncate mr-2">{c.name}</p>
        </a>
        {/if}
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

export function select(c) {
    // If edit form is open, just cancel edit without selecting anything.
    if (editid != -1) {
        editid = -1;
        return;
    }

    selid = c.currencyid;
}
function onselect(c) {
    select(c);
}
function onedit(c) {
    editid = c.currencyid;
}
function oncreate(c) {
    editid = 0;
    selid = 0;
}
function currencyform_done(e) {
    editid = -1;
}

</script>


