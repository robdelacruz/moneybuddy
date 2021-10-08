<div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="width: 15rem;">
    <h1 class="text-sm font-bold fg-h1 mb-4">User Login</h1>
    <form class="" autocomplete="off" on:submit|preventDefault={onSubmit}>
        <div class="mb-2">
            <input class="block bg-input fg-normal py-1 px-2 w-full" name="username" id="username" type="text" placeholder="Username" bind:value={frm_username}>
        </div>
        <div class="mb-4">
            <input class="block bg-input fg-normal py-1 px-2 w-full" name="password" id="password" type="password" placeholder="Password" bind:value={frm_password}>
        </div>
        <div class="flex flex-row justify-center mb-2">
            <button class="border border-normal py-1 px-2 bg-inputok mr-2">Login</button>
            <a href="/" class="hidden action self-center" on:click|preventDefault="{e => dispatch('cancel')}">Cancel</a>
        </div>
        {#if status != ""}
        <div class="">
            <p class="uppercase italic text-xs">{status}</p>
        </div>
        {/if}
    </form>
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, del} from "./helpers.js";
import * as data from "./data.js";

let svcurl = "/api";
let status = "";
let frm_username = "";
let frm_password = "";
let frm_password2 = "";

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    status = "processing";

    let loginreq = {};
    loginreq.username = frm_username;
    loginreq.password = frm_password;

    let sreq = `${svcurl}/login`;
    let [result, err] = await submit(sreq, "POST", loginreq);
    if (err != null) {
        console.error(err);
        status = err;
        return;
    }
    if (result.error != "") {
        console.error(result.error);
        status = result.error;
        return;
    }

    document.cookie = `user=${result.userid}|${result.username}|${result.sig};path=/`;
    dispatch("submit");
}

function onCancel(e) {
    dispatch("cancel");
}

</script>

