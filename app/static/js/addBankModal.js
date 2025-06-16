document.addEventListener('DOMContentLoaded', () => {
    const modalElement = document.getElementById("addBankModal");
    const form = modalElement?.querySelector("form");
    const toastElement = document.getElementById("toastContainer");
    const toastMessage = document.getElementById("toastMessage");
    const bsToast = new bootstrap.Toast(toastElement, { delay: 10000 });

    if (!modalElement || !form || !toastElement || !toastMessage) {
        console.error("Required DOM elements not found");
        return;
    }

    const bootstrapModal = new bootstrap.Modal(modalElement);

    document.querySelectorAll('.open-add-bank-button').forEach(button => {
        button.addEventListener('click', () => {
            bootstrapModal.show();
        });
    });

    modalElement.addEventListener('hidden.bs.modal', () => {
        form.reset();
    });

    form.addEventListener("submit", function (event) {
        event.preventDefault();

        const formData = new FormData(form);

        fetch("/addBank", {
            method: "POST",
            body: formData
        })
        .then(response => {
            if (!response.ok) {
                throw new Error("Upload failed");
            }
            return response.json();
        })
        .then(data => {
            toastElement.classList.remove("text-bg-danger");
            toastElement.classList.add("text-bg-success");
            toastMessage.textContent = "Bank added successfully!";
            bsToast.show();
            bootstrapModal.hide();
        })
        .catch(error => {
            toastElement.classList.remove("text-bg-success");
            toastElement.classList.add("text-bg-danger");
            toastMessage.textContent = "Error uploading bank.";
            bsToast.show();
            console.error(error);
        });
    });
});
