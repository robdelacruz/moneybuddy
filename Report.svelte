<div class="flex flex-row">
    <div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="width: 40rem;">
    {#if rptdata == null || rptdata.bookrpts.length == 0}
        <p class="fg-dim">No data</p>
    {:else}
        <div class="flex flex-row justify-between border-b border-cell mb-4">
            <div>
                <select class="text-xs fg-normal bg-normal pr-1 mr-4" id="rpttype" name="rpttype" placeholder="Select Report" bind:value={selrptid}>
                    {#each menurpts as menurpt}
                        {#if menurpt.id == selrptid}
                        <option selected value={menurpt.id}>{menurpt.name}</option>
                        {:else}
                        <option value={menurpt.id}>{menurpt.name}</option>
                        {/if}
                    {/each}
                </select>
                <select class="text-xs fg-normal bg-normal pr-1" id="book" name="book" placeholder="Select Book" bind:value={selbookid}>
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
                <select class="text-xs fg-normal bg-normal pr-1" id="currency" name="currency" placeholder="Select Currency" bind:value={selcurrencyid}>
                    {#each currencies as c}
                        {#if c.currencyid == selcurrencyid}
                        <option selected value={c.currencyid}>{c.currency}</option>
                        {:else}
                        <option value={c.currencyid}>{c.currency}</option>
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

    <div class="dim bg-normal fg-normal mb-2 py-2 px-4" style="width: 20rem;">
        <h1 class="text-sm font-bold mb-2">Lorem Ipsum</h1>
        <p class="mb-4">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis consequat est eget est accumsan, eu iaculis tellus fermentum. Nunc pharetra ante feugiat maximus dapibus. Sed dui sapien, hendrerit vel viverra ut, maximus vitae risus. Nunc scelerisque bibendum magna, a faucibus nunc. Suspendisse sapien eros, tincidunt ac ultrices at, hendrerit ac purus. Aliquam erat volutpat. Suspendisse aliquam accumsan ornare. Sed ac leo vitae enim fringilla tristique at sit amet odio. Nunc vel sollicitudin est, vitae commodo purus. Nunc sit amet tellus tincidunt, ultrices quam a, cursus mi. Maecenas sollicitudin vehicula arcu, ut eleifend purus luctus vel. Integer scelerisque mi quis tincidunt laoreet. Fusce dapibus vulputate arcu, eu sollicitudin orci efficitur vel.</p>
    </div>
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, getls, setls} from "./helpers.js";
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

let selbookid = getls("selbookid", "Report", firstbookid(rptdata));
let selrptid = getls("selrptid", "Report", menurpts[0].id);
let selcurrencyid = getls("selcurrencyid", "Report", firstcurrencyid(currencies));

// Remember selections when changed.
$: setls("selbookid", "Report", selbookid);
$: setls("selrptid", "Report", selrptid);
$: setls("selcurrencyid", "Report", selcurrencyid);

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

$: init(userid, selcurrencyid);

async function init(userid, currencyid) {
    let [d, err] = await data.loadRptdata(userid, currencyid);
    rptdata = d;
}

</script>

