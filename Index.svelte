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
    <Accounts bind:this={waccounts} on:select={accounts_select} widgetstate={ui.accountsstate} accounts={model.accounts}/>
    <Transactions bind:this={wtransactions} on:select={transactions_select} account={ui.activeAccount} widgetstate={ui.transactionsstate} />

    <div class="dim bg-normal fg-normal mb-2 py-2 px-4" style="width: 20rem;">
        <h1 class="text-sm font-bold mb-2">Lorem Ipsum</h1>
        <p class="mb-4">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis consequat est eget est accumsan, eu iaculis tellus fermentum. Nunc pharetra ante feugiat maximus dapibus. Sed dui sapien, hendrerit vel viverra ut, maximus vitae risus. Nunc scelerisque bibendum magna, a faucibus nunc. Suspendisse sapien eros, tincidunt ac ultrices at, hendrerit ac purus. Aliquam erat volutpat. Suspendisse aliquam accumsan ornare. Sed ac leo vitae enim fringilla tristique at sit amet odio. Nunc vel sollicitudin est, vitae commodo purus. Nunc sit amet tellus tincidunt, ultrices quam a, cursus mi. Maecenas sollicitudin vehicula arcu, ut eleifend purus luctus vel. Integer scelerisque mi quis tincidunt laoreet. Fusce dapibus vulputate arcu, eu sollicitudin orci efficitur vel.</p>
        <p class="mb-4">Sed lacinia ligula a quam vehicula porttitor. Ut sollicitudin ante mi, convallis rhoncus dolor egestas non. Phasellus non leo volutpat, consectetur tortor ac, cursus leo. Duis quis tortor quis odio consequat consectetur eget vestibulum tellus. Etiam vel lacus neque. Pellentesque sagittis sodales scelerisque. Nam tristique feugiat enim eu fringilla. Proin lorem quam, ornare sit amet lectus in, aliquam volutpat mauris. Curabitur sed libero turpis. Proin vitae sagittis ex. Morbi semper turpis nisl, sed fermentum dolor pellentesque et.</p>
    </div>
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";
import Accounts from "./Accounts.svelte";
import Transactions from "./Transactions.svelte";

let svcurl = "/api";
let waccounts;
let wtransactions;

let model = {};
model.accounts = [];

let ui = {};
ui.activeAccount = null;
ui.activeTransaction = null;
ui.accountsstate = "";
ui.transactionsstate = "";

$: init();
async function init() {
    let [aa, err] = await loadAccounts();
    model.accounts = aa;
}

document.addEventListener("keydown", function(e) {
//    accounts.onEvent(e);
//    transactions.onEvent(e);
});

async function loadAccounts() {
    let sreq = `${svcurl}/accounts`;
    let [aa, err] = await find(sreq);
    if (err != null) {
        console.error(err);
        return [aa, err];
    }
    if (aa == null) {
        aa = [];
        return [aa, null];
    }

    for (let i=0; i < aa.length; i++) {
        let formatter = new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: aa[i].currency,
            minimumFractionDigits: 2
        });
        aa[i].fmtbalance = formatter.format(aa[i].balance);
    }
    return [aa, null];
}

async function loadTransactions(account) {
    if (account == null) {
        return [null, null];
    }

    let sreq = `${svcurl}/account?accountid=${account.accountid}`;
    let [a, err] = await find(sreq);
    if (err != null) {
        console.error(err);
        return err;
    }
    if (a == null) {
        return null;
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

    account.transactions = a.transactions;
    return null;
}

async function accounts_select(e) {
    let account = e.detail;
    await loadTransactions(account);
    ui.activeAccount = account;
    wtransactions.reset();

    ui.accountsstate = "";
    ui.transactionsstate = "dim";
}

function transactions_select(e) {
    let transaction = e.detail;
    ui.activeTransaction = transaction;

    ui.accountsstate = "dim";
    ui.transactionsstate = "";
}

</script>


