<div class="accounts bg-normal fg-normal py-2 px-4">
{#if root == null || root.books.length == 0}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flex flex-row justify-between items-end mb-2">
        <div class="">
            <select class="text-sm font-bold fg-h1 bg-normal pr-2" id="book" name="book" placeholder="Select Book" bind:value={bookid} on:change={onbookchange} on:blur="{e => {}}">
            {#each root.books as b}
                {#if b.bookid == bookid}
                <option selected value={b.bookid}>{b.name}</option>
                {:else}
                <option value={b.bookid}>{b.name}</option>
                {/if}
            {/each}
            </select>
        </div>
        <a class="text-xs pill" href="/" on:click|preventDefault={oncreate}>New</a>
    </div>
    {#if editid != 0}
    <!-- Don't show filter when Create form is visible. -->
        <div class="mb-4">
            <form autocomplete="off" on:submit|preventDefault="{e => {}}">
                <input class="block bg-input fg-normal py-1 px-2 w-full" name="filter" id="accountfilter" type="text" placeholder="Filter" bind:value={frm_filter}>
            </form>
        </div>
    {/if}
    {#if editid == 0}
        <div class="p-2 border-b border-cell">
            <AccountForm book={selbook} account={newaccount} root={root} on:submit={accountform_done} on:cancel={accountform_done} />
        </div>
    {/if}
    {#if display_bb.length > 0}
    <div class="mb-1">
        <h2 class="font-bold fg-h2">Bank Accounts</h2>
    </div>
    {/if}
    {#each display_bb as a (a.accountid)}
        {#if selid == a.accountid && selid != editid}
        <a class="flex flex-row justify-between p-1 border-b border-cell highlight" href="/" on:click|preventDefault="{e => oneditaccount(a)}">
            <p class="flex-grow truncate mr-2">{a.name}</p>
            {#if a.balance >= 0}
            <p class="whitespace-nowrap text-right mr-1">{a.fmtbalance}</p>
            {:else}
            <p class="whitespace-nowrap text-right mr-1">{a.fmtbalance}</p>
            {/if}
        </a>
        {:else}
        <a class="flex flex-row justify-between p-1 border-b border-cell" href="/" on:click|preventDefault="{e => onselectaccount(a)}">
            <p class="flex-grow truncate mr-2">{a.name}</p>
            {#if a.balance >= 0}
            <p class="whitespace-nowrap fg-number-plus text-right mr-1">{a.fmtbalance}</p>
            {:else}
            <p class="whitespace-nowrap fg-number-minus text-right mr-1">{a.fmtbalance}</p>
            {/if}
        </a>
        {/if}
        {#if editid == a.accountid}
        <div class="p-2 border-b border-cell">
            <AccountForm book={selbook} account={a} root={root} on:submit={accountform_done} on:cancel={accountform_done} />
        </div>
        {/if}
    {/each}
    <div class="mb-4">
    </div>
    {#if display_ss.length > 0}
    <div class="mb-1">
        <h2 class="font-bold fg-h2">Stocks</h2>
    </div>
    {/if}
    {#each display_ss as a (a.accountid)}
        {#if selid == a.accountid && selid != editid}
        <a class="flex flex-row justify-between p-1 border-b border-cell highlight" href="/" on:click|preventDefault="{e => oneditaccount(a)}">
            <p class="flex-grow truncate mr-2">{a.name}</p>
            {#if a.balance >= 0}
            <p class="text-right mr-1">{a.fmtbalance}</p>
            {:else}
            <p class="text-right mr-1">{a.fmtbalance}</p>
            {/if}
        </a>
        {:else}
        <a class="flex flex-row justify-between p-1 border-b border-cell" href="/" on:click|preventDefault="{e => onselectaccount(a)}">
            <p class="flex-grow truncate mr-2">{a.name}</p>
            {#if a.balance >= 0}
            <p class="fg-number-plus text-right mr-1">{a.fmtbalance}</p>
            {:else}
            <p class="fg-number-minus text-right mr-1">{a.fmtbalance}</p>
            {/if}
        </a>
        {/if}
        {#if editid == a.accountid}
        <div class="p-2 border-b border-cell">
            <AccountForm book={selbook} account={a} root={root} on:submit={accountform_done} on:cancel={accountform_done} />
        </div>
        {/if}
    {/each}
{/if}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, ifnull, getls, setls} from "./helpers.js";
import * as data from "./data.js";
import AccountForm from "./AccountForm.svelte";

export let root = null;
export let bookid = 0;
export let selaccountid = 0;

let selid = selaccountid;
let editid = -1;
let newaccount = {
    accountid: 0,
    code: "",
    name: "",
    accounttype: 0,
    currencyid: 0,
};

let frm_filter = "";

let selbook = null;
let display_bb = []; // bankaccounts to display
let display_ss = []; // stockaccounts to display

// root + bookid --> selbook
// selbook + frm_filter --> display_bb, display_ss

$: selbook = getSelectedBook(root, bookid);
$: [display_bb, display_ss] = filterAccounts(selbook, frm_filter);

function getSelectedBook(rootdata, bookid) {
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
// Returns [bankaccounts, stockaccounts] where account.name matches sfilter
function filterAccounts(book, sfilter) {
    let bb = [];
    let ss = [];
    if (book == null) {
        return [bb, ss];
    }
    if (book.bankaccounts == null || book.stockaccounts == null) {
        return [bb, ss];
    }
    sfilter = sfilter.trim().toLowerCase();
    if (sfilter == "") {
        return [book.bankaccounts, book.stockaccounts];
    }

    for (let i=0; i < book.bankaccounts.length; i++) {
        let a = book.bankaccounts[i];
        if (a.name.toLowerCase().includes(sfilter)) {
            bb.push(a);
        }
    }
    for (let i=0; i < book.stockaccounts.length; i++) {
        let a = book.stockaccounts[i];
        if (a.name.toLowerCase().includes(sfilter)) {
            ss.push(a);
        }
    }
    return [bb, ss];
}

function onbookchange(e) {
    selid = 0;
    editid = -1;

    dispatch("selectbookid", bookid);
}

export function reset() {
    selid = 0;
    editid = -1;
}

export function selectAccount(account) {
    // If edit form is open, just cancel edit without selecting anything.
    if (editid != -1) {
        editid = -1;
        return;
    }

    selid = account.accountid;
}
function onselectaccount(account) {
    selectAccount(account);
    dispatch("selectaccount", account);
}
function oneditaccount(account) {
    editid = account.accountid;
}
function oncreate(e) {
    editid = 0;
    selid = 0;
}
function accountform_done(e) {
    editid = -1;
}

</script>

