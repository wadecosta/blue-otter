document.addEventListener("DOMContentLoaded", function() {
	const currentLocation = window.location.href;
	const navLinks = document.querySelectorAll(".nav-link");

	navLinks.forEach(link => {
		if (link.href === currentLocation) {
			link.classList.add("active");
		}
	});
});
