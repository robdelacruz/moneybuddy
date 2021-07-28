import {find, submit} from "./helpers.js";

let svcurl = "/api";

export async function loadAccountsReq(sreq) {
    let [aa, err] = await find(sreq);
    if (err != null) {
        return [null, err];
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

export async function loadAccounts() {
    let sreq = `${svcurl}/accounts`;
    return loadAccountsReq(sreq);
}

export async function loadAccountsTxns() {
    let sreq = `${svcurl}/accountstxns`;
    let [aa, err] = await loadAccountsReq(sreq);
    if (err != null) {
        return [null, err];
    }

    for (let i=0; i < aa.length; i++) {
        let a = aa[i];
        addFormattedAmts(a.txns, a.currency);
    }
    return [aa, null];
}

function addFormattedAmts(tt, currency) {
    let formatter = new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: currency,
        minimumFractionDigits: 2
    });

    for (let i=0; i < tt.length; i++) {
        let t = tt[i];
        if (t.amt > 0) {
            t.fmtamt = formatter.format(t.amt);
        } else {
            // Show negative amt as "(123.45)"
            t.fmtamt = `(${formatter.format(Math.abs(t.amt))})`;
        }
    }
}

// request account and its transactions
export async function loadAccount(accountid) {
    let sreq = `${svcurl}/account?id=${accountid}`;
    let [a, err] = await find(sreq);
    if (err != null) {
        return [null, err];
    }
    if (a == null) {
        return [null, null];
    }

    addFormattedAmts(a.txns, a.currency);
    return [a, null];
}

