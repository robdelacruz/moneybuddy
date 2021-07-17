<div class:dim="{widgetstate == 'dim'}" class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="width: 20rem;">
    <h1 class="text-sm font-bold mb-2">Accounts</h1>
{#each accounts as account, i}
    <a href="/" on:click|preventDefault="{e => onselaccount(account, i)}">
        <div class:highlight="{ui.isel == i}" class="flex flex-row justify-between p-1 border-b border-cell">
            <p class="truncate cell-desc">{account.name}</p>
            <p class="fg-dim text-right cell-amt">{account.fmtbalance}</p>
        </div>
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
ui.isel = -1;

export function onEvent(e) {
    if (accounts.length == 0) {
        return;
    }

    if (e.key == "ArrowUp") {
        ui.isel--;
    } else if (e.key == "ArrowDown") {
        ui.isel++;
    } else {
        return;
    }

    if (ui.isel < 0) {
        ui.isel = 0;
    }
    if (ui.isel > accounts.length-1) {
        ui.isel = accounts.length-1;
    }

    dispatch("select", accounts[ui.isel]);
}

function onselaccount(account, i) {
    ui.isel = i;
    dispatch("select", account);
}

</script>

