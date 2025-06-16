document.addEventListener('DOMContentLoaded', () => {
    const modalElement = document.getElementById("addBankAccountModal");

    if (!modalElement) {
        console.error("Modal element not found");
        return;
    }

    const bootstrapModal = new bootstrap.Modal(modalElement);
    const toastElement = document.getElementById("toastContainer");
    const toastMessage = document.getElementById("toastMessage");
    const bsToast = new bootstrap.Toast(toastElement);

    const form = modalElement.querySelector("form");
    const key = sessionStorage.getItem("key");
    const ivElement = document.getElementById("iv");
    const iv = ivElement ? ivElement.value : "";

    // Show modal
    document.querySelectorAll('.open-add-bank-account-button').forEach(button => {
        button.addEventListener('click', () => {
            bootstrapModal.show();
        });
    });

    // Clear form when modal is hidden
    modalElement.addEventListener('hidden.bs.modal', () => {
        form.reset(); // Clears all fields
        // You could also clear validation errors here if any
    });

    // Handle form submission
    form.addEventListener("submit", function (event) {
        event.preventDefault();

        const bankInput = form.querySelector("#bank");
        const amountInput = form.querySelector("#amount");

        const bankValue = bankInput.value;
        const amountValue = amountInput.value;

        const encryptedAmount = encryptText(amountValue, key, iv);

        const xhr = new XMLHttpRequest();
        xhr.open("POST", "/addBankAccount");
        xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8");

        const body = JSON.stringify({
            bank: bankValue.toString(),
            amount: encryptedAmount.toString()
        });

        xhr.onreadystatechange = function () {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status === 200) {
                    toastElement.classList.remove("text-bg-danger");
                    toastElement.classList.add("text-bg-success");
                    toastMessage.textContent = "Bank account added successfully!";
                    bsToast.show();
                    form.reset();
                    bootstrapModal.hide();
                } else {
                    toastElement.classList.remove("text-bg-success");
                    toastElement.classList.add("text-bg-danger");
                    toastMessage.textContent = "Failed to add bank account.";
                    bsToast.show();
                }
            }
        };

        xhr.send(body);
    });
});
