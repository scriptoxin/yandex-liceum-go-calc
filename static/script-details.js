const urlParams = new URLSearchParams(window.location.search);
const exprId = urlParams.get('id');

async function loadExpressionDetails() {
  const response = await fetch(`${API_BASE}/expressions/${exprId}`);
  const data = await response.json();

  const container = document.getElementById('expression-details');
  container.innerHTML = `
        <p><b>ID:</b> ${data.expression.id}</p>
        <p><b>Выражение:</b> ${data.expression.expression}</p>
        <p><b>Статус:</b> <span class="status ${data.expression.status}">${
    data.expression.status
  }</span></p>
        ${
          data.expression.result
            ? `<p><b>Результат:</b> ${data.expression.result}</p>`
            : ''
        }
        <p><b>Время создания:</b> ${new Date(
          data.expression.created_at
        ).toLocaleString()}</p>
    `;
}

loadExpressionDetails();
setInterval(loadExpressionDetails, 2000);
