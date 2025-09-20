// Modal functionality
document.addEventListener("DOMContentLoaded", function () {
  const modal = document.getElementById("contactModal");
  const openModalCard = document.getElementById("openModalCard");
  const closeModalBtn = document.getElementById("closeModalBtn");
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

  // Open Add Contact modal when the add card is clicked
  openModalCard.addEventListener("click", function () {
    resetForm();
    document.getElementById("modalTitle").innerText = "Add Contact";
    form.setAttribute("hx-post", "/contacts"); // use POST for new
    modal.style.display = "flex";
  });

  // Close modal when X is clicked
  closeModalBtn.addEventListener("click", function () {
    modal.style.display = "none";
    resetForm();
  });

  // Close modal when clicking outside the modal content
  window.addEventListener("click", function (event) {
    if (event.target === modal) {
      modal.style.display = "none";
      resetForm();
    }
  });

  // Close after submit (keeps UX smooth)
  window.closeModal = function () {
    setTimeout(() => {
      modal.style.display = "none";
      resetForm();
    }, 300);
  };

  // Open modal for editing
  window.openEditModal = function (
    id,
    contactType,
    firstName,
    lastName,
    email,
    phone,
  ) {
    modal.style.display = "flex";
    document.getElementById("modalTitle").innerText = "Edit Contact";

    // Set the hidden ID field
    document.getElementById("contactId").value = id;

    document.getElementById("contactType").value = contactType;
    document.getElementById("firstName").value = firstName;
    document.getElementById("lastName").value = lastName;
    document.getElementById("email").value = email;
    document.getElementById("phone").value = phone;

    // Update form attributes for editing
    form.setAttribute("hx-post", `/contacts/${id}`);
    form.setAttribute("hx-target", `#contact-${id}`);
    form.setAttribute("hx-swap", "outerHTML");
  };

  // Reset form for "Add"
  function resetForm() {
    form.reset();
    document.getElementById("contactId").value = ""; // Clear the ID

    // Reset form attributes for adding new contact
    form.setAttribute("hx-post", "/contacts");
    form.setAttribute("hx-target", "#contacts-list");
    form.setAttribute("hx-swap", "afterbegin");
  }
});
