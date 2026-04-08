const productForm = document.getElementById("product-form"); // Берём форму добавления товара.
const productNameInput = document.getElementById("product-name"); // Берём поле названия.
const productCategoryInput = document.getElementById("product-category"); // Берём поле категории.
const productPriceInput = document.getElementById("product-price"); // Берём поле цены.
const formError = document.getElementById("form-error"); // Берём блок ошибок формы.
const reloadButton = document.getElementById("reload-button"); // Берём кнопку обновления списка.
const messageBox = document.getElementById("message"); // Берём общий блок сообщений.
const emptyState = document.getElementById("empty-state"); // Берём блок пустого состояния.
const tableContainer = document.getElementById("table-container"); // Берём контейнер таблицы.
const productsBody = document.getElementById("products-body"); // Берём тело таблицы, куда будем вставлять строки.

productForm.addEventListener("submit", handleCreateProduct); // Подписываемся на отправку формы.
reloadButton.addEventListener("click", loadProducts); // Подписываемся на нажатие кнопки обновления списка.

loadProducts(); // Сразу загружаем список товаров при открытии страницы.

async function loadProducts() { // Функция загружает товары с API.
  hideMessage(); // Скрываем старые сообщения перед новой загрузкой.

  try { // Начинаем блок try для безопасной работы с fetch.
    const response = await fetch("/products"); // Делаем GET-запрос к API на получение списка товаров.

    if (!response.ok) { // Если сервер ответил неуспешным статусом...
      throw new Error(await getErrorMessage(response, "Не удалось загрузить каталог")); // ...создаём ошибку с текстом.
    }

    const products = await response.json(); // Преобразуем тело ответа из JSON в JavaScript-массив.
    renderProducts(products); // Передаём массив товаров в функцию отрисовки таблицы.
  } catch (error) { // Если что-то пошло не так...
    renderProducts([]); // ...очищаем таблицу.
    showMessage(error.message, "error"); // ...показываем сообщение об ошибке.
  }
}

function renderProducts(products) { // Функция рисует таблицу товаров.
  productsBody.innerHTML = ""; // Очищаем текущие строки таблицы.

  if (products.length === 0) { // Если товаров нет...
    tableContainer.classList.add("hidden"); // ...скрываем таблицу.
    emptyState.classList.remove("hidden"); // ...показываем пустое состояние.
    return; // Прерываем выполнение функции.
  }

  tableContainer.classList.remove("hidden"); // Показываем таблицу.
  emptyState.classList.add("hidden"); // Скрываем блок пустого состояния.

  for (const product of products) { // Проходим по каждому товару.
    const row = document.createElement("tr"); // Создаём новую строку таблицы.
    const price = Number(product.price)
    const priceText = Number.isFinite(price) ? price.toFixed(2) : "_";
    row.innerHTML = ` 
      <td>${product.id}</td> 
      <td>${product.name}</td> 
      <td>${product.category}</td> 
      <td>${priceText}</td> 
      <td class="actions"> 
        <button type="button" class="button button-secondary" data-action="edit">Изменить</button> 
        <button type="button" class="button button-danger" data-action="delete">Удалить</button> 
      </td>
    `; // Заполняем строку товара: id, name, category, price и кнопки действий.

    const editButton = row.querySelector('[data-action="edit"]'); // Находим кнопку редактирования внутри строки.
    const deleteButton = row.querySelector('[data-action="delete"]'); // Находим кнопку удаления внутри строки.

    editButton.addEventListener("click", () => handleEditProduct(product)); // Навешиваем обработчик на редактирование.
    deleteButton.addEventListener("click", () => handleDeleteProduct(product)); // Навешиваем обработчик на удаление.

    productsBody.appendChild(row); // Добавляем готовую строку в таблицу.
  }
}

async function handleCreateProduct(event) { // Функция срабатывает при отправке формы создания товара.
  event.preventDefault(); // Отменяем стандартное поведение формы, чтобы страница не перезагружалась.
  clearFormError(); // Очищаем старые ошибки формы.

  const name = productNameInput.value.trim(); // Читаем и очищаем название товара.
  const category = productCategoryInput.value.trim(); // Читаем и очищаем категорию товара.
  const price = Number(productPriceInput.value); // Читаем цену и превращаем её в число.

  if (name === "") { // Если название пустое...
    showFormError("Введите название товара."); // ...показываем ошибку.
    return; // Прерываем выполнение функции.
  }

  if (category === "") { // Если категория пустая...
    showFormError("Введите категорию товара."); // ...показываем ошибку.
    return; // Прерываем выполнение функции.
  }

  if (!Number.isFinite(price) || price <= 0) { // Если цена не является корректным числом больше нуля...
    showFormError("Введите корректную цену больше нуля."); // ...показываем ошибку.
    return; // Прерываем выполнение функции.
  }

  try { // Начинаем безопасный блок отправки запроса.
    const response = await fetch("/products", { // Делаем POST-запрос на создание товара.
      method: "POST", // Указываем HTTP-метод POST.
      headers: { // Описываем заголовки запроса.
        "Content-Type": "application/json" // Сообщаем серверу, что отправляем JSON.
      },
      body: JSON.stringify({ name, category, price }) // Превращаем объект с данными в JSON-строку.
    });

    if (!response.ok) { // Если сервер ответил ошибкой...
      throw new Error(await getErrorMessage(response, "Не удалось добавить товар")); // ...создаём ошибку.
    }

    productForm.reset(); // Очищаем форму после успешного создания товара.
    showMessage("Товар успешно добавлен.", "success"); // Показываем сообщение об успехе.
    await loadProducts(); // Перезагружаем таблицу товаров.
  } catch (error) { // Если произошла ошибка...
    showFormError(error.message); // ...показываем её в блоке ошибок формы.
  }
}

async function handleEditProduct(product) { // Функция редактирует существующий товар.
  const newName = window.prompt("Введите новое название товара:", product.name); // Просим ввести новое название.
  if (newName === null) { // Если пользователь нажал Отмена...
    return; // ...прерываем функцию.
  }

  const newCategory = window.prompt("Введите новую категорию товара:", product.category); // Просим ввести новую категорию.
  if (newCategory === null) { // Если пользователь нажал Отмена...
    return; // ...прерываем функцию.
  }

  const newPriceRaw = window.prompt("Введите новую цену товара:", String(product.price)); // Просим ввести новую цену.
  if (newPriceRaw === null) { // Если пользователь нажал Отмена...
    return; // ...прерываем функцию.
  }

  const name = newName.trim(); // Очищаем новое название.
  const category = newCategory.trim(); // Очищаем новую категорию.
  const price = Number(newPriceRaw); // Превращаем введённую цену в число.

  if (name === "") { // Если название пустое...
    showMessage("Название товара не может быть пустым.", "error"); // ...показываем ошибку.
    return; // Прерываем выполнение функции.
  }

  if (category === "") { // Если категория пустая...
    showMessage("Категория товара не может быть пустой.", "error"); // ...показываем ошибку.
    return; // Прерываем выполнение функции.
  }

  if (!Number.isFinite(price) || price <= 0) { // Если цена невалидная...
    showMessage("Цена должна быть больше нуля.", "error"); // ...показываем ошибку.
    return; // Прерываем выполнение функции.
  }

  try { // Начинаем безопасный блок обновления.
    const response = await fetch(`/products/${product.id}`, { // Делаем PUT-запрос на обновление товара.
      method: "PUT", // Указываем HTTP-метод PUT.
      headers: { // Описываем заголовки запроса.
        "Content-Type": "application/json" // Сообщаем серверу, что отправляем JSON.
      },
      body: JSON.stringify({ name, category, price }) // Превращаем новые данные в JSON-строку.
    });

    if (!response.ok) { // Если сервер ответил неуспешно...
      throw new Error(await getErrorMessage(response, "Не удалось обновить товар")); // ...создаём ошибку.
    }

    showMessage("Товар успешно обновлён.", "success"); // Показываем сообщение об успехе.
    await loadProducts(); // Перезагружаем список товаров.
  } catch (error) { // Если произошла ошибка...
    showMessage(error.message, "error"); // ...показываем сообщение об ошибке.
  }
}

async function handleDeleteProduct(product) { // Функция удаляет товар.
  const confirmed = window.confirm(`Удалить товар "${product.name}"?`); // Просим пользователя подтвердить удаление.

  if (!confirmed) { // Если пользователь отказался...
    return; // ...прерываем функцию.
  }

  try { // Начинаем безопасный блок удаления.
    const response = await fetch(`/products/${product.id}`, { // Делаем DELETE-запрос.
      method: "DELETE" // Указываем HTTP-метод DELETE.
    });

    if (!response.ok) { // Если сервер ответил ошибкой...
      throw new Error(await getErrorMessage(response, "Не удалось удалить товар")); // ...создаём ошибку.
    }

    showMessage("Товар удалён.", "success"); // Показываем сообщение об успехе.
    await loadProducts(); // Перезагружаем список товаров.
  } catch (error) { // Если произошла ошибка...
    showMessage(error.message, "error"); // ...показываем сообщение об ошибке.
  }
}

function showMessage(text, type) { // Функция показывает общее сообщение пользователю.
  messageBox.textContent = text; // Записываем текст сообщения.
  messageBox.className = `message ${type}`; // Назначаем CSS-классы по типу сообщения.
}

function hideMessage() { // Функция скрывает блок общих сообщений.
  messageBox.textContent = ""; // Очищаем текст сообщения.
  messageBox.className = "message hidden"; // Возвращаем скрытое состояние.
}

function showFormError(text) { // Функция показывает ошибку формы.
  formError.textContent = text; // Записываем текст ошибки.
  formError.classList.remove("hidden"); // Показываем блок ошибки.
}

function clearFormError() { // Функция очищает ошибку формы.
  formError.textContent = ""; // Удаляем текст ошибки.
  formError.classList.add("hidden"); // Скрываем блок ошибки.
}

async function getErrorMessage(response, fallbackMessage) { // Функция пытается достать текст ошибки из JSON-ответа API.
  try { // Начинаем блок try.
    const data = await response.json(); // Читаем JSON-ответ сервера.
    if (data && data.error) { // Если в ответе есть поле error...
      return data.error; // ...возвращаем его.
    }
    return fallbackMessage; // Если поля error нет, возвращаем запасное сообщение.
  } catch { // Если JSON вообще не удалось прочитать...
    return fallbackMessage; // ...возвращаем запасное сообщение.
  }
}