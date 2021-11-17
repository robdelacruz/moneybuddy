<div class="journal-container">
    <Accounts bind:this={waccounts} root={root} bookid={selbookid} selaccountid={selaccountid} on:selectbookid={accounts_selectbookid} on:selectaccount={accounts_selectaccount} />
    <Txns bind:this={wtxns} root={root} bookid={selbookid} accountid={selaccountid} on:selectaccount={txns_selectaccount} on:selecttxn={txns_selecttxn} />
    <TxnDetail txn={seltxn} />
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, getlsInt, setlsInt} from "./helpers.js";

import Accounts from "./Accounts.svelte";
import Txns from "./Txns.svelte";
import TxnDetail from "./TxnDetail.svelte";

export let root = null;
let waccounts;
let wtxns;

let selbookid = getlsInt("Journal", "selbookid", firstbookid(root));
let selaccountid = getlsInt("Journal", "selaccountid", 0);
let seltxn = null;

$: selbookid = getlsInt("Journal", "selbookid", firstbookid(root));

// Refresh seltxn whenever root data changed.
$: if (seltxn != null) {
    seltxn = findTxn(root, selbookid, selaccountid, seltxn.txnid);
}

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

// Look for txn in rootdata. 
function findTxn(rootdata, bookid, accountid, txnid) {
    if (txnid == 0 || bookid == 0 || accountid == 0) {
        return null;
    }
    let book = null;
    let account = null;
    for (let i=0; i < rootdata.books.length; i++) {
        if (bookid == rootdata.books[i].bookid) {
            book = rootdata.books[i];
            break;
        }
    }
    if (book == null) {
        return null;
    }
    for (let i=0; i < book.bankaccounts.length; i++) {
        if (accountid == book.bankaccounts[i].accountid) {
            account = book.bankaccounts[i];
            break;
        }
    }
    if (account == null) {
        for (let i=0; i < book.stockaccounts.length; i++) {
            if (accountid == book.stockaccounts[i].accountid) {
                account = book.stockaccounts[i];
                break;
            }
        }
    }
    if (account == null) {
        return null;
    }
    for (let i=0; i < account.txns.length; i++) {
        if (txnid == account.txns[i].txnid) {
            return account.txns[i];
        }
    }
    return null;
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
function txns_selecttxn(e) {
    let t = e.detail;
    seltxn = t;
}

</script>
