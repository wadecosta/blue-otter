document.addEventListener('DOMContentLoaded', (event) => {
	
	let modal = document.getElementById("addBankModal");
	if (!modal) {
		console.error("Modal element not found");
	}

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

	let form = modal.querySelector("form");
	form.addEventListener("submit", function(event) {
		event.preventDefault();

		const formData = new FormData(form);

		/* Send the form data (including the file */
		fetch("/addBank", {
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
