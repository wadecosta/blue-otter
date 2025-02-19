document.addEventListener('DOMContentLoaded', (event) => {
        
        let modal = document.getElementById("addBankAccountModal");
        if (!modal) {
                console.error("Modal element not found");
                return;
        }

        let span = document.querySelector(".closeAddBankAccount");

	let key = sessionStorage.getItem("key");
    	let iv = document.getElementById("iv").value;

        document.querySelectorAll('.open-add-bank-account-button').forEach(button => {
                button.onclick = function() {
                        modal.showModal();
                };
        });

        span.onclick = function() {
                modal.close();
        };

        modal.addEventListener('click', function(event) {
                if (event.target === modal) {
                        modal.close();
                }
        });

        let form = modal.querySelector("form");

	form.addEventListener("submit", function(event) {
    		event.preventDefault();

		let bankInput = modal.querySelector("#bank");
		let amountInput = modal.querySelector("#amount");

		let bankValue = bankInput.value;
        	let amountValue = amountInput.value;

		let encryptedAmount = encryptText(amountValue, key, iv);

		const xhr = new XMLHttpRequest();
		xhr.open("POST", "/addBankAccount");
		xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8");

		const body = JSON.stringify({
			bank: bankValue.toString(),
			amount: encryptedAmount.toString()
		});

		console.log(body);

		xhr.onreadystatechange = function() {
			if (xhr.readyState === XMLHttpRequest.DONE) {
				if (xhr.status === 200) {
					window.location.reload();
				} else {
					console.error("Failed to add sticky:", xhr,status, xhr.statusText);
				}
			}
		};

		xhr.send(body);
		modal.close();

        });
});
