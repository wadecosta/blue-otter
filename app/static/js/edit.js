    document.addEventListener("DOMContentLoaded", function() {
        // Select all paragraph elements with IDs starting with "Card-CardBank-"
        var paragraphs = document.querySelectorAll("p[id^='Card-CardBank-']");

        // Iterate over each paragraph element
        paragraphs.forEach(function(paragraph) {
            // Add click event listener to each paragraph
            paragraph.addEventListener("click", function() {
                // Create a new textbox element
                var textbox = document.createElement("input");
                textbox.type = "text";
                // Set the value of the textbox to the current text content of the paragraph
                textbox.value = paragraph.textContent;
                // Set the ID of the textbox
                textbox.id = "Textbox-CardBank-" + paragraph.id.split("-")[3];
                // Replace the paragraph element with the textbox
                paragraph.parentNode.replaceChild(textbox, paragraph);

                // Add blur event listener to the textbox
                textbox.addEventListener("blur", function() {
                    // Create a new paragraph element
                    var newParagraph = document.createElement("p");
                    // Set the text content of the new paragraph to the value of the textbox
                    newParagraph.textContent = textbox.value;
                    // Set the ID of the new paragraph
                    newParagraph.id = paragraph.id;
                    // Replace the textbox with the new paragraph
                    textbox.parentNode.replaceChild(newParagraph, textbox);
                });

                // Add click event listener to the textbox
                textbox.addEventListener("click", function(event) {
                    // Prevent event from propagating to the parent elements
                    event.stopPropagation();
                });
            });
        });
    });
