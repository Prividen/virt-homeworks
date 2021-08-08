# Домашняя работа по занятию "7.2. Облачные провайдеры и синтаксис Терраформ."

> 1. Регистрация в aws и знакомство с основами (необязательно, но крайне желательно).
> - Создайте IAM политику для терраформа c правами
>   * AmazonEC2FullAccess
>   * AmazonS3FullAccess
>   * AmazonDynamoDBFullAccess
>   * AmazonRDSFullAccess
>   * CloudWatchFullAccess
>   * IAMFullAccess

```
mak@test-xu20:~$ aws iam list-users
{
    "Users": [
        {
            "Path": "/",
            "UserName": "terraform_agent",
            "UserId": "AIDAX7SACLKN56BK3MKWD",
            "Arn": "arn:aws:iam::548816444059:user/terraform_agent",
            "CreateDate": "2021-07-28T21:39:17+00:00"
        }
    ]
}

mak@test-xu20:~$ aws iam list-attached-user-policies --user-name terraform_agent |jq -r '.AttachedPolicies[].PolicyName'
AmazonRDSFullAccess
AmazonEC2FullAccess
IAMFullAccess
AmazonS3FullAccess
CloudWatchFullAccess
AmazonDynamoDBFullAccess
```

> - Создайте, остановите и удалите ec2 инстанс (любой с пометкой `free tier`) через веб интерфейс.

```
mak@test-xu20:~$ aws ec2 describe-instances |jq '.Reservations[].Instances[] |{Region:.Placement.AvailabilityZone,InstanceId,ImageId,InstanceType}'
{
  "Region": "us-west-2b",
  "InstanceId": "i-0bf6a19b9f0f2549a",
  "ImageId": "ami-083ac7c7ecf9bb9b0",
  "InstanceType": "t2.micro"
}
```

> В виде результата задания приложите вывод команды `aws configure list`.
```
mak@test-xu20:~$ aws configure list
      Name                    Value             Type    Location
      ----                    -----             ----    --------
   profile                <not set>             None    None
access_key     ****************ND6Q              env    
secret_key     ****************ZMjq              env    
    region                us-west-2      config-file    ~/.aws/config
```

---
> 2. Созданием ec2 через терраформ.
> В качестве результата задания предоставьте:
> - Ответ на вопрос: при помощи какого инструмента (из разобранных на прошлом занятии) можно создать свой образ ami?

У нас образами виртуалок `packer` [занимается](https://www.packer.io/docs/builders/amazon). Хотя, вроде можно создать пустой инстанс, отпровизионить его Ансиблем, и сохранить как ami.

> - Ссылку на репозиторий с исходной конфигурацией терраформа.

https://github.com/Prividen/devops-netology/tree/d9beb269fd95605bade07a4db68ea24d21c2ccc1/terraform

```
mak@test-xu20:~/terraform$ terraform apply -auto-approve

...

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:

account_id = "548816444059"
aws_region = "us-west-2"
caller_user = "AIDAX7SACLKN56BK3MKWD"
instance_private_ip = "172.31.9.65"
instance_public_ip = "34.219.115.77"
instance_subnet_id = "subnet-a92649f4"
```


