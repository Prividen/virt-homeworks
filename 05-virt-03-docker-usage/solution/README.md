# Домашняя работа по занятию "5.3. Контейнеризация на примере Docker"

> 1. Посмотрите на сценарий ниже и ответьте на вопрос:
> "Подходит ли в этом сценарии использование докера? Или лучше подойдет виртуальная машина, физическая машина? Или возможны разные варианты?"
> Детально опишите и обоснуйте свой выбор.

> - Высоконагруженное монолитное java веб-приложение; 

Если высоконагруженность его зашкаливает, то лучше физическая машина. Ява замечательно потребит все замеченные ресурсы.

> - Go-микросервис для генерации отчетов;
> - Nodejs веб-приложение;

Контейнеры, Go просто создан для них. И Nodejs хорошо туда упаковывается.
 
> - Мобильное приложение c версиями для Android и iOS;

Если имеется ввиду код для другой архитектуры, то виртуалка с трансляцией инструкций, типа qemu. 
Или если серверная часть, то в контейнер наверное.

> - База данных postgresql используемая, как кэш;
  
Базу данных на физический хост. 

> - Шина данных на базе Apache Kafka;
> - Очередь для Logstash на базе Redis;
> - Elastic stack для реализации логирования продуктивного веб-приложения - три ноды elasticsearch, два logstash и две ноды kibana;
> - Мониторинг-стек на базе prometheus и grafana;
  
Все это хорошо для контейнеров.

> - Mongodb, как основное хранилище данных для java-приложения;
  
Если вон для того вон зашкаливающе высоконагруженного из первого пункта, то лучше на физический хост. Из-за неё поди и тормозит всё.  
Или если особой нагрузки нет, то и в контейнере неплохо жить будет.

> - Jenkins-сервер.

Контейнер, или виртуалка. В зависимости от масштабов, наверное.

---
> 2. Сценарий выполения задачи:
> - создайте свой репозиторий на докерхаб; 
> - выберете любой образ, который содержит апачи веб-сервер;
> - создайте свой форк образа;
> - реализуйте функциональность: 
>запуск веб-сервера в фоне с индекс-страницей, содержащей HTML-код ниже: ... 

Про "в фоне" не совсем понятно, имеется ввиду опция докера -d?
Что-нибудь вроде ```docker run -d -p 8080:80 --rm prividen/apache-netology```

> Опубликуйте созданный форк в своем репозитории и предоставьте ответ в виде ссылки на докерхаб-репо.

https://hub.docker.com/r/prividen/apache-netology  
И [Dockerfile](Dockerfile)

---
> Задача 3 
> - Запустите первый контейнер из образа centos c любым тэгом в фоновом режиме, подключив папку info из текущей рабочей директории на хостовой машине в /share/info контейнера;
```
$ docker run -d -it --name container-1 -v $(readlink -e ./info):/share/info centos:8 
74f67249eb8d2798f3d1ba604333951e6a42ff0d90f602f392f40ea5c26b1ab7
```  

> - Запустите второй контейнер из образа debian:latest в фоновом режиме, подключив папку info из текущей рабочей директории на хостовой машине в /info контейнера;

```
$ docker run -d -it --name container-2 -v $(readlink -e ./info):/info debian:latest
88c9f00115095ece4f69a659db5d8592b81c7ea130ef707d3c27736dfee4406e
```
```
$ docker ps
CONTAINER ID   IMAGE           COMMAND       CREATED              STATUS              PORTS     NAMES
88c9f0011509   debian:latest   "bash"        6 seconds ago        Up 5 seconds                  container-2
74f67249eb8d   centos:8        "/bin/bash"   About a minute ago   Up About a minute             container-1
```

> - Подключитесь к первому контейнеру с помощью exec и создайте текстовый файл любого содержания в /share/info ;

```
$ docker exec -it container-1 /bin/bash
[root@74f67249eb8d /]# echo "AnyAnyAny" > /share/info/test.txt
```  

> - Добавьте еще один файл в папку info на хостовой машине;

```
$ echo "Hello from host" > info/host.txt
```

> - Подключитесь во второй контейнер и отобразите листинг и содержание файлов в /info контейнера.

```
$ docker exec -it container-2 /bin/bash
root@88c9f0011509:/# ls -l /info/
total 8
-rw-rw-r-- 1 1000 1000 16 Jun 24 00:15 host.txt
-rw-r--r-- 1 root root 10 Jun 24 00:13 test.txt
root@88c9f0011509:/# cat /info/test.txt 
AnyAnyAny
root@88c9f0011509:/# cat /info/host.txt 
Hello from host
```