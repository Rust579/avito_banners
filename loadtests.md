# Результаты нагрузочного тестирования

Нагрузочное тестирование проводил с помощью Apache Jmeter<br/>
Использовал повторяющиеся запросы, которые можно отправлять подряд с одинаковыми параметрами<br/>
Результаты представлены с нагрузкой, при которой сервер успевает обработать запрос до появления следующего (кроме /update-banner с RPS=500)

- /get-banner с флагом поиска актуальной версии (напрямую из базы)

Нагрузка 100 RPS<br/>
![3.PNG](loadtests%2Fget-banner%20use_last_revision%20%3D%20true%2F3.PNG)
![1.PNG](loadtests%2Fget-banner%20use_last_revision%20%3D%20true%2F1.PNG)
![2.PNG](loadtests%2Fget-banner%20use_last_revision%20%3D%20true%2F2.PNG)

- /get-banner (получение из кэша)

Нагрузка 1000 RPS<br/>
![3.PNG](loadtests%2Fget-banner%20use_last_revision%20%3D%20false%2F3.PNG)
![1.PNG](loadtests%2Fget-banner%20use_last_revision%20%3D%20false%2F1.PNG)
![2.PNG](loadtests%2Fget-banner%20use_last_revision%20%3D%20false%2F2.PNG)

- /get-banners

Нагрузка 500 RPS<br/>
![3.PNG](loadtests%2Fget-banners%2F3.PNG)
![1.PNG](loadtests%2Fget-banners%2F1.PNG)
![2.PNG](loadtests%2Fget-banners%2F2.PNG)

- /update-banner

Нагрузка 500 RPS<br/>
![6.PNG](loadtests%2Fupdate-banner%2F6.PNG)
![4.PNG](loadtests%2Fupdate-banner%2F4.PNG)
![5.PNG](loadtests%2Fupdate-banner%2F5.PNG)

Нагрузка 100 RPS<br/>
![3.PNG](loadtests%2Fupdate-banner%2F3.PNG)
![1.PNG](loadtests%2Fupdate-banner%2F1.PNG)
![2.PNG](loadtests%2Fupdate-banner%2F2.PNG)

- /get-banner-versions

Нагрузка 1000 RPS<br/>
![3.PNG](loadtests%2Fget-banner-versions%2F3.PNG)
![1.PNG](loadtests%2Fget-banner-versions%2F1.PNG)
![2.PNG](loadtests%2Fget-banner-versions%2F2.PNG)
