<div class="setup section-container">
{#if root == null}
    <p class="fg-dim">No data</p>
{:else}
    <div class="flexrow border-b border-cell mb-4">
        <select class="text-xs" id="setupmenu" name="setupmenu" placeholder="Select" bind:value={selmenuid}>
            {#each menuitems as menuitem}
                {#if menuitem.id == selmenuid}
                <option selected value={menuitem.id}>{menuitem.name}</option>
                {:else}
                <option value={menuitem.id}>{menuitem.name}</option>
                {/if}
            {/each}
        </select>
    </div>
    {#if selmenuid == "books"}
        <SetupBooks root={root} userid={userid} />
    {:else if selmenuid == "currencies"}
        <SetupCurrencies root={root} userid={userid} />
    {:else if selmenuid == "user"}
        <SetupUser userid={userid} />
    {/if}
{/if}
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();
import {find, submit, getls, setls} from "./helpers.js";
import * as data from "./data.js";
import SetupBooks from "./SetupBooks.svelte";
import SetupCurrencies from "./SetupCurrencies.svelte";
import SetupUser from "./SetupUser.svelte";

export let userid = 0;
export let root = null;

let menuitems = [
    {id: "books", name: "Books"},
    {id: "currencies", name: "Currencies"},
    {id: "user", name: "User"}
];

let selmenuid = getls("Setup", "selmenuid", menuitems[0].id);

// Remember selections when changed.
$: setls("Setup", "selmenuid", selmenuid);

</script>

