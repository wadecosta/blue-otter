document.addEventListener('DOMContentLoaded', () => {
    const modalElement = document.getElementById("addCDModal");
    if (!modalElement) {
        console.error("Modal element not found");
        return;
    }

    const bootstrapModal = new bootstrap.Modal(modalElement);
    const form = modalElement.querySelector("form");
    const key = sessionStorage.getItem("key");
    const iv = document.getElementById("iv")?.value;

    const toastElement = document.getElementById("toastContainer");
    const toastMessage = document.getElementById("toastMessage");
    const bsToast = new bootstrap.Toast(toastElement, { delay: 10000 }); // 10 sec

    // Show modal on button click
    document.querySelectorAll('.open-CD-button').forEach(button => {
        button.addEventListener('click', () => {
            bootstrapModal.show();
        });
    });

    // Reset form when modal is closed
    modalElement.addEventListener('hidden.bs.modal', () => {
        form.reset();
    });

    // Handle form submission
    form.addEventListener("submit", (event) => {
        event.preventDefault();

        const bank = form.querySelector("#bank").value;
        const startDate = form.querySelector("#startDate").value;
        const deposit = form.querySelector("#deposit").value;
        const term = form.querySelector("#term").value;
        const apy = form.querySelector("#apy").value;

        const encryptedStartDate = encryptText(startDate, key, iv);
        const encryptedDeposit = encryptText(deposit, key, iv);
        const encryptedTerm = encryptText(term, key, iv);
        const encryptedApy = encryptText(apy, key, iv);

        const body = JSON.stringify({
            bank: bank.toString(),
            startDate: encryptedStartDate.toString(),
            deposit: encryptedDeposit.toString(),
            term: encryptedTerm.toString(),
            apy: encryptedApy.toString()
        });

        const xhr = new XMLHttpRequest();
        xhr.open("POST", "/addCD");
        xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8");

        xhr.onreadystatechange = () => {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status === 200) {
                    toastElement.classList.remove("text-bg-danger");
                    toastElement.classList.add("text-bg-success");
                    toastMessage.textContent = "CD added successfully!";
                    bsToast.show();

                    form.reset();
                    bootstrapModal.hide();
                    setTimeout(() => window.location.reload(), 500); // Reload after modal closes
                } else {
                    toastElement.classList.remove("text-bg-success");
                    toastElement.classList.add("text-bg-danger");
                    toastMessage.textContent = "Failed to add CD.";
                    bsToast.show();
                }
            }
        };

        xhr.send(body);
    });
});

