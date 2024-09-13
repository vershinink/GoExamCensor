### GoExamCensor

Сервис проверки комментариев к новостным статьям на наличие недопустимых выражений. Практика на курсе "Go-разработчик" от SkillFactory. Часть итогового проекта курса.

Для запуска нужно установить путь к файлу конфига в переменную окружения `CENSOR_CONFIG_PATH`. Остальные входные данные указываются в файле конфига. Список недопустимых выражений также находится в файле конфига.

Сам файл конфига `config.yaml` лежит в каталоге config.

**Сделано:**

- Логирование в stdout через пакет slog стандартной библиотеки Go.
- REST API метод проверка комментария на наличие недопустимых выражений.
- Тесты для всех основных пакетов приложения.
- Использование контекстов при работе сервера.
- Использоваие middleware для трассировки запросов и логирования.
- Завершение работы приложения по сигналу прерывания с использованием graceful shutdown.
- Сборка и запуск сервиса в Docker контейнере.

**Методы:**

- POST `/` , проверяет содержимое тела запроса на наличие слов из списка в конфиге. Возвращает либо код 200, либо 400. В теле должен быть JSON с обязательным полем `content` , где содержится текст комментария.
