console.log("Ran addBankModal");

document.addEventListener('DOMContentLoaded', (event) => {
	
	let modal = document.getElementById("addBankModal");
	if (!modal) {
		console.error("Modal element not found");
	}

	console.log("Found modal");

	let span = document.querySelector(".closeAddBank");

	document.querySelectorAll('.open-add-bank-button').forEach(button => {
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
});
