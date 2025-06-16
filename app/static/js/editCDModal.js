document.addEventListener('DOMContentLoaded', () => {
    const modalElement = document.getElementById('editCDModal');
    if (!modalElement) {
        console.error("editCDModal not found");
        return;
    }

    const bootstrapModal = new bootstrap.Modal(modalElement);
    const form = modalElement.querySelector("form");

    const key = sessionStorage.getItem("key");
    const ivElement = document.getElementById("iv");
    const iv = ivElement ? ivElement.value : "";

    let id;
    let encryptedBank, encryptedStartDate, encryptedDeposit, encryptedTerm, encryptedApy;

    document.querySelectorAll('.open-edit-CD-button').forEach(button => {
        button.addEventListener('click', () => {
            id = button.getAttribute("data-id");

            const bankID = button.getAttribute("data-bank");
            const bankElement = modalElement.querySelector("#bank");

            if (bankElement) {
                bankElement.value = bankID;
                if (bankElement.value !== bankID) {
                    for (let option of bankElement.options) {
                        if (option.value === bankID) {
                            option.selected = true;
                            break;
                        }
                    }
                }
            }

            encryptedStartDate = button.getAttribute("data-start-date");
            encryptedDeposit = button.getAttribute("data-deposit");
            encryptedTerm = button.getAttribute("data-term");
            encryptedApy = button.getAttribute("data-apy");

            const decryptedStartDate = decryptText(encryptedStartDate, key, iv);
            const decryptedDeposit = decryptText(encryptedDeposit, key, iv);
            const decryptedTerm = decryptText(encryptedTerm, key, iv);
            const decryptedApy = decryptText(encryptedApy, key, iv);

            modalElement.querySelector("#startDate").value = decryptedStartDate;
            modalElement.querySelector("#deposit").value = decryptedDeposit;
            modalElement.querySelector("#term").value = decryptedTerm;
            modalElement.querySelector("#apy").value = decryptedApy;

            bootstrapModal.show();
        });
    });

    modalElement.addEventListener('hidden.bs.modal', () => {
        form.reset(); // Clear all fields when closed
    });

    // Add form submission logic here as needed
});

