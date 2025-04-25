document.addEventListener('DOMContentLoaded', (event) => {
	let modal = document.getElementById('addCDModal');
	let span = document.querySelector('.closeAddCD');

	/* Needed for encryption/decryption */
	let key = sessionStorage.getItem("key");
	let iv = document.getElementById("iv").value;
	
	let id;
	let encryptedBank;
	let encryptedStartDate;
	let encryptedDeposit;
	let encryptedTerm;
	let encryptedApy;

	document.querySelectorAll('.open-edit-CD-button').forEach(button => {
		button.onclick = function() {
			id = this.getAttribute("data-id");

			let bankID = this.getAttribute("data-bank");
			let bankElement = document.getElementById("bank");

			console.log("Bank ID:", bankID);
			
			if (bankElement) {
            console.log("Dropdown options:", [...bankElement.options].map(opt => opt.value));

            // Attempt direct value assignment
            bankElement.value = bankID;

            // If direct assignment fails, manually select the option
            if (bankElement.value !== bankID) {
                console.warn(`Direct selection failed for "${bankID}", selecting manually.`);
                for (let option of bankElement.options) {
                    if (option.value === bankID) {
                        option.selected = true;
                        console.log(`Selected option: ${option.value}`);
                        break;
                    }
                }
            }
        }


			encryptedStartDate = this.getAttribute("data-start-date");
			decryptedStartDate = decryptText(encryptedStartDate, key, iv);

			encryptedDeposit = this.getAttribute("data-deposit");
			decryptedDeposit = decryptText(encryptedDeposit, key, iv);

			encryptedTerm = this.getAttribute("data-term");
			decryptedTerm = decryptText(encryptedTerm, key, iv);

			encryptedApy = this.getAttribute("data-apy");
			decryptedApy = decryptText(encryptedApy, key, iv);


			document.getElementById("startDate").value = decryptedStartDate;
            		document.getElementById("deposit").value = decryptedDeposit;
            		document.getElementById("term").value = decryptedTerm;
            		document.getElementById("apy").value = decryptedApy;

			modal.showModal();
		}
	});

	span.onclick = function() {
		modal.close();
	};
	
	modal.addEventListener('click', function(event) {
		if (event.target === modal) {
			modal.close();
		}
	});
});
