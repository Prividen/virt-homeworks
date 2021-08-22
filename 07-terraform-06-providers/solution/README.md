# Домашняя работа по занятию "7.6. Написание собственных провайдеров для Terraform."

> Задача 1. 
> 1. Найдите, где перечислены все доступные `resource` и `data_source`, приложите ссылку на эти строки в коде на 
гитхабе.

Они перечислены в файле aws/provider.go:
- --> [`data_source`](https://github.com/hashicorp/terraform-provider-aws/blob/2b12ec179f2616975ce0afe67b454dce7368a4ed/aws/provider.go#L186)
- --> [`resource`](https://github.com/hashicorp/terraform-provider-aws/blob/2b12ec179f2616975ce0afe67b454dce7368a4ed/aws/provider.go#L456)

> 2. Для создания очереди сообщений SQS используется ресурс `aws_sqs_queue` у которого есть параметр `name`. 
> - С каким другим параметром конфликтует `name`? Приложите строчку кода, в которой это указано.

Конфликтует с name_prefix, ```ConflictsWith: []string{"name_prefix"},``` ([код](https://github.com/hashicorp/terraform-provider-aws/blob/2b12ec179f2616975ce0afe67b454dce7368a4ed/aws/resource_aws_sqs_queue.go#L99))

> - Какая максимальная длина имени?
> - Какому регулярному выражению должно подчиняться имя? 

Как сказано в `website/docs/r/sqs_queue.html.markdown`, <cite>Queue names must be made up of only uppercase and lowercase ASCII letters, numbers, underscores, and hyphens, and must be between 1 and 80 characters long</cite>

То же самое [написано в `aws/resource_aws_sqs_queue.go`:](https://github.com/hashicorp/terraform-provider-aws/blob/2b12ec179f2616975ce0afe67b454dce7368a4ed/aws/resource_aws_sqs_queue.go#L415)  
```re = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,80}$`)``` 


---

> Задача 2. (Не обязательно) 
> - Проделайте все шаги создания провайдера.

Я извиняюсь, но на него какое-то дикое количество времени угрохалось с непривычки. Провайдер простенький, читает конфиги из JSON (я про это на лекции как раз спрашивал) 
 
> - В виде результата приложение ссылку на исходный код.

https://github.com/Prividen/terraform-provider-jsonconf

> - Попробуйте скомпилировать провайдер, если получится то приложите снимок экрана с командой и результатом компиляции.   

```
[mak@mak-ws terraform-provider-jsonconf]$ make install
go build -o terraform-provider-jsonconf
mkdir -p ~/.terraform.d/plugins/github.com/prividen/jsonconf/0.1.13/linux_amd64
mv terraform-provider-jsonconf ~/.terraform.d/plugins/github.com/prividen/jsonconf/0.1.13/linux_amd64

[mak@mak-ws terraform-provider-jsonconf]$ cd examples/
[mak@mak-ws examples]$ terraform init

Initializing the backend...

Initializing provider plugins...
- Finding github.com/prividen/jsonconf versions matching "~> 0.1"...
- Installing github.com/prividen/jsonconf v0.1.13...
- Installed github.com/prividen/jsonconf v0.1.13 (unauthenticated)
...
Terraform has been successfully initialized!
...

[mak@mak-ws examples]$ terraform apply -auto-approve

Changes to Outputs:
  + datacenters = {
      + 0 = "Matrix4"
      + 1 = "eq.ash"
      + 2 = "M9"
      + 3 = "eq.la"
      + 4 = "CloudVSP.bj"
      + 5 = "l3.lima"
      + 6 = "CMC.vn"
      + 7 = "eq.HK"
      + 8 = "gs.SG"
      + 9 = "gs.SYD"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real
infrastructure.

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
...
```

К сожалению, не получилось добавить имя файла как конфигурационный параметр к провайдеру. 
При попытке взять его сразу из `data_source_jsonconf.go` (`dataSourceJsonConfRead`, что-нибудь вроде ```config = d.Get("config").(string)```) провайдер падает с трейсом, а написать полноценную функцию `providerConfigure` в `provider.go` я ниасилил, надо найти какой-нибудь пример, где без http-клиента обходятся... 