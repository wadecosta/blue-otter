document.addEventListener('DOMContentLoaded', (event) => {
	
	let modal = document.getElementById("addCreditCardModal");
	if (!modal) {
		console.error("Modal element not found");
	}

	let span = document.querySelector(".closeAddCreditCard");

	document.querySelectorAll('.open-add-credit-card-button').forEach(button => {
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

		const formData = new FormData(form);

		/* Send the form data (including the file */
		fetch("/addCard", {
			method: "POST",
			body: formData,
		})
		.then(response => response.json())
		.then(data => {
			alert("Upload Successful");
			modal.close();
		})
		.catch(error => {
			alert("Error uploading file");
			console.error(error);
		});
	});
});
