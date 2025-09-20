const modal = document.getElementById("contactModal");
const openBtn = document.getElementById("openModalBtn");
const closeBtn = document.getElementById("closeModalBtn");
const form = document.getElementById("contactForm");

// Open modal
openBtn.onclick = () => {
  modal.style.display = "flex";
};

// Close modal
closeBtn.onclick = () => {
  modal.style.display = "none";
  form.reset();
};

// Close click outside modal
window.onclick = (e) => {
  if (e.target === modal) {
    modal.style.display = "none";
    form.reset();
  }
};

// Close after submit
function closeModal() {
  setTimeout(() => {
    modal.style.display = "none";
    form.reset();
  }, 300);
}
