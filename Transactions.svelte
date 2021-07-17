<div class:dim="{widgetstate == 'dim'}" class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="width: 40rem;">
{#if accountid == 0}
    <p class="fg-dim">Select Account</p>
{:else if ui.account == null}
    <p class="fg-dim">Account not found</p>
{:else}
    <h1 class="text-sm font-bold mb-2">{ui.account.name}: Transactions</h1>
    {#each ui.account.transactions as t, i}
    <a href="/" on:click|preventDefault="{e => onseltransaction(t, i)}">
        {#if ui.isel == i}
            <div class="flex flex-row flex-start p-1 border-b border-cell highlight">
                <p class="highlight cell-date">{t.date.substring(0, 10)}</p>
                <p class="highlight truncate cell-desc">{t.desc}</p>
                {#if t.amt >= 0}
                <p class="highlight text-right cell-amt">{t.fmtamt}</p>
                <p class="highlight text-right cell-amt"></p>
                {:else}
                <p class="highlight text-right cell-amt"></p>
                <p class="highlight text-right cell-amt">{t.fmtamt}</p>
                {/if}
            </div>
        {:else}
            <div class="flex flex-row flex-start p-1 border-b border-cell">
                <p class="fg-dim cell-date">{t.date.substring(0, 10)}</p>
                <p class="truncate cell-desc">{t.desc}</p>
                {#if t.amt >= 0}
                <p class="fg-dim text-right cell-amt">{t.fmtamt}</p>
                <p class="fg-dim text-right cell-amt"></p>
                {:else}
                <p class="fg-dim text-right cell-amt"></p>
                <p class="fg-dim text-right cell-amt">{t.fmtamt}</p>
                {/if}
            </div>
        {/if}
    </a>
    {/each}
{/if}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";

export let widgetstate = "";
export let accountid = 0;

let svcurl = "/api";
let ui = {};
ui.account = null;
ui.isel = -1;
ui.status = "";

$: init(accountid);

async function init(accountid) {
    console.log("init()");
    ui.status = "";
    ui.account = null;

    if (accountid == 0) {
        return;
    }

    let sreq = `${svcurl}/account?accountid=${accountid}`;
    let [a, err] = await find(sreq);
    if (err != null) {
        console.error(err);
        ui.status = "Server error while fetching account";
        return;
    }
    if (a == null) {
        ui.status = "Account not found";
        return;
    }

    let formatter = new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: a.currency,
        minimumFractionDigits: 2
    });

    for (let i=0; i < a.transactions.length; i++) {
        let t = a.transactions[i];
        if (t.amt > 0) {
            a.transactions[i].fmtamt = formatter.format(t.amt);
        } else {
            // Show negative amt as "(123.45)"
            a.transactions[i].fmtamt = `(${formatter.format(Math.abs(t.amt))})`;
        }
    }
    ui.account = a;
    ui.isel = -1;
}

function dispatchAction(action, entryid) {
    dispatch("action", {
        action: action,
        itemid: entryid,
    });
}

export function onEvent(e) {
    if (ui.account == null || ui.account.transactions.length == 0) {
        return;
    }

    if (e.key == "ArrowUp") {
        ui.isel--;
    } else if (e.key == "ArrowDown") {
        ui.isel++;
    }

    if (ui.isel < 0) {
        ui.isel = 0;
    }
    if (ui.isel > ui.account.transactions.length-1) {
        ui.isel = ui.account.transactions.length-1;
    }
}

function onseltransaction(transaction, i) {
    ui.isel = i;
    dispatch("select", transaction);
}

</script>
