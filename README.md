# testJob
Для трудоустройства

Возможны проблемой с SMTP yandex банит, как спам. Можно тестовую почту создать, а так все работает.
Настройка SMTP 
1 Войти в директори cd ./go
2 Открыть файл auth.go в любом удобном редвакторе и найти функцию sendEmail vim ./auth.go зайтем нажать i
3 Заменить 
auth := smtp.PlainAuth("", "iisusnawin@yandex.ru", "testAuthJob02", "smtp.yandex.ru") 
значения
smtp.PlainAuth("", "Ваша почта", "логин", "smtp сервер")
и
err := smtp.SendMail("smtp.yandex.ru:25", auth, "iisusnawin@yandex.ru", []string{email}, []byte(message))
по аналоги переменные не трогать
по оконнчанию редактирования нажмите esc и shift + : затем введите wq и нажмите enter
4(a) После go run ./*.go  . Помните вы директори go ещё
4(b) После go build -a. Помните вы директори go ещё. Флаг "а" необязательный очень, но нужен если решите пересобрать файл.
По окнчанию пояаится файл "go" в консоли запустите ./go
5 В браузере http://localhost:4888/ 

Так же файл может работать в режиме демона но думаю, что это ненужно.
Мой номер +79518825061 вайбер и ватсап
