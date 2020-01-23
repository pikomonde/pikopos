package entity

// CompanyStatus is used for listing paid tier of the company account
type CompanyStatus int

const (
	// CompanyStatusFree is for free account
	CompanyStatusFree CompanyStatus = iota
	// CompanyStatusFreeAds is for free account with ads
	CompanyStatusFreeAds
	// CompanyStatusPaid01 is for paid account tier 01
	CompanyStatusPaid01
)

func (cs CompanyStatus) String() string {
	return [...]string{"free", "free-ads", "paid-01"}[cs]
}

// Company is used to contains company data
type Company struct {
	ID       int
	Username string `max:"16"`
	Name     string `max:"32"`
	Status   CompanyStatus
}

// CompanyDetail contains details of a company
type CompanyDetail struct {
	ID          int
	CompanyID   int
	Address     *string `max:"128"`
	City        *string `max:"64"`
	Province    *string `max:"32"`
	Postal      *string `max:"5"`
	Lat         *float64
	Lon         *float64
	PhoneNumber *string `max:"16"`
	WhatsApp    *string `max:"16"`
	Email       *string `max:"48"`
	Website     *string `max:"64"`
	Twitter     *string `max:"48"`
	Facebook    *string `max:"48"`
}
