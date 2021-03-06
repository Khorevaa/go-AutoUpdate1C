package update

import (
	"fmt"

	"github.com/khorevaa/go-AutoUpdate1C/logging"

	"github.com/khorevaa/go-v8runner"
	//"github.com/sirupsen/logrus"
)

type Обновление struct {
	*v8runner.Конфигуратор

	СтрокаПодключения                string
	Пользователь                     string
	Пароль                           string
	ФайлОбновления                   string
	ВыполнитьЗагрузкуВместоОбновения bool
	ИспользоватьПолныйДистрибутив    bool
	НаСервере                        bool
	УправляемыйРежим                 bool

	ФайлОбработкиПриЗапуске string

	СеансыПользователей *УправлениеСеансамиПользователей

	ДополнительныеПараметры []string

	logging.Log // *log.Entry
}

func (о *Обновление) УстановитьЛог(log logging.Log) {

	о.Log = log.Context(logging.LogFeilds{
		"СтрокаПодключения": о.СтрокаПодключения,
		"Пользователь":      о.Пользователь,
	})

}

func НовоеОбновление(строкаПодключения, пользователь, пароль string) (о *Обновление) {
	о = &Обновление{
		СтрокаПодключения: строкаПодключения,
		Пользователь:      пользователь,
		Пароль:            пароль,
		Конфигуратор:      v8runner.НовыйКонфигуратор(),
	}

	о.УстановитьКлючСоединенияСБазой(о.СтрокаПодключения)
	return

}

func (о *Обновление) ВыполнитьОбновление() (e error) {

	if о.ВыполнитьЗагрузкуВместоОбновения {
		e = о.выполнитьЗагрузкуКонфигурации()
	} else {
		e = о.выполнитьОбновлениеКонфигурации()
	}

	if e != nil {
		о.WithError(e).Warning("Ошибка выполнения обновления")
	}

	return

}

func (о *Обновление) ВыполнитьВРежимеПредприятия(КомандаПриЗапуске, ФайлОбработки string, privileged bool) (e error) {

	log := о.Log.Context(logging.LogFeilds{
		"КомандаПриЗапуске": КомандаПриЗапуске,
		"ФайлОбработки":     ФайлОбработки,
		"privileged":        privileged,
	})

	log.Infof("Выполняю запуск информационной базы в режиме 1С.Предприятие")

	var ДополнительныеПараметрыВыполнения []string

	copy(ДополнительныеПараметрыВыполнения, о.ДополнительныеПараметры)

	if len(КомандаПриЗапуске) != 0 {
		ДополнительныеПараметрыВыполнения = append(ДополнительныеПараметрыВыполнения, fmt.Sprintf("/C%s", КомандаПриЗапуске))
	}

	if len(ФайлОбработки) != 0 {
		ДополнительныеПараметрыВыполнения = append(ДополнительныеПараметрыВыполнения, fmt.Sprintf("/Execute%s", ФайлОбработки))
	}

	if privileged {
		ДополнительныеПараметрыВыполнения = append(ДополнительныеПараметрыВыполнения, "/UsePrivilegedMode")
	}

	e = о.ЗапуститьВРежимеПредприятия(о.УправляемыйРежим, ДополнительныеПараметрыВыполнения...)

	if e != nil {
		log.WithError(e).Warning("Ошибка выполнения запуска в режиме предприятия")
	}

	return

}

func (о *Обновление) выполнитьЗагрузкуКонфигурации() (e error) {
	log := о.Log.Context(logging.LogFeilds{
		"ФайлОбновления": о.ФайлОбновления,
	})

	log.Infof("Загрузка конфигурации в информационную базу")

	e = о.ЗагрузитьКонфигурациюИзФайла(о.ФайлОбновления) // Не сработает.., о.ДополнительныеПараметры...)

	if e != nil {
		return
	}

	log.Infof("Обновление конфигурации в информационной базе")
	e = о.ОбновитьКонфигурациюБазыДанных(false, о.НаСервере, false, о.ДополнительныеПараметры...)

	return
}

func (о *Обновление) выполнитьОбновлениеКонфигурации() (e error) {

	log := о.Log.Context(logging.LogFeilds{
		"ФайлОбновления":                о.ФайлОбновления,
		"ИспользоватьПолныйДистрибутив": о.ИспользоватьПолныйДистрибутив,
		"НаСервере":                     о.НаСервере,
	})

	log.Infof("Выполняю обнволение информационной базы из поставки")

	e = о.ОбновитьКонфигурацию(о.ФайлОбновления, о.ИспользоватьПолныйДистрибутив, true, о.НаСервере, false, о.ДополнительныеПараметры...)

	return
}

type УправлениеСеансамиПользователей struct {
}
