package db

import "time"

// Enum for different progress states of an application
const (
	Rejected = iota
	InSubmission
	InProcess
	Confirmed
	Running
	CostsPending
	CostsInProcess
	Done
)

// Enum for different kinds of applications
const (
	SchoolEvent = iota
	Training
	Seminar
	Conference
	Course
	Miscellaneous
	OtherReason
	Careleave
	ServiceMandate
	MedicalAppointment
	Other
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
	UUID string `json:"uuid" example:"693aa616-9895-418b-8904-765f0f6d26a4"`
	// The name on how this Application should be referenced by
	Name string `json:"name" example:"Sommersportwoche"`
	// The kind of this Application (for more see the Enum for the kinds of Application on this level only Training, SchoolEvent and OtherReason is applicable, the sub kinds should be used in the further detail section of the corresponding site)
	Kind int `json:"kind" example:"0"`
	// The Reasoning of this Application (there is none if this isn't of the type Miscellaneous)
	MiscellaneousReason string `json:"miscellaneous_reason" example:"Guter Grund"`
	// The Progress of this Application in filing (for more see the Enum for the Progress)
	Progress int `json:"progress" example:"3"`
	// the time the underlying event of this Application starts
	StartTime time.Time `json:"start_time"`
	// the time the underlying event of this Application ends
	EndTime time.Time `json:"end_time"`
	// Other Notes regarding this Application
	Notes string `json:"notes" example:"Wichtig ist, dass wir die Reise bewilligen lassen!"`
	// The starting address of this Application
	StartAddress string `json:"start_address" example:"TGM, Wexstraße 19-23, 1220 Wien"`
	// The Destination Address of this Application
	DestinationAddress string `json:"destination_address" example:"Karl Hönck Heim, Kärnten"`
	// The timestamp this application was changed last
	LastChanged time.Time `json:"last_changed"`
	// Further Details if this is of the kind SchoolEvent, if not this will be empty
	SchoolEventDetails SchoolEventDetails `json:"school_event_details"`
	// Further Details if this is of the kind Training, if not this will be empty
	TrainingDetails TrainingDetails `json:"training_details"`
	// Further Details if this is of the kind of any other, if not this will be empty
	OtherReasonDetails OtherReasonDetails `json:"other_reason_details"`
	// The regarding BusinessTripApplication for each teacher
	BusinessTripApplications []BusinessTripApplication `json:"business_trip_applications"`
	// The regarding TravelInvoice for each teacher
	TravelInvoices []TravelInvoice `json:"travel_invoices"`
}

// SchoolEventDetails are details an Application has if it is of the kind of SchoolEvent
type SchoolEventDetails struct {
	// The participating classes
	Classes []string `json:"classes" example:"5BHIT,5AHIT,5CHIT,5DHIT"`
	// The amount of male students
	AmountMaleStudents int `json:"amount_male_students" example:"17"`
	// The amount of female students
	AmountFemaleStudents int `json:"amount_female_students" example:"0"`
	// The duration of the event in days
	DurationInDays int `json:"duration_in_days" example:"2"`
	// Details of each teacher participating in the SchoolEvent
	Teachers []SchoolEventTeacherDetails `json:"teachers"`
}

// SchoolEventTeacherDetails are details of each teacher participating in a SchoolEvent
type SchoolEventTeacherDetails struct {
	// The full name of a teacher
	Name string `json:"name" example:"Stefan Zakall"`
	// The short name (abbrevation) of a teacher
	Shortname string `json:"shortname" example:"szakall"`
	// The teacher will be attending the SchoolEvent from
	AttendanceFrom time.Time `json:"attendance_from"`
	// The teacher will be attend the SchoolEvent till
	AttendanceTill time.Time `json:"attendance_till"`
	// The group number
	Group int `json:"group" example:"1"`
	// Where the teacher starts their travel from
	StartAddress string `json:"start_address" example:"TGM, Wexstraße 19-23, 1220 Wien"`
	// Where the teacher will meet with the group to travel together
	MeetingPoint string `json:"meeting_point" example:"TGM, Wexstraße 19-23, 1220 Wien"`
	// The role of each teacher (Leader or Companion)
	Role int `json:"role" example:"0"`
}

// TrainingDetails are details an Application has if it is of the kind of Training
type TrainingDetails struct {
	// The kind of Training
	Kind int `json:"kind" example:"2"`
	// if its miscellaneous a reasoning for the Training
	MiscellaneousReason string `json:"miscellaneous_reason" example:"Ein sonstiger Grund"`
	// the ph number of the teacher
	PH int `json:"ph" example:"938503154"`
	// The company who organizes the Training
	Organizer string `json:"organizer" example:"Accenture"`
	// the teacher this application belongs to
	Filer string `json:"filer" example:"Stefan Zakall"`
}

// OtherReasonDetails are details an Application has if it isnt a Training or SchoolEvent
type OtherReasonDetails struct {
	// The kind of other Reason this Application is filed
	Kind int `json:"kind" example:"7"`
	// The title if the other reason is a ServiceMandate
	ServiceMandateTitle string `json:"service_mandate_title" example:"Dienstverrichtung"`
	// the gz number if the other reason is a ServiceMandate
	ServiceMandateGZ int `json:"service_mandate_gz"`
	// the reasoning if the other reason is of kind Miscellaneous
	MiscellaneousReason string `json:"miscellaneous_reason" example:"Ein guter Grund"`
	// the teacher this application belongs to
	Filer string `json:"filer" example:"Stefan Zakall"`
}

// A BusinessTripApplication represents one Business Trip Application belonging to an Application for each teacher
type BusinessTripApplication struct {
	// The id (counting upwards) of this BusinessTripApplication regarding to the uid
	ID int `json:"id" example:"1"`
	// Surname of the Teacher
	Surname string `json:"surname" example:"Zakall"`
	// Name of the Teacher
	Name string `json:"name" example:"Stefan"`
	// Degree of the Teacher
	Degree string `json:"degree" example:"DI"`
	// Title of the Teacher
	Title string `json:"title" example:"Prof"`
	// The staffnr of the regarding teacher
	Staffnr int `json:"staffnr" example:"938503154"`
	// The time the trip begins
	TripBeginTime time.Time `json:"trip_begin_time"`
	// The time the trip ends
	TripEndTime time.Time `json:"trip_end_time"`
	// The time the service begins
	ServiceBeginTime time.Time `json:"service_begin_time"`
	// The time the service ends
	ServiceEndTime time.Time `json:"service_end_time"`
	// The trip goal (address)
	TripGoal string `json:"trip_goal" example:"Technisches Museum Wien"`
	// The purpose of travelling
	TravelPurpose string `json:"travel_purpose" example:"Lehrausgang"`
	// The travel mode (see the regarding Enum for this)
	TravelMode int `json:"travel_mode" example:"6"`
	// The starting point (see the regarding Enum: OwnApartment or Office
	StartingPoint int `json:"starting_point" example:"1"`
	// The end point (see the regarding Enum: OwnApartment or Office)
	EndPoint int `json:"end_point" example:"1"`
	// The reasoing behind the trip application
	Reasoning string `json:"reasoning" example:"Lehrausgang ins technische Museum"`
	// The name of other participants of this trip
	OtherParticipants []string `json:"other_participants" example:"Markus Schabel,Gottfried Koppensteiner"`
	// the confirmation of the first bonus mile clause
	BonusMileConfirmation1 bool `json:"bonus_mile_confirmation_1" example:"true"`
	// the confirmation of the second bonus mile clause
	BonusMileConfirmation2 bool `json:"bonus_mile_confirmation_2" example:"true"`
	// whether the travel costs are paid by someone else
	TravelCostsPaidBySomeone bool `json:"travel_costs_paid_by_someone" example:"true"`
	// whether the staying costs are paid by someone else
	StayingCostsPaidBySomeone bool `json:"staying_costs_paid_by_someone" example:"true"`
	// if some costs are paid by someone else by whom
	PaidByWhom string `json:"paid_by_whom" example:"Technologenverband"`
	// other costs which appeared
	OtherCosts float32 `json:"other_costs" example:"2.42"`
	// the total estimated costs
	EstimatedCosts float32 `json:"estimated_costs" example:"25.32"`
	// the date this application is filed
	DateApplicationFiled time.Time `json:"date_application_filed"`
	// the date this application is approved
	DateApplicationApproved time.Time `json:"date_application_approved"`
	// the referee checking this application
	Referee string `json:"referee"`
	// whether a business card was emitted outwards
	BusinessCardEmittedOutward bool `json:"business_card_emitted_outward" example:"false"`
	// whether a business card was emitted on the return
	BusinessCardEmittedReturn bool `json:"business_card_emitted_return" example:"false"`
}

// A TravelInvoice represents one Travel Invoice belonging to an Application for each teacher
type TravelInvoice struct {
	// The id (counting upwards) of this TravelInvoice regarding to the uid
	ID int `json:"id" example:"1"`
	// Surname of the Teacher
	Surname string `json:"surname" example:"Zakall"`
	// Name of the Teacher
	Name string `json:"name" example:"Stefan"`
	// Degree of the Teacher
	Degree string `json:"degree" example:"DI"`
	// Title of the Teacher
	Title string `json:"title" example:""`
	// The time the trip begins
	TripBeginTime time.Time `json:"trip_begin_time"`
	// The time the trip ends
	TripEndTime time.Time `json:"trip_end_time"`
	// The granted travel costs
	TravelCostsPreGrant float32 `json:"travel_costs_pre_grant" example:"0"`
	// The personnel number of the teacher
	Staffnr int `json:"staffnr" example:"938503154"`
	// the starting point of the trip
	StartingPoint string `json:"starting_point" example:"TGM, Wexstraße 19-23, 1220 Wien"`
	// the end point of the trip
	EndPoint string `json:"end_point" example:"Hauptbahnhof, 1220 Wien"`
	// the clerk maintaining and checking this application
	Clerk string `json:"clerk"`
	// the reviewer reviewing the approval of this application
	Reviewer string `json:"reviewer"`
	// the zi number
	ZI int `json:"zi"`
	// the date this application was filed
	FilingDate time.Time `json:"filing_date"`
	// the date this application was approved
	ApprovalDate time.Time `json:"approval_date"`
	// the mode how daily charges are handled
	DailyChargesMode int `json:"daily_charges_mode" example:"1"`
	// the amount the daily charges should be shortened
	ShortenedAmount float32 `json:"shortened_amount" example:"0"`
	// the mode how nightly charges are handled
	NightlyChargesMode int `json:"nightly_charges_mode" example:"1"`
	// the amount of breakfasts
	Breakfasts int `json:"breakfasts" example:"2"`
	// the amount of lunches
	Lunches int `json:"lunches" example:"3"`
	// the amount of dinners
	Dinners int `json:"dinners" example:"4"`
	// whether the teacher got a official business card
	OfficialBusinessCardGot bool `json:"official_business_card_got" example:"true"`
	// whether the teacher got a travel grant
	TravelGrant bool `json:"travel_grant" example:"false"`
	// whether the teacher got a replacement for an advantage card
	ReplacementForAdvantageCard bool `json:"replacement_for_advantage_card" example:"false"`
	// whether the teacher got a replacement for a train card class 2
	ReplacementForTrainCardClass2 bool `json:"replacement_for_train_card_class_2" example:"false"`
	// whether the teacher got a kilometre allowance
	KilometreAllowance bool `json:"kilometre_allowance" example:"true"`
	// the regarding kilometre amount
	KilometreAmount float32 `json:"kilometre_amount" example:"25.12"`
	// whether the participants of the trip are counted and clearly indicated
	NRAndIndicationsOfParticipants bool `json:"nr_and_indications_of_participants" example:"true"`
	// whether the travel costs are clearly cited
	TravelCostsCited bool `json:"travel_costs_cited" example:"false"`
	// whether there aren't any travel costs
	NoTravelCosts bool `json:"no_travel_costs" example:"true"`
	// the regarding calculation
	Calculation Calculation `json:"calculation"`
}

// Calculation represent the calc field in a TravelInvoice
type Calculation struct {
	// the id of this calculation
	ID int `json:"id" example:"1"`
	// rows of this calculation
	Rows []Row `json:"rows"`
	// the sum of all travel costs
	SumTravelCosts float32 `json:"sum_travel_costs"`
	// the sum of all daily charges
	SumDailyCharges float32 `json:"sum_daily_charges"`
	// the sum of all nightly charges
	SumNightlyCharges float32 `json:"sum_nightly_charges"`
	// the sum of all additional costs
	SumAdditionalCosts float32 `json:"sum_additional_costs"`
	// the sum of all sums
	SumOfSums float32 `json:"sum_of_sums"`
}

// A Row in a Calculation
type Row struct {
	// The row nr
	NR int `json:"nr" example:"1"`
	// The date this Row refers to
	Date time.Time `json:"date"`
	// the begin time this Row refers to
	Begin time.Time `json:"begin"`
	// the end time this Row refers to
	End time.Time `json:"end"`
	// the kind of costs this row describes (see costs enum)
	KindsOfCost []int `json:"kind_of_cost" example:"1,2,3"`
	// the amount of kilometres this row refers to
	Kilometres float32 `json:"kilometres" example:"4.32"`
	// the travelCosts this Row conducts
	TravelCosts float32 `json:"travel_costs" example:"4.32"`
	// the dailyCharges this Row conducts
	DailyCharges float32 `json:"daily_charges" example:"4.32"`
	// the nightlyCharges this Row conducts
	NightlyCharges float32 `json:"nightly_charges" example:"4.32"`
	// the additionalCosts this Row conducts
	AdditionalCosts float32 `json:"additional_costs" example:"4.32"`
	// the sum of all costs in this Row
	Sum float32 `json:"sum" example:"4.32"`
}

// Teacher includes further information of a teacher (which isnt saved in the LDAP-instance)
type Teacher struct {
	// the uuid of this Teacher
	UUID string `json:"uuid" example:"3fcf7f67-e0ed-4339-99b4-a6765aaa3dc4"`
	// the short name of the Teacher
	Short string `json:"short" example:"szakall"`
	// the longname (firstname + sirname) of the Teacher
	Longname string `json:"longname" example:"Stefan Zakall"`
	// Superuser (total admin) of this software
	SuperUser bool `json:"super_user" example:"true"`
	// whether this Teacher as av rights
	AV bool `json:"av" example:"true"`
	// whether this Teacher as administration rights
	Administration bool `json:"administration" example:"true"`
	// whether this Teacher as pek rights
	PEK bool `json:"pek" example:"true"`
	// Degree of the Teacher
	Degree string `json:"degree" example:"DI"`
	// Title of the Teacher
	Title string `json:"title" example:"Prof"`
	// The Staffnr of the regarding teacher
	Staffnr int `json:"staffnr" example:"938503154"`
	// The Group number
	Group int `json:"group" example:"1"`
	// The StartingAddresses of the teacher
	StartingAddresses []string `json:"starting_addresses" example:"Zuhause 1,Zuhause 2"`
	// The TripGoals the teacher visited before
	TripGoals []string `json:"trip_goals" example:"Karl Hönck Heim,PH Wien,Landesgericht St. Pölten"`
	// The Departments this teacher belongs to
	Departments []string `json:"departments" example:"HIT,HBG"`
	// The Untis abbrevation of the teacher
	Untis string `json:"untis" example:"ZAKS"`
}
