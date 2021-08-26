<div class="flex flex-row">
    <ReportsMenu bind:this={wreportsmenu} rptdata={rptdata} bookid={selbookid} on:selectbookid={reportsmenu_selectbookid} on:selectrpt={reportsmenu_selectrpt} />
    {#if selrptid == "summaryrpt"}
        <SummaryRpt bind:this={wsummaryrpt} rptdata={rptdata} bookid={selbookid} />
    {:else if selrptid == "somethingelse"}
        <p>something else</p>
    {/if}

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

let rptdata = null;
let wreportsmenu;
let wsummaryrpt;

let selbookid = 1;
let selrptid = "";

init();
async function init() {
    let [d, err] = await data.loadRptdata();
    rptdata = d;
}

function reportsmenu_selectbookid(e) {
    selbookid = e.detail;
}
function reportsmenu_selectrpt(e) {
    selrptid = e.detail;
    console.log(selrptid);
}


</script>

