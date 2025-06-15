package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	// Разбиваем строку на слайс строк через разделитель ','
	sliceData := strings.Split(data, ",")

	// Проверяем на целостность получаемой строки
	if len(sliceData) < 2 {
		return 0, 0, nil
	}

	// Преобразовываем первое значение в слайсе в тип int и проверяем на отсутствие шагов
	steps, err := strconv.Atoi(sliceData[0])
	if err != nil {
		return 0, 0, nil
	}
	if steps <= 0 {
		return 0, 0, nil
	}

	// Преобразовываем второе значение в слайсе в тип time.Duration
	duration, err := time.ParseDuration(sliceData[1])
	if err != nil {
		return 0, 0, nil
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}

	// Если шаги отстутствуют, возвращаем пустую строку
	if steps == 0 {
		return ""
	}

	// Вычисляем дистанцию в метрах и переводим значение в километры
	distance := (float64(steps) * stepLength) / mInKm

	// Вычисляем количество калорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)

}
