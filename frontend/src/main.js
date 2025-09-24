import { mount } from "svelte";
import AppSvelte from "./App.svelte";

const app = mount(AppSvelte, {
  target: document.querySelector("#app"),
  // props: { some: 'property' }
});
