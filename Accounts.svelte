<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="min-width: 20rem;">
    <h1 class="text-sm font-bold mb-2">Accounts</h1>
{#each accounts as account, i}
    {#if ui.editid == account.accountid}
    <div class="p-2 border-b border-cell">
        <AccountForm account={account} on:update={accountform_update} on:cancel={accountform_cancel} />
    </div>
    {:else if ui.selid == account.accountid}
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
import AccountForm from "./AccountForm.svelte";

export let accounts = [];

let svcurl = "/api";
let ui = {};
ui.selid = 0;
ui.editid = 0;
ui.mode = "";

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
    ui.editid = 0;
    dispatch("select", account);
}
function oneditaccount(account) {
    ui.editid = account.accountid;
    dispatch("edit", account);
}

function accountform_update(e) {
    ui.editid = 0;

    let updatedAccount = e.detail;
    for (let i=0; i < accounts.length; i++) {
        if (accounts[i].accountid == updatedAccount.accountid) {
            accounts[i] = updatedAccount;
        }
    }
}
function accountform_cancel(e) {
    ui.editid = 0;
}

</script>

