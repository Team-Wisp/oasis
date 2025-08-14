// static/script.js
const emailInput = document.getElementById("email");
const sendOtpBtn = document.getElementById("send-otp");
const otpInput = document.getElementById("otp");
const verifyOtpBtn = document.getElementById("verify-otp");
const passwordInput = document.getElementById("password");
const confirmInput = document.getElementById("confirm");
const createAccountBtn = document.getElementById("create-account");

const stepEmail = document.getElementById("step-email");
const stepOtp = document.getElementById("step-otp");
const stepPassword = document.getElementById("step-password");
const feedback = document.getElementById("feedback");

const PARENT_ORIGIN = (() => {
  try { return new URL(document.referrer).origin || "*"; } catch { return "*"; }
})();

function showStep(step) {
  [stepEmail, stepOtp, stepPassword].forEach(el => el.classList.add("hidden"));
  step.classList.remove("hidden");
}
function setFeedback(msg, error = false) {
  feedback.textContent = msg;
  feedback.classList.remove("hidden");
  feedback.classList.toggle("text-green-400", !error);
  feedback.classList.toggle("text-red-400", error);
}
async function sha256(msg) {
  const buf = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(msg));
  return [...new Uint8Array(buf)].map(x => x.toString(16).padStart(2, "0")).join("");
}

sendOtpBtn.onclick = async () => {
  const email = emailInput.value.trim();
  if (!email.includes("@")) return setFeedback("Invalid email", true);

  sendOtpBtn.disabled = true; sendOtpBtn.classList.add("opacity-70","cursor-not-allowed");
  setFeedback("Verifying domain...");

  try {
    const domainRes = await fetch("/verify-domain", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email })
    });
    if (!domainRes.ok) throw new Error(await domainRes.text());

    setFeedback("Domain verified. Sending OTP...");
    await fetch("/send-otp", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email })
    });

    showStep(stepOtp);
    setFeedback("OTP sent to your email.");
  } catch (err) {
    setFeedback(err.message || "Verification failed", true);
  } finally {
    sendOtpBtn.disabled = false; sendOtpBtn.classList.remove("opacity-70","cursor-not-allowed");
  }
};

verifyOtpBtn.onclick = async () => {
  const otp = otpInput.value.trim();
  const emailHash = await sha256(emailInput.value.trim());
  if (!otp) return setFeedback("Enter the OTP", true);

  verifyOtpBtn.disabled = true; verifyOtpBtn.classList.add("opacity-70","cursor-not-allowed");

  try {
    const otpRes = await fetch("/verify-otp", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: emailHash, otp })
    });
    const otpResult = await otpRes.json();
    if (!otpResult.verified) throw new Error("Invalid OTP");

    showStep(stepPassword);
    setFeedback("OTP verified. Create your password.");
  } catch (err) {
    setFeedback(err.message || "OTP verification failed", true);
  } finally {
    verifyOtpBtn.disabled = false; verifyOtpBtn.classList.remove("opacity-70","cursor-not-allowed");
  }
};

createAccountBtn.onclick = async () => {
  const pw = passwordInput.value;
  const confirm = confirmInput.value;
  if (pw !== confirm) return setFeedback("Passwords don’t match", true);

  createAccountBtn.disabled = true; createAccountBtn.classList.add("opacity-70","cursor-not-allowed");

  try {
    const email = emailInput.value.trim();
    const domain = email.split("@")[1];
    const pwHash = await sha256(pw);
    const emailHash = await sha256(email);

    const res = await fetch("/create-account", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: emailHash, password: pwHash, domain })
    });

    if (res.status === 409) {
      window.parent.postMessage({ event: "OASIS_USER_EXISTS" }, PARENT_ORIGIN);
      return;
    }
    if (!res.ok) throw new Error(await res.text());

    window.parent.postMessage({ event: "OASIS_SIGNUP_SUCCESS" }, PARENT_ORIGIN);
  } catch (err) {
    setFeedback(err.message || "Account creation failed", true);
  } finally {
    createAccountBtn.disabled = false; createAccountBtn.classList.remove("opacity-70","cursor-not-allowed");
  }
};
