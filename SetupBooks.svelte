{#if root == null || root.books.length == 0}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flex flex-row justify-between items-end mb-2">
        <h2 class="text-sm font-bold fg-h1 bg-normal">Setup Books</h2>
        <a class="text-xs pill" href="/" on:click|preventDefault={oncreate}>Add Book</a>
    </div>
    {#if editid == 0}
        <div class="p-2 border-b border-cell mb-2">
            <BookForm book={newbook} userid={userid} on:submit={bookform_done} on:cancel={bookform_done} />
        </div>
    {/if}
    {#each root.books as b (b.bookid)}
        {#if selid == b.bookid && selid != editid}
        <a class="flex flex-row justify-start p-1 border-b border-cell highlight" href="/" on:click|preventDefault="{e => onedit(b)}">
            <p class="flex-grow truncate mr-2">{b.name}</p>
        </a>
        {:else}
        <a class="flex flex-row justify-start p-1 border-b border-cell" href="/" on:click|preventDefault="{e => onselect(b)}">
            <p class="flex-grow truncate mr-2">{b.name}</p>
        </a>
        {/if}
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

export function select(b) {
    // If edit form is open, just cancel edit without selecting anything.
    if (editid != -1) {
        editid = -1;
        return;
    }

    selid = b.bookid;
}
function onselect(b) {
    select(b);
}
function onedit(b) {
    editid = b.bookid;
}
function oncreate(e) {
    editid = 0;
}
function bookform_done(e) {
    editid = -1;
}

</script>


