<div class="journal-container">
    <Accounts bind:this={waccounts} root={root} bookid={selbookid} selaccountid={selaccountid} on:selectbookid={accounts_selectbookid} on:selectaccount={accounts_selectaccount} />
    <Txns bind:this={wtxns} root={root} bookid={selbookid} accountid={selaccountid} on:selectaccount={txns_selectaccount}/>
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, getlsInt, setlsInt} from "./helpers.js";

import Accounts from "./Accounts.svelte";
import Txns from "./Txns.svelte";

export let root = null;
let waccounts;
let wtxns;

let selbookid = getlsInt("Journal", "selbookid", firstbookid(root));
let selaccountid = getlsInt("Journal", "selaccountid", 0);

$: selbookid = getlsInt("Journal", "selbookid", firstbookid(root));

function firstbookid(rootdata) {
    if (rootdata == null || rootdata.books.length == 0) {
        return 0;
    }
    // Find first active book's id.
    let bookid = 0;
    for (let i=0; i < rootdata.books.length; i++) {
        if (rootdata.books[i].active == 1) {
            bookid = rootdata.books[i].bookid;
            break;
        }
    }
    return bookid;
}

function accounts_selectbookid(e) {
    selbookid = e.detail;
    selaccountid = 0;
    wtxns.reset();

    setlsInt("Journal", "selbookid", selbookid);
    setlsInt("Journal", "selaccountid", selaccountid);
}

function accounts_selectaccount(e) {
    let a = e.detail;
    selaccountid = a.accountid;
    setlsInt("Journal", "selaccountid", selaccountid);
}

function txns_selectaccount(e) {
    let a = e.detail;
    selaccountid = a.accountid;
    setlsInt("Journal", "selaccountid", selaccountid);

    waccounts.selectAccount(a);
}

</script>
