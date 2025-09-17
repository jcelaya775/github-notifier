import { Events } from "@wailsio/runtime";
import { GreetService } from "./bindings/github.com/jcelaya775/github-notifier/index.js";

const resultElement = document.getElementById("result");
const timeElement = document.getElementById("time");

window.doGreet = () => {
  let name = document.getElementById("name").value;
  if (!name) {
    name = "anonymous";
  }
  GreetService.Greet(name).then((result) => {
    resultElement.innerText = result;
  }).catch((err) => {
    console.log(err);
  });
};

// TODO: Explore later (read docs)
// Events.On("keydown", async (e) => {
//     if (e.key === 'Escape') {
//         // This removes focus from the window
//         await WindowSetFocusable(false);
//
//         // Optionally re-enable focus after a short delay so the user can click back in
//         setTimeout(() => {
//             WindowSetFocusable(true);
//         }, 100);
//     }
// })

document.addEventListener("keydown", async (e) => {
  console.log("keydown:", e.key);
  if (e.key === "Escape") {
    console.log("Escape pressed, emitting event to backend");
    await Events.Emit({
      name: "escape-pressed",
    });
  }
});

Events.On("time", (time) => {
  timeElement.innerText = time.data;
});
