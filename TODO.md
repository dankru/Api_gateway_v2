# API

1. Добавь ручки, масимально просто в одном main.go, хранить данные в map или слайс
    1. `POST` `/user` request body{name, age, ...}, resp id (uuid), 404, 500, 201 (created) 
    2. `GET` `/user/:id` resp body{id, name, age, ...}, 404, 500, 200
    3. `DELETE` `/user/:id` 404, 500, 200
    4. `PUT` `/user` resp body{id, name, age, ...}, 404, 500, 200

Протыкай в postman что робит

1. слайс vs массив, структура слайс, что происходит при append (cap)
2. map, бакеты, коллизии, эвакуации, чек что такое swiss table в новой версии
3. interface, что под капотом, solid, опп в golang чек примеры
4. Приведение типов