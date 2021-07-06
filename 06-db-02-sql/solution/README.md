# Домашняя работа по занятию "6.2. SQL"


> 1. Используя docker поднимите инстанс PostgreSQL (версию 12) c 2 volume, 
в который будут складываться данные БД и бэкапы.  
Приведите получившуюся команду или docker-compose манифест.
 
```
mak@test-xu20:~$ docker run --rm -d --name pg -v pg-db:/var/lib/postgresql/data -v pg-backups:/backups -e POSTGRES_PASSWORD=MyPassw0rd postgres:12 
310e2169eb54abdd589ee61b319502de65491bc643964b8ce62bb908428e3c83
mak@test-xu20:~$ docker exec -it pg psql -U postgres -P pager=off
psql (12.7 (Debian 12.7-1.pgdg100+1))
Type "help" for help.

postgres=#
```

---
> 2. В БД из задачи 1: 
> - создайте пользователя test-admin-user и БД test_db

```
postgres=# create user "test-admin-user";
CREATE ROLE
postgres=# create database "test_db";
CREATE DATABASE
```

- в БД test_db создайте таблицу orders и clients (спeцификация таблиц ниже)
  
```
postgres=# \c test_db
You are now connected to database "test_db" as user "postgres".

test_db=# create table orders (id serial primary key, "наименование" text, "цена" int);
CREATE TABLE

test_db=# create table clients (id serial primary key, "фамилия" text, "страна проживания" text, "заказ" int, constraint "orders_fk" foreign key ("заказ") references orders (id));
CREATE TABLE
test_db=# create index i_country on clients("страна проживания");
CREATE INDEX
```

> - предоставьте привилегии на все операции пользователю test-admin-user на таблицы БД test_db

```
test_db=# grant all on orders,clients to "test-admin-user";
GRANT
```

> - создайте пользователя test-simple-user  
> - предоставьте пользователю test-simple-user права на SELECT/INSERT/UPDATE/DELETE данных таблиц БД test_db

```
test_db=# create user "test-simple-user";
CREATE ROLE
test_db=# grant SELECT, INSERT, UPDATE, DELETE on orders,clients to "test-simple-user" ;
GRANT
```

> Приведите:
> - итоговый список БД после выполнения пунктов выше,

```
test_db=# \l
                                 List of databases
   Name    |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges   
-----------+----------+----------+------------+------------+-----------------------
 postgres  | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 template1 | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
           |          |          |            |            | postgres=CTc/postgres
 test_db   | postgres | UTF8     | en_US.utf8 | en_US.utf8 | 
```

> - описание таблиц (describe)

```
test_db=# \d orders
                               Table "public.orders"
    Column    |  Type   | Collation | Nullable |              Default               
--------------+---------+-----------+----------+------------------------------------
 id           | integer |           | not null | nextval('orders_id_seq'::regclass)
 наименование | text    |           |          | 
 цена         | integer |           |          | 
Indexes:
    "orders_pkey" PRIMARY KEY, btree (id)
Referenced by:
    TABLE "clients" CONSTRAINT "orders_fk" FOREIGN KEY ("заказ") REFERENCES orders(id)

test_db=# \d clients
                                  Table "public.clients"
      Column       |  Type   | Collation | Nullable |               Default               
-------------------+---------+-----------+----------+-------------------------------------
 id                | integer |           | not null | nextval('clients_id_seq'::regclass)
 фамилия           | text    |           |          | 
 страна проживания | text    |           |          | 
 заказ             | integer |           |          | 
Indexes:
    "clients_pkey" PRIMARY KEY, btree (id)
    "i_country" btree ("страна проживания")
Foreign-key constraints:
    "orders_fk" FOREIGN KEY ("заказ") REFERENCES orders(id)
```

> - SQL-запрос для выдачи списка пользователей с правами над таблицами test_db
> - список пользователей с правами над таблицами test_db

```
test_db=# select distinct(grantee) "User", table_catalog "Database", table_name "Table", array (select privilege_type from information_schema.role_table_grants t2 where t2.grantee=t1.grantee and t2.table_name=t1.table_name) "Privileges" from information_schema.role_table_grants t1 where table_catalog='test_db' and table_schema='public';
       User       | Database |  Table  |                        Privileges                         
------------------+----------+---------+-----------------------------------------------------------
 postgres         | test_db  | clients | {TRIGGER,REFERENCES,TRUNCATE,DELETE,UPDATE,SELECT,INSERT}
 postgres         | test_db  | orders  | {TRIGGER,REFERENCES,TRUNCATE,DELETE,UPDATE,SELECT,INSERT}
 test-admin-user  | test_db  | clients | {TRIGGER,REFERENCES,TRUNCATE,DELETE,UPDATE,SELECT,INSERT}
 test-admin-user  | test_db  | orders  | {TRIGGER,REFERENCES,TRUNCATE,DELETE,UPDATE,SELECT,INSERT}
 test-simple-user | test_db  | clients | {DELETE,UPDATE,SELECT,INSERT}
 test-simple-user | test_db  | orders  | {DELETE,UPDATE,SELECT,INSERT}
(6 rows)
```

---
> 3. Используя SQL синтаксис - наполните таблицы следующими тестовыми данными:  
> ...  
> Используя SQL синтаксис:
> - вычислите количество записей для каждой таблицы 
> - приведите в ответе:
>    - запросы 
>    - результаты их выполнения.

```
test_db=# insert into orders ("наименование", "цена") values ('Шоколад', 10), ('Принтер', 3000), ('Книга', 500), ('Монитор', 7000), ('Гитара', 4000);
INSERT 0 5

test_db=# insert into clients ("фамилия", "страна проживания") values ('Иванов Иван Иванович', 'USA'), ('Петров Петр Петрович', 'Canada'), ('Иоганн Себастьян Бах', 'Japan'), ('Ронни Джеймс Дио', 'Russia'), ('Ritchie Blackmore', 'Russia');
INSERT 0 5

test_db=# select count(*) from clients;
 count 
-------
     5
(1 row)

test_db=# select count(*) from orders;
 count 
-------
     5
(1 row)
```

---
> 4. Часть пользователей из таблицы clients решили оформить заказы из таблицы orders.
> Используя foreign keys свяжите записи из таблиц, согласно таблице:  
> ...  
> Приведите SQL-запросы для выполнения данных операций.

```
test_db=# update clients set "заказ"=(select id from orders where "наименование"='Книга') where "фамилия"='Иванов Иван Иванович';
UPDATE 1
test_db=# update clients set "заказ"=(select id from orders where "наименование"='Монитор') where "фамилия"='Петров Петр Петрович';
UPDATE 1
test_db=# update clients set "заказ"=(select id from orders where "наименование"='Гитара') where "фамилия"='Иоганн Себастьян Бах';
UPDATE 1
```

> Приведите SQL-запрос для выдачи всех пользователей, которые совершили заказ, а также вывод данного запроса.
```
test_db=# select c."фамилия" "ФИО", o."наименование" "Заказ" from clients c right outer join orders o on c."заказ"=o.id where c."заказ" is not null;
         ФИО          |  Заказ  
----------------------+---------
 Иванов Иван Иванович | Книга
 Петров Петр Петрович | Монитор
 Иоганн Себастьян Бах | Гитара
(3 rows)
```

---
> 5. Получите полную информацию по выполнению запроса выдачи всех пользователей из задачи 4 
(используя директиву EXPLAIN). Приведите получившийся результат и объясните что значат полученные значения.

```
test_db=# explain analyze select c."фамилия" "ФИО", o."наименование" "Заказ" from clients c right outer join orders o on c."заказ"=o.id where c."заказ" is not null;
                                                    QUERY PLAN                                                     
-------------------------------------------------------------------------------------------------------------------
 Hash Join  (cost=37.00..57.23 rows=806 width=64) (actual time=0.042..0.047 rows=3 loops=1)
   Hash Cond: (c."заказ" = o.id)
   ->  Seq Scan on clients c  (cost=0.00..18.10 rows=806 width=36) (actual time=0.016..0.018 rows=3 loops=1)
         Filter: ("заказ" IS NOT NULL)
         Rows Removed by Filter: 2
   ->  Hash  (cost=22.00..22.00 rows=1200 width=36) (actual time=0.014..0.015 rows=5 loops=1)
         Buckets: 2048  Batches: 1  Memory Usage: 17kB
         ->  Seq Scan on orders o  (cost=0.00..22.00 rows=1200 width=36) (actual time=0.005..0.008 rows=5 loops=1)
 Planning Time: 0.190 ms
 Execution Time: 0.081 ms
```
https://explain.depesz.com/s/HF3R

Тут происходит последовательное чтение обоих таблиц, clients с дополнительным фильтром 
`"заказ" IS NOT NULL`, и объединение результатов по условию  `c."заказ" = o.id`.  

Если попросить Postgres сделать анализ таблиц, то план запроса может измениться. Например, вместо 
последовательного сканирования таблицы orders будут использоваться индексы: 
```
                                                         QUERY PLAN                                                          
-----------------------------------------------------------------------------------------------------------------------------
 Nested Loop  (cost=0.15..21.57 rows=3 width=65) (actual time=0.030..0.042 rows=3 loops=1)
   ->  Seq Scan on clients c  (cost=0.00..1.05 rows=3 width=37) (actual time=0.015..0.017 rows=3 loops=1)
         Filter: ("заказ" IS NOT NULL)
         Rows Removed by Filter: 2
   ->  Index Scan using orders_pkey on orders o  (cost=0.15..6.84 rows=1 width=36) (actual time=0.005..0.005 rows=1 loops=3)
         Index Cond: (id = c."заказ")
 Planning Time: 0.424 ms
 Execution Time: 0.077 ms
```
https://explain.depesz.com/s/HjuY


---
> 6. Создайте бэкап БД test_db и поместите его в volume, предназначенный для бэкапов (см. Задачу 1).
> Остановите контейнер с PostgreSQL (но не удаляйте volumes).
> Поднимите новый пустой контейнер с PostgreSQL.
> Восстановите БД test_db в новом контейнере.  
> Приведите список операций, который вы применяли для бэкапа данных и восстановления. 

```
mak@test-xu20:~$ docker exec -it pg bash
root@310e2169eb54:/# pg_dump -U postgres -Fc test_db  >/backups/test_db.dump
root@310e2169eb54:/# exit
mak@test-xu20:~$ docker stop pg
pg
mak@test-xu20:~$ docker run --rm -d --name pg_new -v pg-db-new:/var/lib/postgresql/data -v pg-backups:/backups -e POSTGRES_PASSWORD=MyPassw0rd postgres:12 
de0700a8ac1f071ba98496e73532efaf45723d9ef9c094fee79b1b4eae6aaf15
mak@test-xu20:~$ docker exec -it pg_new bash
root@de0700a8ac1f:/# createdb -U postgres test_db
root@de0700a8ac1f:/# createuser -U postgres test-admin-user
root@de0700a8ac1f:/# createuser -U postgres test-simple-user
root@de0700a8ac1f:/# pg_restore -U postgres -d test_db /backups/test_db.dump
root@de0700a8ac1f:/# psql -U postgres test_db -c "select c."фамилия" "ФИО", o."наименование" "Заказ" from clients c right outer join orders o on c."заказ"=o.id where c."заказ" is not null"
         ФИО          |  Заказ  
----------------------+---------
 Иванов Иван Иванович | Книга
 Петров Петр Петрович | Монитор
 Иоганн Себастьян Бах | Гитара
(3 rows)
```
Получилось, всё работает.

Можно было бы бакапить сразу с пользователями, если их много, и с паролями:
```
root@310e2169eb54:/# pg_dumpall -U postgres -r > /backups/test_db.sql 
root@310e2169eb54:/# pg_dump -U postgres test_db >>/backups/test_db.sql
...
root@de0700a8ac1f:/# createdb -U postgres test_db
root@de0700a8ac1f:/# cat /backups/test_db1.sql |psql -U postgres
```
