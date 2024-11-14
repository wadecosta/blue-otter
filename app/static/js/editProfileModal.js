document.addEventListener("DOMContentLoaded", function() {
	const editProfileButton = document.getElementById("editProfileButton");
	const editProfileModal = document.getElementById("editProfileModal");
	const closeEditProfile = document.querySelector(".closeEditProfile");

	let email_data = document.getElementById("email").value;
	let username_data = document.getElementById("username").value;

	editProfileButton.addEventListener("click", () => {
		username.value = username_data;
		email.value = email_data;

		editProfileModal.showModal();
	});

	closeEditProfile.addEventListener("click", () => {
		editProfileModal.close();
	});
});
