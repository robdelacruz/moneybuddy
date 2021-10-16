<div class="flex flex-row justify-between border-b border-normal fg-normal pb-1 mb-2">
    <div>
        <h1 class="inline text-sm ml-1 mr-2"><a href="/">Money Buddy</a></h1>
        <a href="about.html" class="mr-2">About</a>
    </div>
    <div>
        {#if username != ""}
        <a class="inline mr-1" href="/">{username}</a>
        <a class="inline mr-1" href="/" on:click|preventDefault={onlogout}>Logout</a>
        {:else}
        <a class="inline mr-1" href="/" on:click|preventDefault={onlogin}>Login</a>
        {/if}
    </div>
</div>

{#if mode == ""}
<Tablinks bind:this={tablinks} links="journal|Journal;report|Report;setup|Setup" bind:sel={tabsel} />
    {#if tabsel == "journal"}
        <Journal bind:this={wjournal} root={root} />
    {:else if tabsel == "report"}
        <Report bind:this={wreport} userid={userid} currencies={currencies} />
    {:else if tabsel == "setup"}
        <Setup bind:this={wsetup} userid={userid} root={root} />
    {/if}
{:else if mode == "login"}
<Tablinks bind:this={tablinks} links="login|Login;signup|Sign Up" bind:sel={tabsel} />
    {#if tabsel == "login"}
    <UserLogin bind:this={wuserlogin} on:submit={resetlogin} />
    {:else if tabsel == "signup"}
    <UserSignup bind:this={wusersignup} on:submit={resetlogin} />
    {/if}
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, subscribe, getls, setls} from "./helpers.js";
import * as data from "./data.js";
import Tablinks from "./Tablinks.svelte";
import Journal from "./Journal.svelte";
import Report from "./Report.svelte";
import Setup from "./Setup.svelte";
import UserLogin from "./UserLogin.svelte";
import UserSignup from "./UserSignup.svelte";

let wjournal;
let wreport;
let wsetup;
let wuserlogin;
let wusersignup;

let root = null;
let tablinks;
let tabsel = "";
let mode = "";

let userid = 0;
let username = "";
let sig = "";

let currencies = [];
$: if (root != null) currencies = root.currencies;

[userid, username, sig] = currentSession();

$: resetTabs(userid);
$: subscribeRootdata(userid);

// Remember selections when changed.
$: setls("tabsel", "Index", tabsel);


async function subscribeRootdata(userid) {
    if (userid == 0) {
        root = null;
        return;
    }

    let [rootdata, err] = await data.loadRootdata(userid);
    root = rootdata;

    // Subscribe to data changes.
    let sreq = `/api/subscriberoot?userid=${userid}`;
    console.log(`Subscribing rootdata...`);
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

// Returns [userid, username, sig]
// Reads user cookie of the format:
//   user=<username>|<user signature>;
//   Ex: user=robdelacruz|abc123
function currentSession() {
    let userid = 0;
    let suserid = "";
    let username = "";
    let sig = "";

    let cookies = document.cookie.split(";");
    for (let i=0; i < cookies.length; i++) {
        let cookie = cookies[i].trim();
        let [k,v] = cookie.split("=");
        if (k != "user") {
            continue;
        }
        if (v == undefined) {
            v = "";
        }
        [suserid, username, sig] = v.split("|");
        if (suserid == "") {
            userid = 0;
        } else {
            userid = parseInt(suserid, 10);
            if (userid == NaN) {
                userid = 0;
            }
        }
        if (username == undefined) {
            username = "";
        }
        if (sig == undefined) {
            sig = "";
        }
        break;
    }
    return [userid, username, sig];
}

function onlogin(e) {
    mode = "login";
    tabsel = "login";
}
function onlogout(e) {
    userid = "";
    username = "";
    sig = "";
    document.cookie = `user=;path=/`;

    mode = "login";
    tabsel = "login";
}

function resetlogin() {
    localStorage.clear();
    [userid, username, sig] = currentSession();
}

function resetTabs(userid) {
    // Not logged in, show login prompt.
    if (userid == 0) {
        mode = "login";
        tabsel = "login";
        return;
    }

    // Logged in, show previously selected tab ("journal" is default).
    mode = "";
    tabsel = getls("tabsel", "Index", "journal");
}

</script>


