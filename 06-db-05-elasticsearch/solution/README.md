# Домашняя работа по занятию "6.5. Elasticsearch"

> 1. Используя докер образ ...
> - составьте Dockerfile-манифест для elasticsearch

```
FROM centos:7

RUN cd /opt && \
	curl https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.13.3-linux-x86_64.tar.gz \
		> elasticsearch-7.13.3-linux-x86_64.tar.gz && \
	curl https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.13.3-linux-x86_64.tar.gz.sha512 \
		> elasticsearch-7.13.3-linux-x86_64.tar.gz.sha512 && \
	sha512sum -c elasticsearch-7.13.3-linux-x86_64.tar.gz.sha512 && \
	tar -xzf elasticsearch-7.13.3-linux-x86_64.tar.gz && \
	rm -f elasticsearch-7.13.3-linux-x86_64.tar.gz* && \
	useradd -r es_user && \
	chown -R es_user /opt/elasticsearch-7.13.3/

WORKDIR /opt/elasticsearch-7.13.3/

RUN	mkdir -p /var/lib/elasticsearch && \
	chown -R es_user /var/lib/elasticsearch && \
	echo 'node.name: netology_test' >> config/elasticsearch.yml && \
	echo 'path.data: /var/lib/elasticsearch' >> config/elasticsearch.yml && \
	echo 'network.host: 0.0.0.0' >> config/elasticsearch.yml && \
	echo 'discovery.type: single-node' >> config/elasticsearch.yml && \
	echo 'xpack.security.enabled: false' >> config/elasticsearch.yml && \
        echo 'path.repo: /opt/elasticsearch-7.13.3/snapshots' >> config/elasticsearch.yml && \
	mkdir snapshots && chown es_user snapshots

EXPOSE 9200
USER es_user
ENTRYPOINT ["bin/elasticsearch"]
```

> - соберите docker-образ и сделайте `push` в ваш docker.io репозиторий

https://hub.docker.com/r/prividen/elasticsearch

> - запустите контейнер из получившегося образа и выполните запрос пути `/` c хост-машины

```
mak@test-xu20:~/docker/060501$ docker run --rm -p 9200:9200 -d --name es --network es-netw -v es-data:/var/lib/elasticsearch -v es-snashots:/opt/elasticsearch-7.13.3/snapshots prividen/elasticsearch
73ba71cc0bb5a09032040cf780e7d889fd1f75a896298ba57dcfdb73b66a3516

mak@test-xu20:~$ curl localhost:9200/
{
  "name" : "netology_test",
  "cluster_name" : "elasticsearch",
  "cluster_uuid" : "UMc3_6PDTrOZNPGx_XCgLg",
  "version" : {
    "number" : "7.13.3",
    "build_flavor" : "default",
    "build_type" : "tar",
    "build_hash" : "5d21bea28db1e89ecc1f66311ebdec9dc3aa7d64",
    "build_date" : "2021-07-02T12:06:10.804015202Z",
    "build_snapshot" : false,
    "lucene_version" : "8.8.2",
    "minimum_wire_compatibility_version" : "6.8.0",
    "minimum_index_compatibility_version" : "6.0.0-beta1"
  },
  "tagline" : "You Know, for Search"
}
```


---
> 2. добавьте в `elasticsearch` 3 индекса, в соответствии со таблицей:
...

```
mak@test-xu20:~$ curl -X PUT "localhost:9200/ind-1?pretty" -H 'Content-Type: application/json' -d'{"settings": {"number_of_replicas": 0,"number_of_shards": 1}}'
{
  "acknowledged" : true,
  "shards_acknowledged" : true,
  "index" : "ind-1"
}
mak@test-xu20:~$ curl -X PUT "localhost:9200/ind-2?pretty" -H 'Content-Type: application/json' -d'{"settings": {"number_of_replicas": 1,"number_of_shards": 2}}'
{
  "acknowledged" : true,
  "shards_acknowledged" : true,
  "index" : "ind-2"
}
mak@test-xu20:~$ curl -X PUT "localhost:9200/ind-3?pretty" -H 'Content-Type: application/json' -d'{"settings": {"number_of_replicas": 2,"number_of_shards": 4}}'
{
  "acknowledged" : true,
  "shards_acknowledged" : true,
  "index" : "ind-3"
}
```

> Получите список индексов и их статусов, используя API и **приведите в ответе** на задание.

```
mak@test-xu20:~$ curl localhost:9200/_cat/indices/ind*
green  open ind-1 qAEwxaYSRTy87bw2vWNBsw 1 0 0 0 208b 208b
yellow open ind-3 d8zujqpeROGexbSTb5lFlQ 4 2 0 0 832b 832b
yellow open ind-2 mWfiLoU7T1SESqy2jViOBg 2 1 0 0 416b 416b
```

> Получите состояние кластера `elasticsearch`, используя API.

```mak@test-xu20:~$ curl localhost:9200/_cluster/health?pretty
{
  "cluster_name" : "elasticsearch",
  "status" : "yellow",
  "timed_out" : false,
  "number_of_nodes" : 1,
  "number_of_data_nodes" : 1,
  "active_primary_shards" : 13,
  "active_shards" : 13,
  "relocating_shards" : 0,
  "initializing_shards" : 0,
  "unassigned_shards" : 10,
  "delayed_unassigned_shards" : 0,
  "number_of_pending_tasks" : 0,
  "number_of_in_flight_fetch" : 0,
  "task_max_waiting_in_queue_millis" : 0,
  "active_shards_percent_as_number" : 56.52173913043478
}
```
Ещё там есть `_cluster/state`, но оно длинное.


> Как вы думаете, почему часть индексов и кластер находится в состоянии yellow?

Потому что мы создали индексы, которые хотят больше реплик, чем у нас есть. Реплик не хватает, индексы-кластер желтые и несчастные.

> Удалите все индексы.

```
mak@test-xu20:~$ curl -X DELETE "localhost:9200/ind*?pretty"
{
  "acknowledged" : true
}
```

---
> 3. Создайте директорию `{путь до корневой директории с elasticsearch в образе}/snapshots`.
Используя API [зарегистрируйте](https://www.elastic.co/guide/en/elasticsearch/reference/current/snapshots-register-repository.html#snapshots-register-repository) 
данную директорию как `snapshot repository` c именем `netology_backup`.  
**Приведите в ответе** запрос API и результат вызова API для создания репозитория.

```
PUT /_snapshot/netology_backup
{
  "type": "fs",
  "settings": {
    "location": "/opt/elasticsearch-7.13.3/snapshots",
    "compress": true
  }
}

{
  "acknowledged" : true
}
```

Создайте индекс `test` с 0 реплик и 1 шардом и **приведите в ответе** список индексов.

```
PUT /test
{
  "settings": {
    "index": {
      "number_of_replicas": 0,
      "number_of_shards": 1
    }
  }
}
```

```
GET /_cat/indices/test*
green open test qneID9vFTEGPgxqDSHyTIA 1 0 0 0 208b 208b
```

> [Создайте `snapshot`](https://www.elastic.co/guide/en/elasticsearch/reference/current/snapshots-take-snapshot.html) 
состояния кластера `elasticsearch`.  
**Приведите в ответе** список файлов в директории со `snapshot`ами.

```
./meta-NgywM2StQnKU0uTN2F0v_g.dat
./snap-NgywM2StQnKU0uTN2F0v_g.dat
./index-0
./indices
./indices/nrYmptd9S5GrjQp3FuQkYw
./indices/nrYmptd9S5GrjQp3FuQkYw/0
./indices/nrYmptd9S5GrjQp3FuQkYw/0/snap-NgywM2StQnKU0uTN2F0v_g.dat
./indices/nrYmptd9S5GrjQp3FuQkYw/0/index-URa3ELY5Skqq1mD3tLULcQ
./indices/nrYmptd9S5GrjQp3FuQkYw/meta-fz0ZoXoBsyw6DX3ceXO7.dat
./index.latest
```

> Удалите индекс `test` и создайте индекс `test-2`. **Приведите в ответе** список индексов.

```
GET /_cat/indices/test*
green open test-2 1Lb_A_I4Q92oo_MSrn9vEQ 4 0 0 0 832b 832b
```

> [Восстановите](https://www.elastic.co/guide/en/elasticsearch/reference/current/snapshots-restore-snapshot.html) состояние
кластера `elasticsearch` из `snapshot`, созданного ранее. 
**Приведите в ответе** запрос к API восстановления и итоговый список индексов.

```
POST /_snapshot/netology_backup/snap-1/_restore
{
  "accepted" : true
}

GET /_cat/indices/test*
green open test-2 1Lb_A_I4Q92oo_MSrn9vEQ 4 0 0 0 832b 832b
green open test   CIlzgRgVQ8iyOHdf-Ydo0A 1 0 0 0 208b 208b
```

