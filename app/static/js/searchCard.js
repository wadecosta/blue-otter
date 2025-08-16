document.addEventListener("DOMContentLoaded", function () {
	const input = document.getElementById("searchCardInput");
	const cardList = document.getElementById("cardList");

	let items = Array.from(cardList.getElementsByClassName("card-item"));
	let activeIndex = -1;

	// Sort items alphabetically
	items.sort((a, b) => {
		const nameA = a.dataset.name.toLowerCase();
		const nameB = b.dataset.name.toLowerCase();
		return nameA.localeCompare(nameB);
	});

	// Re-append sorted
	items.forEach(item => cardList.appendChild(item));

	function updateHighlightIndex(visibleItems) {
		visibleItems.forEach((item, i) => {
			item.classList.toggle("active", i === activeIndex);
		});
	}

	input.addEventListener("input", function () {
		const query = this.value.trim().toLowerCase();
		activeIndex = -1;

		items.forEach(item => {
			const nameDiv = item.querySelector(".card-name");
			if (!nameDiv) return;

			const originalName = item.dataset.name;
			const nameLower = originalName.toLowerCase();

			if (nameLower.includes(query)) {
				if (query === "") {
					nameDiv.innerHTML = originalName;
				} else {
					const regex = new RegExp(`(${query})`, "gi");
					nameDiv.innerHTML = originalName.replace(regex, `<mark>$1</mark>`);
				}
				item.style.display = "flex";
			} else {
				item.style.display = "none";
			}
			item.classList.remove("active");
		});
	});

	input.addEventListener("keydown", function (e) {
		const visibleItems = items.filter(i => i.style.display !== "none");

		if (visibleItems.length === 0) return;

		if (e.key === "ArrowDown") {
			e.preventDefault();
			activeIndex = (activeIndex + 1) % visibleItems.length;
			updateHighlightIndex(visibleItems);
			visibleItems[activeIndex].scrollIntoView({ behavior: "smooth", block: "center" });
		}

		if (e.key === "ArrowUp") {
			e.preventDefault();
			activeIndex = (activeIndex - 1 + visibleItems.length) % visibleItems.length;
			updateHighlightIndex(visibleItems);
			visibleItems[activeIndex].scrollIntoView({ behavior: "smooth", block: "center" });
		}

		if (e.key === "Enter" && activeIndex >= 0) {
			e.preventDefault();
			const selected = visibleItems[activeIndex];

			if (selected) {
				const deleteBtn = selected.querySelector("button");
				if (deleteBtn) { 
					deleteBtn.focus();
				}
			}
		}
	});
});
