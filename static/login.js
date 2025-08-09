// static/login.js
async function sha256(msg) {
  const buf = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(msg));
  return [...new Uint8Array(buf)].map(x => x.toString(16).padStart(2, "0")).join("");
}

const $ = (id) => document.getElementById(id);
const emailEl = $("email");
const passwordEl = $("password");
const feedbackEl = $("feedback");
const btn = $("login");

// Derive parent origin from referrer (safer than "*")
const PARENT_ORIGIN = (() => {
  try { return new URL(document.referrer).origin || "*"; } catch { return "*"; }
})();

function showError(msg) {
  feedbackEl.textContent = msg;
  feedbackEl.classList.remove("hidden");
}

btn.onclick = async () => {
  const email = emailEl.value.trim();
  const password = passwordEl.value.trim();
  feedbackEl.classList.add("hidden");

  if (!email || !password) return showError("Email and password are required.");

  btn.disabled = true;
  btn.classList.add("opacity-70", "cursor-not-allowed");

  try {
    const emailHash = await sha256(email);
    const passwordHash = await sha256(password);

    const res = await fetch("/verify-login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: emailHash, password: passwordHash })
    });

    if (!res.ok) throw new Error(await res.text() || "Login failed");

    const data = await res.json();
    window.parent.postMessage({
      event: "OASIS_LOGIN_SUCCESS",
      data: { token: data.token }
    }, PARENT_ORIGIN);
  } catch (err) {
    showError(err.message || "Something went wrong");
  } finally {
    btn.disabled = false;
    btn.classList.remove("opacity-70", "cursor-not-allowed");
  }
};
