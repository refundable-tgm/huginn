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
	UUID string
	// The name on how this Application should be referenced by
	Name string
	// The kind of this Application (for more see the Enum for the kinds of Application)
	Kind int
	// The Reasoning of this Application (there is none if this isn't of the type Miscellaneous)
	MiscellaneousReason string
	// The Progress of this Application in filing (for more see the Enum for the Progress)
	Progress int
	// the time the underlying event of this Application starts
	StartTime time.Time
	// the time the underlying event of this Application ends
	EndTime time.Time
	// Other Notes regarding this Application
	Notes string
	// The starting address of this Application
	StartAddress string
	// The Destination Address of this Application
	DestinationAddress string
	// The timestamp this application was changed last
	LastChanged time.Time
	// Further Details if this is of the kind SchoolEvent, if not this will be empty
	SchoolEventDetails SchoolEventDetails
	// Further Details if this is of the kind Training, if not this will be empty
	TrainingDetails TrainingDetails
	// Further Details if this is of the kind of any other, if not this will be empty
	OtherReasonDetails OtherReasonDetails
	// The regarding BusinessTripApplication for each teacher
	BusinessTripApplications []BusinessTripApplication
	// The regarding TravelInvoice for each teacher
	TravelInvoices []TravelInvoice
}

// Details an Application has if it is of the kind of SchoolEvent
type SchoolEventDetails struct {
	// The participating classes
	Classes []string
	// The amount of male students
	AmountMaleStudents int
	// The amount of female students
	AmountFemaleStudents int
	// The duration of the event in days
	DurationInDays int
	// Details of each teacher participating in the SchoolEvent
	Teachers []SchoolEventTeacherDetails
}

// The details of each teacher participating in a SchoolEvent
type SchoolEventTeacherDetails struct {
	// The full name of a teacher
	Name string
	// The short name (abbrevation) of a teacher
	Shortname string
	// The teacher will be attending the SchoolEvent from
	AttendanceFrom time.Time
	// The teacher will be attend the SchoolEvent till
	AttendanceTill time.Time
	// The group number
	Group int
	// Where the teacher starts their travel from
	StartAddress string
	// Where the teacher will meet with the group to travel together
	MeetingPoint string
	// The role of each teacher (Leader or Companion)
	Role int
}

// Details an Application has if it is of the kind of Training
type TrainingDetails struct {
	// The kind of Training
	Kind int
	// if its miscellaneous a reasoning for the Training
	MiscellaneousReason string
	// the personnell number of the teacher
	PH int
	// The company who organizes the Training
	Organizer string
}

// Details an Application has if it isnt a Training or SchoolEvent
type OtherReasonDetails struct {
	// The kind of other Reason this Application is filed
	Kind int
	// The title if the other reason is a ServiceMandate
	ServiceMandateTitle string
	// the gz number if the other reason is a ServiceMandate
	ServiceMandateGZ int
	// the reasoning if the other reason is of kind Miscellaneous
	MiscellaneousReason string
}

// A BusinessTripApplication represents one Business Trip Application belonging to an Application for each teacher
type BusinessTripApplication struct {
	// The id (counting upwards) of this BusinessTripApplication regarding to the uid
	ID int
	// The staffnr of the regarding teacher
	Staffnr int
	// The time the trip begins
	TripBeginTime time.Time
	// The time the trip ends
	TripEndTime time.Time
	// The time the service begins
	ServiceBeginTime time.Time
	// The time the service ends
	ServiceEndTime time.Time
	// The trip goal (address)
	TripGoal string
	// The purpose of travelling
	TravelPurpose string
	// The travel mode (see the regarding Enum for this)
	TravelMode int
	// The starting point (see the regarding Enum: OwnApartment or Office
	StartingPoint int
	// The end point (see the regarding Enum: OwnApartment or Office)
	EndPoint int
	// The reasoing behind the trip application
	Reasoning string
	// The name of other participants of this trip
	OtherParticipants []string
	// the confirmation of the first bonus mile clause
	BonusMileConfirmation1 bool
	// the confirmation of the second bonus mile clause
	BonusMileConfirmation2 bool
	// whether the travel costs are payed by someone else
	TravelCostsPayedBySomeone bool
	// whether the staying costs are payed by someone else
	StayingCostsPayedBySomeone bool
	// if some costs are payed by someone else by whom
	PayedByWhom string
	// other costs which appeared
	OtherCosts float32
	// the total estimated costs
	EstimatedCosts float32
	// the date this application is filed
	DateApplicationFiled time.Time
	// the date this application is approved
	DateApplicationApproved time.Time
	// the referee checking this application
	Referee string
	// whether a business card was emitted outwards
	BusinessCardEmittedOutward bool
	// whether a business card was emitted on the return
	BusinessCardEmittedReturn bool
}

// A TravelInvoice represents one Travel Invoice belonging to an Application for each teacher
type TravelInvoice struct {
	// The id (counting upwards) of this TravelInvoice regarding to the uid
	ID int
	// The time the trip begins
	TripBeginTime time.Time
	// The time the trip ends
	TripEndTime time.Time
	// The personell number of the teacher
	Staffnr int
	// the starting point of the trip
	StartingPoint string
	// the end point of the trip
	EndPoint string
	// the clerk maintaining and checking this application
	Clerk string
	// the reviewer reviewing the approval of this application
	Reviewer string
	// the travel mode (see the regarding enum)
	TravelMode int
	// the zi number
	ZI int
	// the date this application was filed
	FilingDate time.Time
	// the date this application was approved
	ApprovalDate time.Time
	// the mode how daily charges are handled
	DailyChargesMode int
	// the amount the daily charges should be shortened
	ShortenedAmount int
	// the mode how nightly charges are handled
	NightlyChargesMode int
	// the amount of breakfasts
	Breakfasts int
	// the amount of lunches
	Lunches int
	// the amount of dinners
	Dinners int
	// whether the teacher got a official business card
	OfficialBusinessCardGot bool
	// whether the teacher got a travel grant
	TravelGrant bool
	// whether the teacher got a replacement for an advantage card
	ReplacementForAdvantageCard bool
	// whether the teacher got a replacement for a train card class 2
	ReplacementForTrainCardClass2 bool
	// whether the teacher got a kilometre allowance
	KilometreAllowance bool
	// the regarding kilometre amount
	KilometreAmount float32
	// whether the participants of the trip are counted and clearly indicated
	NRAndIdicationsOfParticipants bool
	// whether the travel costs are clearly cited
	TravelCostsCited bool
	// whether there aren't any travel costs
	NoTravelCosts bool
	// the regarding calculation
	Calculation Calculation
}

// The calculations in a TravelInvoice
type Calculation struct {
	// the id of this calculation
	ID int
	// rows of this calculation
	Rows []Row
	// the sum of all travel costs
	SumTravelCosts float32
	// the sum of all daily charges
	SumDailyCharges float32
	// the sum of all nightly charges
	SumNightlyCharges float32
	// the sum of all additional costs
	SumAdditionalCosts float32
	// the sum of all sums
	SumOfSums float32
}

// The specific Row in a Calculation
type Row struct {
	// The row nr
	NR int
	// The date this Row refers to
	Date time.Time
	// the begin time this Row refers to
	Begin time.Time
	// the end time this Row refers to
	End time.Time
	// the amount of kilometres this row refers to
	Kilometres float32
	// the travelCosts this Row conducts
	TravelCosts float32
	// the dailyCharges this Row conducts
	DailyCharges float32
	// the nightlyCharges this Row conducts
	NightlyCharges float32
	// the additionalCosts this Row conducts
	AdditionalCosts float32
	// the sum of all costs in this Row
	Sum float32
}

// Further information of a Teacher (which isnt saved in the LDAP-instance)
type Teacher struct {
	// the uuid of this Teacher
	UUID string
	// the short name of the Teacher
	Short string
	// the longname (firstname + sirname) of the Teacher
	Longname string
	// whether this Teacher as av rights
	AV bool
	// whether this Teacher as administration rights
	Administration bool
	// whether this Teacher as pek rights
	PEK bool
}
