<div class="report bg-normal fg-normal py-2 px-4">
{#if rptdata == null || rptdata.bookrpts.length == 0}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flexrow justify-between border-b border-cell mb-4">
        <div class="flexrow gap-4">
            <select class="text-xs" id="rpttype" name="rpttype" placeholder="Select Report" bind:value={selrptid}>
                {#each menurpts as menurpt}
                    {#if menurpt.id == selrptid}
                    <option selected value={menurpt.id}>{menurpt.name}</option>
                    {:else}
                    <option value={menurpt.id}>{menurpt.name}</option>
                    {/if}
                {/each}
            </select>
            <select class="text-xs" id="book" name="book" placeholder="Select Book" bind:value={selbookid}>
                {#each rptdata.bookrpts as bookrpt}
                    {#if bookrpt.bookid == selbookid}
                    <option selected value={bookrpt.bookid}>{bookrpt.bookname}</option>
                    {:else}
                    <option value={bookrpt.bookid}>{bookrpt.bookname}</option>
                    {/if}
                {/each}
            </select>
        </div>
        <div>
            <select class="text-xs" id="currency" name="currency" placeholder="Select Currency" bind:value={selcurrencyid}>
                {#each currencies as c}
                    {#if c.currencyid == selcurrencyid}
                    <option selected value={c.currencyid}>{c.name}</option>
                    {:else}
                    <option value={c.currencyid}>{c.name}</option>
                    {/if}
                {/each}
            </select>
        </div>
    </div>

    {#if selrptid == "summaryrpt"}
    <SummaryRpt bind:this={wsummaryrpt} rptdata={rptdata} bookid={selbookid} />
    {:else if selrptid == "somethingelse"}
        <p>something else</p>
    {/if}
{/if}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, getls, getlsInt, setls, setlsInt} from "./helpers.js";
import * as data from "./data.js";
import SummaryRpt from "./SummaryRpt.svelte";

export let userid = 0;
export let currencies = [];
let rptdata = null;
let wsummaryrpt;

let menurpts = [
    {id: "summaryrpt", name: "Summary Report"},
    {id: "robrpt", name: "Rob Report"}
];

let selrptid = getls("Report", "selrptid", menurpts[0].id);
let selbookid = 0;
let selcurrencyid = getlsInt("Report", "selcurrencyid", firstcurrencyid(currencies));

// Remember selections when changed.
$: setls("Report", "selrptid", selrptid);
$: if (selbookid != 0) setlsInt("Report", "selbookid", selbookid);
$: setlsInt("Report", "selcurrencyid", selcurrencyid);

$: loadrptdata(userid, selcurrencyid);

async function loadrptdata(userid, currencyid) {
    let [d, err] = await data.loadRptdata(userid, currencyid);
    rptdata = d;
    selbookid = getlsInt("Report", "selbookid", firstbookid(rptdata));
}

function firstbookid(rptdata) {
    if (rptdata == null || rptdata.bookrpts.length == 0) {
        return 0;
    }
    return rptdata.bookrpts[0].bookid;
}

function firstcurrencyid(currencies) {
    if (currencies == null || currencies.length == 0) {
        return 0;
    }
    return currencies[0].currencyid;
}

</script>

