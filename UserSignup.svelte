<div class="userpanel webpanel">
    <h1 class="h2 mb-4">Create Your Account</h1>
    <form class="" autocomplete="off" on:submit|preventDefault={onSubmit}>
        <div class="mb-2">
            <input class="input w-full" name="username" id="username" type="text" placeholder="Username" bind:value={frm_username}>
        </div>
        <div class="mb-2">
            <input class="input w-full" name="password" id="password" type="password" placeholder="Password" bind:value={frm_password}>
        </div>
        <div class="mb-4">
            <input class="input w-full" name="password2" id="password2" type="password" placeholder="Re-enter Password" bind:value={frm_password2}>
        </div>
        <div class="flexrow justify-center mb-2">
            <button class="btn bg-inputok mr-2">Sign Up</button>
            <a href="/" class="hidden action self-center" on:click|preventDefault="{e => dispatch('cancel')}">Cancel</a>
        </div>
        {#if status != ""}
        <div class="">
            <p class="prompt">{status}</p>
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
    if (frm_password != frm_password2) {
        status = "passwords don't match";
        return;
    }

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

