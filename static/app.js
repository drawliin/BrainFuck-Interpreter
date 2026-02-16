const input = document.getElementById("bf-input");
const output = document.getElementById("output");
const statusBadge = document.getElementById("status");
const count = document.getElementById("op-count");
const statOps = document.getElementById("stat-ops");
const statIgnored = document.getElementById("stat-ignored");
const statDepth = document.getElementById("stat-depth");
const statOutput = document.getElementById("stat-output");

const runBtn = document.getElementById("run-btn");
const clearBtn = document.getElementById("clear-btn");
const sampleBtn = document.getElementById("sample-btn");
const chips = document.querySelectorAll(".chip");
const barRows = document.querySelectorAll(".bar-row");

const helloWorld =
  "++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>.";
const commands = "><+-.,[]";

function setStatus(kind, text) {
  statusBadge.className = `status ${kind}`;
  statusBadge.textContent = text;
}

function analyzeCode(code) {
  const stats = {
    valid: 0,
    ignored: 0,
    maxDepth: 0,
    byCommand: {
      ">": 0,
      "<": 0,
      "+": 0,
      "-": 0,
      ".": 0,
      ",": 0,
      "[": 0,
      "]": 0,
    },
  };

  let depth = 0;
  for (const char of code) {
    if (!commands.includes(char)) {
      stats.ignored++;
      continue;
    }

    stats.valid++;
    stats.byCommand[char]++;

    if (char === "[") {
      depth++;
      if (depth > stats.maxDepth) {
        stats.maxDepth = depth;
      }
    } else if (char === "]" && depth > 0) {
      depth--;
    }
  }

  return stats;
}

function updateInsights() {
  const stats = analyzeCode(input.value);
  statOps.textContent = String(stats.valid);
  statIgnored.textContent = String(stats.ignored);
  statDepth.textContent = String(stats.maxDepth);

  const maxCount = Math.max(1, ...Object.values(stats.byCommand));
  barRows.forEach((row) => {
    const cmd = row.dataset.cmd;
    const value = stats.byCommand[cmd] || 0;
    const percentage = (value / maxCount) * 100;
    const fill = row.querySelector(".bar-fill");
    const label = row.querySelector(".bar-value");

    fill.style.width = `${percentage}%`;
    label.textContent = String(value);
  });
}

function updateCounter() {
  count.textContent = `${input.value.length} chars`;
  updateInsights();
}

async function runProgram() {
  const code = input.value.trim();
  if (!code) {
    setStatus("error", "No code");
    output.textContent = "Write or paste Brainfuck code first.";
    return;
  }

  setStatus("running", "Running");
  output.textContent = "Executing...";
  runBtn.disabled = true;

  try {
    const res = await fetch("/api/run", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ code }),
    });

    const data = await res.json();
    if (!res.ok) {
      setStatus("error", "Error");
      output.textContent = data.error || "Execution failed.";
      statOutput.textContent = "0";
      return;
    }

    const result = data.output || "";
    setStatus("ok", "Success");
    output.textContent = result || "(no output)";
    statOutput.textContent = String(result.length);
  } catch (err) {
    setStatus("error", "Server down");
    output.textContent = "Could not reach backend. Start server with: go run .";
    statOutput.textContent = "0";
  } finally {
    runBtn.disabled = false;
  }
}

runBtn.addEventListener("click", runProgram);

clearBtn.addEventListener("click", () => {
  input.value = "";
  output.textContent = "Your output will appear here.";
  setStatus("idle", "Idle");
  statOutput.textContent = "0";
  updateCounter();
  input.focus();
});

sampleBtn.addEventListener("click", () => {
  input.value = helloWorld;
  updateCounter();
  input.focus();
});

chips.forEach((chip) => {
  chip.addEventListener("click", () => {
    input.value = chip.dataset.code;
    updateCounter();
    input.focus();
  });
});

input.addEventListener("input", updateCounter);

input.addEventListener("keydown", (event) => {
  if ((event.ctrlKey || event.metaKey) && event.key === "Enter") {
    event.preventDefault();
    runProgram();
  }
});

updateCounter();
