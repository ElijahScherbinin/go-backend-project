#!/bin/bash

# Проверяем количество аргументов
if [ $# -ne 1 ]; then
    echo "Использование: $0 <пароль>"
    exit 1
fi

PASSWORD=$1

execute_sql() {
    local sql_file=$1
    
    echo "Выполнение $sql_file..."
    
    if [ ! -f "$sql_file" ]; then
        echo "Файл $sql_file не найден"
        return 1
    fi
    
    PGPASSWORD=$PASSWORD psql -h localhost -p 5432 -U admin user_service < "$sql_file"
    
    if [ $? -ne 0 ]; then
        echo "Ошибка при выполнении $sql_file"
        return 1
    fi
    
    echo "Скрипт $sql_file выполнен успешно"
}

echo "Начало выполнения скриптов"
execute_sql "$(pwd)/user-service/script/sql/create_tables.sql"
execute_sql "$(pwd)/user-service/script/sql/create_users.sql"
echo "Все скрипты выполнены"
