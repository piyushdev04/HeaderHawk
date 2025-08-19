// Set current year
document.getElementById('year').textContent = new Date().getFullYear();

const report = JSON.parse(localStorage.getItem("hh_report"));
const backendBase = "https://headerhawk.onrender.com"; // Update if deployed elsewhere

const reportSection = document.getElementById("report-section");
const list = document.getElementById("issues-list");
const noIssues = document.getElementById("no-issues");
const targetURLEl = document.getElementById("target-url");
const scoreEl = document.getElementById("score");
const gradeEl = document.getElementById("grade");
const viewJSON = document.getElementById("view-json");

if (!report) {
  reportSection.innerHTML = `
    <div class="empty-state">
      <h3>No report available</h3>
      <p>Please run a scan from the homepage.</p>
    </div>
  `;
} else {
  // Fill main info
  const targetURL = report.target_url || "Unknown";
  targetURLEl.textContent = targetURL;
  scoreEl.textContent = report.score ?? "-";
  gradeEl.textContent = report.grade ?? "-";

  // Link to JSON API
  viewJSON.href = `${backendBase}/api/scan?url=${encodeURIComponent(targetURL)}`;

  // Populate issues
  if (report.findings && report.findings.length > 0) {
    report.findings.forEach(f => {
      const li = document.createElement("li");
      li.className = "list-item";
      li.innerHTML = `
        <div class="list-line">
          <span class="chip warn">${f.header}</span>
          <span class="issue">${f.issue}</span>
        </div>
        <p class="advice">${f.advice}</p>
      `;
      list.appendChild(li);
    });
  } else {
    noIssues.style.display = "block";
  }
}