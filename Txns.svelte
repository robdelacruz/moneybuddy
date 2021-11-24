<div class="txns bg-normal fg-normal py-2 px-4">
{#if root == null || selbook == null || displayaccount == null}
    <p class="fg-dim">Select Account</p>
{:else}
    <div class="flex flex-row justify-between items-end mb-2">
        <div class="flex flex-row">
            <select class="text-sm font-bold fg-h1 bg-normal pr-2 mr-2" id="selectaccount" name="selectaccount" placeholder="Select Account" bind:value={accountid} on:change={onaccountchange} on:blur="{e => {}}">
                {#each bookaccounts as a}
                    {#if a.accountid == accountid}
                    <option selected value={a.accountid}>{a.name}</option>
                    {:else}
                    <option value={a.accountid}>{a.name}</option>
                    {/if}
                {/each}
            </select>
            {#if displayaccount.balance >= 0}
            <p class="text-sm self-end fg-number-plus">{displayaccount.fmtbalance}</p>
            {:else}
            <p class="text-sm self-end fg-number-minus">{displayaccount.fmtbalance}</p>
            {/if}
        </div>
        <a class="text-xs pill" href="/" on:click|preventDefault={oncreate}>New</a>
    </div>
    {#if editid != 0}
    <!-- Don't show filter when Create form is visible. -->
        <div class="mb-2">
            <form autocomplete="off" on:submit|preventDefault="{e => {}}">
                <input class="block bg-input fg-normal py-1 px-2 w-full" name="filter" id="txnfilter" type="text" placeholder="Filter" bind:value={frm_filter} bind:this={input_filter}>
            </form>
        </div>
    {/if}
    {#if editid == 0}
        <div class="p-2 border-b border-cell">
            <TxnForm txn={newtxn} account={displayaccount} on:submit={txnform_done} on:cancel={txnform_done} />
        </div>
    {/if}
    {#each displaytxns as t (t.txnid)}
        <a class="txnrow" class:sel="{selid == t.txnid}" href="/" on:click|preventDefault="{e => onclicktxn(t)}">
            <p class="cell-date">{t.date.substring(0, 10)}</p>
            <p class="cell-refno">{t.ref}</p>
            <a class="cell-tag" class:sel="{selnoteid == t.txnid}" href="/" on:click|preventDefault="{e => onclicktag(e, t)}">
            {#if t.memo != ""}
            <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M0 10V2l2-2h8l10 10-10 10L0 10zm4.5-4a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3z"/></svg>
            {/if}
            </a>
            <p class="cell-desc">{t.desc}</p>
            {#if t.amt < 0}
            <p class="cell-amt fg-number-minus">{t.fmtamt}</p>
            {:else}
            <p class="cell-amt fg-number-plus">{t.fmtamt}</p>
            {/if}
        </a>
        {#if selnoteid == t.txnid && t.memo != ""}
        <div class="txnrow">
            <p class="cell-date"></p>
            <p class="cell-refno"></p>
            <p class="cell-tag" href="/"></p>
            <div class="cell-desc">
                <div class="memonote memo">
                    {@html textToHtml(t.memo)}
                </div>
            </div>
            <p class="cell-amt"></p>
        </div>
        {/if}
        {#if editid == t.txnid}
        <div class="p-2 border-b border-cell">
            <TxnForm txn={t} account={displayaccount} on:submit={txnform_done} on:cancel={txnform_done} />
        </div>
        {/if}
    {/each}
{/if}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, textToHtml} from "./helpers.js";
import * as data from "./data.js";
import TxnForm from "./TxnForm.svelte";

export let root = null;
export let bookid = 1;
export let accountid = 0;

let svcurl = "/api";
let selid = 0;
let editid = -1;
let selnoteid = 0;
let newtxn = {
    txnid: 0,
    accountid: 0,
    date: "",
    desc: "",
    amt: null,
    ref: "",
    memo: "",
};

let frm_filter = "";
let input_filter = null;

let bookaccounts = [];
let selbook = null;
let displayaccount = null;
let displaytxns = [];

// root + bookid --> selbook
$: selbook = getSelectedBook(root, bookid);

// selbook --> bookaccounts
$: bookaccounts = getBookAccounts(selbook);

// selbook + accountid --> displayaccount
$: displayaccount = getBookAccount(selbook, accountid);

// displayaccount + frm_filter --> displaytxns
$: displaytxns = filterTxns(displayaccount, frm_filter);

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        // If there's a selection, signal that it has been deselected.
        if (selid != 0 && editid == -1) {
            dispatch("selecttxn", null);
        }
        // Don't remove the selection if edit form is active.
        if (editid == -1) {
            reset();
        }
    }
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

function getBookAccounts(book) {
    if (book == null) {
        return null;
    }
    let aa = [];
    for (let i=0; i < book.bankaccounts.length; i++) {
        aa.push(book.bankaccounts[i]);
    }
    for (let i=0; i < book.stockaccounts.length; i++) {
        aa.push(book.stockaccounts[i]);
    }
    return aa;
}

function getBookAccount(book, accountid) {
    if (book == null || accountid == 0) {
        return null;
    }
    for (let i=0; i < book.bankaccounts.length; i++) {
        let a = book.bankaccounts[i];
        if (a.accountid == accountid) {
            return a;
        }
    }
    for (let i=0; i < book.stockaccounts.length; i++) {
        let a = book.stockaccounts[i];
        if (a.accountid == accountid) {
            return a;
        }
    }
    return null;
}

function filterTxns(account, sfilter) {
    if (account == null) {
        return [];
    }
    if (account.txns == null) {
        return [];
    }
    sfilter = sfilter.trim().toLowerCase();
    if (sfilter == "") {
        return account.txns;
    }
    let tt = [];
    for (let i=0; i < account.txns.length; i++) {
        let t = account.txns[i];
        if (t.desc.toLowerCase().includes(sfilter) || 
            t.ref.toLowerCase().includes(sfilter)  ||
            t.memo.toLowerCase().includes(sfilter)) {
                tt.push(t);
        }
    }
    return tt;
}

function onaccountchange(e) {
    let a = getBookAccount(selbook, accountid);
    if (a == null) {
        return;
    }
    selid = 0;
    editid = -1;
    dispatch("selectaccount", a);
}

export function reset() {
    selid = 0;
    editid = -1;
    selnoteid = 0;
}

function onclicktxn(txn) {
    // If txn already selected, edit it.
    if (selid == txn.txnid && selid != editid) {
        editid = txn.txnid;
        return;
    }

    // txn not selected

    // If edit form is open, just cancel edit without selecting anything.
    if (editid != -1) {
        editid = -1;
        return;
    }

    selid = txn.txnid;
    dispatch("selecttxn", txn);
}
function oncreate(e) {
    if (accountid == 0) {
        return;
    }
    newtxn.accountid = accountid;
    newtxn.date = new Date().toISOString(),
    editid = 0;
    selid = 0;
}
function txnform_done(e) {
    editid = -1;
}

function onclicktag(e, txn) {
    e.stopPropagation();
    if (selnoteid == txn.txnid) {
        selnoteid = 0;
        return;
    }
    selnoteid = txn.txnid;
}

</script>
