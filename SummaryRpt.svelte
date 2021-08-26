<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="width: 40rem;">
{#if rptdata == null || selbookrpt == null}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flex flex-row justify-between mb-2">
        <h1 class="text-sm font-bold fg-h1 bg-normal">{selbookrpt.summaryrpt.heading}</h1>
    </div>

    {#each selbookrpt.summaryrpt.rptitems as ri}
    <div class="flex flex-row flex-start p-1 border-b border-cell">
        <p class="cell-desc">{ri.caption}</p>
        <p class="text-right cell-amt">{ri.val}</p>
    </div>
    {/each}
{/if}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();

export let rptdata = null;
export let bookid = 1;

let selbookrpt = null;

// rptdata + bookid --> selbookrpt
$: selbookrpt = getSelectedBookRpt(rptdata, bookid);

function getSelectedBookRpt(rptdata, bookid) {
    if (rptdata == null) {
        return null;
    }
    let bookrpt = null;
    for (let i=0; i < rptdata.bookrpts.length; i++) {
        if (bookid == rptdata.bookrpts[i].bookid) {
            bookrpt = rptdata.bookrpts[i];
            break;
        }
    }
    return bookrpt
}

</script>

