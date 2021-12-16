{#if root == null || root.books.length == 0}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flexrow justify-between mb-2">
        <h2 class="h2">Setup Books</h2>
        <a class="pill text-xs" href="/" on:click|preventDefault={oncreate}>New</a>
    </div>
    {#if editid == 0}
        <div class="p-2 border-b border-cell mb-2">
            <BookForm book={newbook} userid={userid} on:submit={bookform_done} on:cancel={bookform_done} />
        </div>
    {/if}
    {#each root.books as b (b.bookid)}
        <a class="tblrow" class:sel="{selid == b.bookid}" href="/" on:click|preventDefault="{e => onclickrow(b, b.bookid)}">
            <p class="cell-desc">{b.name}</p>
        </a>
        {#if editid == b.bookid}
        <div class="p-2 border-b border-cell">
            <BookForm book={b} userid={userid} on:submit={bookform_done} on:cancel={bookform_done} />
        </div>
        {/if}
    {/each}
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";
import * as data from "./data.js";
import BookForm from "./BookForm.svelte";

export let root = null;
export let userid = 0;

let selid = 0;
let editid = -1;
let newbook = {
    bookid: 0,
    booktype: data.USERBOOK,
    name: "",
    userid: userid,
};

export function reset() {
    selid = 0;
    editid = -1;
}

function onclickrow(item, itemid) {
    // If row already selected, edit it.
    if (selid == itemid && selid != editid) {
        editid = itemid;
        return;
    }

    // row not selected

    // If edit form is open, just cancel edit without selecting anything.
    if (editid != -1) {
        editid = -1;
        return;
    }

    selid = itemid;
}
function oncreate(e) {
    editid = 0;
    selid = 0;
}
function bookform_done(e) {
    editid = -1;
}

</script>


