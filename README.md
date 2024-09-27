# UDP service

    client.go: Программа получающая от сервера информацию с файла

    server.go: Программа передающая информацию с файла клиентам

# UDP service v1 использование:
## Получение go.mod файла:

    Для клиента в папке client
    go mod init client.go

>

    Для сервера в папке server
    go mod init server.go

## Запуск скомпилированного файла

    Для клиента в папке client
    ./client ip:порт имя_файла
>   ip:порт: ip-адрес и порт **(1024 - 65535)** для подключения к серверу<br>
>   имя_файла: имя файла для сохранения полученной информации

>

    Для сервера в папке server
    ./server порт имя файла

>   порт: порт для запуска сервера на нём<br>
>   имя_файла: имя файла для чтения информации из него

# UDP service v2 usage:
## Получение go.mod файла:

    Для клиента в папке client
    go mod init client.go

>

    Для сервера в папке server
    go mod init server.go

## Запуск скомпилированного файла

    Для клиента в папке client
    ./client ip:порт имя_файла

    Далее следует писать 0 или 1 для получения того же среза байт или следующего, соответственно.

>   ip:порт: ip-адрес и порт **(1024 - 65535)** для подключения к серверу<br>
>   имя_файла: имя файла для сохранения полученной информации

>

    Для сервера в папке server
    ./server порт имя файла

>   порт: порт для запуска сервера на нём<br>
>   имя_файла: имя файла для чтения информации из него


_Я не думаю, что это кто-то прочитает_
_Почему я этим занимаюсь..._