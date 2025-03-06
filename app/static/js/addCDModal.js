document.addEventListener('DOMContentLoaded', (event) => {
	
	let modal = document.getElementById("addCDModal");
	if (!modal) {
		console.error("Modal element not found");
		return;
	}

	let span = document.querySelector(".closeAddCD");

	let key = sessionStorage.getItem("key");
	let iv = document.getElementById("iv").value;

	document.querySelectorAll('.open-CD-button').forEach(button => {
		button.onclick = function() {
			modal.showModal();
		}
	});

	span.onclick = function() {
		modal.close();
	}

	modal.addEventListener('click', function(event) {
		if (event.target === modal) {
			modal.close();
		}
	});

	let form = modal.querySelector("form");

	form.addEventListener("submit", function(event) {
		event.preventDefault();

		let bankInput = modal.querySelector("#bank");
		let startDateInput = modal.querySelector("#startDate");
		let depositInput = modal.querySelector("#deposit");
		let termInput = modal.querySelector("#term");
		let aprInput = modal.querySelector("#apy");

		let bankValue = bankInput.value;
		let startDateValue = startDateInput.value;
		let depositValue = depositInput.value;
		let termValue = termInput.value;
		let apyValue = aprInput.value;

		let encryptedStartDate = encryptText(startDateValue, key, iv);
		let encryptedDeposit = encryptText(depositValue, key, iv);
		let encryptedTerm = encryptText(termValue, key, iv);
		let encryptedApy = encryptText(apyValue, key, iv);

		const xhr = new XMLHttpRequest();
		xhr.open("POST", "/addCD");
		xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8");

		const body = JSON.stringify({
			bank: bankValue.toString(),
			startDate: encryptedStartDate.toString(),
			deposit: encryptedDeposit.toString(),
			term: encryptedTerm.toString(),
			apy: encryptedApy.toString()
		});

		xhr.onreadystatechange = function() {
			if(xhr.readyState === XMLHttpRequest.DONE) {
				if(xhr.status === 200) {
					window.location.reload();
				} else {
					console.error("Failed to add CD:", xhr.status, xhr.statusText);
				}
			}
		};

		xhr.send(body);
		modal.close();

	});
});
