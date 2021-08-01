import {find, submit} from "./helpers.js";

let svcurl = "/api";

export async function loadRootdata() {
    let sreq = `${svcurl}/rootdata`;
    let [rootdata, err] = await find(sreq);
    if (err != null) {
        return [null, err];
    }

    for (let i=0; i < rootdata.accounts.length; i++) {
        let a = rootdata.accounts[i];
        addFormattedAmts(a);
    }
    return [rootdata, null];
}

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
        let a = aa[i];
        addFormattedAmts(a);
    }
    return [aa, null];
}

export async function loadAccounts() {
    let sreq = `${svcurl}/accounts`;
    return loadAccountsReq(sreq);
}

export async function loadAccountsTxns() {
    let sreq = `${svcurl}/accountstxns`;
    return loadAccountsReq(sreq);
}

export async function loadCurrencies() {
    let sreq = `${svcurl}/currencies`;
    let [cc, err] = await find(sreq);
    if (err != null) {
        return [null, err];
    }
    if (cc == null) {
        cc = [];
        return [cc, null];
    }

    return [cc, null];
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

    addFormattedAmts(a);
    return [a, null];
}

function createFormatter(currency) {
    if (currency == null || currency == "") {
        currency = "USD";
    }
    let formatter;
    try {
        formatter = new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: currency,
            minimumFractionDigits: 2
        });
    } catch(e) {
        formatter = new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: "USD",
            minimumFractionDigits: 2
        });
    }
    return formatter;
}

// Set account.fmtbalance and account's txns'.fmtamt to currency amount format.
export function addFormattedAmts(account) {
    let formatter = createFormatter(account.currency);
    account.fmtbalance = formatter.format(account.balance);

    for (let i=0; i < account.txns.length; i++) {
        let t = account.txns[i];
        if (t.amt > 0) {
            t.fmtamt = formatter.format(t.amt);
        } else {
            // Show negative amt as "(123.45)"
            t.fmtamt = `(${formatter.format(Math.abs(t.amt))})`;
        }
    }
}

export function formattedAmt(amt, currency) {
    let formatter = createFormatter(currency);

    // Show negative amt as "(123.45)"
    if (amt < 0) {
        return `(${formatter.format(Math.abs(amt))})`;
    }
    return formatter.format(amt);
}

