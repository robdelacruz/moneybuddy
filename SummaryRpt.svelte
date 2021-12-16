{#if rptdata == null || selbookrpt == null}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flexrow mb-2">
        <h1 class="text-sm font-bold fg-h2">{selbookrpt.summaryrpt.heading}</h1>
    </div>

    {#each selbookrpt.summaryrpt.rptitems as ri}
        {#if ri.caption == "" && ri.cols.length == 0}
            <div class="py-2"></div>
        {:else}
            <div class="rptrow">
            {#if ri.caption.startsWith("# ")}
                <h2 class="cell-desc font-bold fg-h2">{ri.caption.substring(2)}</h2>
            {:else}
                <p class="cell-desc">{ri.caption}</p>
            {/if}
            {#each ri.cols as col}
                {#if col.amt}
                    {#if col.amt >= 0}
                    <p class="cell-amt whitespace-nowrap fg-number-plus">{data.formattedAmt(col.amt, col.currencyname)}</p>
                    {:else}
                    <p class="cell-amt whitespace-nowrap fg-number-minus">{data.formattedAmt(col.amt, col.currencyname)}</p>
                    {/if}
                {:else if typeof col == "string"}
                    <p class="cell-caption fg-h3 text-right">{col}</p>
                {/if}
            {/each}
            </div>
        {/if}
    {/each}
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import * as data from "./data.js";

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

