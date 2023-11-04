// Select DOM elements 
const form = document.getElementById('search-form');
const resourceNameInput = document.getElementById('resource-name');
const resourceTypeSelect = document.getElementById('resource-type');
const resultsDiv = document.getElementById('search-results');
const spinner = document.querySelector('.spinner-border');

// Form submit handler
form.addEventListener('submit', (e) => {

  e.preventDefault();
  
  // Show loading spinner
  spinner.hidden = false;

  // Get form values
  const resourceName = resourceNameInput.value;
  const resourceType = resourceTypeSelect.value;

    // Make POST request to API
    const url = `search/?resource_name=${resourceName}&resource_type=${resourceType}`;
    fetch(url)
    .then(response => response.json())
    .then(data => {

        // Hide loading spinner
        spinner.hidden = true;

        // Display results
        resultsDiv.innerText = JSON.stringify(data);

    });

});