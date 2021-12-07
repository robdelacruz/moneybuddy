{#if rptdata == null || selbookrpt == null}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flex flex-row justify-between mb-2">
        <h1 class="text-sm font-bold fg-h2">{selbookrpt.summaryrpt.heading}</h1>
    </div>

    {#each selbookrpt.summaryrpt.rptitems as ri}
        {#if ri.caption == ""}
            <div class="flex flex-row py-2">
            </div>
        {:else if ri.caption.startsWith("# ")}
            <div class="flex flex-row p-1">
                <h2 class="cell-desc font-bold fg-h2">{ri.caption.substring(2)}</h2>
            </div>
        {:else}
        <div class="flex flex-row justify-between p-1 border-b border-cell">
            <p class="cell-desc">{ri.caption}</p>
            {#if ri.currencyamt.amt >= 0}
            <p class="whitespace-nowrap fg-number-plus text-right cell-amt">{data.formattedAmt(ri.currencyamt.amt, ri.currencyamt.currency.currency)}</p>
            {:else}
            <p class="whitespace-nowrap fg-number-minus text-right cell-amt">{data.formattedAmt(ri.currencyamt.amt, ri.currencyamt.currency.currency)}</p>
            {/if}
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

