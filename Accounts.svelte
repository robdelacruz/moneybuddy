<div class:dim="{widgetstate == 'dim'}" class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="min-width: 20rem;">
    <h1 class="text-sm font-bold mb-2">Accounts</h1>
{#each accounts as account, i}
    <a href="/" on:click|preventDefault="{e => onselaccount(account)}">
    {#if ui.selid == account.accountid}
        <div class="flex flex-row justify-between p-1 border-b border-cell highlight">
            <p class="truncate cell-desc mr-2">{account.name}</p>
            <p class="fg-dim text-right cell-amt mr-1">{account.fmtbalance}</p>
            <a class="cell-icon pl-1 py-1" href="/" on:click|preventDefault="{e => onseleditaccount(account)}">
                <svg class="w-2 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M12.3 3.7l4 4L4 20H0v-4L12.3 3.7zm1.4-1.4L16 0l4 4-2.3 2.3-4-4z"/></svg>
            </a>
        </div>
    {:else}
        <div class="flex flex-row justify-between p-1 border-b border-cell">
            <p class="truncate cell-desc mr-2">{account.name}</p>
            <p class="fg-dim text-right cell-amt mr-1">{account.fmtbalance}</p>
            <p class="cell-icon pl-1 py-1"></p>
        </div>
    {/if}
    </a>
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

function onseleditaccount(account) {
    dispatch("edit", account);
    console.log("onseleditaccount");
}

</script>

