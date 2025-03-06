document.addEventListener('DOMContentLoaded', (event) => {
    console.log("Script loaded successfully");

    let modal = document.getElementById("addStickyModal");
    if (!modal) {
        console.error("Modal element not found");
        return;
    }

    let span = document.querySelector(".closeAddSticky");

    let key = sessionStorage.getItem("key");
    let iv = document.getElementById("iv").value;

    console.log("Modal element:", modal);

    document.querySelectorAll('.open-add-sticky-button').forEach(button => {
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

	let titleInput = modal.querySelector("#title");
	let descriptionInput = modal.querySelector("#description");

	let titleValue = titleInput.value;
        let descriptionValue = descriptionInput.value;

	console.log(titleValue);
        console.log(descriptionValue);

	let encryptedTitle = encryptText(titleValue, key, iv);
	let encryptedDescription = encryptText(descriptionValue, key, iv);

	console.log(encryptedTitle);
	console.log(encryptedDescription);

	const xhr = new XMLHttpRequest();
	xhr.open("POST", "/addSticky");
	xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8");

	const body = JSON.stringify({
		title: encryptedTitle.toString(),
		description: encryptedDescription.toString()
	});

	xhr.onreadystatechange = function() {
		if (xhr.readyState === XMLHttpRequest.DONE) {
			if (xhr.status === 200) {
				window.location.reload();
			} else {
				console.error("Failed to add sticky:", xhr.status, xhr.statusText);
			}
		}
	};

	xhr.send(body);
	modal.close();
    });
});
