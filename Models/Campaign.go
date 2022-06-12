package Models

import (
	"fmt"
	"reflect"
	"seentech/RECR/Utils"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Campaign struct {
	ID                   primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                 string               `json:"name" binding:"required"`
	Serial               string               `json:"serial"`
	Desc                 string               `json:"desc"`
	Type                 string               `json:"type" binding:"required"` // Online || Event
	Event                CampaignEvent        `json:"event,omitempty" bson:"event,omitempty"`
	Online               CampaignOnline       `json:"online,omitempty" bson:"online,omitempty"`
	Budget               float64              `json:"budget"`
	ActualCost           float64              `json:"actualcost"`
	Status               string               `json:"status"` // Draft || Review || Confirmed || Started || Closed || Canceled
	ExpectedStartingDate primitive.DateTime   `json:"expectedstartingdate,omitempty" bson:"expectedstartingdate,omitempty"`
	ExpectedEndingDate   primitive.DateTime   `json:"expectedendingdate,omitempty" bson:"expectedendingdate,omitempty"`
	Products             []primitive.ObjectID `json:"products,omitempty" bson:"products,omitempty" binding:"required"`
	ActualStartingDate   primitive.DateTime   `json:"actualstartingdate,omitempty" bson:"actualstartingdate,omitempty"`
	ActualEndingDate     primitive.DateTime   `json:"actualendingdate,omitempty" bson:"actualendingdate,omitempty"`
	OwnerRef             primitive.ObjectID   `json:"ownerref,omitempty" bson:"ownerref,omitempty"`
	ReviewerRef          primitive.ObjectID   `json:"reviewerref,omitempty" bson:"reviewerref,omitempty"`
	Team                 []primitive.ObjectID `json:"team,omitempty" bson:"team,omitempty"`
}

type CampaignEvent struct {
	Place    string   `json:"place"`
	Files    []string `json:"files"`
	FilesExt []string `json:"filesext"`
	Others   string   `json:"others"`
}

type CampaignOnline struct {
	ChannelRef primitive.ObjectID `json:"channelref,omitempty" bson:"channelref,omitempty"`
	Files      []string           `json:"files" bson:"files"`
	FilesExt   []string           `json:"filesext"`
	Others     string             `json:"others"`
}

type CampaignOnlinePopulated struct {
	ChannelRef Channel  `json:"channelref,omitempty" bson:"channelref,omitempty"`
	Files      []string `json:"files" bson:"files"`
	FilesExt   []string `json:"filesext"`
	Others     string   `json:"others"`
}

func (obj *CampaignOnlinePopulated) CloneFrom(other CampaignOnline) {
	obj.ChannelRef = Channel{}
	obj.Files = other.Files
	obj.FilesExt = other.FilesExt
	obj.Others = other.Others
}

type CampaignPopulated struct {
	ID                   primitive.ObjectID      `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                 string                  `json:"name" binding:"required"`
	Serial               string                  `json:"serial"`
	Desc                 string                  `json:"desc"`
	Type                 string                  `json:"type" binding:"required"`
	Event                CampaignEvent           `json:"event,omitempty" bson:"event,omitempty"`
	Online               CampaignOnlinePopulated `json:"online,omitempty" bson:"online,omitempty"`
	Budget               float64                 `json:"budget"`
	ActualCost           float64                 `json:"actualcost"`
	Status               string                  `json:"status" binding:"required"`
	ExpectedStartingDate primitive.DateTime      `json:"expectedstartingdate,omitempty" bson:"expectedstartingdate,omitempty"`
	ExpectedEndingDate   primitive.DateTime      `json:"expectedendingdate,omitempty" bson:"expectedendingdate,omitempty"`
	Products             []Product               `json:"products,omitempty" bson:"products,omitempty" binding:"required"`
	ActualStartingDate   primitive.DateTime      `json:"actualstartingdate,omitempty" bson:"actualstartingdate,omitempty"`
	ActualEndingDate     primitive.DateTime      `json:"actualendingdate,omitempty" bson:"actualendingdate,omitempty"`
	OwnerRef             User                    `json:"ownerref,omitempty" bson:"ownerref,omitempty"`
	ReviewerRef          User                    `json:"reviewerref,omitempty" bson:"reviewerref,omitempty"`
	Team                 []User                  `json:"team,omitempty" bson:"team,omitempty"`
}

type CampaignSearch struct {
	IDIsUsed                        bool                 `json:"idisused,omitempty" bson:"_id,omitempty"`
	ID                              primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	NameIsUsed                      bool                 `json:"nameisused" binding:"required"`
	Name                            string               `json:"name" binding:"required"`
	SerialIsUsed                    bool                 `json:"serialisused"`
	Serial                          string               `json:"serial"`
	Desc                            string               `json:"desc"`
	DescIsUsed                      bool                 `json:"descisused"`
	Type                            string               `json:"type" binding:"required"`
	TypeIsUsed                      bool                 `json:"typeisused" binding:"required"`
	Event                           CampaignEvent        `json:"event,omitempty" bson:"event,omitempty"`
	Online                          CampaignOnline       `json:"online,omitempty" bson:"online,omitempty"`
	Budget                          float64              `json:"budget"`
	ActualCost                      float64              `json:"actualcost"`
	StatusIsUsed                    bool                 `json:"statusisused,omitempty"`
	Status                          string               `json:"status" binding:"required"`
	ExpectedStartingDateRangeIsUsed bool                 `json:"expectedstartingdaterangeisused,omitempty" bson:"expectedstartingdaterangeisused,omitempty"`
	ExpectedStartingDateFrom        primitive.DateTime   `json:"expectedstartingdatefrom,omitempty" bson:"expectedstartingdatefrom,omitempty"`
	ExpectedStartingDateTo          primitive.DateTime   `json:"expectedstartingdateto,omitempty" bson:"expectedstartingdateto,omitempty"`
	ExpectedEndingDateRangeIsUsed   bool                 `json:"expectedendingdaterangeisused,omitempty" bson:"expectedendingdaterangeisused,omitempty"`
	ExpectedEndingDateFrom          primitive.DateTime   `json:"expectedendingdatefrom,omitempty" bson:"expectedendingdatefrom,omitempty"`
	ExpectedEndingDateTo            primitive.DateTime   `json:"expectedendingdateto,omitempty" bson:"expectedendingdateto,omitempty"`
	Products                        []primitive.ObjectID `json:"products,omitempty" bson:"products,omitempty" binding:"required"`
	ActualStartingDateRangeIsUsed   bool                 `json:"actualstartingdaterangeisused,omitempty" bson:"actualstartingdaterangeisused,omitempty"`
	ActualStartingDateFrom          primitive.DateTime   `json:"actualstartingdatefrom,omitempty" bson:"actualstartingdatefrom,omitempty"`
	ActualStartingDateTo            primitive.DateTime   `json:"actualstartingdateto,omitempty" bson:"actualstartingdateto,omitempty"`
	ActualEndingDateRangeIsUsed     bool                 `json:"actualendingdaterangeisused,omitempty" bson:"actualendingdaterangeisused,omitempty"`
	ActualEndingDateFrom            primitive.DateTime   `json:"actualendingdatefrom,omitempty" bson:"actualendingdatefrom,omitempty"`
	ActualEndingDateTo              primitive.DateTime   `json:"actualendingdateto,omitempty" bson:"actualendingdateto,omitempty"`
	OwnerRef                        primitive.ObjectID   `json:"ownerref,omitempty" bson:"ownerref,omitempty"`
	ReviewerRef                     primitive.ObjectID   `json:"reviewerref,omitempty" bson:"reviewerref,omitempty"`
	Team                            []primitive.ObjectID `json:"team,omitempty" bson:"team,omitempty"`
}

const (
	StatusDraft     string = "Draft"
	StatusReview           = "Review"
	StatusConfirmed        = "Confirmed"
	StatusStarted          = "Started"
	StatusClosed           = "Closed"
	StatusCanceled         = "Canceled"
)

const (
	TypeOnline string = "Online"
	TypeEvent         = "Event"
)

func (obj Campaign) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Desc, validation.Required),
		validation.Field(&obj.Type, validation.Required, validation.In("Event", "Online")),
		validation.Field(&obj.Status, validation.Required, validation.In("Draft", "Review", "Confirmed", "Started", "Closed", "Canceled")))
}

func (obj CampaignEvent) ValidateEvent() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Place, validation.Required))
}

func (obj CampaignOnline) ValidateOnline() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.ChannelRef, validation.Required))
}

func (obj *CampaignPopulated) CloneFrom(other Campaign) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.Serial = other.Serial
	obj.Desc = other.Desc
	obj.Type = other.Type
	obj.Event = other.Event
	obj.Online = CampaignOnlinePopulated{}
	obj.Budget = other.Budget
	obj.ActualCost = other.ActualCost
	obj.Status = other.Status
	obj.ExpectedStartingDate = other.ExpectedStartingDate
	obj.ExpectedEndingDate = other.ExpectedEndingDate
	obj.Products = []Product{}
	obj.ActualStartingDate = other.ActualStartingDate
	obj.ActualEndingDate = other.ActualEndingDate
	obj.OwnerRef = User{}
	obj.ReviewerRef = User{}
	obj.Team = []User{}

}

func (obj CampaignSearch) GetCampaignSearchBSONObj() bson.M {
	self := bson.M{}

	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.SerialIsUsed {
		self["serial"] = obj.Serial
	}

	if obj.TypeIsUsed {
		self["type"] = obj.Type
	}

	if obj.DescIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Desc)
		self["desc"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.ExpectedStartingDateRangeIsUsed {
		self["expectedstartingdate"] = bson.M{
			"$gte": obj.ExpectedStartingDateFrom,
			"$lte": obj.ExpectedStartingDateTo,
		}
	}

	if obj.ExpectedEndingDateRangeIsUsed {
		self["expectedendingdate"] = bson.M{
			"$gte": obj.ExpectedEndingDateFrom,
			"$lte": obj.ExpectedEndingDateTo,
		}
	}

	if obj.ActualStartingDateRangeIsUsed {
		self["actualstartingdate"] = bson.M{
			"$gte": obj.ActualStartingDateFrom,
			"$lte": obj.ActualStartingDateTo,
		}
	}

	if obj.ActualEndingDateRangeIsUsed {
		self["actualendingdate"] = bson.M{
			"$gte": obj.ActualEndingDateFrom,
			"$lte": obj.ActualEndingDateTo,
		}
	}

	return self
}

func (obj Campaign) GetCampaignModificationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()
	invalidFieldNames := []string{"ID", "Serial"}

	for i := 0; i < valueOfObj.NumField(); i++ {
		if Utils.ArrayStringContains(invalidFieldNames, typeOfObj.Field(i).Name) {
			continue
		}
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}
	return self
}
