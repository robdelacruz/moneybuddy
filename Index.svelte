<div class="maincontainer p-2">
    <div class="header flex flex-row justify-between border-b border-normal fg-normal pb-1 mb-2">
        <div class="flex flex-row items-center">
            <div class="fg-h2">
                <svg class="h-4 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M10 20a10 10 0 1 1 0-20 10 10 0 0 1 0 20zm1-5h1a3 3 0 0 0 0-6H7.99a1 1 0 0 1 0-2H14V5h-3V3H9v2H8a3 3 0 1 0 0 6h4a1 1 0 1 1 0 2H6v2h3v2h2v-2z"/></svg>
            </div>
            <h1 class="font-bold text-sm ml-1 mr-4">Money Buddy</h1>
            <a href="about.html" class="self-end mr-2">About</a>
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

    <div class="main py-2">
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
    </div>

    <div class="footer flex flex-row text-xs fg-dim p-4">
        <a class="mr-2" href="https://github.com/robdelacruz/moneybuddy" target="_blank">source</a>
        <a class="" href="https://github.com/robdelacruz/moneybuddy" target="_blank">donate</a>
    </div>
</div>

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
$: setls("Index", "tabsel", tabsel);


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
    tabsel = getls("Index", "tabsel", "journal");
}

</script>


