document.addEventListener('DOMContentLoaded', () => {
	const modalElement = document.getElementById("warningModal");
	const form = modalElement?.querySelector("form");

	let id;
	let type;

	if (!modalElement) {
		console.error("Required DOM elements not found");
		return;
	} 

	const bootstrapModal = new bootstrap.Modal(modalElement);

	document.querySelectorAll('.open-warning-modal-button').forEach(button => {
		button.addEventListener('click', function () {
			id = this.getAttribute("data-id");
			type = this.getAttribute("data-type");

			bootstrapModal.show();
		});
	});

	/* Reset form when modal is closed */
	modalElement.addEventListener('hiden.bs.modal', () => {
		form.reset();
	});

	/* Submit form */
	form.addEventListener("submit", (event) => {
		event.preventDefault();

		let message = `Removing id=${id} and type=${type}`;
		console.log(message);

		postDelete(type, id);
	});
});

function postDelete(type, id) {

	if (type === "sticky") {
		if (!id || isNaN(id) ) {
			console.error(`Invalid id: ${id}`);
			return;
		}
		postDeleteSticky(id);	
	}
	/* TODO : Impliment same feature for other types */
	else {
		console.log(`TYPE=${type} has not been implemented.`);
	}
}
