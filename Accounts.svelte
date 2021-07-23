<div class:dim="{widgetstate == 'dim'}" class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="min-width: 20rem;">
    <h1 class="text-sm font-bold mb-2">Accounts</h1>
{#each accounts as account, i}
    {#if ui.selid == account.accountid}
    <a class="flex flex-row justify-between p-1 border-b border-cell highlight" href="/" on:click|preventDefault="{e => oneditaccount(account)}">
        <p class="truncate cell-desc mr-2">{account.name}</p>
        <p class="fg-dim text-right cell-amt mr-1">{account.fmtbalance}</p>
    </a>
    {:else}
    <a class="flex flex-row justify-between p-1 border-b border-cell" href="/" on:click|preventDefault="{e => onselaccount(account)}">
        <p class="truncate cell-desc mr-2">{account.name}</p>
        <p class="fg-dim text-right cell-amt mr-1">{account.fmtbalance}</p>
    </a>
    {/if}
{/each}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";

export let widgetstate = "";
export let accounts = [];

let svcurl = "/api";
let ui = {};
ui.selid = 0;

export function reset() {
    ui.selid = 0;
}

export function onEvent(e) {
    if (e.key == "ArrowUp") {
    } else if (e.key == "ArrowDown") {
    } else {
        return;
    }
}

function onselaccount(account) {
    ui.selid = account.accountid;
    dispatch("select", account);
}

function oneditaccount(account) {
    dispatch("edit", account);
    console.log("oneditaccount");
}

</script>

