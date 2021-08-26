<div class="flex flex-row justify-between border-b border-normal fg-normal pb-1 mb-2">
    <div>
        <h1 class="inline text-sm ml-1 mr-2"><a href="/">Accounts Buddy</a></h1>
        <a href="about.html" class="mr-2">About</a>
    </div>
    <div>
        <a class="inline mr-1" href="/">robdelacruz</a>
        <a class="inline mr-1" href="/">Logout</a>
    </div>
</div>

<Tablinks bind:this={tablinks} links="journal|Journal;report|Report;setup|Setup" sel="journal" on:sel={tablinks_sel} />

{#if tabsel == "journal"}
    <Journal bind:this={wjournal} root={root} />
{:else if tabsel == "report"}
    <Report bind:this={wreport} />
{:else if tabsel == "setup"}
    <p class="fg-normal">setup</p>
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, subscribe} from "./helpers.js";
import * as data from "./data.js";
import Tablinks from "./Tablinks.svelte";
import Journal from "./Journal.svelte";
import Report from "./Report.svelte";

let root = null;
let tablinks;
let tabsel = "journal";

let wjournal;
let wreport;

init();
async function init() {
    let [rootdata, err] = await data.loadRootdata();
    root = rootdata;

    // Subscribe to data changes.
    let qwho = getqs("who");
    let sreq = `/api/subscriberoot?who=${qwho}`;
    console.log(`Subscribing (${qwho})...`);
    await subscribe(sreq, "json", function(rootdata, err) {
        if (err != null) {
            console.log("Error receiving root data...");
            console.error(err);
            return;
        }

        console.log("Received root data...");
        console.log(rootdata);

        for (let i=0; i < rootdata.books.length; i++) {
            let b = rootdata.books[i];
            data.formatBookAmts(b);
        }

        root = rootdata;
    });
}

function getqs(q) {
    let url = new URL(document.location.href);
    let v = url.searchParams.get(q);
    if (v == null) {
        v = "";
    }
    return v;
}

document.addEventListener("keyup", function(e) {
    if (wjournal) {
        wjournal.postEvent(e);
    }
});

function tablinks_sel(e) {
    tabsel = e.detail;
}

</script>


