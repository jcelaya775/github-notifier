import { Events } from "@wailsio/runtime";
import { GitHubAPIService, GreetService } from "./bindings/github.com/jcelaya775/github-notifier/index.js";

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

document.addEventListener("keydown", async (e) => {
  console.log("keydown:", e.key);
  if (e.key === "Escape") {
    console.log("Escape pressed, emitting event to backend");
    await Events.Emit({
      name: "escape-pressed",
    });
  }
});

(async () => {
  const notifications = await GitHubAPIService.GetNotifications();
  console.log("Notifications:", notifications);
})();

Events.On("time", (time) => {
  timeElement.innerText = time.data;
});
