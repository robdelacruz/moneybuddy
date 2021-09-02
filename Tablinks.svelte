<div class="flex flex-row fg-normal text-sm">
{#each tabs as tab}
    {#if tab.id == sel}
        <a href="/" class="rounded-t-md py-1 px-3 bg-normal" on:click|preventDefault='{e => onlink(tab.id)}'>{tab.caption}</a>
    {:else}
        <a href="/" class="py-1 px-3" on:click|preventDefault='{e => onlink(tab.id)}'>{tab.caption}</a>
    {/if}
{/each}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();

// Ex. "entries|Entries;images|Images;files|Files"
export let links = "";
export let sel = "";
let tabs = [];

let ll = links.split(";");
for (let i=0; i < ll.length; i++) {
    let ss = ll[i].split("|");
    if (ss.length != 2) {
        continue;
    }
    let id = ss[0].trim();
    let caption = ss[1];
    tabs.push({id: id, caption: caption});

    // Make the first link active by default.
    if (sel == "") {
        sel = id;
    }
}

function onlink(id) {
    sel = id;
    dispatch("sel", id);
}

</script>

