<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="width: 20rem;">
    <h1 class="text-sm font-bold mb-2">Accounts</h1>
    <ul>
{#each ui.accounts as account, i}
    {#if ui.isel == i}
    <li class="flex flex-row justify-between p-1 border-b border-cell highlight">
        <p class="truncate cell-desc">{account.name}</p>
        <p class="fg-dim text-right cell-amt">{account.fmtbalance}</p>
    </li>
    {:else}
    <li class="flex flex-row justify-between p-1 border-b border-cell">
        <p class="truncate cell-desc">{account.name}</p>
        <p class="fg-dim text-right cell-amt">{account.fmtbalance}</p>
    </li>
    {/if}
{/each}
    </ul>
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";
import {currentSession} from "./helpers.js";
let session = currentSession();

export let userid = 0;
if (userid == 0) {
    userid = session.userid;
}

let svcurl = "/api";
let ui = {};
ui.accounts = [];
ui.isel = 0;
ui.status = "";

init(userid);

async function init(userid) {
    ui.status = "";

    let sreq = `${svcurl}/accounts?userid=${userid}`;
    let [aa, err] = await find(sreq);
    if (err != null) {
        console.error(err);
        ui.status = "Server error while fetching accounts";
    }
    if (aa == null) {
        aa = [];
    }

    for (let i=0; i < aa.length; i++) {
        let formatter = new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: aa[i].currency,
            minimumFractionDigits: 2
        });
        aa[i].fmtbalance = formatter.format(aa[i].balance);
    }
    ui.accounts = aa;

    if (aa.length > 0) {
        ui.isel = 0;
    }
}

function dispatchAction(action, entryid) {
    dispatch("action", {
        action: action,
        itemid: entryid,
    });
}

export function onEvent(e) {
    if (e.key == "ArrowUp") {
        ui.isel--;
    } else if (e.key == "ArrowDown") {
        ui.isel++;
    }

    if (ui.isel < 0) {
        ui.isel = 0;
    }
    if (ui.isel > ui.accounts.length-1) {
        ui.isel = ui.accounts.length-1;
    }
}

</script>

