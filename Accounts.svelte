<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="min-width: 20rem;">
{#if root == null || ui.accounts == null || ui.selbook == null}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flex flex-row justify-between items-end mb-2">
        <div class="flex-grow">
            <select class="text-sm font-bold fg-normal bg-normal pr-2" id="book" name="book" placeholder="Select Book" bind:value={ui.frm.bookid}>
                {#each root.books as book}
                <option value={book.bookid}>{book.name}</option>
                {/each}
            </select>
        </div>
        <a class="text-xs pill" href="/" on:click|preventDefault={oncreate}>Create</a>
    </div>
    {#if ui.editid != 0}
    <!-- Don't show filter when Create form is visible. -->
        <div class="mb-2">
            <form autocomplete="off" on:submit|preventDefault="{e => {}}">
                <input class="block bg-input fg-normal py-1 px-2 w-full" name="filter" id="filter" type="text" placeholder="Filter" bind:value={ui.frm.filter} on:keyup={processFilter}>
            </form>
        </div>
    {/if}
    {#if ui.editid == 0}
        <div class="p-2 border-b border-cell">
            <AccountForm account={ui.newaccount} currencies={root.currencies} on:submit={accountform_done} on:cancel={accountform_done} />
        </div>
    {/if}
    {#each ui.accounts as account (account.accountid)}
        {#if ui.editid == account.accountid}
        <div class="p-2 border-b border-cell">
            <AccountForm account={account} currencies={root.currencies} on:submit={accountform_done} on:cancel={accountform_done} />
        </div>
        {:else if ui.selid == account.accountid}
        <a class="flex flex-row justify-between p-1 border-b border-cell highlight" href="/" on:click|preventDefault="{e => oneditaccount(account)}">
            <p class="flex-grow truncate mr-2">{account.name}</p>
            <p class="fg-dim text-right mr-1">{account.fmtbalance}</p>
        </a>
        {:else}
        <a class="flex flex-row justify-between p-1 border-b border-cell" href="/" on:click|preventDefault="{e => onselaccount(account)}">
            <p class="flex-grow truncate mr-2">{account.name}</p>
            <p class="fg-dim text-right mr-1">{account.fmtbalance}</p>
        </a>
        {/if}
    {/each}
{/if}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";
import * as data from "./data.js";
import AccountForm from "./AccountForm.svelte";

export let root = null;

let svcurl = "/api";
let ui = {};
ui.accounts = null;
ui.selbook = null;

$: init();

$: if (root != null) {
    for (let i=0; i < root.books.length; i++) {
        if (ui.frm.bookid == root.books[i].bookid) {
            ui.selbook = root.books[i];
            break;
        }
    }
    processFilter();
}


function init() {
    ui.selid = 0;
    ui.editid = -1;
    ui.newaccount = {
        accountid: 0,
        code: "",
        name: "",
        accounttype: 0,
        currencyid: 0,
    };

    ui.frm = {};
    ui.frm.filter = "";
    ui.frm.bookid = 1;
}

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
    // If edit form is open, just cancel edit without selecting anything.
    if (ui.editid != -1) {
        ui.editid = -1;
        return;
    }

    ui.selid = account.accountid;
    dispatch("select", account);
}
function oneditaccount(account) {
    ui.editid = account.accountid;
}
function oncreate(e) {
    ui.editid = 0;
}
function accountform_done(e) {
    ui.editid = -1;
}

function processFilter() {
    if (ui.selbook == null) {
        ui.accounts = null;
        return;
    }
    let filter = ui.frm.filter.trim();
    if (filter == "") {
        ui.accounts = ui.selbook.accounts;
        return;
    }
    ui.accounts = filterAccounts(ui.selbook.accounts, filter);
}
function filterAccounts(accounts, filter) {
    filter = filter.toLowerCase();
    let aa = [];
    for (let i=0; i < accounts.length; i++) {
        let a = accounts[i];
        if (a.name.toLowerCase().includes(filter)) {
            aa.push(a);
        }
    }
    return aa;
}
</script>

