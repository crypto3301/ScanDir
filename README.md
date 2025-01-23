# SnapScan

**SnapScan** — это утилита, написанная на языке Golang, предназначенная для создания снимков (snapshots) директорий. Она обеспечивает целостность данных, рассчитывая хэши файлов. Для удобства использования утилита имеет графический интерфейс (GUI).

---

## 🛠 **Особенности**
- **Снимки директорий**: Создание снимков директорий для фиксации их текущего состояния.
- **Целостность данных**: Вычисление хэшей файлов (SHA-256) для проверки целостности данных.
- **Графический интерфейс**: Удобная навигация с помощью клавиатуры.

---

## 🎮 **Управление**
Навигация по интерфейсу осуществляется клавишами:
- **Tab**: Переключение между окнами.
- **Стрелки**: Перемещение по списку файлов.
- **Enter**: Выбор файла или выполнение команды.
- **Пробел**: Аналог клавиши Esc для возврата на предыдущий экран.

---

## 🪟 **Окна интерфейса**

### **LOG Window**
- **Назначение**: Логирование всех выполненных команд и их результатов.
- **Использование**: Отображает журнал действий в реальном времени.

### **INFO Window**
- **Назначение**: Работа со снимками.
- **Навигация**:
  - Используйте стрелки для перемещения по списку файлов.
  - Нажмите **Enter**, чтобы выбрать файл или выполнить действие.
  - Нажмите **Space**, чтобы вернуться назад (аналог Esc).

### **Terminal Window**
- **Назначение**: Ввод и выполнение команд.
- **Использование**: Введите команды в терминале для взаимодействия с утилитой.

---

## 🚀 **Установка**

1. Убедитесь, что у вас установлен Go. Если нет, скачайте и установите его с [golang.org](https://golang.org/).
2. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/crypto3301/ScanDir.git
