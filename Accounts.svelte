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
    {#if display_aa.length > 0}
    <div class="mb-1">
        <h2 class="font-bold fg-h2">Bank Accounts</h2>
    </div>
    {/if}
    {#each display_aa as a, i (a.accountid)}
        {#if i == ifirststock}
        <div class="mt-4 mb-1">
            <h2 class="font-bold fg-h2">Securities</h2>
        </div>
        {/if}
        <a class="accountrow" draggable="true" data-accountid="{a.accountid}" data-seq="{a.seq}" class:sel="{selid == a.accountid}" class:detail="{expandids.has(a.accountid)}" href="/" on:click|preventDefault="{e => onclickaccount(a)}">
            <p class="cell-desc">{a.name}</p>
            {#if a.balance >= 0}
            <p class="cell-amt" class:fg-number-plus="{selid != a.accountid}">{a.fmtbalance}</p>
            {:else}
            <p class="cell-amt" class:fg-number-minus="{selid != a.accountid}">{a.fmtbalance}</p>
            {/if}
            {#if a.ref != "" || a.memo != ""}
            <a class="cell-detailicon" href="/" on:click|preventDefault="{e => onclickdetail(e, a)}">
                {#if expandids.has(a.accountid)}
                <!-- chevron up arrow -->
                <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M10.707 7.05L10 6.343 4.343 12l1.414 1.414L10 9.172l4.243 4.242L15.657 12z"/></svg>
                {:else}
                <!-- chevron down arrow -->
                <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z"/></svg>
                {/if}
            </a>
            {:else}
            <p class="cell-detailicon"></p>
            {/if}
        </a>

        {#if expandids.has(a.accountid) && (a.ref != "" || a.memo != "")}
        <div class="accountrow">
            <div class="flex-grow flex-shrink">
                {#if a.ref != ""}
                <p class="fg-dim mb-1">{a.ref}</p>
                {/if}
                {#if a.memo != ""}
                <div class="memonote memo">
                    {@html textToHtml(a.memo)}
                </div>
                {/if}
            </div>
        </div>
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
import {find, submit, ifnull, getls, setls, textToHtml} from "./helpers.js";
import * as data from "./data.js";
import AccountForm from "./AccountForm.svelte";

export let root = null;
export let bookid = 0;
export let selaccountid = 0;

let selid = selaccountid;
let editid = -1;
let expandids = new Set();
let newaccount = {
    accountid: 0,
    code: "",
    name: "",
    accounttype: 0,
    currencyid: 0,
    ref: "",
    memo: "",
};

let frm_filter = "";

let selbook = null;
let display_aa = []; // bankaccounts and stockaccounts matching filter
let ifirststock = 0;

// root + bookid --> selbook
// selbook + frm_filter --> display_aa

$: selbook = getSelectedBook(root, bookid);
$: [display_aa, ifirststock] = filterAccounts(selbook, frm_filter);

document.addEventListener("dragstart", function(e) {
    let el = e.target;
    if (!el.classList.contains("accountrow")) {
        return;
    }
    e.dataTransfer.setData("text/plain", el.dataset.accountid);
    console.log(`dragstart: ${el.dataset.accountid}`);
});
document.addEventListener("dragover", function(e) {
    e.preventDefault();

    let el = e.target.closest(".accountrow");
    if (el == null) {
        e.dataTransfer.dropEffect = "none";
        return;
    }
    e.dataTransfer.dropEffect = "move";
});
document.addEventListener("drop", function(e) {
    e.preventDefault();

    let el = e.target.closest(".accountrow");
    if (el == null) {
        return;
    }

    let saccountid = e.dataTransfer.getData("text/plain");
    let accountid = parseInt(saccountid, 10);
    let targetseq = parseInt(el.dataset.seq, 10);
    console.log(`dropped ${accountid} to ${targetseq}`);
});


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
// Returns matching accounts where name, ref, or memo matches sfilter
// [aa, i] returned where:
//   aa = array of matching accounts
//   i  = index within aa of first stock account
//        (if no stock accounts, index returned will be out of range of aa)
function filterAccounts(book, sfilter) {
    let aa = [];
    let ifirststock = 0;
    if (book == null) {
        return [aa, ifirststock];
    }
    if (book.bankaccounts == null || book.stockaccounts == null) {
        return [aa, ifirststock];
    }
    sfilter = sfilter.trim().toLowerCase();
    if (sfilter == "") {
        aa = aa.concat(book.bankaccounts, book.stockaccounts);
        ifirststock = book.bankaccounts.length;
        return [aa, ifirststock];
    }

    for (let i=0; i < book.bankaccounts.length; i++) {
        let a = book.bankaccounts[i];
        if (a.name.toLowerCase().includes(sfilter) ||
            a.ref.toLowerCase().includes(sfilter) ||
            a.memo.toLowerCase().includes(sfilter)) {
            aa.push(a);
        }
    }

    ifirststock = aa.length;

    for (let i=0; i < book.stockaccounts.length; i++) {
        let a = book.stockaccounts[i];
        if (a.name.toLowerCase().includes(sfilter) ||
            a.ref.toLowerCase().includes(sfilter) ||
            a.memo.toLowerCase().includes(sfilter)) {
            aa.push(a);
        }
    }
    return [aa, ifirststock];
}

function onbookchange(e) {
    selid = 0;
    editid = -1;

    dispatch("selectbookid", bookid);
}

export function reset() {
    selid = 0;
    editid = -1;
    expandids.clear();
}

export function selectAccount(a) {
    // If edit form is open, just cancel edit without selecting anything.
    if (editid != -1) {
        editid = -1;
        return;
    }

    selid = a.accountid;
}
function onclickaccount(a) {
    // If account already selected, edit it.
    if (selid == a.accountid && selid != editid) {
        editid = a.accountid;
        return;
    }

    selectAccount(a);
    dispatch("selectaccount", a);
}
function oneditaccount(a) {
    editid = a.accountid;
}
function oncreate(e) {
    editid = 0;
    selid = 0;
}
function accountform_done(e) {
    editid = -1;
}

function onclickdetail(e, a) {
    e.stopPropagation();
    if (expandids.has(a.accountid)) {
        expandids.delete(a.accountid);
        expandids = expandids;
        return;
    }
    expandids.add(a.accountid);
    expandids = expandids;
}

</script>

