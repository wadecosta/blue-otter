document.addEventListener('DOMContentLoaded', () => {
  const modalElement = document.getElementById("editStickyModal");
  const bsModal = new bootstrap.Modal(modalElement);

  const form = document.getElementById("editStickyForm");
  const titleInput = document.getElementById("title");
  const descriptionInput = document.getElementById("description");

  const key = sessionStorage.getItem("key");
  const iv = document.getElementById("iv").value;

  let id, encryptedTitle, encryptedDescription;

  // Handle open buttons
  document.querySelectorAll('.open-edit-sticky-button').forEach(button => {
    button.addEventListener('click', () => {
      id = button.getAttribute("data-id");
      encryptedTitle = button.getAttribute("data-title");
      encryptedDescription = button.getAttribute("data-description");

      const decryptedTitle = decryptText(encryptedTitle, key, iv);
      const decryptedDescription = decryptText(encryptedDescription, key, iv);

      titleInput.value = decryptedTitle;
      descriptionInput.value = decryptedDescription;

      bsModal.show();
    });
  });

  // Reset form on modal hide
  modalElement.addEventListener('hidden.bs.modal', () => {
    form.reset();
  });

  // Show toast helper
  function showToast(message, isSuccess = true) {
    const toastEl = document.getElementById("toastContainer");
    const toastMsg = document.getElementById("toastMessage");

    toastEl.classList.remove("text-bg-success", "text-bg-danger");
    toastEl.classList.add(isSuccess ? "text-bg-success" : "text-bg-danger");
    toastMsg.textContent = message;

    const toast = new bootstrap.Toast(toastEl, { delay: 10000 });
    toast.show();
  }

  // Form submission
  form.addEventListener("submit", function (event) {
    event.preventDefault();

    const changedTitle = titleInput.value;
    const changedDescription = descriptionInput.value;

    const changedEncryptedTitle = encryptText(changedTitle, key, iv);
    const changedEncryptedDescription = encryptText(changedDescription, key, iv);

    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/editSticky");
    xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8");

    const body = JSON.stringify({
      id: id.toString(),
      old_sticky_title: encryptedTitle,
      old_sticky_description: encryptedDescription,
      new_sticky_title: changedEncryptedTitle.toString(),
      new_sticky_description: changedEncryptedDescription.toString()
    });

    xhr.onreadystatechange = function () {
      if (xhr.readyState === XMLHttpRequest.DONE) {
        if (xhr.status === 200) {
          showToast("Sticky note updated successfully!");
          setTimeout(() => window.location.reload(), 1000);
        } else {
          showToast("Failed to update sticky note.", false);
          console.error("Failed to save changes:", xhr.status, xhr.statusText);
        }
      }
    };

    xhr.send(body);
    bsModal.hide();
  });
});
