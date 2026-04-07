const productForm = document.getElementById("product-form");
const productNameInput = document.getElementById("product-name");
const formError = document.getElementById("form-error");
const reloadButton = document.getElementById("reload-button");
const messageBox = document.getElementById("message");
const emptyState = document.getElementById("empty-state");
const tableContainer = document.getElementById("table-container");
const productsBody = document.getElementById("products-body");

productForm.addEventListener("submit", handleCreateProduct);
reloadButton.addEventListener("click", loadProducts);

loadProducts();

async function loadProducts() {
  hideMessage();

  try {
    const response = await fetch("/products");

    if (!response.ok) {
      throw new Error(await getErrorMessage(response, "Не удалось загрузить каталог"));
    }

    const products = await response.json();
    renderProducts(products);
  } catch (error) {
    renderProducts([]);
    showMessage(error.message, "error");
  }
}

function renderProducts(products) {
  productsBody.innerHTML = "";

  if (products.length === 0) {
    tableContainer.classList.add("hidden");
    emptyState.classList.remove("hidden");
    return;
  }

  tableContainer.classList.remove("hidden");
  emptyState.classList.add("hidden");

  for (const product of products) {
    const row = document.createElement("tr");

    row.innerHTML = `
      <td>${product.id}</td>
      <td>${product.name}</td>
      <td class="actions">
        <button type="button" class="button button-secondary" data-action="edit">
          Изменить
        </button>
        <button type="button" class="button button-danger" data-action="delete">
          Удалить
        </button>
      </td>
    `;

    const editButton = row.querySelector('[data-action="edit"]');
    const deleteButton = row.querySelector('[data-action="delete"]');

    editButton.addEventListener("click", () => handleEditProduct(product));
    deleteButton.addEventListener("click", () => handleDeleteProduct(product));

    productsBody.appendChild(row);
  }
}

async function handleCreateProduct(event) {
  event.preventDefault();
  clearFormError();

  const name = productNameInput.value.trim();

  if (name === "") {
    showFormError("Введите название товара.");
    return;
  }

  try {
    const response = await fetch("/products", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ name })
    });

    if (!response.ok) {
      throw new Error(await getErrorMessage(response, "Не удалось добавить товар"));
    }

    productForm.reset();
    showMessage("Товар успешно добавлен.", "success");
    await loadProducts();
  } catch (error) {
    showFormError(error.message);
  }
}

async function handleEditProduct(product) {
  const newName = window.prompt("Введите новое название товара:", product.name);

  if (newName === null) {
    return;
  }

  const trimmedName = newName.trim();

  if (trimmedName === "") {
    showMessage("Название товара не может быть пустым.", "error");
    return;
  }

  try {
    const response = await fetch(`/products/${product.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ name: trimmedName })
    });

    if (!response.ok) {
      throw new Error(await getErrorMessage(response, "Не удалось обновить товар"));
    }

    showMessage("Название товара обновлено.", "success");
    await loadProducts();
  } catch (error) {
    showMessage(error.message, "error");
  }
}

async function handleDeleteProduct(product) {
  const confirmed = window.confirm(`Удалить товар "${product.name}"?`);

  if (!confirmed) {
    return;
  }

  try {
    const response = await fetch(`/products/${product.id}`, {
      method: "DELETE"
    });

    if (!response.ok) {
      throw new Error(await getErrorMessage(response, "Не удалось удалить товар"));
    }

    showMessage("Товар удалён.", "success");
    await loadProducts();
  } catch (error) {
    showMessage(error.message, "error");
  }
}

function showMessage(text, type) {
  messageBox.textContent = text;
  messageBox.className = `message ${type}`;
}

function hideMessage() {
  messageBox.textContent = "";
  messageBox.className = "message hidden";
}

function showFormError(text) {
  formError.textContent = text;
  formError.classList.remove("hidden");
}

function clearFormError() {
  formError.textContent = "";
  formError.classList.add("hidden");
}

async function getErrorMessage(response, fallbackMessage) {
  try {
    const data = await response.json();

    if (data && data.error) {
      return data.error;
    }

    return fallbackMessage;
  } catch {
    return fallbackMessage;
  }
}