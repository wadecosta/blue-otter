function getValue() { 
	// Get the input element by its ID 
	let inputField = document.getElementById("key");
          
	// Get the value of the input field 
	let value = inputField.value; 
        console.log(value);
	sessionStorage.setItem("key", value);

	/*	
	// Simulate a mouse click:
	//window.location.href = "/dashboard";
	// Create a form element
	var form = document.createElement("form");

	// Set the form's method to POST
	form.method = "POST";

	// Set the form's action to the desired URL route ("/vault")
	form.action = "/vault";

	// Append the form to the document body
	document.body.appendChild(form);

	// Submit the form
	form.submit();
	*/
} 
