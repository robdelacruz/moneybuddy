import {find, submit} from "./helpers.js";

let svcurl = "/api";

export let BANKACCOUNT = 0;
export let STOCKACCOUNT = 1;

export let USERBOOK = 0;
export let SYSTEMBOOK = 1;

export async function loadRootdata(userid) {
    let sreq = `${svcurl}/rootdata?userid=${userid}`;
    let [rootdata, err] = await find(sreq);
    if (err != null) {
        return [null, err];
    }

    for (let i=0; i < rootdata.books.length; i++) {
        let b = rootdata.books[i];
        formatBookAmts(b);
    }
    return [rootdata, null];
}
export function formatBookAmts(b) {
    for (let i=0; i < b.bankaccounts.length; i++) {
        addFormattedAmts(b.bankaccounts[i]);
    }
    for (let i=0; i < b.stockaccounts.length; i++) {
        addFormattedAmts(b.stockaccounts[i]);
    }
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
    if (account.accounttype == 0) {
        addCurrencyFormattedAmts(account);
    } else {
        addUnitFormattedAmts(account);
    }
}
// Add 'fmtamt' property showing currency amount.
function addCurrencyFormattedAmts(account) {
    let formatter = createFormatter(account.currency.name);
    account.fmtbalance = formatter.format(account.balance);

    for (let i=0; i < account.txns.length; i++) {
        let t = account.txns[i];
        t.fmtamt = formatter.format(t.amt);
    }
}
// Add 'fmtamt' property showing non-currency amount.
function addUnitFormattedAmts(account) {
    let formatter = createFormatter(account.currency.name);
    account.fmtbalance = formatter.format(account.balance);

    for (let i=0; i < account.txns.length; i++) {
        let t = account.txns[i];
        t.fmtamt = t.amt.toLocaleString("en-US", {minimumFractionDigits: 2});
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

export async function loadRptdata(userid, currencyid) {
    let sreq = `${svcurl}/rptdata?userid=${userid}&currencyid=${currencyid}`;
    let [rptdata, err] = await find(sreq);
    return [rptdata, null];
}

