{#if userid == 0}
    <p class="fg-dim">Select User</p>
{:else}
    <form class="" on:submit|preventDefault={onSubmit}>
        <div class="mb-2">
            <input class="block bg-input fg-normal py-1 px-2 w-full" name="password" id="password" type="password" placeholder="Current Password" bind:value={frm_password}>
        </div>
        <div class="mb-2">
            <input class="block bg-input fg-normal py-1 px-2 w-full" name="newpassword" id="newpassword" type="password" placeholder="New Password" bind:value={frm_newpassword}>
        </div>
        <div class="flex flex-row justify-left mb-2">
            <button class="border border-normal py-1 px-2 bg-inputok mr-2">Update</button>
            <a href="/" class="action self-center" on:click|preventDefault="{e => dispatch('cancel')}">Cancel</a>
        </div>
        {#if status != ""}
        <div class="">
            <p class="uppercase italic text-xs">{status}</p>
        </div>
        {/if}
    </form>
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, del} from "./helpers.js";
import * as data from "./data.js";

export let userid = 0;

let svcurl = "/api";
let status = "";
let frm_password = "";
let frm_newpassword = "";

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

async function onSubmit(e) {
    status = "processing";

    let passwordreq = {};
    passwordreq.userid = userid;
    passwordreq.password = frm_password;
    passwordreq.newpassword = frm_newpassword;

    let sreq = `${svcurl}/password`;
    let [result, err] = await submit(sreq, "POST", passwordreq);
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

