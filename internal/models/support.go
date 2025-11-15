package models

type CategoryContacting string

const (
	RequestBug              CategoryContacting = "Bug"
	RequestOffer            CategoryContacting = "Offer"
	RequestProductComplaint CategoryContacting = "Product Complaint"
)

type StatusContacting string

const (
	RequestOpened StatusContacting = "Opened"
	RequestClosed StatusContacting = "Closed"
	RequestAtWork StatusContacting = "At Work"
)

type Support struct {
	ID              int
	UserID          int
	CategoryRequest CategoryContacting
	StatusRequest   StatusContacting
	Message         string
}
