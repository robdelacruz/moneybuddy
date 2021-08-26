<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="width: 20rem;">
{#if rptdata == null}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flex flex-row justify-start items-end mb-4">
        <div class="flex-grow">
            <select class="text-sm font-bold fg-h1 bg-normal pr-2" id="book" name="book" placeholder="Select Book" bind:value={bookid} on:change={onbookchange} on:blur="{e => {}}">
                {#each rptdata.bookrpts as b}
                <option value={b.bookid}>{b.bookname}</option>
                {/each}
            </select>
        </div>
    </div>

    <div class="mb-2">
        <h2 class="text-sm font-bold">Reports Menu</h2>
    </div>

    {#each menurpts as menurpt, i}
        {#if seli == i}
        <a class="flex flex-row justify-between p-1 border-b border-cell highlight" href="/" on:click|preventDefault="{e => onselectrpt(menurpt, i)}">
            <p class="flex-grow truncate mr-2">{menurpt.name}</p>
            <p></p>
        </a>
        {:else}
        <a class="flex flex-row justify-between p-1 border-b border-cell" href="/" on:click|preventDefault="{e => onselectrpt(menurpt, i)}">
            <p class="flex-grow truncate mr-2">{menurpt.name}</p>
            <p></p>
        </a>
        {/if}
    {/each}
{/if}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();

export let rptdata = null;
export let bookid = 1;

let menurpts = [
    {id: "summaryrpt", name: "Summary"},
    {id: "robrpt", name: "Rob Report"}
];
let seli = 0;

function onbookchange(e) {
    dispatch("selectbookid", bookid);
}

function onselectrpt(menurpt, irpt) {
    seli = irpt;
    dispatch("selectrpt", menurpt.id);
}

</script>

