# Домашняя работа 6 — конкурентность и generic-типизация в Go

Домашка закрепляет занятие 6: goroutines, channels, `sync.WaitGroup`, data race, `sync.Mutex`, `close`, `range`, `select`, отмену через `context.Context` и базовые generics.

## Разделы

| Раздел | Пакет | Количество задач |
|---|---|---:|
| Goroutines | `internal/goroutines` | 5 |
| Channels | `internal/channels` | 5 |
| WaitGroup | `internal/waitgroup` | 5 |
| Race и Mutex | `internal/racemutex` | 5 |
| Close и Range | `internal/closerange` | 5 |
| Select | `internal/selectflow` | 5 |
| Context | `internal/contextflow` | 5 |
| Generics | `internal/generics` | 5 |

Всего: 40 функций и методов — по 5 коротких задач в каждом разделе.

Для каждой задачи предусмотрено не менее 10 основных тестовых сценариев. У каждой функции в `docs/task.md` есть одинаковые опоры: «Что тренируем», «Как рассуждать», «Мини-пример» и «Частая ошибка». Дополнительный `TestExample` проверяет итоговый вывод раздела, а integration-тест запускает все восемь `cmd`.

## Как не потеряться в конкурентности

Не пытайся писать весь раздел сразу. В конкурентном коде ошибка часто выглядит не как обычное неверное число, а как зависание, data race или goroutine, которая не может завершиться.

Работай по одной функции:

1. Прочитай её описание в `docs/task.md`.
2. Нарисуй словами: кто запускает goroutine, кто ждёт, кто владеет channel и кто может остановить работу.
3. Реализуй только одну функцию.
4. Запусти тесты текущего раздела.
5. После блока с общей памятью обязательно запусти `make test-race`.

Внутри каждого раздела сначала реши задачи 1–3 — они закрепляют основной механизм. Затем переходи к задачам 4–5: они повторяют тот же принцип в новом контексте и помогают проверить, что решение не было запомнено только для одного примера.

Рекомендуемый маршрут:

| Этап | Что выполнить | Команда |
|---|---|---|
| 1 | Goroutines | `make test-goroutines` |
| 2 | Channels | `make test-channels` |
| 3 | WaitGroup | `make test-waitgroup` |
| 4 | Race и Mutex | `make test-race-mutex`, затем `make test-race` |
| 5 | Close и Range | `make test-close-range` |
| 6 | Select | `make test-select` |
| 7 | Context | `make test-context` |
| 8 | Generics | `make test-generics` |
| 9 | Все unit- и integration-тесты | `make test` |
| 10 | Полная проверка | `make ci` |

В стартовом проекте TODO ещё не реализованы, поэтому тесты и GitHub Actions сначала красные. Это ожидаемо. Заглушки специально не должны зависать: они позволяют компилировать проект и постепенно делать зелёным каждый раздел.

## Важно

- Не меняй сигнатуры функций и не редактируй тесты.
- Не используй `time.Sleep` для ожидания завершения goroutine.
- Вызывай `WaitGroup.Add` до запуска goroutine, а `Done` — через `defer` внутри неё.
- Общие изменяемые данные защищай одним и тем же Mutex и при записи, и при чтении.
- Channel закрывает отправитель, который точно знает, что новых значений не будет.
- Не закрывай channel со стороны получателя и не отправляй значение после `close`.
- `context.Context` передавай первым аргументом; не сохраняй его в структуре и не передавай `nil`.
- Проверяй ошибки context через `errors.Is(err, context.Canceled)` и `errors.Is(err, context.DeadlineExceeded)`.
- В generic-функции выбирай минимальное constraint: `any`, если операции над значениями не нужны; `comparable`, если используешь `==`, `!=` или ключи map.
- Не заменяй type parameters на `any` с последующим type assertion: это убирает compile-time проверку, ради которой используется generics.

## Самопроверка

```bash
make compile
make fmt
make fmt-check
make test-goroutines
make test-channels
make test-waitgroup
make test-race-mutex
make test-close-range
make test-select
make test-context
make test-generics
make test-unit
make test-integration
make test-race
make ci
```
