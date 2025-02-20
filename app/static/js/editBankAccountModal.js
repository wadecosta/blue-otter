document.addEventListener('DOMContentLoaded', (event) => {
	let modal = document.getElementById('editBankAccountModal');
	let span = document.querySelector(".closeEditBankAccount");
	let amount = document.getElementById("bankAccountAmount");

	/* Needed for encryption/decryption */
	let key = sessionStorage.getItem("key");
	let iv = document.getElementById("iv").value;

	let id;
	let encryptedAmount;

	document.querySelectorAll('.open-edit-bank-account-button').forEach(button => {
		button.onclick = function() {
			id = this.getAttribute("data-id");

			encryptedAmount = this.getAttribute("data-amount");
			decryptedAmount = decryptText(encryptedAmount, key, iv);
			
			bankAccountAmount.value = decryptedAmount;

			modal.showModal();
		}
	});

	span.onclick = function() {
		modal.close()
	}

	modal.addEventListener('click', function(event) {
		if (event.target === modal) {
			modal.close();
		}
	});

	let form = modal.querySelector("form");
	form.addEventListener("submit", function(event) {
		event.preventDefault();

		let changedAmount = bankAccountAmount.value;
		
		let changedEncryptedAmount = encryptText(changedAmount, key, iv);

		const xhr = new XMLHttpRequest();
		xhr.open("POST", "/editBankAccount");
		xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8");

		const body = JSON.stringify({
			id: id.toString(),
			old_bank_account_amount: encryptedAmount,
			new_bank_account_amount: changedEncryptedAmount.toString()
		});

		xhr.onreadystatechange = function() {
			if(xhr.readyState === XMLHttpRequest.DONE) {
				if(xhr.status === 200) {
					window.location.reload();
				} else {
					console.error("Failed to save changes:", xhr.status, xhr.statusText);
				}
			}
		};

		xhr.send(body);
		modal.close();
	});
});
