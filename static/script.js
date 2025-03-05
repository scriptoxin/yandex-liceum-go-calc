const API_BASE = 'http://localhost:8080/api/v1';

// Отправка нового выражения
async function submitExpression() {
  const exprInput = document.getElementById('expression');
  const response = await fetch(`${API_BASE}/calculate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      expression: exprInput.value,
    }),
  });

  if (response.ok) {
    exprInput.value = '';
    loadExpressions();
  } else {
    alert('Ошибка при отправке выражения');
  }
}

// Загрузка списка выражений
async function loadExpressions() {
  const response = await fetch(`${API_BASE}/expressions`);
  const data = await response.json();

  const list = document.getElementById('expressions-list');
  list.innerHTML = '';

  data.expressions.forEach((expr) => {
    const item = document.createElement('div');
    item.className = 'expression-item';
    item.innerHTML = `
            <div>
                <b>ID:</b> ${expr.id}<br>
                <b>Выражение:</b> ${expr.expression}
            </div>
            <div>
                <div class="status ${expr.status}">${expr.status}</div>
                ${
                  expr.result
                    ? `<div class="result">= ${expr.result}</div>`
                    : ''
                }
            </div>
        `;

    item.onclick = () => (window.location = `/expression.html?id=${expr.id}`);
    list.appendChild(item);
  });
}

// Автообновление каждые 5 секунд
setInterval(loadExpressions, 5000);
window.onload = loadExpressions;
