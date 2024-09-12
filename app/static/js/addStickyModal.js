document.addEventListener('DOCContentLoaded', (event) => {
	let modal = document.getElementById("addStickyModal");
	let span = document.querySelector(".closeAddSticky");

	let key = sessionStorage.getItem("key");
        let iv = document.getElementById("iv").value;

	document.querySelectorAdd('.open-add-sticky-button').forEach(button => {
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
