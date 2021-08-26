<div class="flex flex-row">
    <Accounts bind:this={waccounts} root={root} bookid={selbookid} on:selectbookid={accounts_selectbookid} on:selectaccount={accounts_selectaccount} />
    <Txns bind:this={wtxns} root={root} bookid={selbookid} accountid={selaccountid} on:selectaccount={txns_selectaccount}/>

    <div class="dim bg-normal fg-normal mb-2 py-2 px-4" style="width: 20rem;">
        <h1 class="text-sm font-bold mb-2">Lorem Ipsum</h1>
        <p class="mb-4">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis consequat est eget est accumsan, eu iaculis tellus fermentum. Nunc pharetra ante feugiat maximus dapibus. Sed dui sapien, hendrerit vel viverra ut, maximus vitae risus. Nunc scelerisque bibendum magna, a faucibus nunc. Suspendisse sapien eros, tincidunt ac ultrices at, hendrerit ac purus. Aliquam erat volutpat. Suspendisse aliquam accumsan ornare. Sed ac leo vitae enim fringilla tristique at sit amet odio. Nunc vel sollicitudin est, vitae commodo purus. Nunc sit amet tellus tincidunt, ultrices quam a, cursus mi. Maecenas sollicitudin vehicula arcu, ut eleifend purus luctus vel. Integer scelerisque mi quis tincidunt laoreet. Fusce dapibus vulputate arcu, eu sollicitudin orci efficitur vel.</p>
        <p class="mb-4">Sed lacinia ligula a quam vehicula porttitor. Ut sollicitudin ante mi, convallis rhoncus dolor egestas non. Phasellus non leo volutpat, consectetur tortor ac, cursus leo. Duis quis tortor quis odio consequat consectetur eget vestibulum tellus. Etiam vel lacus neque. Pellentesque sagittis sodales scelerisque. Nam tristique feugiat enim eu fringilla. Proin lorem quam, ornare sit amet lectus in, aliquam volutpat mauris. Curabitur sed libero turpis. Proin vitae sagittis ex. Morbi semper turpis nisl, sed fermentum dolor pellentesque et.</p>
    </div>
</div>

<script>
import {onMount, createEventDispatcher} from "svelte";
let dispatch = createEventDispatcher();

import Accounts from "./Accounts.svelte";
import Txns from "./Txns.svelte";

export let root = null;
let waccounts;
let wtxns;

let selbookid = 1;
let selaccountid = null;

function accounts_selectbookid(e) {
    selbookid = e.detail;
    selaccountid = null;
    wtxns.reset();
}

function accounts_selectaccount(e) {
    let a = e.detail;
    selaccountid = a.accountid;
}

function txns_selectaccount(e) {
    let a = e.detail;
    waccounts.selectAccount(a);
}

export function postEvent(e) {
    waccounts.postEvent(e);
}

</script>
