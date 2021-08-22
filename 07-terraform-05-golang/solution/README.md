# Домашняя работа по занятию "7.5. Основы golang"

> Задача 1. Установите golang. Воспользуйтесь инструкций с официального сайта: [https://golang.org/](https://golang.org/).

Красивое...

---
> Напишите программу для перевода метров в футы (1 фут = 0.3048 метр). Можно запросить исходные данные 
у пользователя, а можно статически задать в коде.

--> [0705-1.go](0705-1/0705-1.go)
```
$ go run 0705-1
Please enter meters value: 12
12.000000 meters is 39.370079 foots
```
 
> Напишите программу, которая найдет наименьший элемент в любом заданном списке, например:  
    ```
    x := []int{48,96,86,68,57,82,63,70,37,34,83,27,19,97,9,17,}
    ```
> 
--> [0705-2.go](0705-2/0705-2.go)   
```
$ go run 0705-2
X array: [48 96 86 68 57 82 63 70 37 34 83 27 19 97 9 17]
Minimal value of x array is: 9
Truly 9!
```

> Напишите программу, которая выводит числа от 1 до 100, которые делятся на 3. То есть `(3, 6, 9, …)`.

--> [0705-3.go](0705-3/0705-3.go)
```
$ go run 0705-3
[3 6 9 12 15 18 21 24 27 30 33 36 39 42 45 48 51 54 57 60 63 66 69 72 75 78 81 84 87 90 93 96 99]
```

---
> 4. Создайте тесты для функций из предыдущего задания. 

--> [0705-1_test.go](0705-1/0705-1_test.go)  
--> [0705-2_test.go](0705-2/0705-2_test.go)  
--> [0705-3_test.go](0705-3/0705-3_test.go)  

```
$ go test 0705*
ok  	0705-1	(cached)
ok  	0705-2	(cached)
ok  	0705-3	(cached)
```
