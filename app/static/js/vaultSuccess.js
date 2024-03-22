function getValue() { 
	// Get the input element by its ID 
	let inputField = document.getElementById("key");
	let AEShash = document.getElementById("hash");
          
	// Get the value of the input field 
	let value = inputField.value; 
        console.log(value);
	sessionStorage.setItem("key", value);

	let hash = AEShash.value;
	console.log(hash);

	// Simulate a mouse click:
	window.location.href = "/dashboard";
} 
