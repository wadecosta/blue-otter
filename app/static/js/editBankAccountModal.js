document.addEventListener('DOMContentLoaded', () => {
    const modalElement = document.getElementById('editBankAccountModal');
    if (!modalElement) {
        console.error("Modal element not found");
        return;
    }

    const bootstrapModal = new bootstrap.Modal(modalElement);
    const form = modalElement.querySelector("form");
    const amountInput = form.querySelector("#bankAccountAmount");

    const key = sessionStorage.getItem("key");
    const iv = document.getElementById("iv")?.value;

    let id;
    let encryptedAmount;

    const toastElement = document.getElementById("toastContainer");
    const toastMessage = document.getElementById("toastMessage");
    const bsToast = new bootstrap.Toast(toastElement, { delay: 10000 });

    // Open modal and populate input
    document.querySelectorAll('.open-edit-bank-account-button').forEach(button => {
        button.addEventListener('click', function () {
            id = this.getAttribute("data-id");
            encryptedAmount = this.getAttribute("data-amount");

            const decryptedAmount = decryptText(encryptedAmount, key, iv);
            amountInput.value = decryptedAmount;

            bootstrapModal.show();
        });
    });

    // Reset form when modal is closed
    modalElement.addEventListener('hidden.bs.modal', () => {
        form.reset();
    });

    // Submit form
    form.addEventListener("submit", (event) => {
        event.preventDefault();

        const changedAmount = amountInput.value;
        const changedEncryptedAmount = encryptText(changedAmount, key, iv);

        const body = JSON.stringify({
            id: id.toString(),
            old_bank_account_amount: encryptedAmount,
            new_bank_account_amount: changedEncryptedAmount.toString()
        });

        const xhr = new XMLHttpRequest();
        xhr.open("POST", "/editBankAccount");
        xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8");

        xhr.onreadystatechange = function () {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status === 200) {
                    toastElement.classList.remove("text-bg-danger");
                    toastElement.classList.add("text-bg-success");
                    toastMessage.textContent = "Bank account updated successfully!";
                    bsToast.show();

                    form.reset();
                    bootstrapModal.hide();
                    setTimeout(() => window.location.reload(), 500);
                } else {
                    toastElement.classList.remove("text-bg-success");
                    toastElement.classList.add("text-bg-danger");
                    toastMessage.textContent = "Failed to update bank account.";
                    bsToast.show();
                }
            }
        };

        xhr.send(body);
    });
});

