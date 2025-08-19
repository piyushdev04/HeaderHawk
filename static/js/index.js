// Set current year
document.getElementById('year').textContent = new Date().getFullYear();

const form = document.getElementById("scan-form");
const input = document.getElementById("url");

const backendBase = "https://headerhawk.onrender.com"; // Update if deployed elsewhere

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  const targetURL = input.value.trim();
  if (!targetURL) return;

  form.querySelector("button").disabled = true;
  form.querySelector("button").textContent = "Scanning…";

  try {
    const res = await fetch(`${backendBase}/api/scan?url=${encodeURIComponent(targetURL)}`);
    if (!res.ok) throw new Error(`Scan failed: ${res.statusText}`);
    const report = await res.json();

    // Save report in localStorage
    localStorage.setItem("hh_report", JSON.stringify(report));

    // Redirect to report page
    window.location.href = "/report.html";
  } catch (err) {
    console.error(err);
    alert("Error scanning the site. Check the console for details.");
  } finally {
    form.querySelector("button").disabled = false;
    form.querySelector("button").textContent = "Scan";
  }
});