
# Mangabuffdl CLI

Утилита для скачивания манги. Парсит страницы глав, извлекает изображения и сохраняет их в структурированные папки.

## Особенности

- Скачивание глав по диапазону номеров
- Автоматическое создание папок для глав
- Поддержка томов манги

## Установка

1. Скачивайте последний [релиз отсюда](https://github.com/bdrlv/mangabuffdl/releases) 

## Использование

```bash
./mb.exe -u URL [-v N] [-s X] [-e Y]
```

### Флаги

| Флаг       | Обязательный | Описание                          | По умолчанию |
|------------|--------------|-----------------------------------|--------------|
| `-u`       | Да           | Базовый URL манги                 | -            |
| `-v`       | Нет          | Номер тома                        | 1            |
| `-s`       | Нет          | Стартовый номер главы             | 1            |
| `-e`       | Нет          | Конечный номер главы              | 1            |

### Примеры

Скачать том 1, главы 1-5:
```bash
./mb.exe -u "https://mangabuff.ru/manga/ya-budu-korolem-v-etoi-zhizni" -s 1 -e 5
```

Скачать том 3, главу 10:
```bash
./mb.exe -u "https://mangabuff.ru/manga/ya-budu-korolem-v-etoi-zhizni" -v 3 -s 10 -e 10
```

Скачать все главы тома 2:
```bash
./mb.exe -u "https://mangabuff.ru/manga/ya-budu-korolem-v-etoi-zhizni" -v 2 -s 1 -e 999
```

## Структура файлов

После скачивания файлы сохраняются в следующей структуре:
```
manga_name/
  Chapter 1/
    1.jpg
    2.jpg
    ...
  Chapter 2/
    ...
```

## Ограничения

- Работает только с mangabuff.ru


## Отказ от ответственности

Данная утилита разработана исключительно в образовательных и исследовательских целях. Автор программного обеспечения:

1. **Не несёт ответственности за содержимое**  
   - За любой контент, скачанный пользователями с помощью этой утилиты  
   - За материалы, размещённые на сайте-источнике (mangabuff.ru, иных и связанных доменах)  
   - За возможное нарушение авторских прав третьих лиц при использовании утилиты

2. **Не поддерживает и не одобряет пиратство**  
   - Утилита не предназначена для незаконного распространения защищённого авторским правом контента  
   - Разработчик не призывает и не поощряет действия, нарушающие интеллектуальные права  
   - Любые упоминания сторонних ресурсов не являются рекомендацией к их использованию

3. **Ограничения использования**  
   - Программа должна использоваться ТОЛЬКО для доступа к контенту, распространяемому на законных основаниях  
   - Пользователи обязаны самостоятельно проверять правовой статус скачиваемых материалов  
   - Запрещено использовать утилиту для массового скачивания или создания архивов без явного разрешения правообладателей

4. **Позиция в отношении авторских прав**  
   - Разработчик уважает права правообладателей и призывает пользователей делать то же самое  
   - Любой скачанный материал должен быть удалён в течение 24 часов после ознакомления  
   - Автор не несёт ответственности за использование программы в противоречии с местным законодательством

5. **Технические ограничения**  
   - Утилита не является официальным клиентом сайта-источника  
   - Разработчик не гарантирует стабильную работу программы в связи с возможными изменениями структуры сайта  
   - Любые совпадения с коммерческими продуктами случайны

**ВАЖНО**: Используя данную программу, вы подтверждаете что:  
- Имеете законные права на скачиваемый контент  
- Не будете использовать полученные материалы в коммерческих целях  
- Принимаете полную ответственность за последствия использования утилиты

Автор оставляет за собой право вносить изменения в данный отказ без предварительного уведомления. Последняя версия документа всегда доступна в репозитории проекта.

---

**Поддержка**: Для сообщений об ошибках используйте [Issues](https://github.com/bdrlv/mangabuffdl/issues)
```
