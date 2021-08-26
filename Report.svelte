<div class="flex flex-row">
    <div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="width: 40rem;">
    {#if rptdata == null}
        <p class="fg-dim">No data</p>
    {:else}
        <div class="flex flex-row justify-between border-b border-cell mb-4">
            <div>
                <select class="text-xs fg-normal bg-normal pr-1 mr-4" id="rpttype" name="rpttype" placeholder="Select Report" bind:value={selrptid}>
                    {#each menurpts as menurpt}
                    <option value={menurpt.id}>{menurpt.name}</option>
                    {/each}
                </select>
                <select class="text-xs fg-normal bg-normal pr-1" id="book" name="book" placeholder="Select Book" bind:value={selbookid}>
                    {#each rptdata.bookrpts as bookrpt}
                    <option value={bookrpt.bookid}>{bookrpt.bookname}</option>
                    {/each}
                </select>
            </div>
            <div>
                <select class="text-xs fg-normal bg-normal pr-1" id="currency" name="currency" placeholder="Select Currency" bind:value={selcurrencyid}>
                    {#each currencies as c}
                    <option value={c.currencyid}>{c.currency}</option>
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
import * as data from "./data.js";
import ReportsMenu from "./ReportsMenu.svelte";
import SummaryRpt from "./SummaryRpt.svelte";

export let currencies = [];
let rptdata = null;
let wsummaryrpt;

let menurpts = [
    {id: "summaryrpt", name: "Summary Report"},
    {id: "robrpt", name: "Rob Report"}
];

let selbookid = 1;
let selrptid = menurpts[0].id;
let selcurrencyid = 1;

$: init(selcurrencyid);

async function init(currencyid) {
    let [d, err] = await data.loadRptdata(currencyid);
    rptdata = d;
}

</script>

