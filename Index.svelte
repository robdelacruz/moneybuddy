<script>
import EventBox from "./EventBox.svelte";

export let name;
let eventbox;

document.addEventListener("keydown", function(e) {
    if (e.key >= 'a' && e.key <= 'e') {
        console.log(`e.key: '${e.key}', e.code: '${e.code}'`);
    }
    eventbox.setEvent(e);
});
</script>

<p class="hello">Hello {name}</p>

<EventBox bind:this={eventbox} />

