package pdf

import (
	"fmt"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/refundable-tgm/huginn/db"
	"github.com/refundable-tgm/huginn/untis"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const BasePath = "/vol/files/"
const ClassAbsenceFormFileName = "class_absence_form_%v.pdf"
const TeacherAbsenceFormFileName = "teacher_absence_form_%v.pdf"
const CompensationForEducationalSupportFileName = "compensation_for_educational_support.pdf"

func GeneratePDFEnvironment(app db.Application) (string, error) {
	dirname := app.UUID
	path := filepath.Join(BasePath, dirname)
	err := os.MkdirAll(path, os.ModePerm)
	return path, err
}

func GenerateAbsenceFormForClass(path, username string, app db.Application) error {
	client := untis.GetClient(username)
	defer client.Close()
	if app.Kind != db.SchoolEvent {
		return fmt.Errorf("this pdf can only be generated for school events")
	}
	for _, class := range app.SchoolEventDetails.Classes {
		m := pdf.NewMaroto(consts.Portrait, consts.A4)
		m.SetPageMargins(10, 15, 10)

		m.RegisterHeader(func() {
			m.Row(20, func() {
				m.Col(3, func() {
					_ = m.FileImage("TGM_Logo.png", props.Rect{
						Left:    0,
						Top:     0,
						Center:  false,
						Percent: 100,
					})
				})

				m.ColSpace(2)

				m.Col(2, func() {
					m.Text("Abwesenheitsmeldung eines Jahrgangs", props.Text{
						Align:  consts.Center,
						Family: consts.Helvetica,
						Size:   12,
					})
				})

				m.ColSpace(2)

				m.Col(3, func() {
					m.QrCode("https://refundable.tech/viewer?uuid="+app.UUID, props.Rect{
						Left:    27,
						Top:     0,
						Center:  false,
						Percent: 100,
					})
				})
			})
		})

		m.SetDefaultFontFamily(consts.Helvetica)
		m.Line(3.0)
		m.Row(3, func() {})
		m.Row(10, func() {
			m.Col(12, func() {
				m.Text("Allgemeine Informationen:", props.Text{
					Top:   3,
					Align: consts.Left,
					Style: consts.Bold,
				})
			})
		})
		leader := ""
		companion := ""
		for _, teacher := range app.SchoolEventDetails.Teachers {
			if teacher.Role == db.Leader {
				leader = teacher.Name
			} else if teacher.Role == db.Companion {
				companion = companion + teacher.Name + ", "
			}
		}
		companion = companion[0 : len(companion)-2]

		m.Row(10, func() {
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Jahrgang:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				m.Text(class, props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Lehrkraft:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				m.Text(leader, props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
		})
		m.Line(1.0)

		m.Row(10, func() {
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Anzahl m/w Schüler/innen:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				m.Text(strconv.Itoa(app.SchoolEventDetails.AmountMaleStudents)+" / "+strconv.Itoa(app.SchoolEventDetails.AmountFemaleStudents), props.Text{
					Top:   2.5,
					Align: consts.Center,
				})
			})
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Begleitpersonen:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				m.Text(companion, props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
		})
		m.Line(1.0)
		m.Row(10, func() {
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Von:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				start := app.StartTime
				weekday := getWeekday(int(start.Weekday()))
				month := strconv.Itoa(int(start.Month()))
				if len(month) == 1 {
					month = "0" + month
				}
				day := strconv.Itoa(start.Day())
				if len(day) == 1 {
					day = "0" + day
				}
				year := start.Year()
				hour := strconv.Itoa(start.Hour())
				if len(hour) == 1 {
					hour = "0" + hour
				}
				minute := strconv.Itoa(start.Minute())
				if len(minute) == 1 {
					minute = "0" + minute
				}
				m.Text(fmt.Sprintf("%v, %v.%v.%d %v:%v", weekday, day, month, year, hour, minute), props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Bis:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				end := app.EndTime
				weekday := getWeekday(int(end.Weekday()))
				month := strconv.Itoa(int(end.Month()))
				if len(month) == 1 {
					month = "0" + month
				}
				day := strconv.Itoa(end.Day())
				if len(day) == 1 {
					day = "0" + day
				}
				year := end.Year()
				hour := strconv.Itoa(end.Hour())
				if len(hour) == 1 {
					hour = "0" + hour
				}
				minute := strconv.Itoa(end.Minute())
				if len(minute) == 1 {
					minute = "0" + minute
				}
				m.Text(fmt.Sprintf("%v, %v.%v.%d %v:%v", weekday, day, month, year, hour, minute), props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
		})
		m.Line(1.0)
		m.Row(10, func() {
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Anmerkungen:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
			})
		})
		m.Row(15, func() {
			m.Col(12, func() {
				m.Text(app.Notes, props.Text{
					Top:   2.5,
					Align: consts.Left,
					Style: consts.Italic,
				})
			})
		})
		m.Line(2.0)
		m.Row(3, func() {})
		m.Row(10, func() {
			m.Col(12, func() {
				m.Text("Schulveranstaltung:", props.Text{
					Top:   3,
					Align: consts.Left,
					Style: consts.Bold,
				})
			})
		})

		m.Row(10, func() {
			m.Col(12, func() {
				m.Col(6, func() {
					m.Text("Veranstaltung:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				m.Text(app.Name, props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
		})
		m.Line(1.0)
		m.Row(10, func() {
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Treffpunkt:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				m.Text(app.StartAddress, props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Uhrzeit:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				start := app.StartTime
				hour := strconv.Itoa(start.Hour())
				if len(hour) == 1 {
					hour = "0" + hour
				}
				minute := strconv.Itoa(start.Minute())
				if len(minute) == 1 {
					minute = "0" + minute
				}
				m.Text(fmt.Sprintf("%v:%v", hour, minute), props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
		})

		m.Line(1.0)
		m.Row(10, func() {
			m.Col(12, func() {
				m.Col(6, func() {
					m.Text("Dauer der Veranstaltung:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				dayString := ""
				if app.SchoolEventDetails.DurationInDays == 1 {
					dayString = "1-tägig (002)"
				} else if app.SchoolEventDetails.DurationInDays > 3 {
					dayString = "mehr als 3-tägig (004)"
				} else {
					dayString = "2-3-tägig (003)"
				}
				m.Text(dayString, props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
		})

		m.Line(2.0)
		m.Row(3, func() {})
		m.Row(10, func() {
			m.Col(12, func() {
				m.Text("Supplierungen:", props.Text{
					Top:   3,
					Align: consts.Left,
					Style: consts.Bold,
				})
			})
		})
		tableStrings := make([][]string, 0)
		lessons, err := client.GetTimetableOfClass(app.StartTime, app.EndTime, class)
		if err != nil {
			return err
		}
		for _, lesson := range lessons {
			date := lesson.Start
			month := strconv.Itoa(int(date.Month()))
			if len(month) == 1 {
				month = "0" + month
			}
			day := strconv.Itoa(date.Day())
			if len(day) == 1 {
				day = "0" + day
			}
			year := date.Year()
			beginLesson := untis.GetLessonNrByStart(lesson.Start)
			endLesson := untis.GetLessonNrByEnd(lesson.End)
			hourString := ""
			if beginLesson == endLesson {
				hourString = fmt.Sprintf("%d.", beginLesson)
			} else {
				hourString = fmt.Sprintf("%d. - %d.", beginLesson, endLesson)
			}
			rooms := ""
			for _, room := range lesson.Rooms {
				rooms = room + ", "
			}
			rooms = rooms[0 : len(rooms)-2]
			row := []string{"", class,
				fmt.Sprintf("%v.%v.%d", day, month, year),
				fmt.Sprintf(hourString),
				rooms,
			}
			tableStrings = append(tableStrings, row)
		}
		m.TableList([]string{"H/R/E", "Jahrgang", "Datum", "Stunde", "Saal", "LK Entf.", "LK Supp.", "Paraphe"}, tableStrings)

		m.Line(2.0)
		m.Row(3, func() {})
		m.Row(10, func() {
			m.Col(12, func() {
				m.Text("Kenntnisnahme:", props.Text{
					Top:   3,
					Align: consts.Left,
					Style: consts.Bold,
				})
			})
		})
		ackStrings := make([][]string, 6)
		ackStrings[0][0] = "AV"
		ackStrings[1][0] = "AV"
		ackStrings[2][0] = "WL"
		ackStrings[3][0] = "Begleitperson"
		ackStrings[4][0] = "Ersteller/in"
		ackStrings[5][0] = "UNTIS Eintragung"
		m.TableList([]string{"Stelle", "Datum", "Paraphe"}, ackStrings)

		err = m.OutputFileAndClose(path + fmt.Sprintf(ClassAbsenceFormFileName, class))
		if err != nil {
			return fmt.Errorf("could not save pdf: %v", err)
		}
	}
	return nil
}

func GenerateCompensationForEducationalSupport(path string, app db.Application) error {
	if app.Kind != db.SchoolEvent {
		return fmt.Errorf("this pdf can only be generated for school events")
	}
	var leader db.SchoolEventTeacherDetails
	companions := make([]db.SchoolEventTeacherDetails, 0)
	teachers := app.SchoolEventDetails.Teachers
	for _, teacher := range teachers {
		if teacher.Role == db.Leader {
			leader = teacher
		} else {
			companions = append(companions, teacher)
		}
	}
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(3, func() {
				_ = m.FileImage("TGM_Logo.png", props.Rect{
					Left:    0,
					Top:     0,
					Center:  false,
					Percent: 100,
				})
			})

			m.ColSpace(2)

			m.Col(2, func() {
				m.Text("Abgeltung für pädagogische Betreeung gemäß §63a", props.Text{
					Align:  consts.Center,
					Family: consts.Helvetica,
					Size:   12,
				})
			})

			m.ColSpace(2)

			m.Col(3, func() {
				m.QrCode("https://refundable.tech/viewer?uuid="+app.UUID, props.Rect{
					Left:    27,
					Top:     0,
					Center:  false,
					Percent: 100,
				})
			})
		})
	})

	m.SetDefaultFontFamily(consts.Helvetica)
	m.Line(3.0)
	m.Row(3, func() {})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Allgemeine Informationen:", props.Text{
				Top:   3,
				Align: consts.Left,
				Style: consts.Bold,
			})
		})
	})

	m.Row(8, func() {
		m.Col(12, func() {
			m.Text("Formular ist vom Leiter bzw. der Leiterin der Schulveranstaltung mit der Reiserechnung in der PEK abzugeben", props.Text{Style: consts.Bold})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Col(6, func() {
				m.Text("Veranstaltung:", props.Text{
					Top:   2.5,
					Align: consts.Left,
				})
			})
			m.Text(app.Name, props.Text{
				Top:   2.5,
				Align: consts.Center,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Col(6, func() {
				m.Text("Datum:", props.Text{
					Top:   2.5,
					Align: consts.Left,
				})
			})
			start := app.StartTime
			sweekday := getWeekday(int(start.Weekday()))
			smonth := strconv.Itoa(int(start.Month()))
			if len(smonth) == 1 {
				smonth = "0" + smonth
			}
			sday := strconv.Itoa(start.Day())
			if len(sday) == 1 {
				sday = "0" + sday
			}
			syear := start.Year()
			shour := strconv.Itoa(start.Hour())
			if len(shour) == 1 {
				shour = "0" + shour
			}
			sminute := strconv.Itoa(start.Minute())
			if len(sminute) == 1 {
				sminute = "0" + sminute
			}
			end := app.EndTime
			eweekday := getWeekday(int(end.Weekday()))
			emonth := strconv.Itoa(int(end.Month()))
			if len(emonth) == 1 {
				emonth = "0" + emonth
			}
			eday := strconv.Itoa(end.Day())
			if len(eday) == 1 {
				eday = "0" + eday
			}
			eyear := end.Year()
			ehour := strconv.Itoa(end.Hour())
			if len(ehour) == 1 {
				ehour = "0" + ehour
			}
			eminute := strconv.Itoa(end.Minute())
			if len(eminute) == 1 {
				eminute = "0" + eminute
			}
			m.Text(fmt.Sprintf("%v, %v.%v.%d %v:%v - %v, %v.%v.%d %v:%v",
				sweekday, sday, smonth, syear, shour, sminute,
				eweekday, eday, emonth, eyear, ehour, eminute),
				props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
		})
	})
	m.Line(2.0)
	m.Row(5, func() {})
	m.Row(10, func() {
		m.Col(6, func() {
			m.Col(3, func() {
				m.Text("Leitung:", props.Text{
					Top:   2.5,
					Align: consts.Left,
				})
			})
			m.Text(leader.Name, props.Text{
				Top:   2.5,
				Align: consts.Center,
				Style: consts.Italic,
			})
		})
		m.Col(6, func() {
			m.Col(3, func() {
				m.Text("Verwendungsgruppe:", props.Text{
					Top:   2.5,
					Align: consts.Left,
				})
			})
			m.Text(fmt.Sprintf("L%d", leader.Group), props.Text{
				Top:   2.5,
				Align: consts.Center,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(6, func() {
			m.Col(3, func() {
				m.Text("Beginn:", props.Text{
					Top:   2.5,
					Align: consts.Left,
				})
			})
			start := leader.AttendanceFrom
			weekday := getWeekday(int(start.Weekday()))
			month := strconv.Itoa(int(start.Month()))
			if len(month) == 1 {
				month = "0" + month
			}
			day := strconv.Itoa(start.Day())
			if len(day) == 1 {
				day = "0" + day
			}
			year := start.Year()
			hour := strconv.Itoa(start.Hour())
			if len(hour) == 1 {
				hour = "0" + hour
			}
			minute := strconv.Itoa(start.Minute())
			if len(minute) == 1 {
				minute = "0" + minute
			}
			m.Text(fmt.Sprintf("%v, %v.%v.%d %v:%v", weekday, day, month, year, hour, minute), props.Text{
				Top:   2.5,
				Align: consts.Center,
				Style: consts.Italic,
			})
		})
		m.Col(6, func() {
			m.Col(3, func() {
				m.Text("Ende:", props.Text{
					Top:   2.5,
					Align: consts.Left,
				})
			})
			end := leader.AttendanceTill
			weekday := getWeekday(int(end.Weekday()))
			month := strconv.Itoa(int(end.Month()))
			if len(month) == 1 {
				month = "0" + month
			}
			day := strconv.Itoa(end.Day())
			if len(day) == 1 {
				day = "0" + day
			}
			year := end.Year()
			hour := strconv.Itoa(end.Hour())
			if len(hour) == 1 {
				hour = "0" + hour
			}
			minute := strconv.Itoa(end.Minute())
			if len(minute) == 1 {
				minute = "0" + minute
			}
			m.Text(fmt.Sprintf("%v, %v.%v.%d %v:%v", weekday, day, month, year, hour, minute), props.Text{
				Top:   2.5,
				Align: consts.Center,
				Style: consts.Italic,
			})
		})
	})
	m.Line(3.0)
	m.Row(3, func() {})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Pädagogisch-inhaltliche Betreuung:", props.Text{
				Top:   3,
				Align: consts.Left,
				Style: consts.Bold,
			})
		})
	})

	tableString := make([][]string, 0)
	for _, teacher := range teachers {
		start := app.StartTime
		sweekday := getWeekday(int(start.Weekday()))
		smonth := strconv.Itoa(int(start.Month()))
		if len(smonth) == 1 {
			smonth = "0" + smonth
		}
		sday := strconv.Itoa(start.Day())
		if len(sday) == 1 {
			sday = "0" + sday
		}
		syear := start.Year()
		shour := strconv.Itoa(start.Hour())
		if len(shour) == 1 {
			shour = "0" + shour
		}
		sminute := strconv.Itoa(start.Minute())
		if len(sminute) == 1 {
			sminute = "0" + sminute
		}
		end := app.EndTime
		eweekday := getWeekday(int(end.Weekday()))
		emonth := strconv.Itoa(int(end.Month()))
		if len(emonth) == 1 {
			emonth = "0" + emonth
		}
		eday := strconv.Itoa(end.Day())
		if len(eday) == 1 {
			eday = "0" + eday
		}
		eyear := end.Year()
		ehour := strconv.Itoa(end.Hour())
		if len(ehour) == 1 {
			ehour = "0" + ehour
		}
		eminute := strconv.Itoa(end.Minute())
		if len(eminute) == 1 {
			eminute = "0" + eminute
		}
		row := []string{
			teacher.Name,
			fmt.Sprintf("L%d", teacher.Group),
			fmt.Sprintf("%v, %v.%v.%d %v:%v", sweekday, sday, smonth, syear, shour, sminute),
			fmt.Sprintf("%v, %v.%v.%d %v:%v", eweekday, eday, emonth, eyear, ehour, eminute),
		}
		tableString = append(tableString, row)
	}
	m.TableList([]string{"Name", "Verwendungsgruppe", "Beginn", "Ende"}, tableString)

	m.Row(15, func() {})
	m.Line(1.0)
	m.Row(8, func() {
		m.Col(12, func() {
			m.Text("Datum und Unterschrift des Leiters der Schulveranstaltung")
		})
	})
	m.Row(5, func() {})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("1. Dem Lehrer gebührt für die Teilnahme an mindestens zweitägigen Schulveranstaltungen" +
				" mit Nächtigung, sofern er die pädagogisch-inhaltliche Betreuung" +
				" einer Schülergruppe innehat, eine Abgeltung.", props.Text{Size: 8})
		})
	})
	m.Row(5, func() {})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("2. Weiters gebührt dem Leiter einer mindestens viertägigen Schulveranstaltung als" +
				"Abgeltung die Einrechnung in die Lehrverpflichtung von 4.55 WE in jener Woche in der die" +
				"Schulveranstaltung endet.", props.Text{Size: 8})
		})
	})

	err := m.OutputFileAndClose(path + CompensationForEducationalSupportFileName)
	if err != nil {
		return fmt.Errorf("could not save pdf: %v", err)
	}
	return nil
}

func GenerateAbsenceFormForTeacher(path string, app db.Application) {

}

func GenerateTravelInvoice(path string, app db.Application) {

}

func GenerateBusinessTripApplication(path string, app db.Application) {

}

func getWeekday(weekday int) string {
	switch weekday {
	case int(time.Monday):
		return "Montag"
	case int(time.Tuesday):
		return "Dienstag"
	case int(time.Wednesday):
		return "Mittwoch"
	case int(time.Thursday):
		return "Donnerstag"
	case int(time.Friday):
		return "Freitag"
	case int(time.Saturday):
		return "Samstag"
	case int(time.Sunday):
		return "Sonntag"
	}
	return ""
}
