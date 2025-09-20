const modal = document.getElementById("contactModal");
const openBtn = document.getElementById("openModalBtn");
const closeBtn = document.getElementById("closeModalBtn");
const form = document.getElementById("contactForm");

// Open Add Contact modal
openBtn.onclick = () => {
  resetForm();
  document.getElementById("modalTitle").innerText = "Add Contact";
  form.setAttribute("hx-post", "/contacts"); // use POST for new
  form.removeAttribute("hx-put");
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

// Close after submit
function closeModal() {
  setTimeout(() => {
    modal.style.display = "none";
    resetForm();
  }, 300);
}

// Reset form
function resetForm() {
  form.reset();
  document.getElementById("contactId").value = "";
}

// Open modal in Edit mode
function openEditModal(id, contactType, firstName, lastName, email, phone) {
  modal.style.display = "flex";
  document.getElementById("modalTitle").innerText = "Edit Contact";

  document.getElementById("contactId").value = id;
  document.getElementById("contactType").value = contactType; // ✅ dropdown value
  document.getElementById("firstName").value = firstName;
  document.getElementById("lastName").value = lastName;
  document.getElementById("email").value = email;
  document.getElementById("phone").value = phone;

  form.setAttribute("hx-put", `/contacts/${id}`);
  form.removeAttribute("hx-post");
  form.setAttribute("hx-target", `#contact-${id}`);
  form.setAttribute("hx-swap", "outerHTML");
}
