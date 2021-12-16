{#if userid == 0}
    <p class="fg-dim">Select User</p>
{:else}
    <div class="flexrow justify-between mb-2">
        <h2 class="h2">User Settings</h2>
        <div></div>
    </div>
    <a class="tblrow" href="/" on:click|preventDefault="{e => onselect(1)}">Change Password</a>
    {#if selid == 1}
    <div class="p-2 border-b border-cell">
        <UserPassword userid={userid} on:submit={form_done} on:cancel={form_done} />
    </div>
    {/if}
    <a class="tblrow" href="/" on:click|preventDefault="{e => onselect(2)}">Delete User</a>
    {#if selid == 2}
    <div class="p-2 border-b border-cell">
        <UserDel userid={userid} on:submit={form_done} on:cancel={form_done} />
    </div>
    {/if}
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";
import * as data from "./data.js";
import UserPassword from "./UserPassword.svelte";
import UserDel from "./UserDel.svelte";

export let userid = 0;

let selid = 0;

export function reset() {
    selid = 0;
}

function onselect(nrow) {
    // Cancel any open form.
    if (selid != 0) {
        selid = 0;
        return;
    }

    selid = nrow;
}
function form_done(e) {
    selid = 0;
}

</script>


