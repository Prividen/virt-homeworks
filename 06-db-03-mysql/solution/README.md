# Домашняя работа по занятию "6.3. MySQL"

> 1. Используя docker поднимите инстанс MySQL (версию 8). Данные БД сохраните в volume.

```
mak@test-xu20:~$ docker run --rm -d --name mysql -v mysql-db:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=MyPassw0rd mysql:8
```

> Найдите команду для выдачи статуса БД и **приведите в ответе** из ее вывода версию сервера БД.

```
mysql> status
--------------
mysql  Ver 8.0.25 for Linux on x86_64 (MySQL Community Server - GPL)
...
Server version:		8.0.25 MySQL Community Server - GPL
```

> Подключитесь к восстановленной БД и получите список таблиц из этой БД.
```
mak@test-xu20:~$ docker exec -it -e LANG=en_US.utf8 mysql mysql -pMyPassw0rd --skip-pager test_db
mysql> show tables;
+-------------------+
| Tables_in_test_db |
+-------------------+
| orders            |
+-------------------+
1 row in set (0.00 sec)
```

> **Приведите в ответе** количество записей с `price` > 300.

``` 
mysql> select count(*) from orders where price > 300;
+----------+
| count(*) |
+----------+
|        1 |
+----------+
```

---

> 2. Создайте пользователя test в БД c паролем test-pass, используя:
> - плагин авторизации mysql_native_password
> - срок истечения пароля - 180 дней 
> - количество попыток авторизации - 3 
> - максимальное количество запросов в час - 100
> - аттрибуты пользователя:
>    - Фамилия "Pretty"
>    - Имя "James"

```
mysql> create user test IDENTIFIED WITH mysql_native_password BY 'test-pass' with MAX_QUERIES_PER_HOUR 100 PASSWORD EXPIRE INTERVAL 180 DAY FAILED_LOGIN_ATTEMPTS 3 ATTRIBUTE '{"Имя": "James", "Фамилия": "Pretty"}';
Query OK, 0 rows affected (0.01 sec)
```
Для поддержки ввода UTF8/кириллицы пришлось немножко локали в контейнер доустановить/сконфигурировать.

> Предоставьте привелегии пользователю `test` на операции SELECT базы `test_db`.

```
mysql> grant SELECT on test_db.* to 'test';
Query OK, 0 rows affected (0.00 sec)
```
   
>Используя таблицу INFORMATION_SCHEMA.USER_ATTRIBUTES получите данные по пользователю `test` и 
**приведите в ответе к задаче**.

```
mysql> select * from INFORMATION_SCHEMA.USER_ATTRIBUTES where user='test';
+------+------+-------------------------------------------------+
| USER | HOST | ATTRIBUTE                                       |
+------+------+-------------------------------------------------+
| test | %    | {"Имя": "James", "Фамилия": "Pretty"}           |
+------+------+-------------------------------------------------+
1 row in set (0.00 sec)
```

---
> 3. Установите профилирование `SET profiling = 1`.
> Исследуйте, какой `engine` используется в таблице БД `test_db` и **приведите в ответе**.

```
mysql> SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES where table_name='orders';
+------------+--------+
| TABLE_NAME | ENGINE |
+------------+--------+
| orders     | InnoDB |
+------------+--------+
```

> Измените `engine` и **приведите время выполнения и запрос на изменения из профайлера в ответе**:
> - на `MyISAM`
> - на `InnoDB`

```
mysql> SHOW PROFILES;
+----------+------------+----------------------------------+
| Query_ID | Duration   | Query                            |
+----------+------------+----------------------------------+
|        1 | 0.02834000 | alter table orders engine=MyISAM |
|        2 | 0.02701925 | alter table orders engine=InnoDB |
+----------+------------+----------------------------------+
2 rows in set, 1 warning (0.00 sec)
```

---
> 4. Изучите файл `my.cnf` в директории /etc/mysql.
> Измените его согласно ТЗ (движок InnoDB):
> - Скорость IO важнее сохранности данных
> - Нужна компрессия таблиц для экономии места на диске
> - Размер буффера с незакомиченными транзакциями 1 Мб
> - Буффер кеширования 30% от ОЗУ
> - Размер файла логов операций 100 Мб
>
> Приведите в ответе измененный файл `my.cnf`.

```
$ grep -v '^#' my.cnf


[mysqld]
pid-file        = /var/run/mysqld/mysqld.pid
socket          = /var/run/mysqld/mysqld.sock
datadir         = /var/lib/mysql
secure-file-priv= NULL

innodb_flush_log_at_trx_commit = 2
innodb_file_per_table = 1
innodb_log_buffer_size = 1M
innodb_buffer_pool_size = 9830M
innodb_buffer_pool_chunk_size = 9830M
innodb_buffer_pool_instances = 1
innodb_log_file_size = 100M


!includedir /etc/mysql/conf.d/
```

