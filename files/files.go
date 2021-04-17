package files

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/refundable-tgm/huginn/db"
	"github.com/refundable-tgm/huginn/untis"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const BasePath = "/vol/files/"
const TemplatePath = "/go/src/huginn/excel_template/"
const UploadFolderName = "upload/"
const ReceiptFileName = "%d_%v_receipt.pdf"
const ExcelTemplateTravelInvoicePath = "reiserechnung.xlsx"
const ExcelTemplateBusinessTripApplicationPath = "dienstreiseantrag.xlsx"
const ClassAbsenceFormFileName = "class_absence_form_%v.pdf"
const TeacherAbsenceFormFileName = "teacher_absence_form_%v.pdf"
const CompensationForEducationalSupportFileName = "compensation_for_educational_support.pdf"
const BusinessTripApplicationPDFFileName = "business_trip_application_%v.pdf"
const TravelInvoicePDFFileName = "travel_invoice_%v.pdf"
const BusinessTripApplicationExcelFileName = "business_trip_application_%v.xlsx"
const TravelInvoiceExcelFileName = "travel_invoice_%v.xlsx"

const CheckedCheckBox = "☑"
const UncheckedCheckBox = "☐"
const Sheet = "Sheet1"
const URL = "https://refundable.tech/viewer?uuid="

const (
	BTAWorkplace                          = "A1"
	BTASurname                            = "A4"
	BTAName                               = "O4"
	BTADegree                             = "Y4"
	BTATitle                              = "AF4"
	BTATel                                = "AP4"
	BTAPNR1                               = "M7"
	BTAPNR2                               = "N7"
	BTAPNR3                               = "O7"
	BTAPNR4                               = "P7"
	BTAPNR5                               = "Q7"
	BTAPNR6                               = "R7"
	BTAPNR7                               = "S7"
	BTAPNR8                               = "T7"
	BTAVGr                                = "AA7"
	BTAEGr                                = "AE7"
	BTADKI                                = "AI7"
	BTAGSt                                = "AM7"
	BTAESt                                = "AQ7"
	BTAFeeLevel                           = "AV7"
	BTATripBeginDate                      = "T10"
	BTATripBeginTime                      = "AC10"
	BTATripEndDate                        = "AN10"
	BTATripEndTime                        = "AW10"
	BTAServiceBeginDate                   = "T12"
	BTAServiceBeginTime                   = "AC12"
	BTAServiceEndDate                     = "AN12"
	BTAServiceEndTime                     = "AW12"
	BTADestination                        = "H15"
	BTATravelReasoning                    = "A18"
	BTACheckOfficialBusinessCardClass2    = "B20"
	BTACheckPassenger                     = "B22"
	BTACheckOfficialBusinessCardClass1    = "B24"
	BTACheckTravelGrant                   = "P20"
	BTACheckFlight                        = "P22"
	BTACheckCheapFlight                   = "AD20"
	BTACheckTrainClass2                   = "AD22"
	BTACheckOwnCar                        = "AD24"
	BTACheckSleepTrain                    = "AR20"
	BTACheckBus                           = "AR22"
	BTACheckStartAddressOffice            = "K27"
	BTACheckStartAddressOwnApartment      = "K29"
	BTACheckEndAddressOffice              = "Z27"
	BTACheckEndAddressOwnApartment        = "Z29"
	BTAReasoning                          = "A33"
	BTAOtherParticipants                  = "A35"
	BTACheckBonusMiles1                   = "L37"
	BTACheckBonusMiles2                   = "L40"
	BTACheckTravelCostsPayedBySomeoneYes  = "K45"
	BTACheckTravelCostsPayedBySomeoneNo   = "F45"
	BTACheckStayingCostsPayedBySomeoneYes = "W45"
	BTACheckStayingCostsPayedBySomeoneNo  = "R45"
	BTAPayedByWhom                        = "AE45"
	BTAOtherCosts                         = "J47"
	BTAEstimatedCosts                     = "AN47"
	BTAApprovalDate                       = "Y56"
	BTAFilingDate                         = "AQ62"
	BTACheckBusinessCardEmittedOutward    = "AO69"
	BTACheckBusinessCardEmittedReturn     = "BA69"
	BTAReferee                            = "AP65"
)

const (
	TIWorkplace                                   = "A1"
	TITripBeginYear1                              = "S5"
	TITripBeginYear2                              = "T5"
	TITripBeginYear3                              = "U5"
	TITripBeginYear4                              = "V5"
	TITripBeginMonth1                             = "W5"
	TITripBeginMonth2                             = "X5"
	TITripBeginDay1                               = "Y5"
	TITripBeginDay2                               = "Z5"
	TITripBeginHour1                              = "AA5"
	TITripBeginHour2                              = "AB5"
	TITripBeginMinute1                            = "AC5"
	TITripBeginMinute2                            = "AD5"
	TITripEndYear1                                = "S7"
	TITripEndYear2                                = "T7"
	TITripEndYear3                                = "U7"
	TITripEndYear4                                = "V7"
	TITripEndMonth1                               = "W7"
	TITripEndMonth2                               = "X7"
	TITripEndDay1                                 = "Y7"
	TITripEndDay2                                 = "Z7"
	TITripEndHour1                                = "AA7"
	TITripEndHour2                                = "AB7"
	TITripEndMinute1                              = "AC7"
	TITripEndMinute2                              = "AD7"
	TITravelCostsGrant                            = "AU4"
	TIExtraAmount                                 = "AU6"
	TIZI                                          = "CG3"
	TIFilingDate                                  = "CO5"
	TISurname                                     = "A9"
	TIName                                        = "N9"
	TIDegree                                      = "X9"
	TITitle                                       = "AE9"
	TIPNR1                                        = "L12"
	TIPNR2                                        = "M12"
	TIPNR3                                        = "N12"
	TIPNR4                                        = "O12"
	TIPNR5                                        = "P12"
	TIPNR6                                        = "Q12"
	TIPNR7                                        = "R12"
	TIPNR8                                        = "S12"
	TIStartingAddress                             = "AD12"
	TIDestinationAddress                          = "AD13"
	TIClerk                                       = "CA15"
	TIReviewer                                    = "DG15"
	TICheckOfficialBusinessCardGot                = "B17"
	TICheckTravelGrant                            = "L17"
	TICheckReplacementForAdvantageCard            = "V17"
	TICheckReplacementForTrainCardClass2          = "AE17"
	TICheckKilometreAllowance                     = "AQ17"
	TIKilometreAmount                             = "BD19"
	TICheckNRAndIndicationsOfParticipants         = "BR17"
	TICheckTravelCostsCited                       = "CL17"
	TICheckNoTravelCosts                          = "DE17"
	TICheckDailyChargesType1                      = "B20"
	TICheckDailyChargesType2                      = "L20"
	TICheckToBeShortened                          = "V20"
	TIToBeShortenedAmount                         = "W21"
	TIBreakfasts                                  = "AE21"
	TILunches                                     = "AL21"
	TIDinners                                     = "AX21"
	TICheckNightlyChargesProofNeededForCharges    = "BL20"
	TICheckNightlyChargesNoProofNeeded            = "CH20"
	TICheckNightlyChargesNoClaimForNightlyCharges = "DD20"
	TIFirstConstantRow                            = 28
	TIFirstGeneratedRow                           = 33
	TICalcNrColumn                                = "A"
	TICalcDayColumn                               = "D"
	TICalcBeginColumn                             = "G"
	TICalcEndColumn                               = "K"
	TICalcKindOfFeeColumn                         = "O"
	TICalcKilometreAmountColumn                   = "AO"
	TICalcTravelCostsColumn                       = "AY"
	TICalcDailyChargesColumn                      = "BI"
	TICalcNightlyChargesColumn                    = "BS"
	TICalcAdditionalCostsColumn                   = "CH"
	TICalcSumColumn                               = "CY"
)

func GenerateFileEnvironment(app db.Application) (string, error) {
	dirname := app.UUID
	path := filepath.Join(BasePath, dirname)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	path = filepath.Join(path, UploadFolderName)
	err = os.MkdirAll(path, os.ModePerm)
	return path, err
}

func GenerateAbsenceFormForClass(path, username string, app db.Application) ([]string, error) {
	paths := make([]string, 0)
	client := untis.GetClient(username)
	defer client.Close()
	if app.Kind != db.SchoolEvent {
		return nil, fmt.Errorf("this pdf can only be generated for school events")
	}
	for _, class := range app.SchoolEventDetails.Classes {
		m := pdf.NewMaroto(consts.Portrait, consts.A4)
		m.SetPageMargins(10, 15, 10)

		m.RegisterHeader(func() {
			m.Row(20, func() {
				m.Col(3, func() {
					_ = m.FileImage("../assets/TGM_Logo.png", props.Rect{
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
					m.QrCode(URL+app.UUID, props.Rect{
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
			return nil, err
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
				rooms = rooms + room + ", "
			}
			rooms = rooms[0 : len(rooms)-2]
			supp := ""
			for _, teach := range lesson.Teachers {
				supp = supp + teach + ", "
			}
			supp = supp[0 : len(rooms)-2]

			row := []string{"", class,
				fmt.Sprintf("%v.%v.%d", day, month, year),
				fmt.Sprintf(hourString),
				rooms,
				leader + ", " + companion,
				supp,
				"",
			}
			tableStrings = append(tableStrings, row)
		}
		m.TableList([]string{"H/R/E", "Jahrgang", "Datum", "Stunde", "Saal", "LK Supp.", "LK Entf.", "Paraphe"}, tableStrings)

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

		savePath := filepath.Join(path, fmt.Sprintf(ClassAbsenceFormFileName, class))
		err = m.OutputFileAndClose(savePath)
		paths = append(paths, savePath)
		if err != nil {
			return nil, fmt.Errorf("could not save pdf: %v", err)
		}
	}
	return paths, nil
}

func GenerateCompensationForEducationalSupport(path string, app db.Application) (string, error) {
	if app.Kind != db.SchoolEvent {
		return "", fmt.Errorf("this pdf can only be generated for school events")
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
				_ = m.FileImage("../assets/TGM_Logo.png", props.Rect{
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
				m.QrCode(URL+app.UUID, props.Rect{
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
			m.Text("1. Dem Lehrer gebührt für die Teilnahme an mindestens zweitägigen Schulveranstaltungen"+
				" mit Nächtigung, sofern er die pädagogisch-inhaltliche Betreuung"+
				" einer Schülergruppe innehat, eine Abgeltung.", props.Text{Size: 8})
		})
	})
	m.Row(5, func() {})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("2. Weiters gebührt dem Leiter einer mindestens viertägigen Schulveranstaltung als"+
				"Abgeltung die Einrechnung in die Lehrverpflichtung von 4.55 WE in jener Woche in der die"+
				"Schulveranstaltung endet.", props.Text{Size: 8})
		})
	})
	savePath := filepath.Join(path, CompensationForEducationalSupportFileName)
	err := m.OutputFileAndClose(savePath)
	if err != nil {
		return "", fmt.Errorf("could not save pdf: %v", err)
	}
	return savePath, nil
}

func GenerateAbsenceFormForTeacher(path, username, teacher string, app db.Application) (string, error) {
	client := untis.GetClient(username)
	defer client.Close()
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(3, func() {
				_ = m.FileImage("../assets/TGM_Logo.png", props.Rect{
					Left:    0,
					Top:     0,
					Center:  false,
					Percent: 100,
				})
			})

			m.ColSpace(2)

			m.Col(2, func() {
				m.Text("Abwesenheitsmeldung eines Lehrers", props.Text{
					Align:  consts.Center,
					Family: consts.Helvetica,
					Size:   12,
				})
			})

			m.ColSpace(2)

			m.Col(3, func() {
				m.QrCode(URL+app.UUID, props.Rect{
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

	m.Row(10, func() {
		m.Col(12, func() {
			m.Col(6, func() {
				m.Text("Name:", props.Text{
					Top:   2.5,
					Align: consts.Left,
				})
			})
			m.Text(username, props.Text{
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
			m.Text("Abwesenheitsgrund:", props.Text{
				Top:   3,
				Align: consts.Left,
				Style: consts.Bold,
			})
		})
	})

	if app.Kind == db.SchoolEvent {
		m.Row(10, func() {
			m.Row(10, func() {
				m.Col(12, func() {
					m.Col(6, func() {
						m.Text("Schulveranstaltung:", props.Text{
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
		})
	} else if app.Kind == db.Training {
		m.Row(10, func() {
			m.Col(12, func() {
				m.Text("Fortbildung:", props.Text{
					Top:   3,
					Align: consts.Left,
					Style: consts.Italic,
				})
			})
		})
		m.Row(10, func() {
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Titel der Fortbildung:", props.Text{
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
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("PH-Zahl:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				m.Text(strconv.Itoa(app.TrainingDetails.PH), props.Text{
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
					m.Text("Art der Veranstaltung:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				st := ""
				switch app.TrainingDetails.Kind {
				case db.Seminar:
					st = "Seminar"
					break
				case db.Conference:
					st = "Tagung"
					break
				case db.Course:
					st = "Lehrgang"
				case db.Miscellaneous:
					st = "Sonstiger Grund: " + app.TrainingDetails.MiscellaneousReason
				}
				m.Text(st, props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
			m.Col(6, func() {
				m.Col(3, func() {
					m.Text("Veranstalter:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				m.Text(app.TrainingDetails.Organizer, props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
		})
	} else if app.Kind == db.OtherReason {
		m.Row(10, func() {
			m.Col(12, func() {
				m.Col(6, func() {
					m.Text("Anderer Grund:", props.Text{
						Top:   2.5,
						Align: consts.Left,
					})
				})
				st := ""
				switch app.OtherReasonDetails.Kind {
				case db.Careleave:
					st = "Pflegefreistellung"
					break
				case db.ServiceMandate:
					st = "Dienstauftrag"
					break
				case db.MedicalAppointment:
					st = "Arzttermin"
				case db.Miscellaneous:
					st = "Sonstige Gründe"
				}
				m.Text(st, props.Text{
					Top:   2.5,
					Align: consts.Center,
					Style: consts.Italic,
				})
			})
		})
		if app.OtherReasonDetails.Kind == db.ServiceMandate {
			m.Line(1)
			m.Row(10, func() {
				m.Col(6, func() {
					m.Col(3, func() {
						m.Text("GZ:", props.Text{
							Top:   2.5,
							Align: consts.Left,
						})
					})
					m.Text(strconv.Itoa(app.OtherReasonDetails.ServiceMandateGZ), props.Text{
						Top:   2.5,
						Align: consts.Center,
						Style: consts.Italic,
					})
				})
				m.Col(6, func() {
					m.Col(3, func() {
						m.Text("Titel:", props.Text{
							Top:   2.5,
							Align: consts.Left,
						})
					})
					m.Text(app.OtherReasonDetails.ServiceMandateTitle, props.Text{
						Top:   2.5,
						Align: consts.Center,
						Style: consts.Italic,
					})
				})
			})
		} else if app.OtherReasonDetails.Kind == db.Other {
			m.Line(1)
			m.Row(10, func() {
				m.Col(12, func() {
					m.Col(6, func() {
						m.Text("Grund:", props.Text{
							Top:   2.5,
							Align: consts.Left,
						})
					})
					m.Text(app.OtherReasonDetails.MiscellaneousReason, props.Text{
						Top:   2.5,
						Align: consts.Center,
						Style: consts.Italic,
					})
				})
			})
		}
	}

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
	lessons, err := client.GetTimetableOfSpecificTeacher(app.StartTime, app.EndTime, teacher)
	if err != nil {
		return "", err
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
			rooms = rooms + room + ", "
		}
		rooms = rooms[0 : len(rooms)-2]
		classes := ""
		for _, class := range lesson.Classes {
			classes = classes + class + ", "
		}
		classes = classes[0 : len(rooms)-2]
		row := []string{"", classes,
			fmt.Sprintf("%v.%v.%d", day, month, year),
			fmt.Sprintf(hourString),
			rooms,
			"",
			username,
			"",
		}
		tableStrings = append(tableStrings, row)
	}
	m.TableList([]string{"H/R/E", "Jahrgang", "Datum", "Stunde", "Saal", "LK Supp.", "LK Entf.", "Paraphe"}, tableStrings)

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
	ackStrings := make([][]string, 4)
	ackStrings[0][0] = "AV"
	ackStrings[2][0] = "WL"
	ackStrings[4][0] = "Ersteller/in"
	ackStrings[5][0] = "UNTIS Eintragung"
	m.TableList([]string{"Stelle", "Datum", "Paraphe"}, ackStrings)

	savePath := filepath.Join(path, fmt.Sprintf(TeacherAbsenceFormFileName, username))
	err = m.OutputFileAndClose(savePath)
	if err != nil {
		return "", fmt.Errorf("could not save pdf: %v", err)
	}
	return savePath, nil
}

func GenerateTravelInvoice(path, short string, app db.TravelInvoice, uuid string) (string, error) {
	m := pdf.NewMaroto(consts.Landscape, consts.A4)
	m.SetPageMargins(10, 15, 10)

	m.Row(15, func() {
		m.Col(3, func() {
			_ = m.FileImage("../assets/TGM_Logo.png", props.Rect{
				Left:    0,
				Top:     0,
				Center:  true,
				Percent: 100,
			})
		})

		m.ColSpace(2)

		m.Col(2, func() {
			m.Text("Reiserechnung Inland", props.Text{
				Align:  consts.Center,
				Family: consts.Helvetica,
				Size:   12,
			})
		})

		m.ColSpace(2)

		m.Col(3, func() {
			m.QrCode(URL+uuid, props.Rect{
				Left:    27,
				Top:     0,
				Center:  true,
				Percent: 100,
			})
		})
	})
	m.Line(4.0)
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Familienname:", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text(app.Surname, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Vorname:", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text(app.Name, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Akademischer Grad:", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text(app.Degree, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Amtstitel:", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text(app.Title, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Row(10, func() {
		m.Col(1, func() {
			m.Text("Beginn:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text(app.TripBeginTime.Format("02. 01. 2006 15:04"), props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(1, func() {
			m.Text("Ende:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text(app.TripEndTime.Format("02. 01. 2006 15:04"), props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Reisekostenvorschuss:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text(strconv.FormatFloat(float64(app.TravelCostsPreGrant), 'f', 2, 32) + " €", props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Anzahl der Beilagen:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		extras, _ := ioutil.ReadDir(filepath.Join(path, UploadFolderName, short))
		m.Col(1, func() {
			m.Text(strconv.Itoa(len(extras)), props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Row(10, func() {
		m.Col(1, func() {
			m.Text("Personalnr:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text(strconv.Itoa(app.Staffnr), props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(1, func() {
			m.Text("Bearbeiter:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text(app.Clerk, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(1, func() {
			m.Text("Prüfer:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text(app.Reviewer, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(1, func() {
			m.Text("Eingelangt:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text(app.FilingDate.Format("02. 01. 2006 15:04"), props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Ausgangsort:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(4, func() {
			m.Text(app.StartingPoint, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Zielort:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(4, func() {
			m.Text(app.EndPoint, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(2.0)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Zusätzliche Daten:", props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Bold,
			})
		})
	})
	data := ""
	if app.OfficialBusinessCardGot {
		data = data + "Amtl. Businesskarte erhalten;   "
	}
	if app.TravelGrant {
		data = data + "Beförderungszuschuss;   "
	}
	if app.ReplacementForAdvantageCard {
		data = data + "Ersatz für Vorteilscard (Beleg erford.);   "
	}
	if app.ReplacementForTrainCardClass2 {
		data = data + "Ersatz für Bahnfahrt 2. Kl (Beleg erford.);   "
	}
	if app.KilometreAllowance {
		data = data + "Amtl. Kilometergeld für eigenen PKW (" +
			strconv.FormatFloat(float64(app.KilometreAmount), 'f', 2, 32) + " km);   "
	}
	if app.NRAndIndicationsOfParticipants {
		data = data + "Anzahl und namentliche Angabe der Mitfahrer;   "
	}
	if app.TravelCostsCited {
		data = data + "Angeführte andere Reisekosten (nur gegen Beleg);   "
	}
	if app.NoTravelCosts {
		data = data + "Keine Reisekosten;   "
	}
	data = data[0:len(data) - 4]
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(data, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Row(10, func() {
		m.Col(1, func() {
			m.Text("Tagesgebühr:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		day := ""
		switch app.DailyChargesMode {
		case db.DailyChargesType1:
			day = "Tarif I"
			break
		case db.DailyChargesType2:
			day = "Tarif II"
			break
		case db.ToBeShortened:
			day = "zu kürzen um " + strconv.FormatFloat(float64(app.ShortenedAmount), 'f', 2, 32)
		}
		m.Col(2, func() {
			m.Text(day, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(5, func() {
			m.Text(strconv.Itoa(app.Breakfasts) + " Frühstück; " +
				strconv.Itoa(app.Lunches) + " Mittagessen; " +
				strconv.Itoa(app.Dinners) + " Abendessen",
				props.Text{
					Top:   2.5,
					Align: consts.Middle,
					Style: consts.Italic,
				})
		})
		m.Col(2, func() {
			m.Text("Nächtigungsgeb.:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		night := ""
		switch app.NightlyChargesMode {
		case db.ProofNeededForCharges:
			night = "mit Nachweis"
			break
		case db.NoProofNeeded:
			night = "ohne Nachweis"
			break
		case db.NoClaimForNightlyCharges:
			night = "kein Anspruch"
		}
		m.Col(2, func() {
			m.Text(night, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(2.0)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Berechnungsblatt", props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Bold,
			})
		})
	})
	header := []string{"Nr.", "Tag", "Beginn", "Ende", "Art Gebühren",
		"Kilometer", "Reisek.", "Tagk.", "Nachtk.", "Nebenk.", "Summe"}
	content := make([][]string, 0)
	for _, r := range app.Calculation.Rows {
		row := make([]string, len(header))
		row[0] = strconv.Itoa(r.NR)
		row[1] = r.Date.Format("02.01")
		row[2] = r.Begin.Format("15:04")
		row[3] = r.End.Format("15:04")
		geb := ""
		for _, kind := range r.KindsOfCost {
			switch kind {
			case db.TravelCosts:
				geb = geb + "Reisekosten, "
				row[5] = strconv.FormatFloat(float64(r.Kilometres), 'f', 2, 32)
				row[6] = strconv.FormatFloat(float64(r.TravelCosts), 'f', 2, 32)
				break
			case db.DailyCharges:
				geb = geb + "Tagesgebühr, "
				row[7] = strconv.FormatFloat(float64(r.DailyCharges), 'f', 2, 32)
				break
			case db.NightlyCharges:
				geb = geb + "Nächtigungsgebühr, "
				row[8] = strconv.FormatFloat(float64(r.NightlyCharges), 'f', 2, 32)
				break
			case db.AdditionalCosts:
				geb = geb + "Nebenkosten, "
				row[9] = strconv.FormatFloat(float64(r.AdditionalCosts), 'f', 2, 32)
				break
			}
		}
		geb = geb[0:len(geb) - 2]
		row[4] = geb
		row[10] = strconv.FormatFloat(float64(r.Sum), 'f', 2, 32)
		content = append(content, row)
	}
	content = append(content, []string{"","","","","Summe:","",
		strconv.FormatFloat(float64(app.Calculation.SumTravelCosts), 'f', 2, 32),
		strconv.FormatFloat(float64(app.Calculation.SumDailyCharges), 'f', 2, 32),
		strconv.FormatFloat(float64(app.Calculation.SumNightlyCharges), 'f', 2, 32),
		strconv.FormatFloat(float64(app.Calculation.SumAdditionalCosts), 'f', 2, 32),
		strconv.FormatFloat(float64(app.Calculation.SumOfSums), 'f', 2, 32),
		})
	m.TableList(header, content, props.TableList{
		Align: consts.Center,
		HeaderProp: props.TableListContent{
			GridSizes: []uint{1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1},
		},
		ContentProp: props.TableListContent{
			GridSizes: []uint{1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1},
		},
		Line: true,
	})
	m.Row(5, func() {
		m.Col(4, func() {
			m.Text("Die sachliche Richtigkeit wird besätigt:", props.Text{
				Size: 7,
				Top: 2.5,
				Align: consts.Center,
			})
		})
		m.ColSpace(4)
		m.Col(4, func() {
			m.Text("gem. § 37 RGV 55: für die Richtigkeit der Angaben:", props.Text{
				Size: 7,
				Top: 2.5,
				Align: consts.Center,
			})
		})
	})
	m.Row(15, func() {})
	m.Line(1.0)
	m.Row(5, func() {
		m.Col(4, func() {
			m.Text("(Datum, Unterschrift der/s Anweisungsberechtigten)", props.Text{
				Size: 7,
				Top: 2.5,
				Align: consts.Center,
			})
		})
		m.ColSpace(4)
		m.Col(4, func() {
			m.Text("(Datum, Unterschrift der/s Rechnungslegers/in)", props.Text{
				Size: 7,
				Top: 2.5,
				Align: consts.Center,
			})
		})
	})

	savePath := filepath.Join(path, fmt.Sprintf(TravelInvoicePDFFileName, short))
	err := m.OutputFileAndClose(savePath)
	if err != nil {
		return "", fmt.Errorf("could not save pdf: %v", err)
	}
	return savePath, nil
}

func GenerateBusinessTripApplication(path, short string, app db.BusinessTripApplication, uuid string) (string, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(3, func() {
				_ = m.FileImage("../assets/TGM_Logo.png", props.Rect{
					Left:    0,
					Top:     0,
					Center:  false,
					Percent: 100,
				})
			})

			m.ColSpace(2)

			m.Col(2, func() {
				m.Text("Dienstreiseantrag Inland", props.Text{
					Align:  consts.Center,
					Family: consts.Helvetica,
					Size:   12,
				})
			})

			m.ColSpace(2)

			m.Col(3, func() {
				m.QrCode(URL+uuid, props.Rect{
					Left:    27,
					Top:     0,
					Center:  false,
					Percent: 100,
				})
			})
		})
	})
	m.Line(4.0)
	m.Row(10, func() {
		m.Col(1, func() {
			m.Text("Name:", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(4, func() {
			m.Text(app.Surname + " " + app.Name, props.Text{
				Top:   2.5,
				Align: consts.Center,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Akad. Grad:", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text(app.Degree, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Amtstitel:", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(1, func() {
			m.Text(app.Title, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Personalnr.:", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			m.Text(strconv.Itoa(app.Staffnr), props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Reiseziel:", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(5, func() {
			m.Text(app.TripGoal, props.Text{
				Top:   2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Dienstreise", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text("Beginn:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			m.Text(app.TripBeginTime.Format("02. 01. 2006 15:04 Uhr"), props.Text{
				Top: 2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Ende:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			m.Text(app.TripEndTime.Format("02. 01. 2006 15:04 Uhr"), props.Text{
				Top: 2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Dienstverrichtung", props.Text{
				Top:   2.5,
				Align: consts.Left,
			})
		})
		m.Col(2, func() {
			m.Text("Beginn:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			m.Text(app.ServiceBeginTime.Format("02. 01. 2006 15:04 Uhr"), props.Text{
				Top: 2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(2, func() {
			m.Text("Ende:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			m.Text(app.ServiceEndTime.Format("02. 01. 2006 15:04 Uhr"), props.Text{
				Top: 2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Reisezweck:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(10, func() {
			m.Text(app.TravelPurpose, props.Text{
				Top: 2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Reiseart:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(10, func() {
			travel := ""
			switch app.TravelMode {
			case db.OfficialBusinessCardClass2:
				travel = "Amtl BUSINESSKARTE 2. Kl"
				break
			case db.Passenger:
				travel = "MITFAHRER/INNEN"
				break
			case db.OfficialBusinessCardClass1:
				travel = "Amtl. BUSINESSKARTE / BAHNVERRECHNUNG 1. Kl - (Begründung erford.)"
				break
			case db.TravelGrant:
				travel = "BEFÖRDERUNGSZUSCHUSS"
				break
			case db.Flight:
				travel = "FLUG"
				break
			case db.TrainClass2:
				travel = "BAHN 2. Kl. - (Beleg erford.)"
				break
			case db.CheapFlight:
				travel = "BILLIGFLUG"
				break
			case db.OwnCar:
				travel = "EIGENER PKW - (Begründung erford.)"
				break
			case db.SleepTrain:
				travel = "SCHLAFWAGEN"
				break
			case db.Bus:
				travel = "BUS - (Beleg erford.)"
			}
			m.Text(travel, props.Text{
				Top: 2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(3, func() {
			m.Text("Ausgangspunkt:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			starting := ""
			if app.StartingPoint == db.Office {
				starting = "Dienststelle"
			} else if app.StartingPoint == db.OwnApartment {
				starting = "Wohnung"
			}
			m.Text(starting, props.Text{
				Top: 2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
		m.Col(3, func() {
			m.Text("Endpunkt:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			ending := ""
			if app.EndPoint == db.Office {
				ending = "Dienststelle"
			} else if app.EndPoint == db.OwnApartment {
				ending = "Wohnung"
			}
			m.Text(ending, props.Text{
				Top: 2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Begründung:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(10, func() {
			m.Text(app.Reasoning, props.Text{
				Top: 2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(10, func() {
		m.Col(3, func() {
			m.Text("Sonstige Teilnehmer/innen:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(9, func() {
			part := ""
			for _, t := range app.OtherParticipants {
				part = part + t + ", "
			}
			part = part[0:len(part) - 2]
			m.Text(part, props.Text{
				Top: 2.5,
				Align: consts.Left,
				Style: consts.Italic,
			})
		})
	})
	m.Line(1.0)
	m.Row(15, func() {
		m.Col(6, func() {
			text := ""
			if app.BonusMileConfirmation1 {
				text = "Ich bestätige, dass ich anlässlich von Dienstreisen im Rahmen personenbezogener" +
					" Bonusprogramme erworbene Prämien nicht privat in Anspruch nehme."
			} else {
				text = "1. nicht bestätigt"
			}
			m.Text(text, props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(6, func() {
			text := ""
			if app.BonusMileConfirmation1 {
				text = "Für die Dienstreise verwende ich auf meine Meilenkonto" +
					" gutgeschriebene, dienstlich erworbene Meilen."
			} else {
				text = "2. nicht bestätigt"
			}
			m.Text(text, props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
	})
	m.Row(10, func() {
		text := ""
		if app.TravelCostsPayedBySomeone && app.StayingCostsPayedBySomeone {
			text = "Es werden Aufenthaltskosten und Reisekosten von " + app.PayedByWhom + " getragen"
		} else if app.TravelCostsPayedBySomeone && !app.StayingCostsPayedBySomeone {
			text = "Es werden Reisekosten von " + app.PayedByWhom + " getragen"
		} else if !app.TravelCostsPayedBySomeone && app.StayingCostsPayedBySomeone {
			text = "Es werden Aufenthaltskosten von " + app.PayedByWhom + " getragen"
		} else {
			text = "Es werden keine Kosten von anderer Stelle getragen"
		}
		m.Col(12, func() {
			m.Text(text, props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
	})
	m.Row(10, func() {
		m.Col(3, func() {
			m.Text("Sonstige Kosten:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			m.Text(strconv.FormatFloat(float64(app.OtherCosts), 'f', 2, 32) + " €", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			m.Text("Geschätzte Kosten:", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
		m.Col(3, func() {
			m.Text(strconv.FormatFloat(float64(app.EstimatedCosts), 'f', 2, 32) + " €", props.Text{
				Top: 2.5,
				Align: consts.Left,
			})
		})
	})
	m.Row(15, func() {})
	m.Line(1.0)
	m.Row(5, func() {
		m.Col(4, func() {
			m.Text("Antragsteller/in", props.Text{
				Size: 7,
				Top: 2.5,
				Align: consts.Center,
			})
		})
		m.Col(8, func() {
			m.Text("Instituts-/Abteilungsleiter/in", props.Text{
				Size: 7,
				Top: 2.5,
				Align: consts.Center,
			})
		})
	})
	m.Line(10.0)
	m.Row(10, func() {
		m.Text("Die vorstehend beantragte Dienstreise wird mit " +
			app.DateApplicationApproved.Format("02. 01. 2006") + " genehmigt.", props.Text{
			Top: 2.5,
			Align: consts.Left,
		})
	})
	m.Row(15, func() {})
	m.Line(1.0)
	m.Row(5, func() {
		m.Col(4, func() {
			m.Text("Ort, Datum", props.Text{
				Size: 7,
				Top: 2.5,
				Align: consts.Center,
			})
		})
		m.ColSpace(4)
		m.Col(4, func() {
			m.Text("Unterschrift", props.Text{
				Size: 7,
				Top: 2.5,
				Align: consts.Center,
			})
		})
	})
	m.Line(10.0)
	m.Row(10, func() {
		m.ColSpace(8)
		m.Col(4, func() {
			m.Text("Eingabedatum: " + app.DateApplicationFiled.Format("02. 01. 2006"), props.Text{
				Top: 2.5,
				Align: consts.Right,
			})
		})
	})
	m.Row(10, func() {
		m.ColSpace(8)
		m.Col(4, func() {
			m.Text("Referent/in: " + app.Referee, props.Text{
				Top: 2.5,
				Align: consts.Right,
			})
		})
	})
	if app.BusinessCardEmittedOutward || app.BusinessCardEmittedReturn {
		m.Row(10, func() {
			m.ColSpace(6)
			text := ""
			if app.BusinessCardEmittedReturn && app.BusinessCardEmittedOutward {
				text = "Hin- und Rückfahrt"
			} else if app.BusinessCardEmittedOutward {
				text = "Hinfahrt"
			} else if app.BusinessCardEmittedReturn {
				text = "Rückfahrt"
			}
			m.Col(6, func() {
				m.Text("Businesskarte bei " + text + " ausgefolgt.", props.Text{
					Top: 2.5,
					Align: consts.Right,
				})
			})
		})
	}
	savePath := filepath.Join(path, fmt.Sprintf(BusinessTripApplicationPDFFileName, short))
	err := m.OutputFileAndClose(savePath)
	if err != nil {
		return "", fmt.Errorf("could not save pdf: %v", err)
	}
	return savePath, nil
}

func GenerateTravelInvoiceExcel(path, short string, app db.TravelInvoice) (string, error) {
	sourceF, err := os.Stat(filepath.Join(TemplatePath, ExcelTemplateTravelInvoicePath))
	if err != nil {
		return "", err
	}
	if !sourceF.Mode().IsRegular() {
		return "", fmt.Errorf("template file isn't regular")
	}
	source, err := os.Open(filepath.Join(TemplatePath, ExcelTemplateTravelInvoicePath))
	if err != nil {
		return "", err
	}
	defer source.Close()
	newPath := filepath.Join(path, fmt.Sprintf(TravelInvoiceExcelFileName, short))
	dest, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer dest.Close()
	_, err = io.Copy(dest, source)
	if err != nil {
		return "", err
	}
	excel, err := excelize.OpenFile(newPath)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIWorkplace, "tgm - Schule der Technik")
	if err != nil {
		return "", err
	}
	byear := app.TripBeginTime.Format("2006")
	bmonth := app.TripBeginTime.Format("01")
	bday := app.TripBeginTime.Format("02")
	bhour := app.TripBeginTime.Format("15")
	bminute := app.TripBeginTime.Format("04")
	err = excel.SetCellValue(Sheet, TITripBeginYear1, byear[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginYear2, byear[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginYear3, byear[2])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginYear4, byear[3])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginMonth1, bmonth[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginMonth2, bmonth[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginDay1, bday[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginDay2, bday[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginHour1, bhour[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginHour2, bhour[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginMinute1, bminute[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripBeginMinute2, bminute[1])
	if err != nil {
		return "", err
	}
	eyear := app.TripEndTime.Format("2006")
	emonth := app.TripEndTime.Format("01")
	eday := app.TripEndTime.Format("02")
	ehour := app.TripEndTime.Format("15")
	eminute := app.TripEndTime.Format("04")
	err = excel.SetCellValue(Sheet, TITripEndYear1, eyear[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndYear2, eyear[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndYear3, eyear[2])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndYear4, eyear[3])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndMonth1, emonth[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndMonth2, emonth[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndDay1, eday[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndDay2, eday[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndHour1, ehour[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndHour2, ehour[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndMinute1, eminute[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITripEndMinute2, eminute[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITravelCostsGrant, app.TravelCostsPreGrant)
	if err != nil {
		return "", err
	}
	extras, _ := ioutil.ReadDir(filepath.Join(path, UploadFolderName, short))
	err = excel.SetCellValue(Sheet, TIExtraAmount, len(extras))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIZI, app.ZI)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIFilingDate, app.FilingDate.Format("02.01.2006"))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TISurname, app.Surname)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIName, app.Name)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIDegree, app.Degree)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TITitle, app.Title)
	if err != nil {
		return "", err
	}
	ph := strconv.Itoa(app.Staffnr)
	if len(ph) != 8 {
		return "", fmt.Errorf("staffnr length isnt 8")
	}
	err = excel.SetCellValue(Sheet, TIPNR1, ph[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIPNR2, ph[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIPNR3, ph[2])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIPNR4, ph[3])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIPNR5, ph[4])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIPNR6, ph[5])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIPNR7, ph[6])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIPNR8, ph[7])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIStartingAddress, app.StartingPoint)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIDestinationAddress, app.EndPoint)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIClerk, app.Clerk)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIReviewer, app.Reviewer)
	if err != nil {
		return "", err
	}

	switch app.DailyChargesMode {
	case db.DailyChargesType1:
		err = excel.SetCellValue(Sheet, TICheckDailyChargesType1, CheckedCheckBox)
		break
	case db.DailyChargesType2:
		err = excel.SetCellValue(Sheet, TICheckDailyChargesType2, CheckedCheckBox)
		break
	case db.ToBeShortened:
		err = excel.SetCellValue(Sheet, TICheckToBeShortened, CheckedCheckBox)
		if err != nil {
			return "", err
		}
		err = excel.SetCellValue(Sheet, TIToBeShortenedAmount, app.ShortenedAmount)
		break
	}
	if err != nil {
		return "", err
	}

	switch app.NightlyChargesMode {
	case db.ProofNeededForCharges:
		err = excel.SetCellValue(Sheet, TICheckNightlyChargesProofNeededForCharges, CheckedCheckBox)
		break
	case db.NoProofNeeded:
		err = excel.SetCellValue(Sheet, TICheckNightlyChargesNoProofNeeded, CheckedCheckBox)
		break
	case db.NoClaimForNightlyCharges:
		err = excel.SetCellValue(Sheet, TICheckNightlyChargesNoClaimForNightlyCharges, CheckedCheckBox)
		break
	}
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIBreakfasts, app.Breakfasts)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TILunches, app.Lunches)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TIDinners, app.Dinners)
	if err != nil {
		return "", err
	}
	if app.OfficialBusinessCardGot {
		err = excel.SetCellValue(Sheet, TICheckOfficialBusinessCardGot, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	if app.TravelGrant {
		err = excel.SetCellValue(Sheet, TICheckTravelGrant, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	if app.ReplacementForAdvantageCard {
		err = excel.SetCellValue(Sheet, TICheckReplacementForAdvantageCard, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	if app.ReplacementForTrainCardClass2 {
		err = excel.SetCellValue(Sheet, TICheckReplacementForTrainCardClass2, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	if app.KilometreAllowance {
		err = excel.SetCellValue(Sheet, TICheckKilometreAllowance, CheckedCheckBox)
		if err != nil {
			return "", err
		}
		err = excel.SetCellValue(Sheet, TIKilometreAmount, app.KilometreAmount)
		if err != nil {
			return "", err
		}
	}
	if app.NRAndIndicationsOfParticipants {
		err = excel.SetCellValue(Sheet, TICheckNRAndIndicationsOfParticipants, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	if app.TravelCostsCited {
		err = excel.SetCellValue(Sheet, TICheckTravelCostsCited, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	if app.NoTravelCosts {
		err = excel.SetCellValue(Sheet, TICheckNoTravelCosts, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	if len(app.Calculation.Rows) > 5 {
		toGen := len(app.Calculation.Rows) - 5
		for i := 0; i < toGen; i++ {
			err = excel.DuplicateRow(Sheet, TIFirstGeneratedRow-1)
			if err != nil {
				return "", err
			}
		}
	}
	for i, row := range app.Calculation.Rows {
		rowNumber := strconv.Itoa(i + TIFirstConstantRow)
		err = excel.SetCellValue(Sheet, TICalcNrColumn+rowNumber, row.NR)
		if err != nil {
			return "", err
		}
		err = excel.SetCellValue(Sheet, TICalcDayColumn+rowNumber, row.Date.Format("02.01"))
		if err != nil {
			return "", err
		}
		err = excel.SetCellValue(Sheet, TICalcBeginColumn+rowNumber, row.Begin.Format("15:04"))
		if err != nil {
			return "", err
		}
		err = excel.SetCellValue(Sheet, TICalcEndColumn+rowNumber, row.End.Format("15:04"))
		if err != nil {
			return "", err
		}
		kinds := ""
		for _, kind := range row.KindsOfCost {
			kindString := ""
			switch kind {
			case db.TravelCosts:
				kindString = "Reisekosten"
				err = excel.SetCellValue(Sheet, TICalcKilometreAmountColumn+rowNumber, row.Kilometres)
				if err != nil {
					return "", err
				}
				err = excel.SetCellValue(Sheet, TICalcTravelCostsColumn+rowNumber, row.TravelCosts)
				if err != nil {
					return "", err
				}
				break
			case db.DailyCharges:
				kindString = "Tagesgebühr"
				err = excel.SetCellValue(Sheet, TICalcDailyChargesColumn+rowNumber, row.DailyCharges)
				if err != nil {
					return "", err
				}
				break
			case db.NightlyCharges:
				kindString = "Nächtigungsgebühr"
				err = excel.SetCellValue(Sheet, TICalcNightlyChargesColumn+rowNumber, row.NightlyCharges)
				if err != nil {
					return "", err
				}
				break
			case db.AdditionalCosts:
				err = excel.SetCellValue(Sheet, TICalcAdditionalCostsColumn+rowNumber, row.AdditionalCosts)
				if err != nil {
					return "", err
				}
				kindString = "Sonstige Nebenkosten"
				break
			}
			kinds = kinds + kindString + ", "
		}
		kinds = kinds[0 : len(kinds)-2]
		err = excel.SetCellValue(Sheet, TICalcKindOfFeeColumn+rowNumber, kinds)
		if err != nil {
			return "", err
		}
		err = excel.SetCellValue(Sheet, TICalcSumColumn+rowNumber, row.Sum)
	}
	row := strconv.Itoa(len(app.Calculation.Rows) + TIFirstConstantRow)
	err = excel.SetCellValue(Sheet, TICalcTravelCostsColumn+row, app.Calculation.SumTravelCosts)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TICalcDailyChargesColumn+row, app.Calculation.SumDailyCharges)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TICalcNightlyChargesColumn+row, app.Calculation.SumNightlyCharges)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TICalcAdditionalCostsColumn+row, app.Calculation.SumAdditionalCosts)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, TICalcSumColumn+row, app.Calculation.SumOfSums)
	if err != nil {
		return "", err
	}
	err = excel.Save()
	return newPath, err
}

func GenerateBusinessTripApplicationExcel(path, short string, app db.BusinessTripApplication) (string, error) {
	sourceF, err := os.Stat(filepath.Join(TemplatePath, ExcelTemplateBusinessTripApplicationPath))
	if err != nil {
		return "", err
	}
	if !sourceF.Mode().IsRegular() {
		return "", fmt.Errorf("template file isn't regular")
	}
	source, err := os.Open(filepath.Join(TemplatePath, ExcelTemplateBusinessTripApplicationPath))
	if err != nil {
		return "", err
	}
	defer source.Close()
	newPath := filepath.Join(path, fmt.Sprintf(BusinessTripApplicationExcelFileName, short))
	dest, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer dest.Close()
	_, err = io.Copy(dest, source)
	if err != nil {
		return "", err
	}
	excel, err := excelize.OpenFile(newPath)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAWorkplace, "tgm - Schule der Technik")
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTASurname, app.Surname)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAName, app.Name)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTADegree, app.Degree)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTATitle, app.Title)
	if err != nil {
		return "", err
	}
	//excel.SetCellValue(Sheet, BTATel, app.Tel)
	ph := strconv.Itoa(app.Staffnr)
	if len(ph) != 8 {
		return "", fmt.Errorf("staff nr doesnt match valid length")
	}
	err = excel.SetCellValue(Sheet, BTAPNR1, ph[0])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAPNR2, ph[1])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAPNR3, ph[2])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAPNR4, ph[3])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAPNR5, ph[4])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAPNR6, ph[5])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAPNR7, ph[6])
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAPNR8, ph[7])
	if err != nil {
		return "", err
	}
	//excel.SetCellValue(Sheet, BTAVGr, app.VGr)
	//excel.SetCellValue(Sheet, BTAEGr, app.EGr)
	//excel.SetCellValue(Sheet, BTADKI, app.DKl)
	//excel.SetCellValue(Sheet, BTAGSt, app.GSt)
	//excel.SetCellValue(Sheet, BTAESt, app.ESt)
	//excel.SetCellValue(Sheet, BTAFeeLevel, app.FeeLevel)
	err = excel.SetCellValue(Sheet, BTATripBeginDate, app.TripBeginTime.Format("02.01.2006"))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTATripBeginTime, app.TripBeginTime.Format("15:04"))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTATripEndDate, app.TripEndTime.Format("02.01.2006"))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTATripEndTime, app.TripEndTime.Format("15:04"))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAServiceBeginDate, app.ServiceBeginTime.Format("02.01.2006"))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAServiceBeginTime, app.ServiceBeginTime.Format("15:04"))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAServiceEndDate, app.ServiceEndTime.Format("02.01.2006"))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAServiceEndTime, app.ServiceEndTime.Format("15:04"))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTADestination, app.TripGoal)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTATravelReasoning, app.TravelPurpose)
	if err != nil {
		return "", err
	}
	switch app.TravelMode {
	case db.OfficialBusinessCardClass2:
		err = excel.SetCellValue(Sheet, BTACheckOfficialBusinessCardClass2, CheckedCheckBox)
		break
	case db.Passenger:
		err = excel.SetCellValue(Sheet, BTACheckPassenger, CheckedCheckBox)
		break
	case db.OfficialBusinessCardClass1:
		err = excel.SetCellValue(Sheet, BTACheckOfficialBusinessCardClass1, CheckedCheckBox)
		break
	case db.TravelGrant:
		err = excel.SetCellValue(Sheet, BTACheckTravelGrant, CheckedCheckBox)
		break
	case db.Flight:
		err = excel.SetCellValue(Sheet, BTACheckFlight, CheckedCheckBox)
		break
	case db.TrainClass2:
		err = excel.SetCellValue(Sheet, BTACheckTrainClass2, CheckedCheckBox)
		break
	case db.CheapFlight:
		err = excel.SetCellValue(Sheet, BTACheckCheapFlight, CheckedCheckBox)
		break
	case db.OwnCar:
		err = excel.SetCellValue(Sheet, BTACheckOwnCar, CheckedCheckBox)
		break
	case db.SleepTrain:
		err = excel.SetCellValue(Sheet, BTACheckSleepTrain, CheckedCheckBox)
		break
	case db.Bus:
		err = excel.SetCellValue(Sheet, BTACheckBus, CheckedCheckBox)
		break
	}
	if err != nil {
		return "", err
	}
	switch app.StartingPoint {
	case db.OwnApartment:
		err = excel.SetCellValue(Sheet, BTACheckStartAddressOwnApartment, CheckedCheckBox)
		break
	case db.Office:
		err = excel.SetCellValue(Sheet, BTACheckStartAddressOffice, CheckedCheckBox)
		break
	}
	if err != nil {
		return "", err
	}
	switch app.EndPoint {
	case db.OwnApartment:
		err = excel.SetCellValue(Sheet, BTACheckEndAddressOwnApartment, CheckedCheckBox)
		break
	case db.Office:
		err = excel.SetCellValue(Sheet, BTACheckEndAddressOffice, CheckedCheckBox)
		break
	}
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAReasoning, app.Reasoning)
	if err != nil {
		return "", err
	}
	participants := ""
	for _, part := range app.OtherParticipants {
		participants = participants + part + ", "
	}
	participants = participants[0 : len(participants)-2]
	err = excel.SetCellValue(Sheet, BTAOtherParticipants, participants)
	if err != nil {
		return "", err
	}
	if app.BonusMileConfirmation1 {
		err = excel.SetCellValue(Sheet, BTACheckBonusMiles1, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	if app.BonusMileConfirmation2 {
		err = excel.SetCellValue(Sheet, BTACheckBonusMiles2, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	if app.TravelCostsPayedBySomeone {
		err = excel.SetCellValue(Sheet, BTACheckTravelCostsPayedBySomeoneYes, CheckedCheckBox)
		if err != nil {
			return "", err
		}
		err = excel.SetCellValue(Sheet, BTACheckTravelCostsPayedBySomeoneNo, UncheckedCheckBox)
		if err != nil {
			return "", err
		}
	} else {
		err = excel.SetCellValue(Sheet, BTACheckTravelCostsPayedBySomeoneNo, CheckedCheckBox)
		if err != nil {
			return "", err
		}
		err = excel.SetCellValue(Sheet, BTACheckTravelCostsPayedBySomeoneYes, UncheckedCheckBox)
		if err != nil {
			return "", err
		}
	}

	if app.StayingCostsPayedBySomeone {
		err = excel.SetCellValue(Sheet, BTACheckStayingCostsPayedBySomeoneYes, CheckedCheckBox)
		if err != nil {
			return "", err
		}
		err = excel.SetCellValue(Sheet, BTACheckStayingCostsPayedBySomeoneNo, UncheckedCheckBox)
		if err != nil {
			return "", err
		}
	} else {
		err = excel.SetCellValue(Sheet, BTACheckStayingCostsPayedBySomeoneNo, CheckedCheckBox)
		if err != nil {
			return "", err
		}
		err = excel.SetCellValue(Sheet, BTACheckStayingCostsPayedBySomeoneYes, UncheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	err = excel.SetCellValue(Sheet, BTAPayedByWhom, app.PayedByWhom)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAOtherCosts, app.OtherCosts)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAEstimatedCosts, app.EstimatedCosts)
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAFilingDate, app.DateApplicationFiled.Format("02.01.2006"))
	if err != nil {
		return "", err
	}
	err = excel.SetCellValue(Sheet, BTAApprovalDate, app.DateApplicationApproved.Format("02.01.2006"))
	if err != nil {
		return "", err
	}
	if app.BusinessCardEmittedOutward {
		err = excel.SetCellValue(Sheet, BTACheckBusinessCardEmittedOutward, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	if app.BusinessCardEmittedReturn {
		err = excel.SetCellValue(Sheet, BTACheckBusinessCardEmittedReturn, CheckedCheckBox)
		if err != nil {
			return "", err
		}
	}
	err = excel.SetCellValue(Sheet, BTAReferee, app.Referee)
	if err != nil {
		return "", err
	}
	err = excel.Save()
	return newPath, err
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