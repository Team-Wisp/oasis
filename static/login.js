async function sha256(msg) {
    const buf = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(msg));
    return [...new Uint8Array(buf)].map(x => x.toString(16).padStart(2, "0")).join("");
  }

  document.getElementById("login").onclick = async () => {
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value.trim();
    const feedback = document.getElementById("feedback");
  
    if (!email || !password) {
      showError("Email and password are required.");
      return;
    }
  
    const emailHash = await sha256(email);
    const passwordHash = await sha256(password);
  
    try {
      const res = await fetch("/verify-login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: emailHash, password: passwordHash })
      });
  
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || "Login failed.");
      }
  
      const data = await res.json();
      window.parent.postMessage({
        event: "OASIS_LOGIN_SUCCESS",
        data: { token: data.token }
      }, "*");// change this to only desert
  
    } catch (err) {
      showError(err.message);
    }
  };
  
  function showError(msg) {
    const feedback = document.getElementById("feedback");
    feedback.textContent = msg;
    feedback.classList.remove("hidden");
  }
  