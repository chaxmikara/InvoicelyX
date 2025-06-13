package models

import "time"

type Invoice struct {
	ID          string    `bson:"_id" json:"id"`
	UserID      string    `bson:"userID" json:"userID"`
	ClientName  string    `bson:"clientName" json:"clientName"`
	ClientEmail string    `bson:"clientEmail" json:"clientEmail"`
	IssueDate   time.Time `bson:"issueDate" json:"issueDate"`
	DueDate     time.Time `bson:"dueDate" json:"dueDate"`
	//invoice item
	//total amount
	Status      string    `bson:"status" json:"status"`
	PdfFilePath string    `bson:"pdfFilePath" json:"pdfFilePath"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}
