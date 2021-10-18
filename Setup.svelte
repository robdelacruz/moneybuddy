<div class="flex flex-row">
    <div class="bg-normal fg-normal mb-2 mr-2 py-2 px-4" style="width: 20rem;">
    {#if root == null}
        <p class="fg-dim">No data</p>
    {:else}
        <div class="flex flex-row border-b border-cell mb-4">
            <select class="text-xs fg-normal bg-normal pr-1" id="setupmenu" name="setupmenu" placeholder="Select" bind:value={selmenuid}>
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

    <div class="dim bg-normal fg-normal mb-2 py-2 px-4" style="width: 20rem;">
        <h1 class="text-sm font-bold mb-2">Lorem Ipsum</h1>
        <p class="mb-4">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis consequat est eget est accumsan, eu iaculis tellus fermentum. Nunc pharetra ante feugiat maximus dapibus. Sed dui sapien, hendrerit vel viverra ut, maximus vitae risus. Nunc scelerisque bibendum magna, a faucibus nunc. Suspendisse sapien eros, tincidunt ac ultrices at, hendrerit ac purus. Aliquam erat volutpat. Suspendisse aliquam accumsan ornare. Sed ac leo vitae enim fringilla tristique at sit amet odio. Nunc vel sollicitudin est, vitae commodo purus. Nunc sit amet tellus tincidunt, ultrices quam a, cursus mi. Maecenas sollicitudin vehicula arcu, ut eleifend purus luctus vel. Integer scelerisque mi quis tincidunt laoreet. Fusce dapibus vulputate arcu, eu sollicitudin orci efficitur vel.</p>
    </div>
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

