package db

import "time"

// Enum for different progress states of a school event
const (
	SERejected = iota
	SEInSubmission
	SEInProcess
	SEConfirmed
	SERunning
	SECostsPending
	SECostsInProcess
	SEDone
)

// Enum for different progress states of a training
const (
	TRejected = iota
	TInProcess
	TConfirmed
	TRunning
	TCostsPending
	TCostsInProcess
	TDone
)

// Enum for different kinds of applications
const (
	Training = iota
	Careleave
	ServiceMandate
	MedicalAppointment
	SchoolEvent
	Seminar
	Conference
	Course
	Miscellaneous
)

// Enum for different roles of a teacher in an application
const (
	Leader = iota
	Companion
)

// Enum for different modes of travel
const (
	OfficialBusinessCardClass2 = iota
	Passenger
	OfficialBusinessCardClass1
	TravelGrant
	Flight
	CheapFlight
	TrainClass2
	OwnCar
	SleepTrain
	Bus
)

// Enum for start and end place
const (
	OwnApartment = iota
	Office
)

// Enum for kinds of costs
const (
	TravelCosts = iota
	DailyCharges
	NightlyCharges
	AdditionalCosts
)

// Enum for modes of daily charges
const (
	DailyChargesType1 = iota
	DailyChargesType2
	ToBeShortened
)


// Enum for modes of nightly charges
const (
	ProofNeededForCharges = iota
	NoProofNeeded
	NoClaimForNightlyCharges
)

// An Application filed by a teacher represents the core group of data in this Application
type Application struct {
	// A generated uuid of this application
	uuid string
	// The name on how this Application should be referenced by
	name string
	// The kind of this Application (for more see the Enum for the kinds of Application)
	kind int
	// The Reasoning of this Application (there is none if this isn't of the type Miscellaneous)
	miscellaneousReason string
	// The Progress of this Application in filing (for more see the Enum for the Progress)
	progress int
	// the time the underlying event of this Application starts
	startTime time.Time
	// the time the underlying event of this Application ends
	endTime time.Time
	// Other Notes regarding this Application
	notes string
	// The starting address of this Application
	startAddress string
	// The Destination Address of this Application
	destinationAddress string
	// The timestamp this application was changed last
	lastChanged time.Time
	// Further Details if this is of the kind SchoolEvent, if not this will be empty
	SchoolEventDetails
	// Further Details if this is of the kind Training, if not this will be empty
	TrainingDetails
	// Further Details if this is of the kind of any other, if not this will be empty
	OtherReasonDetails
	// The regarding BusinessTripApplication for each teacher
	businessTripApplications []BusinessTripApplication
	// The regarding TravelInvoice for each teacher
	travelInvoices []TravelInvoice
}

// Details an Application has if it is of the kind of SchoolEvent
type SchoolEventDetails struct {
	// The participating classes
	classes []string
	// The amount of male students
	amountMaleStudents int
	// The amount of female students
	amountFemaleStudents int
	// The duration of the event in days
	durationInDays int
	// Details of each teacher participating in the SchoolEvent
	teachers []SchoolEventTeacherDetails
}

// The details of each teacher participating in a SchoolEvent
type SchoolEventTeacherDetails struct {
	// The full name of a teacher
	name string
	// The short name (abbrevation) of a teacher
	shortname string
	// The teacher will be attending the SchoolEvent from
	attendanceFrom time.Time
	// The teacher will be attend the SchoolEvent till
	attendanceTill time.Time
	// The group number
	group int
	// Where the teacher starts their travel from
	startAddress string
	// Where the teacher will meet with the group to travel together
	meetingPoint string
	// The role of each teacher (Leader or Companion)
	role int
}

// Details an Application has if it is of the kind of Training
type TrainingDetails struct {
	// The kind of Training
	kind int
	// if its miscellaneous a reasoning for the Training
	miscellaneousReason string
	// the personnell number of the teacher
	ph int
	// The company who organizes the Training
	organizer string
}

// Details an Application has if it isnt a Training or SchoolEvent
type OtherReasonDetails struct {
	// The kind of other Reason this Application is filed
	kind int
	// The title if the other reason is a ServiceMandate
	serviceMandateTitle string
	// the gz number if the other reason is a ServiceMandate
	serviceMandateGZ int
	// the reasoning if the other reason is of kind Miscellaneous
	miscellaneousReason string
}

// A BusinessTripApplication represents one Business Trip Application belonging to an Application for each teacher
type BusinessTripApplication struct {
	// The id (counting upwards) of this BusinessTripApplication regarding to the uid
	id int
	// The staffnr of the regarding teacher
	staffnr int
	// The time the trip begins
	tripBeginTime time.Time
	// The time the trip ends
	tripEndTime time.Time
	// The time the service begins
	serviceBeginTime time.Time
	// The time the service ends
	serviceEndTime time.Time
	// The trip goal (address)
	tripGoal string
	// The purpose of travelling
	travelPurpose string
	// The travel mode (see the regarding Enum for this)
	travelMode int
	// The starting point (see the regarding Enum: OwnApartment or Office
	startingPoint int
	// The end point (see the regarding Enum: OwnApartment or Office)
	endPoint int
	// The reasoing behind the trip application
	reasoning string
	// The name of other participants of this trip
	otherParticipants []string
	// the confirmation of the first bonus mile clause
	bonusMileConfirmation1 bool
	// the confirmation of the second bonus mile clause
	bonusMileConfirmation2 bool
	// whether the travel costs are payed by someone else
	travelCostsPayedBySomeone bool
	// whether the staying costs are payed by someone else
	stayingCostsPayedBySomeone bool
	// if some costs are payed by someone else by whom
	payedByWhom string
	// other costs which appeared
	otherCosts float32
	// the total estimated costs
	estimatedCosts float32
	// the date this application is filed
	dateApplicationFiled time.Time
	// the date this application is approved
	dateApplicationApproved time.Time
	// the referee checking this application
	referee string
	// whether a business card was emitted outwards
	businessCardEmittedOutward bool
	// whether a business card was emitted on the return
	businessCardEmittedReturn bool
}

// A TravelInvoice represents one Travel Invoice belonging to an Application for each teacher
type TravelInvoice struct {
	// The id (counting upwards) of this TravelInvoice regarding to the uid
	id int
	// The time the trip begins
	tripBeginTime time.Time
	// The time the trip ends
	tripEndTime time.Time
	// The personell number of the teacher
	staffnr int
	// the starting point of the trip
	startingPoint string
	// the end point of the trip
	endPoint string
	// the clerk maintaining and checking this application
	clerk string
	// the reviewer reviewing the approval of this application
	reviewer string
	// the travel mode (see the regarding enum)
	travelMode int
	// the zi number
	zi int
	// the date this application was filed
	filingDate time.Time
	// the date this application was approved
	approvalDate time.Time
	// the mode how daily charges are handled
	dailyChargesMode int
	// the amount the daily charges should be shortened
	shortenedAmount int
	// the mode how nightly charges are handled
	nightlyChargesMode int
	// the amount of breakfasts
	breakfasts int
	// the amount of lunches
	lunches int
	// the amount of dinners
	dinners int
	// whether the teacher got a official business card
	officialBusinessCardGot bool
	// whether the teacher got a travel grant
	travelGrant bool
	// whether the teacher got a replacement for an advantage card
	replacementForAdvantageCard bool
	// whether the teacher got a replacement for a train card class 2
	replacementForTrainCardClass2 bool
	// whether the teacher got a kilometre allowance
	kilometreAllowance bool
	// the regarding kilometre amount
	kilometreAmount float32
	// whether the participants of the trip are counted and clearly indicated
	nrAndIdicationsOfParticipants bool
	// whether the travel costs are clearly cited
	travelCostsCited bool
	// whether there aren't any travel costs
	noTravelCosts bool
	// the regarding calculation
	Calculation
}

// The calculations in a TravelInvoice
type Calculation struct {
	// the id of this calculation
	id int
	// rows of this calculation
	rows []Row
	// the sum of all travel costs
	sumTravelCosts float32
	// the sum of all daily charges
	sumDailyCharges float32
	// the sum of all nightly charges
	sumNightlyCharges float32
	// the sum of all additional costs
	sumAdditionalCosts float32
	// the sum of all sums
	sumOfSums float32
}

// The specific Row in a Calculation
type Row struct {
	// The row nr
	nr int
	// The date this Row refers to
	date time.Time
	// the begin time this Row refers to
	begin time.Time
	// the end time this Row refers to
	end time.Time
	// the amount of kilometres this row refers to
	kilometres float32
	// the travelCosts this Row conducts
	travelCosts float32
	// the dailyCharges this Row conducts
	dailyCharges float32
	// the nightlyCharges this Row conducts
	nightlyCharges float32
	// the additionalCosts this Row conducts
	additionalCosts float32
	// the sum of all costs in this Row
	sum float32
}

// Further information of a Teacher (which isnt saved in the LDAP-instance)
type Teacher struct {
	// the uuid of this Teacher
	uuid string
	// the short name of the Teacher
	short string
	// the longname (firstname + sirname) of the Teacher
	longname string
	// whether this Teacher as av rights
	av bool
	// whether this Teacher as pek rights
	pek bool
}
