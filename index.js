import Index from "./Index.svelte";
const index = new Index({
    target: document.querySelector("#container"),
    props: {
        name: "rob"
    }
});

export default index;

