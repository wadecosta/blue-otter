document.addEventListener('DOMContentLoaded', (event) => {
	let modal = document.getElementById("editStickyModal");
	let span = document.querySelector(".closeEditSticky");
	let title = document.getElementById("title");
	   
	/* Needed for encryption/decryption */
	let key = sessionStorage.getItem("key");
	let iv = document.getElementById("iv").value;

	let id;
	let encryptedDescription;
	let encryptedTitle;

	document.querySelectorAll('.open-edit-sticky-button').forEach(button => {
		button.onclick = function() {
			id = this.getAttribute("data-id");

			encryptedTitle = this.getAttribute("data-title");
			decryptedTitle = decryptText(encryptedTitle, key, iv);
			title.value = decryptedTitle;

			encryptedDescription = this.getAttribute("data-description");
			let decryptedDescription = decryptText(encryptedDescription, key, iv);
			description.value = decryptedDescription;

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
		console.log("GOT YOU");

		let changedTitle = title.value;
		let changedDescription = description.value;

		console.log(changedTitle);
		console.log(changedDescription);

		let changedEncryptedTitle = encryptText(changedTitle, key, iv);
		let changedEncryptedDescription = encryptText(changedDescription, key, iv);

		console.log(changedEncryptedTitle);
		console.log(changedEncryptedDescription);
		console.log(id);

		const xhr = new XMLHttpRequest();
		xhr.open("POST", "/editSticky");
		xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8");

		const body = JSON.stringify({
			id: id.toString(),
			old_sticky_description: encryptedDescription,
			old_sticky_title: encryptedTitle.toString(),
			new_sticky_description: changedEncryptedDescription.toString(),
			new_sticky_title: changedEncryptedTitle.toString()
		});

		xhr.onreadystatechange = function() {
			if (xhr.readyState === XMLHttpRequest.DONE) {
				if (xhr.status === 200) {
					window.location.reload();
				} else {
					console.error("Failed to saved changes:", xhr.status, xhr.statusText);
				}
			}
		};

		xhr.send(body);
		modal.close();
	});
});
