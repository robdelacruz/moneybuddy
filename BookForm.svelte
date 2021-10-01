{#if book == null}
    <p class="fg-dim">Select Book</p>
{:else}
    <form class="" on:submit|preventDefault={onSubmit}>
        <div class="mb-2">
            <input class="block bg-input fg-normal py-1 px-2 w-full" name="bookname" id="bookname" type="text" placeholder="Book Name" bind:value={frm_name}>
        </div>
        <div class="mb-2">
            <label class="inline-flex items-center bg-input fg-normal py-1 px-2 w-full" for="bookactive">
                <input class="fg-normal mr-1" name="bookactive" id="bookactive" type="checkbox" bind:checked={frm_active}>
                <span class="">Active</span>
            </label>
        </div>
        {#if mode == ""}
        <div class="flex flex-row justify-between">
            <div>
                {#if book.bookid == 0}
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2">Create</button>
                {:else}
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2">Update</button>
                {/if}
                <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault={onCancel}>Cancel</a>
            </div>
            <div>
                <button class="mx-auto border border-normal py-1 px-2 bg-input" on:click|preventDefault={onDelete}>Delete</button>
            </div>
        </div>
        {:else if mode == "delete"}
        <div class="flex flex-row justify-left">
            <p class="self-center uppercase italic text-xs mr-4">Delete this book?</p>
            <div>
                <button class="mx-auto border border-normal py-1 px-2 bg-inputdel mr-2" on:click|preventDefault={onConfirmDelete}>Delete</button>
                <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault={onCancelDelete}>Cancel</a>
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
import {find, submit, del} from "./helpers.js";
import * as data from "./data.js";

export let book = null;
export let userid = 0;

let svcurl = "/api";
let mode = "";
let status = "";
let frm_name = book.name;
let frm_active = false;
if (book.active > 0) {
    frm_active = true;
}

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    status = "processing";

    let b = {};
    b.bookid = book.bookid;
    b.userid = userid;
    b.name = frm_name;
    b.active = 0;
    console.log(`frm_active = ${frm_active}`);
    if (frm_active == true) {
        b.active = 1;
    }

    let sreq = `${svcurl}/book`;
    let method = "PUT";
    if (b.bookid == 0) {
        method = "POST";
    }
    let err;
    [b, err] = await submit(sreq, method, b);
    if (err != null) {
        console.error(err);
        status = "server error submitting book";
        return;
    }

    status = "";
    book = b;
    dispatch("submit", b);
}

function onDelete(e) {
    mode = "delete";
}

function onCancelDelete(e) {
    mode = "";
}

async function onConfirmDelete(e) {
    status = "processing";

    let sreq = `${svcurl}/book?id=${book.bookid}`;
    let err = await del(sreq);
    if (err != null) {
        console.error(err);
        status = "server error deleting book";
        return;
    }

    status = "";
    dispatch("submit");
}


function onCancel(e) {
    dispatch("cancel");
}

</script>
