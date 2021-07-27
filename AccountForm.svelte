{#if account == null || account.txns == null}
    <p class="fg-dim">Select Account</p>
{:else}
    <form class="">
        <div class="mb-2">
            <input class="block bg-input fg-normal py-1 px-2 w-full" name="accountname" id="accountname" type="text" placeholder="Enter Account Name" value={account.name}>
        </div>
        <div class="flex flex-row mb-2">
            <div class="mr-2 w-1/2">
                <select class="py-1 px-2 bg-input fg-normal w-full" id="accounttype" name="accounttype" placeholder="Account Type">
                    <option value="bankaccount">Bank Account</option>
                    <option value="stock">Stock</option>
                </select>
            </div>
            <div class="w-1/2">
                <select class="py-1 px-2 bg-input fg-normal w-full" id="currency" name="currency">
                    <option value="usd">USD</option>
                    <option value="php">PHP</option>
                </select>
            </div>
        </div>
        <div class="flex flex-row justify-start">
            <div>
                <button class="mx-auto border border-normal py-1 px-2 bg-inputok mr-2" on:click|preventDefault={onUpdate}>Update</button>
                <a href="/" class="mx-auto border-b border-normal pt-1" on:click|preventDefault={onCancel}>Cancel</a>
            </div>
        </div>
    </form>
{/if}

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit} from "./helpers.js";
import * as data from "./data.js";

export let account = null;

let svcurl = "/api";
let ui = {};

document.addEventListener("keydown", function(e) {
    if (e.key == "Escape") {
        dispatch("cancel");
    }
});

function onUpdate(e) {
    dispatch("update", account);
}

function onCancel(e) {
    dispatch("cancel");
}

</script>
