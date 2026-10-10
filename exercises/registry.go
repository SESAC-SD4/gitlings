package exercises

import (
	"gitlings-ko/exercises/ex01"
	"gitlings-ko/exercises/ex02"
	"gitlings-ko/exercises/ex03"
	"gitlings-ko/exercises/ex04"
	"gitlings-ko/exercises/ex05"
	"gitlings-ko/exercises/ex06"
	"gitlings-ko/exercises/ex07"
	"gitlings-ko/exercises/ex08"
	"gitlings-ko/exercises/ex09"
	"gitlings-ko/exercises/ex10"
	"gitlings-ko/exercises/ex11"
	"gitlings-ko/exercises/ex12"
	"gitlings-ko/exercises/ex13"
	"gitlings-ko/internal/exercise"
)

// All은 모든 연습 문제를 ID 순으로 반환한다.
func All() []*exercise.Exercise {
	return []*exercise.Exercise{
		ex01.New(),
		ex02.New(),
		ex03.New(),
		ex04.New(),
		ex05.New(),
		ex06.New(),
		ex07.New(),
		ex08.New(),
		ex09.New(),
		ex10.New(),
		ex11.New(),
		ex12.New(),
		ex13.New(),
	}
}

// ByID는 특정 ID의 문제를 반환한다.
func ByID(id int) *exercise.Exercise {
	for _, ex := range All() {
		if ex.ID == id {
			return ex
		}
	}
	return nil
}
