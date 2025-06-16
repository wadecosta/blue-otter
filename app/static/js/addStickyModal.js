document.addEventListener('DOMContentLoaded', () => {
    const modalElement = document.getElementById("addStickyModal");

    if (!modalElement) {
        console.error("Modal element not found");
        return;
    }

    // Initialize Bootstrap modal instance
    const bootstrapModal = new bootstrap.Modal(modalElement);

    // Open modal when clicking any open-add-sticky-button
    document.querySelectorAll('.open-add-sticky-button').forEach(button => {
        button.addEventListener('click', () => {
            bootstrapModal.show();
        });
    });

    // Optional: Close button inside modal (already handled by data-bs-dismiss="modal")
    // But if you still want custom logic on Cancel:
    const cancelButton = modalElement.querySelector('.btn-danger');
    cancelButton.addEventListener('click', () => {
        bootstrapModal.hide();
    });

    // Get form and handle submission
    const form = modalElement.querySelector("form");
    form.addEventListener("submit", function(event) {
        event.preventDefault();

        const key = sessionStorage.getItem("key");
        const ivElement = document.getElementById("iv");
        const iv = ivElement ? ivElement.value : "";

        const titleInput = modalElement.querySelector("#title");
        const descriptionInput = modalElement.querySelector("#description");

        const titleValue = titleInput.value;
        const descriptionValue = descriptionInput.value;

        console.log(titleValue, descriptionValue);

        const encryptedTitle = encryptText(titleValue, key, iv);
        const encryptedDescription = encryptText(descriptionValue, key, iv);

        const xhr = new XMLHttpRequest();
        xhr.open("POST", "/addSticky");
        xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8");

        const body = JSON.stringify({
            title: encryptedTitle.toString(),
            description: encryptedDescription.toString()
        });

        xhr.onreadystatechange = function () {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status === 200) {
                    window.location.reload();
                } else {
                    console.error("Failed to add sticky:", xhr.status, xhr.statusText);
                }
            }
        };

        xhr.send(body);
        bootstrapModal.hide(); // Close the modal after sending
    });
});

