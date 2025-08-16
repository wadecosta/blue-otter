document.addEventListener("DOMContentLoaded", function () {
	const input = document.getElementById("searchInput");
	const bankList = document.getElementById("bankList");

	let items = Array.from(bankList.getElementsByClassName("bank-item"));
	let activeIndex = -1;

	/* Sort alphabetically */
	items.sort((a, b) => a.dataset.name.localeCompare(b.dataset.name));
	items.forEach(item => bankList.appendChild(item));

	function updateHighlight(index) {
		items.forEach((item, i) => {
			item.classList.toggle("active", i === index && item.style.display !== "none");
		});
	}

	input.addEventListener("input", function () {
		const query = this.value.trim().toLowerCase();
		activeIndex = -1;

		items.forEach(item => {
			const nameDiv = item.querySelector(".bank-name");
			if (!nameDiv) return;

			const originalName = item.dataset.name;
			const nameLower = originalName.toLowerCase();

			if (nameLower.includes(query)) {
				const regex = new RegExp(`(${query})`, "gi");
				nameDiv.innerHTML = query === "" ? originalName : originalName.replace(regex, `<mark>$1</mark>`);
				item.style.display = "flex";
			} else {
				item.style.display = "none";
			}
			item.classList.remove("active");
		});
	});

	input.addEventListener("keydown", function (e) {
		const visibleItems = items.filter(i => i.style.display !== "none");

		if (e.key === "ArrowDown") {
			e.preventDefault();
			activeIndex = (activeIndex + 1) % visibleItems.length;
			updateHighlightIndex(visibleItems);
		} else if (e.key === "ArrowUp") {
			e.preventDefault();
			activeIndex = (activeIndex - 1 + visibleItems.length) % visibleItems.length;
			updateHighlightIndex(visibleItems);
		} else if (e.key === "Enter" && activeIndex >= 0) {
			e.preventDefault();

			const target = visibleItems[activeIndex];

			if (target) {
				target.scrollIntoView({ behavior: "smooth", block: "center" });
				// Optional: click the delete button or another action
				target.querySelector("button").click();
			}
		}
	});

	function updateHighlightIndex(visibleItems) {
		visibleItems.forEach((item, i) => {
			item.classList.toggle("active", i === activeIndex);
		});
	}
});
