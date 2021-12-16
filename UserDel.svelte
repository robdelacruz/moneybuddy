{#if userid == 0}
    <p class="fg-dim">Select User</p>
{:else}
    <form class="" autocomplete="off" on:submit|preventDefault={onSubmit}>
        <div class="mb-2">
            <input class="input w-full" name="username" id="username" type="text" placeholder="Username" bind:value={frm_username}>
        </div>
        <div class="mb-2">
            <input class="input w-full" name="password" id="password" type="password" placeholder="Password" bind:value={frm_password}>
        </div>
        {#if mode == ""}
        <div class="flexrow mb-2">
            <button class="btn bg-inputok mr-2" on:click|preventDefault={onDelete}>Delete User</button>
            <a href="/" class="action" on:click|preventDefault="{e => dispatch('cancel')}">Cancel</a>
        </div>
        {:else if mode == "delete"}
        <div class="">
            <p class="prompt mb-2">This cannot be undone. Confirm Delete?</p>
            <div>
                <button class="btn bg-inputdel mr-2">Delete</button>
                <a href="/" class="action" on:click|preventDefault={onCancelDelete}>Cancel</a>
            </div>
        </div>
        {/if}
        {#if status != ""}
        <div class="">
            <p class="prompt">{status}</p>
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
let mode = "";
let status = "";
let frm_username = "";
let frm_password = "";

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

function onDelete(e) {
    mode = "delete";
}

function onCancelDelete(e) {
    mode = "";
}

async function onSubmit(e) {
    status = "processing";

    // Validate username and password first.
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

    // Delete user.
    sreq = `${svcurl}/user?id=${userid}`;
    err = await del(sreq);
    if (err != null) {
        console.error(err);
        status = "server error deleting user";
        return;
    }

    status = "";
    dispatch("submit");
}

function onCancel(e) {
    dispatch("cancel");
}

</script>

