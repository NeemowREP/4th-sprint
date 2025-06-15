package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разбиваем строку на слайс строк через разделитель ','. Пример получаемой строки - "3456,Ходьба,3h00m"
	sliceData := strings.Split(data, ",")

	// Проверяем на целостность получаемой строки
	if len(sliceData) < 3 {
		return 0, "", 0, fmt.Errorf("ошибка, данные отсутствуют")
	}

	// Преобразовываем первое значение в слайсе в тип int и проверяем на отсутствие шагов
	steps, err := strconv.Atoi(sliceData[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка, преобразование кол-ва шагов невозможно")
	}
	if steps == 0 {
		return 0, "", 0, fmt.Errorf("встань с дивана")
	}

	// Преобразовываем третье значение в слайсе в тип time.Duration
	duration, err := time.ParseDuration(sliceData[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования времени")
	}

	return steps, sliceData[1], duration, nil
}

func distance(steps int, height float64) float64 {

	// Вычисляем длину шага
	stepLength := height * stepLengthCoefficient

	return (float64(steps) * stepLength) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	// Проверяем валидность продолжительности duration
	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию
	dist := distance(steps, height)

	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	// Разбиваем получаемую строку data, пример ввода - "3456,Ходьба,3h00m"
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("ошибка вводных данных")
	}

	// Получаем значения пройденной дистанции и средней скорости
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// Объявляем две переменные для корректной работы switch case
	var calories float64
	var caloriesErr error

	// Проверяем тип тренировки
	switch trainingType {
	case "Бег":
		calories, caloriesErr = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, caloriesErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if caloriesErr != nil {
		return "", caloriesErr
	}

	dur := duration.Hours()

	info := fmt.Sprintf("Тип тренировки: %s\nДлительность: %v ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
trainingType, dur, dist, speed, calories)

return info, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	
	// Проверяем корректность вводимых данных
	switch  {
	case steps <= 0:
		return 0, fmt.Errorf("ошибка, введите корректное значение пройденных шагов")
	case weight <= 0:
		return 0, fmt.Errorf("ошибка, Введите корректное значение вашего веса")
	case height <= 0:
		return 0, fmt.Errorf("ошибка, введите корректное значение вашего роста")
	}

	// Рассчет средней скорости
	speed := meanSpeed(steps, height, duration)

	// Рассчет калорий по формуле
	calories := (weight * speed * duration.Minutes()) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	// Проверяем корректность вводимых данных
	switch  {
	case steps <= 0:
		return 0, fmt.Errorf("ошибка, Введите корректное значение пройденных шагов")
	case weight <= 0:
		return 0, fmt.Errorf("ошибка, Введите корректное значение вашего веса")
	case height <= 0:
		return 0, fmt.Errorf("ошибка, Введите корректное значение вашего роста")
	}

	// Рассчет средней скорости
	speed := meanSpeed(steps, height, duration)
	
	// Рассчет калорий с учетом коэффициента
	calories := (walkingCaloriesCoefficient *((duration.Minutes() * weight * speed) / minInH))

	return calories, nil
}
