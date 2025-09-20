const modal = document.getElementById("contactModal");
const openBtn = document.getElementById("openModalBtn");
const closeBtn = document.getElementById("closeModalBtn");
const form = document.getElementById("contactForm");

// Add this to your script.js to handle HTMX errors
document.body.addEventListener("htmx:responseError", function (evt) {
  console.error(
    "HTMX Error:",
    evt.detail.xhr.status,
    evt.detail.xhr.responseText,
  );
  alert("Error: " + evt.detail.xhr.responseText);
});

// Open Add Contact modal
openBtn.onclick = () => {
  resetForm();
  document.getElementById("modalTitle").innerText = "Add Contact";
  form.setAttribute("hx-post", "/contacts"); // use POST for new
  modal.style.display = "flex";
};

// Close modal
closeBtn.onclick = () => {
  modal.style.display = "none";
  resetForm();
};

// Close click outside modal
window.onclick = (e) => {
  if (e.target === modal) {
    modal.style.display = "none";
    resetForm();
  }
};

// Close after submit (keeps UX smooth)
function closeModal() {
  setTimeout(() => {
    modal.style.display = "none";
    resetForm();
  }, 300);
}

// Reset form for "Add"
function resetForm() {
  form.reset();
  const idInput = document.getElementById("contactId");
  if (idInput) idInput.value = ""; // empty for new

  form.setAttribute("hx-post", "/contacts");
  form.removeAttribute("hx-put"); // keep safe if you used hx-put earlier
  form.setAttribute("hx-target", "#contacts-list");
  form.setAttribute("hx-swap", "afterbegin");
}

// Open modal for editing
function openEditModal(id, contactType, firstName, lastName, email, phone) {
  modal.style.display = "flex";
  document.getElementById("modalTitle").innerText = "Edit Contact";

  const idInput = document.getElementById("contactId");
  if (idInput) idInput.value = id;

  document.getElementById("contactType").value = contactType;
  document.getElementById("firstName").value = firstName;
  document.getElementById("lastName").value = lastName;
  document.getElementById("email").value = email;
  document.getElementById("phone").value = phone;

  // CRITICAL FIX: Set the correct HTMX attributes for update
  form.setAttribute("hx-post", `/contacts/${id}`);
  form.setAttribute("hx-target", `#contact-${id}`);
  form.setAttribute("hx-swap", "outerHTML");

  // Remove any hx-put attributes to avoid confusion
  form.removeAttribute("hx-put");
}
