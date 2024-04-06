// Select DOM elements 
const form = document.getElementById('search-form');
const resourceNameInput = document.getElementById('resource-name');
const resourceTypeSelect = document.getElementById('resource-type');
const resultsDiv = document.getElementById('search-results');
const searchSpinner = document.getElementById('search-spinner');
const resultsCount = document.getElementById('results-count');

// Form submit handler
form.addEventListener('submit', (e) => {

  e.preventDefault();
  
  // Get form values
  const resourceName = resourceNameInput.value;
  const resourceType = resourceTypeSelect.value;

  if (resourceName === '') {
    console.log('Resource name is empty.');
    return; // Do not trigger the event if resourceName is empty
  }

  // Make POST request to API
  const url = `search/?resource_name=${resourceName}&resource_type=${resourceType}`;
  searchSpinner.style.display = 'block';
  resultsCount.style.display = 'block';
  resultsCount.innerHTML = '';
  resultsDiv.innerHTML = '';
  fetch(url)
    .then(response => response.json())
    .then(data => {
      // Build table dynamically
      let table = resultsDiv.appendChild(document.createElement('table'));
      table.style.borderCollapse = 'collapse';
      table.classList.add('table');

      let tr = table.appendChild(document.createElement('tr'));
      const thRowNumber = document.createElement('th');
      thRowNumber.innerText = '#';
      thRowNumber.style.border = '1px solid black';
      thRowNumber.style.padding = '5px';
      tr.appendChild(thRowNumber);

      Object.keys(data.results[0]).forEach((column) => {
        const th = document.createElement('th');
        th.innerText = column;
        th.style.border = '1px solid black';
        th.style.padding = '5px';
        tr.appendChild(th);
      });

      data.results.forEach((result, index) => {
        let tr = table.appendChild(document.createElement('tr'));
        const tdRowNumber = document.createElement('td');
        tdRowNumber.innerText = index + 1;
        tdRowNumber.style.border = '1px solid black';
        tdRowNumber.style.padding = '5px';
        tr.appendChild(tdRowNumber);

        Object.values(result).forEach((value) => {
          const td = document.createElement('td');
          td.innerText = value;
          td.style.border = '1px solid black';
          td.style.padding = '5px';
          td.style.wordBreak = 'break-all';
          td.style.overflowWrap = 'break-word';
          tr.appendChild(td);
        });
      });

      resultsCount.innerText = `Results: ${data.results.length}`;
    })
    .catch((error) => {
      console.log(error);
      resultsCount.innerText = 'Error occurred while fetching results.';
    })
    .finally(() => {
      searchSpinner.style.display = 'none';
    });

});
