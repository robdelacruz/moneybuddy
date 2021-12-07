{#if account == null || root == null || book == null}
    <p class="fg-dim">Select Account</p>
{:else}
    <form class="" autocomplete="off" on:submit|preventDefault={onSubmit}>
    {#if mode == ""}
        {#if frm_accounttype == 0}
        <div class="mb-2">
            <input class="block bg-input fg-normal py-1 px-2 w-full" name="accountname" id="accountname" type="text" placeholder="Account Name" bind:value={frm_name}>
        </div>
        {:else}
        <div class="flex flex-row mb-2">
            <div class="mr-2 w-1/2">
                <input class="block bg-input fg-normal py-1 px-2 w-full" name="accountname" id="accountname" type="text" placeholder="Stock Name" bind:value={frm_name}>
            </div>
            <div class="w-1/2">
                <input class="block bg-input fg-normal py-1 px-2 w-full" name="unitprice" id="unitprice" type="text" placeholder="Unit Price" bind:value={frm_unitprice} on:blur="{e => frm_unitprice = formatnum(frm_unitprice)}">
            </div>
        </div>
        {/if}
        <div class="flex flex-row mb-2">
            <div class="mr-2 w-1/2">
                <select class="py-1 px-2 bg-input fg-normal w-full" id="accounttype" name="accounttype" placeholder="Account Type" bind:value={frm_accounttype}>
                    <option value={data.BANKACCOUNT}>Bank Account</option>
                    <option value={data.STOCKACCOUNT}>Stock</option>
                </select>
            </div>
            <div class="w-1/2">
                <select class="py-1 px-2 bg-input fg-normal w-full" id="currency" name="currency" placeholder="Currency" bind:value={frm_currencyid}>
                    {#each root.currencies as currency}
                    <option value={currency.currencyid}>{currency.name}</option>
                    {/each}
                </select>
            </div>
        </div>
        <div class="mb-2">
            <input class="py-1 px-2 bg-input fg-normal w-full" name="ref" id="ref" type="text" placeholder="Reference No" bind:value={frm_ref}>
        </div>
        <div class="mb-2">
        <textarea class="py-1 px-2 bg-input fg-normal w-full" name="memo" id="memo" placeholder="Memo" rows="4" bind:value={frm_memo}></textarea>
        </div>

        <div class="flex flex-row justify-between items-center">
            <div>
                {#if account.accountid == 0}
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2">Create</button>
                {:else}
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2">Update</button>
                {/if}
                <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault={onCancel}>Cancel</a>
            </div>
            <div class="flex flex-row justify-left">
                {#if account.accountid != 0}
                <button class="border border-normal py-1 px-2 bg-input mr-2" on:click|preventDefault={onMove}>Move...</button>
                {/if}
                {#if book.booktype == data.SYSTEMBOOK}
                <button class="border border-normal py-1 px-2 bg-input" on:click|preventDefault={onDelete}>Delete...</button>
                {/if}
            </div>
        </div>
    {:else if mode == "move"}
        <div class="">
            <p class="mb-2">Move <span class="font-bold">{account.name}</span> to:</p>
            <select class="py-1 px-2 bg-input fg-normal mb-2" id="movebook" name="movebook" bind:value={frm_movebookid}>
                {#each root.books as b}
                <option value={b.bookid}>{b.name}</option>
                {/each}
            </select>
            <div>
                <button class="border border-normal py-1 px-2 bg-inputok mr-2" on:click|preventDefault={onConfirmMove}>Move</button>
                <a href="/" class="border-b border-normal pt-1" on:click|preventDefault={onConfirmCancel}>Cancel</a>
            </div>
        </div>
    {:else if mode == "delete"}
        <div class="">
            <p class="mb-2">Delete <span class="font-bold">{account.name}</span>?</p>
            <div>
                <button class="mx-auto border border-normal py-1 px-2 bg-inputdel mr-2" on:click|preventDefault={onConfirmDelete}>Delete</button>
                <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault={onConfirmCancel}>Cancel</a>
            </div>
        </div>
    {/if}
    {#if status != ""}
        <div class="">
            <p class="uppercase italic text-xs">{status}</p>
        </div>
    {/if}
    </form>
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, del, formatnum} from "./helpers.js";
import * as data from "./data.js";

export let book = null;
export let account = null;
export let root = null;

let svcurl = "/api";
let mode = "";
let status = "";
let frm_name = account.name;
let frm_accounttype = account.accounttype;
let frm_currencyid = account.currencyid;
let frm_unitprice = "";
if (account.unitprice != null) {
    frm_unitprice = formatnum(account.unitprice.toString());
}
let frm_ref = account.ref;
let frm_memo = account.memo;
let frm_movebookid = book.bookid;

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    status = "processing";

    if (frm_unitprice == "") {
        frm_unitprice = "0.00";
    }

    let a = {};
    a.accountid = account.accountid;
    a.bookid = account.bookid;
    a.code = account.code;
    a.name = frm_name;
    a.accounttype = frm_accounttype
    a.currencyid = frm_currencyid;
    a.unitprice = parseFloat(frm_unitprice);
    a.ref = frm_ref;
    a.memo = frm_memo;

    let sreq = `${svcurl}/account?bookid=${book.bookid}`;
    let method = "PUT";
    if (a.accountid == 0) {
        method = "POST";
    }
    let err;
    [a, err] = await submit(sreq, method, a);
    if (err != null) {
        console.error(err);
        status = "server error submitting account";
        return;
    }

    status = "";
    account = a;
    dispatch("submit", a);
}

function onMove(e) {
    mode = "move";
}
function onDelete(e) {
    mode = "delete";
}
function onConfirmCancel(e) {
    mode = "";
}

async function onConfirmDelete(e) {
    status = "processing";

    let sreq = `${svcurl}/account?id=${account.accountid}`;
    let err = await del(sreq);
    if (err != null) {
        console.error(err);
        status = "server error deleting account";
        return;
    }

    status = "";
    dispatch("submit");
}

async function onConfirmMove(e) {
    if (account.accountid == 0) {
        mode = "";
        return;
    }
    if (frm_movebookid == 0) {
        mode = "";
        return;
    }

    let sreq = `${svcurl}/account?bookid=${frm_movebookid}`;
    let method = "PUT";
    let [_, err] = await submit(sreq, method, account);
    if (err != null) {
        console.error(err);
        status = "server error submitting account";
        return;
    }

    status = "";
    dispatch("submit", account);
}

function onCancel(e) {
    dispatch("cancel");
}

</script>
