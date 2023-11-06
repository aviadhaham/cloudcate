// Select DOM elements 
const form = document.getElementById('search-form');
const resourceNameInput = document.getElementById('resource-name');
const resourceTypeSelect = document.getElementById('resource-type');
const resultsDiv = document.getElementById('search-results');
const searchSpinner = document.getElementById('search-spinner');

// Form submit handler
form.addEventListener('submit', (e) => {

  e.preventDefault();
  
  // Get form values
  const resourceName = resourceNameInput.value;
  const resourceType = resourceTypeSelect.value;

    // Make POST request to API
    const url = `search/?resource_name=${resourceName}&resource_type=${resourceType}`;
    searchSpinner.style.display = 'block';
    resultsDiv.innerHTML = '';
    fetch(url)
    .then(response => response.json())
    .then(data => {
        // Build table dynamically
        let table = resultsDiv.appendChild(document.createElement('table'));
        table.style.border = '1px solid black';
        table.style.borderCollapse = 'collapse';

        let tr = table.appendChild(document.createElement('tr'));
        Object.keys(data.results[0]).forEach((column) => {
            const th = document.createElement('th');
            th.innerText = column;
            th.style.border = '1px solid black';
            th.style.padding = '5px';
            tr.appendChild(th);
        });

        data.results.forEach((result) => {
            let tr = table.appendChild(document.createElement('tr'));
            Object.values(result).forEach((value) => {
                const td = document.createElement('td');
                td.innerText = value;
                td.style.border = '1px solid black';
                td.style.padding = '5px';
                tr.appendChild(td);
            });
        });
    })
    .catch((error) => {
      console.log(error);
    })
    .finally(() => {
      searchSpinner.style.display = 'none';
    });

});