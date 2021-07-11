# Домашняя работа по занятию "6.4. PostgreSQL"

> 1. **Найдите и приведите** управляющие команды для:
> - вывода списка БД

```  \l[+]   [PATTERN]      list databases```

> - подключения к БД

```  \c[onnect] {[DBNAME|- USER|- HOST|- PORT|-] | conninfo}```

> - вывода списка таблиц

```  \dt[S+] [PATTERN]      list tables```

> - вывода описания содержимого таблиц

```  \d[S+]  NAME           describe table, view, sequence, or index```

> - выхода из psql

```  \q                     quit psql```

---
> 2. Используя `psql` создайте БД `test_database`.  
> ...  
> Используя таблицу [pg_stats](https://postgrespro.ru/docs/postgresql/12/view-pg-stats), найдите столбец таблицы `orders` 
с наибольшим средним значением размера элементов в байтах.  
> **Приведите в ответе** команду, которую вы использовали для вычисления и полученный результат.

```
test_database=# select attname from pg_stats where schemaname='public' and tablename='orders' order by avg_width desc limit 1;
 attname 
---------
 title
(1 row)
```

---
> 3. Архитектор и администратор БД выяснили, что ваша таблица orders разрослась до невиданных размеров и
поиск по ней занимает долгое время. Вам, как успешному выпускнику курсов DevOps в нетологии предложили
провести разбиение таблицы на 2 (шардировать на orders_1 - price>499 и orders_2 - price<=499).  
> Предложите SQL-транзакцию для проведения данной операции.

```
-- rename old table to another name
test_database=# alter table orders rename to orders_old;
ALTER TABLE

-- create new table where we'll fix indexes
test_database=# create table orders_old_1 (like orders_old including all);
CREATE TABLE

-- recreate primary key to be compatible with partitioning by price
test_database=# alter table orders_old_1 drop constraint orders_old_1_pkey;
ALTER TABLE
test_database=# alter table orders_old_1 add primary key (id,price);
ALTER TABLE

-- create new partitioned table
test_database=# create table orders (like orders_old_1 including all) PARTITION BY RANGE(price);
CREATE TABLE

-- with childs
test_database=# create table orders_1 partition of orders for values from (500) to (~(1::int<<31));
CREATE TABLE
test_database=# create table orders_2 partition of orders for values from (0) to (500);
CREATE TABLE

-- import all data from old table
test_database=# insert into orders select * from orders_old;
INSERT 0 8

-- Checking everything:
test_database=# select * from orders;
 id |        title         | price 
----+----------------------+-------
  1 | War and peace        |   100
  3 | Adventure psql time  |   300
  4 | Server gravity falls |   300
  5 | Log gossips          |   123
  7 | Me and my bash-pet   |   499
  2 | My little database   |   500
  6 | WAL never lies       |   900
  8 | Dbiezdmin            |   501
(8 rows)
-- all data is here, hovewer in another order

test_database=# select * from only orders;
 id | title | price 
----+-------+-------
(0 rows)
-- parent table is empty

test_database=# select * from only orders_1;
 id |       title        | price 
----+--------------------+-------
  2 | My little database |   500
  6 | WAL never lies     |   900
  8 | Dbiezdmin          |   501
(3 rows)

test_database=# select * from only orders_2;
 id |        title         | price 
----+----------------------+-------
  1 | War and peace        |   100
  3 | Adventure psql time  |   300
  4 | Server gravity falls |   300
  5 | Log gossips          |   123
  7 | Me and my bash-pet   |   499
(5 rows)
-- data are splitted by price within child tables

-- test for data inserting:
test_database=# insert into orders (title,price) values ('Some cheap', 3), ('Some expensive', 99999);
INSERT 0 2
test_database=# select * from only orders_1;
 id |       title        | price 
----+--------------------+-------
  2 | My little database |   500
  6 | WAL never lies     |   900
  8 | Dbiezdmin          |   501
 12 | Some expensive     | 99999
(4 rows)

test_database=# select * from only orders_2;
 id |        title         | price 
----+----------------------+-------
  1 | War and peace        |   100
  3 | Adventure psql time  |   300
  4 | Server gravity falls |   300
  5 | Log gossips          |   123
  7 | Me and my bash-pet   |   499
 11 | Some cheap           |     3
(6 rows)
```

> Можно ли было изначально исключить "ручное" разбиение при проектировании таблицы orders?

Дочерние таблицы кто-то должен сделать по какому-то критерию. Возможно, можно как-то продумать нашу архитектуру, чтобы они создавались автоматически - по расписанию, или вот предлагают [специальное расширение](https://github.com/pgpartman/pg_partman)

--- 

> 4. Используя утилиту `pg_dump` создайте бекап БД `test_database`.

```root@06580a9d706a:/# pg_dump -U postgres test_database >/tmp/test_database.sql```

> Как бы вы доработали бэкап-файл, чтобы добавить уникальность значения столбца `title` для таблиц `test_database`?

В самый конец файла можно добавить создание UNIQ-ключа по этому столбцу. 
Проблема только с партиционированной таблицей, там будет ругань на то, что в этом индексе должен 
участвовать столбец, по которому делим. Но, учитывая что уникальность столбца `title` может обеспечиться индексом в 
дочерних таблицах, мы можем добавить этот индекс только для непартиционированных таблиц:

```
root@06580a9d706a:/# for TABLE in $(psql -U postgres test_database -A -t -c \
    "SELECT t.table_name FROM information_schema.tables t right join pg_catalog.pg_class c \
    on t.table_name=c.relname where t.table_schema='public' and c.relkind='r';"); do 
        echo "ALTER TABLE ONLY public.$TABLE ADD CONSTRAINT ${TABLE}_title_key UNIQUE (title);" \
            >> /tmp/test_database.sql ; 
    done
```
