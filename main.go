package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Book описывает книгу в библиотеке
type Book struct {
	ID        int
	Title     string
	Author    string
	Year      int
	Available int // Количество доступных экземпляров
	Total     int // Общее количество экземпляров (для статистики)
}

// Reader описывает читателя библиотеки
type Reader struct {
	ID   int
	Name string
}

// Loan фиксирует факт выдачи книги читателю
type Loan struct {
	BookID   int
	ReaderID int
	Returned bool // false — книга на руках, true — возвращена
}

// 1. Поиск книги по ID (возвращает указатель)
func FindBook(books []Book, id int) *Book {
	for i := range books {
		if books[i].ID == id {
			return &books[i]
		}
	}
	return nil
}

// Поиск читателя по ID (возвращает указатель)
func FindReader(readers []Reader, id int) *Reader {
	for i := range readers {
		if readers[i].ID == id {
			return &readers[i]
		}
	}
	return nil
}

// 2. Выдача книги с использованием указателя
func BorrowBook(book *Book, readerID int, activeLoans map[int][]int) error {
	if book.Available <= 0 {
		return errors.New("нет доступных экземпляров этой книги")
	}

	// Проверяем, не брал ли читатель уже эту книгу (и она не сдана)
	for _, bID := range activeLoans[readerID] {
		if bID == book.ID {
			return errors.New("этот читатель уже взял данную книгу и еще не вернул её")
		}
	}

	// Уменьшаем количество доступных экземпляров по указателю
	book.Available--

	// Регистрируем активную выдачу
	activeLoans[readerID] = append(activeLoans[readerID], book.ID)
	return nil
}

// 3. Возврат книги с использованием указателя
func ReturnBook(book *Book, readerID int, activeLoans map[int][]int) error {
	loansList := activeLoans[readerID]
	foundIdx := -1

	for i, bID := range loansList {
		if bID == book.ID {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		return errors.New("ошибка: этот читатель не получал данную книгу или уже вернул её")
	}

	// Увеличиваем доступное количество по указателю
	book.Available++

	// Удаляем книгу из списка активных выдач читателя
	activeLoans[readerID] = append(loansList[:foundIdx], loansList[foundIdx+1:]...)
	return nil
}

// 4. Поиск книг по названию или автору
func SearchBooks(books []Book, query string) []Book {
	var result []Book
	queryLower := strings.ToLower(query)

	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), queryLower) ||
			strings.Contains(strings.ToLower(book.Author), queryLower) {
			result = append(result, book)
		}
	}
	return result
}

// Функция демонстрации массива [7]string для дней недели (требование пункта 4)
func GetWeeklyReportHeader() string {
	days := [7]string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}
	return fmt.Sprintf("=== ОТЧЕТ БИБЛИОТЕКИ ЗА НЕДЕЛЮ (Начало недели: %s) ===", days[0])
}

func main() {
	// Массив [7]string задействован в логике отчета
	fmt.Println(GetWeeklyReportHeader())

	// Коллекции данных
	var books []Book = []Book{
		{ID: 1, Title: "The Go Programming Language", Author: "Alan A. A. Donovan", Year: 2015, Available: 3, Total: 3},
		{ID: 2, Title: "Concurrency in Go", Author: "Katherine Cox-Buday", Year: 2017, Available: 2, Total: 2},
		{ID: 3, Title: "Clean Code", Author: "Robert Martin", Year: 2008, Available: 1, Total: 1},
	}

	var readers []Reader = []Reader{
		{ID: 101, Name: "Иван Петров"},
		{ID: 102, Name: "Анна Сидорова"},
	}

	var loans []Loan // Полная история выдач и возвратов

	// Карта: map[int]int — индекс книги в срезе по её ID
	bookIndexMap := make(map[int]int)
	for i, b := range books {
		bookIndexMap[b.ID] = i
	}

	// Карта: map[int][]int — идентификаторы активных выдач по читателям (readerID -> []bookID)
	activeLoans := make(map[int][]int)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n============ LIBRARY MANAGER ============")
		fmt.Println("1. Добавить книгу")
		fmt.Println("2. Зарегистрировать читателя")
		fmt.Println("3. Выдать книгу")
		fmt.Println("4. Вернуть книгу")
		fmt.Println("5. Найти книгу")
		fmt.Println("6. Книги читателя")
		fmt.Println("7. Статистика библиотеки")
		fmt.Println("0. Выход")
		fmt.Print("\nВыберите действие: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Ошибка: введите число от 0 до 7.")
			continue
		}

		switch choice {
		case 1:
			// Добавление книги
			fmt.Print("Введите ID книги (уникальный целое число): ")
			idIn, _ := reader.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(idIn))
			if err != nil {
				fmt.Println("Некорректный ID.")
				continue
			}

			// Проверка на дублирующийся ID
			if FindBook(books, id) != nil {
				fmt.Println("Ошибка: книга с таким ID уже существует!")
				continue
			}

			fmt.Print("Введите название книги: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			fmt.Print("Введите автора книги: ")
			author, _ := reader.ReadString('\n')
			author = strings.TrimSpace(author)

			fmt.Print("Введите год издания: ")
			yearIn, _ := reader.ReadString('\n')
			year, _ := strconv.Atoi(strings.TrimSpace(yearIn))

			fmt.Print("Введите количество экземпляров: ")
			qtyIn, _ := reader.ReadString('\n')
			qty, _ := strconv.Atoi(strings.TrimSpace(qtyIn))
			if qty < 1 {
				qty = 1
			}

			newBook := Book{
				ID:        id,
				Title:     title,
				Author:    author,
				Year:      year,
				Available: qty,
				Total:     qty,
			}
			books = append(books, newBook)
			bookIndexMap[id] = len(books) - 1
			fmt.Println("Книга успешно добавлена в библиотеку!")

		case 2:
			// Регистрация читателя
			fmt.Print("Введите ID читателя (число): ")
			idIn, _ := reader.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(idIn))
			if err != nil {
				fmt.Println("Некорректный ID.")
				continue
			}

			if FindReader(readers, id) != nil {
				fmt.Println("Ошибка: читатель с таким ID уже зарегистрирован!")
				continue
			}

			fmt.Print("Введите ФИО читателя: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			readers = append(readers, Reader{ID: id, Name: name})
			fmt.Println("Читатель успешно зарегистрирован!")

		case 3:
			// Выдать книгу
			fmt.Print("ID читателя: ")
			rIn, _ := reader.ReadString('\n')
			readerID, err1 := strconv.Atoi(strings.TrimSpace(rIn))

			fmt.Print("ID книги: ")
			bIn, _ := reader.ReadString('\n')
			bookID, err2 := strconv.Atoi(strings.TrimSpace(bIn))

			if err1 != nil || err2 != nil {
				fmt.Println("Некорректные данные ID.")
				continue
			}

			rObj := FindReader(readers, readerID)
			if rObj == nil {
				fmt.Println("Ошибка: читатель не найден.")
				continue
			}

			bObj := FindBook(books, bookID)
			if bObj == nil {
				fmt.Println("Ошибка: книга не найдена.")
				continue
			}

			// Вызов функции выдачи с указателем на книгу
			err := BorrowBook(bObj, readerID, activeLoans)
			if err != nil {
				fmt.Println("Ошибка выдачи:", err)
			} else {
				loans = append(loans, Loan{BookID: bookID, ReaderID: readerID, Returned: false})
				fmt.Println("\nКнига успешно выдана!")
				fmt.Printf("Название: %s\n", bObj.Title)
				fmt.Printf("Читатель: %s\n", rObj.Name)
				fmt.Printf("Осталось экземпляров: %d\n", bObj.Available)
			}

		case 4:
			// Вернуть книгу
			fmt.Print("ID читателя: ")
			rIn, _ := reader.ReadString('\n')
			readerID, err1 := strconv.Atoi(strings.TrimSpace(rIn))

			fmt.Print("ID книги для возврата: ")
			bIn, _ := reader.ReadString('\n')
			bookID, err2 := strconv.Atoi(strings.TrimSpace(bIn))

			if err1 != nil || err2 != nil {
				fmt.Println("Некорректные данные ID.")
				continue
			}

			bObj := FindBook(books, bookID)
			if bObj == nil {
				fmt.Println("Ошибка: книга не найдена в фонде.")
				continue
			}

			// Вызов функции возврата с указателем на книгу
			err := ReturnBook(bObj, readerID, activeLoans)
			if err != nil {
				fmt.Println("Ошибка возврата:", err)
			} else {
				loans = append(loans, Loan{BookID: bookID, ReaderID: readerID, Returned: true})
				fmt.Println("Книга успешно возвращена в библиотеку!")
				fmt.Printf("Доступно экземпляров теперь: %d\n", bObj.Available)
			}

		case 5:
			// Поиск книг
			fmt.Print("Введите поисковый запрос (название или автор): ")
			query, _ := reader.ReadString('\n')
			query = strings.TrimSpace(query)

			found := SearchBooks(books, query)
			if len(found) == 0 {
				fmt.Println("Книги по вашему запросу не найдены.")
				continue
			}

			fmt.Println("\n--- РЕЗУЛЬТАТЫ ПОИСКА ---")
			for _, b := range found {
				fmt.Printf("ID: %d | Название: \"%s\" | Автор: %s | Год: %d | Доступно: %d/%d\n",
					b.ID, b.Title, b.Author, b.Year, b.Available, b.Total)
			}

		case 6:
			// Просмотр книг читателя
			fmt.Print("Введите ID читателя: ")
			rIn, _ := reader.ReadString('\n')
			readerID, err := strconv.Atoi(strings.TrimSpace(rIn))
			if err != nil {
				fmt.Println("Некорректный ID.")
				continue
			}

			rObj := FindReader(readers, readerID)
			if rObj == nil {
				fmt.Println("Читатель не найден.")
				continue
			}

			activeBookIDs := activeLoans[readerID]
			fmt.Printf("\nКниги на руках у читателя %s:\n", rObj.Name)
			if len(activeBookIDs) == 0 {
				fmt.Println("Нет взятых книг.")
				continue
			}

			for _, bID := range activeBookIDs {
				bObj := FindBook(books, bID)
				if bObj != nil {
					fmt.Printf(" - [ID: %d] \"%s\" (%s)\n", bObj.ID, bObj.Title, bObj.Author)
				}
			}

		case 7:
			// Статистика библиотеки
			totalTitles := len(books)
			totalCopies := 0
			availableCopies := 0

			// Карта для подсчета популярности книг (сколько раз выдавали каждую книгу)
			borrowCountMap := make(map[int]int)
			for _, l := range loans {
				if !l.Returned { // считаем по всей истории выдач или активным
					borrowCountMap[l.BookID]++
				} else {
					borrowCountMap[l.BookID]++
				}
			}

			for _, b := range books {
				totalCopies += b.Total
				availableCopies += b.Available
			}
			loanedCopies := totalCopies - availableCopies

			// Определение самой популярной книги
			topBookTitle := "Нет данных"
			maxLoans := -1
			for bID, count := range borrowCountMap {
				if count > maxLoans {
					maxLoans = count
					bObj := FindBook(books, bID)
					if bObj != nil {
						topBookTitle = bObj.Title
					}
				}
			}

			fmt.Println("\n============ СТАТИСТИКА БИБЛИОТЕКИ ============")
			fmt.Printf("Всего наименований: %d\n", totalTitles)
			fmt.Printf("Всего экземпляров: %d\n", totalCopies)
			fmt.Printf("Доступно: %d\n", availableCopies)
			fmt.Printf("Выдано: %d\n", loanedCopies)
			fmt.Println()
			fmt.Println("Самая популярная книга:")
			if maxLoans > 0 {
				fmt.Printf("%s\n", topBookTitle)
				fmt.Printf("Количество выдач: %d\n", maxLoans)
			} else {
				fmt.Println("Пока нет зафиксированных выдач.")
			}

		case 0:
			fmt.Println("Выход из программы. До свидания!")
			return

		default:
			fmt.Println("Неверный выбор. Пожалуйста, введите число от 0 до 7.")
		}
	}
}
