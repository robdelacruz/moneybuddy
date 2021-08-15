<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="min-width: 20rem;">
{#if root == null}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flex flex-row justify-between items-end mb-2">
        <div class="flex-grow">
            <select class="text-sm font-bold fg-normal bg-normal pr-2" id="book" name="book" placeholder="Select Book" bind:value={frm_bookid} on:change="{e => dispatch('select', null)}" on:blur="{e => {}}">
                {#each root.books as book}
                <option value={book.bookid}>{book.name}</option>
                {/each}
            </select>
        </div>
        <a class="text-xs pill" href="/" on:click|preventDefault={oncreate}>Create</a>
    </div>
    {#if editid != 0}
    <!-- Don't show filter when Create form is visible. -->
        <div class="mb-2">
            <form autocomplete="off" on:submit|preventDefault="{e => {}}">
                <input class="block bg-input fg-normal py-1 px-2 w-full" name="filter" id="filter" type="text" placeholder="Filter" bind:value={frm_filter}>
            </form>
        </div>
    {/if}
    {#if editid == 0}
        <div class="p-2 border-b border-cell">
            <AccountForm account={newaccount} currencies={root.currencies} on:submit={accountform_done} on:cancel={accountform_done} />
        </div>
    {/if}
    {#each displayaccounts as account (account.accountid)}
        {#if editid == account.accountid}
        <div class="p-2 border-b border-cell">
            <AccountForm account={account} currencies={root.currencies} on:submit={accountform_done} on:cancel={accountform_done} />
        </div>
        {:else if selid == account.accountid}
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
let selid = 0;
let editid = -1;
let newaccount = {
    accountid: 0,
    code: "",
    name: "",
    accounttype: 0,
    currencyid: 0,
};

let frm_bookid = 1;
let frm_filter = "";

let selbook = null;
let displayaccounts = [];

// root + frm_bookid --> selbook
// selbook + frm_filter --> displayaccounts

$: selbook = getSelectedBook(root, frm_bookid);
$: displayaccounts = filterAccounts(selbook, frm_filter);

function getSelectedBook(rootdata, bookid) {
    console.log("Accounts.svelte getSelectedBook()");
    if (rootdata == null) {
        return null;
    }
    let b = null;
    for (let i=0; i < rootdata.books.length; i++) {
        if (bookid == rootdata.books[i].bookid) {
            b = rootdata.books[i];
            break;
        }
    }
    return b
}
function filterAccounts(book, sfilter) {
    console.log("Accounts.svelte filterAccounts()");
    if (book == null) {
        return [];
    }
    if (book.accounts == null) {
        return [];
    }
    sfilter = sfilter.trim().toLowerCase();
    if (sfilter == "") {
        return book.accounts;
    }

    let aa = [];
    for (let i=0; i < book.accounts.length; i++) {
        let a = book.accounts[i];
        if (a.name.toLowerCase().includes(sfilter)) {
            aa.push(a);
        }
    }
    return aa;
}

export function reset() {
    selid = 0;
    editid = -1;
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
    if (editid != -1) {
        editid = -1;
        return;
    }

    selid = account.accountid;
    dispatch("select", account);
}
function oneditaccount(account) {
    editid = account.accountid;
}
function oncreate(e) {
    editid = 0;
}
function accountform_done(e) {
    editid = -1;
}

</script>

