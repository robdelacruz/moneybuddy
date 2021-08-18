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

<div class="flex flex-row">
    <Accounts bind:this={waccounts} on:select={accounts_select} root={root}/>
    <Txns bind:this={wtxns} account={ui.activeAccount} />

    <div class="dim bg-normal fg-normal mb-2 py-2 px-4" style="width: 20rem;">
        <h1 class="text-sm font-bold mb-2">Lorem Ipsum</h1>
        <p class="mb-4">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis consequat est eget est accumsan, eu iaculis tellus fermentum. Nunc pharetra ante feugiat maximus dapibus. Sed dui sapien, hendrerit vel viverra ut, maximus vitae risus. Nunc scelerisque bibendum magna, a faucibus nunc. Suspendisse sapien eros, tincidunt ac ultrices at, hendrerit ac purus. Aliquam erat volutpat. Suspendisse aliquam accumsan ornare. Sed ac leo vitae enim fringilla tristique at sit amet odio. Nunc vel sollicitudin est, vitae commodo purus. Nunc sit amet tellus tincidunt, ultrices quam a, cursus mi. Maecenas sollicitudin vehicula arcu, ut eleifend purus luctus vel. Integer scelerisque mi quis tincidunt laoreet. Fusce dapibus vulputate arcu, eu sollicitudin orci efficitur vel.</p>
        <p class="mb-4">Sed lacinia ligula a quam vehicula porttitor. Ut sollicitudin ante mi, convallis rhoncus dolor egestas non. Phasellus non leo volutpat, consectetur tortor ac, cursus leo. Duis quis tortor quis odio consequat consectetur eget vestibulum tellus. Etiam vel lacus neque. Pellentesque sagittis sodales scelerisque. Nam tristique feugiat enim eu fringilla. Proin lorem quam, ornare sit amet lectus in, aliquam volutpat mauris. Curabitur sed libero turpis. Proin vitae sagittis ex. Morbi semper turpis nisl, sed fermentum dolor pellentesque et.</p>
    </div>
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, subscribe} from "./helpers.js";
import * as data from "./data.js";
import Accounts from "./Accounts.svelte";
import Txns from "./Txns.svelte";

let waccounts;
let wtxns;

let root = null;

let ui = {};
ui.activeAccountid = 0;
ui.activeAccount = null;

$: init();
async function init() {
    let [rootdata, err] = await data.loadRootdata();
    root = rootdata;

    //$$: Complex code below - needs rework.

    // Subscribe to data changes and update accounts and txns.
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

        ui.activeAccount = null;
        for (let i=0; i < rootdata.books.length; i++) {
            let b = rootdata.books[i];
            data.formatBookAmts(b);

            for (let j=0; j < b.bankaccounts.length; j++) {
                let a = b.bankaccounts[j];
                if (a.accountid == ui.activeAccountid) {
                    ui.activeAccount = a;
                    break;
                }
            }
            for (let j=0; j < b.stockaccounts.length; j++) {
                let a = b.stockaccounts[j];
                if (a.accountid == ui.activeAccountid) {
                    ui.activeAccount = a;
                    break;
                }
            }
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

function resetTxns() {
    wtxns.reset();
}

document.addEventListener("keyup", function(e) {
    waccounts.postEvent(e);
});

async function accounts_select(e) {
    let err;
    ui.activeAccount = e.detail;
    if (ui.activeAccount != null) {
        ui.activeAccountid = ui.activeAccount.accountid;
    }
    resetTxns();
}

</script>


