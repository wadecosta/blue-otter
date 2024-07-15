function getValue() { 
	// Get the input element by its ID 
	let inputField = document.getElementById("key");

	// Get the value of the input field 
	let value = inputField.value; 
        console.log(value);
	sessionStorage.setItem("key", value);
}
