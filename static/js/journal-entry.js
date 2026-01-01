const tbody = document.querySelector("#tbody");
const debitTotalInput = document.querySelector("input[name='debit-bal']");
const creditTotalInput = document.querySelector("input[name='credit-bal']");
const balanceRow = document.querySelector("#balance-row");

// Recalculate totals without overwriting input values
function recalcTotals() {

  let debitTotal = 0;
  let creditTotal = 0;
  //document.getElementById("save-status").classList.remove("d-none");
  //document.getElementById("save-status").classList.add("btn-primary");
  tbody.querySelectorAll("tr:not(#balance-row)").forEach((row) => {
    const debitInput = row.querySelector(".debit");
    const creditInput = row.querySelector(".credit");
    const accountInput = row.querySelector(".account");
    const descInput = row.querySelector(".description");

    if (!debitInput || !creditInput || !accountInput || !descInput) return;

    const debit = parseFloat(debitInput.value) || 0;
    const credit = parseFloat(creditInput.value) || 0;

    debitTotal += debit;
    creditTotal += credit;

    // Validation: either debit or credit must be filled (not both or neither)
    if (
      (debit > 0 && credit > 0) ||
      (debit === 0 && credit === 0) ||
      !accountInput.value ||
      !descInput.value
    ) {
      row.classList.add("table-warning");
    } else {
      row.classList.remove("table-warning");
    }
  });

  debitTotalInput.value = debitTotal.toFixed(2);
  creditTotalInput.value = creditTotal.toFixed(2);

  if (debitTotal !== creditTotal) {
    balanceRow.classList.remove("table-success");
    balanceRow.classList.add("table-danger");
  } else {
    balanceRow.classList.remove("table-danger");
    balanceRow.classList.add("table-success");
  }
}

// Delegated event handling: recalc totals on input
tbody.addEventListener("input", (e) => {
  if (
    ["debit", "credit", "account", "description"].some((cls) =>
      e.target.classList.contains(cls)
    )
  ) {
    removePostButton();
    document.getElementById("save-status").classList.remove("d-none");
    recalcTotals();
  }
});

// Format debit/credit on blur
tbody.addEventListener(
  "blur",
  (e) => {
    if (
      e.target.classList.contains("debit") ||
      e.target.classList.contains("credit")
    ) {
      const val = parseFloat(e.target.value) || 0;
      e.target.value = val.toFixed(2);
      document.getElementById("save-status").classList.remove("d-none");
      recalcTotals();
    }
  },
  true
);

// Add new line
document.getElementById("add-line").addEventListener("click", () => {
  //alert("Add Line");
  removePostButton();
  const warningRows = tbody.querySelectorAll(".table-warning");
  if (warningRows.length > 0) return; // prevent adding if warnings exist

  const firstRow = tbody.querySelector("tr:not(#balance-row)");
  if (!firstRow) return;

  const clonedRow = firstRow.cloneNode(true);

  // Clear values in cloned row
  clonedRow
    .querySelectorAll("input.debit, input.credit")
    .forEach((input) => (input.value = "0.00"));
  clonedRow
    .querySelectorAll("input.description")
    .forEach((input) => (input.value = ""));
  clonedRow
    .querySelectorAll("select.account")
    .forEach((input) => (input.selectedIndex = 0));

  // Insert before balance row
  tbody.insertBefore(clonedRow, balanceRow);
});

// Form submission validation
document
  .querySelector("form#journal-submit")
  .addEventListener("submit", (event) => {
    let valid = true;

    tbody.querySelectorAll("tr:not(#balance-row)").forEach((row) => {
      const debit = parseFloat(row.querySelector(".debit")?.value) || 0;
      const credit = parseFloat(row.querySelector(".credit")?.value) || 0;
      const account = row.querySelector(".account")?.value;
      const desc = row.querySelector(".description")?.value;

      if (
        !account ||
        !desc ||
        (debit === 0 && credit === 0) ||
        (debit > 0 && credit > 0)
      ) {
        row.classList.add("invalid-row");
        valid = false;
      } else {
        row.classList.remove("invalid-row");
      }
    });

    if (
      parseFloat(debitTotalInput.value) !== parseFloat(creditTotalInput.value)
    ) {
      balanceRow.classList.add("invalid-row");
      valid = false;
    } else {
      balanceRow.classList.remove("invalid-row");
    }

    if (!valid) event.preventDefault();
  });

/* document.addEventListener("DOMContentLoaded", () => {
  const btn = document.getElementById("post-button");
  const form = document.getElementById("journal-post");
  if (!btn) return;
  btn.addEventListener("click", () => {
    form.submit();
  });
}); */

function removePostButton() {
  if (document.querySelector("#post-status")) {
    document.querySelector("#post-status").remove();
  }
}

// Initial calculation on page load
recalcTotals();


document.addEventListener("DOMContentLoaded", () => {
  const confirmBtn = document.getElementById("confirmPostBtn");

  confirmBtn.addEventListener("click", () => {
    // Submit your existing form
    document.querySelector("form#journal-post").submit();

    // Close the modal (optional, but safe)
    const modalEl = document.getElementById("postJournalModal");
    const modal = bootstrap.Modal.getInstance(modalEl);
    if (modal) modal.hide();
  });
});