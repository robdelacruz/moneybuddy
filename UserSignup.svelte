<form class="" on:submit|preventDefault={onSubmit}>
    <div class="mb-2">
        <input class="block bg-input fg-normal py-1 px-2 w-full" name="username" id="username" type="text" placeholder="Username" bind:value={frm_username}>
    </div>
    <div class="mb-2">
        <input class="block bg-input fg-normal py-1 px-2 w-full" name="password" id="password" type="password" placeholder="Password" bind:value={frm_password}>
    </div>
    <div class="mb-2">
        <input class="block bg-input fg-normal py-1 px-2 w-full" name="password2" id="password2" type="password" placeholder="Re-enter Password" bind:value={frm_password2}>
    </div>
    <div class="flex flex-row justify-left">
        <div>
            <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2">Sign Up</button>
            <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault="{e => dispatch('cancel')}">Cancel</a>
        </div>
    </div>
    {#if status != ""}
    <div class="">
        <p class="uppercase italic text-xs">{status}</p>
    </div>
    {/if}
</form>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, del} from "./helpers.js";
import * as data from "./data.js";

let svcurl = "/api";
let mode = "";
let status = "";
let frm_username = "";
let frm_password = "";
let frm_password2 = "";
let frm_newpassword = "";

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    status = "processing";

    let signupreq = {};
    signupreq.username = frm_username;
    signupreq.password = frm_password;

    let sreq = `${svcurl}/signup`;
    let [result, err] = await submit(sreq, "POST", signupreq);
    if (err != null) {
        console.error(err);
        status = err;
        return;
    }
    if (result.error != null) {
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

