package eveapi

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	//CharAccountBalanceURL is the url for the account balance endpoint
	CharAccountBalanceURL = "/char/AccountBalance.xml.aspx"
	//CharSkillQueueURL is the url for the skill queue endpoint
	CharSkillQueueURL = "/char/SkillQueue.xml.aspx"
	//MarketOrdersURL is the url for the market orders endpoint
	MarketOrdersURL = "/char/MarketOrders.xml.aspx"
    //WalletTransactionsURL is the url for the wallet transactions endpoint
	WalletTransactionsURL = "/char/WalletTransactions.xml.aspx"
    //CharacterSheetURL is the url for the character sheet endpoint
    CharacterSheetURL = "/char/CharacterSheet.xml.aspx"

    IndustryJobsURL = "/char/IndustryJobs.xml.aspx"
)

//AccountBalance is defined in corp.go

// CharAccountBalances calls /char/AccountBalance.xml.aspx
// Returns the account balance and any error if occured.
func (api API) CharAccountBalances(charID string) (*AccountBalance, error) {
	output := AccountBalance{}
	arguments := url.Values{}
	arguments.Add("characterID", charID)
	err := api.Call(CharAccountBalanceURL, arguments, &output)
	if err != nil {
		return nil, err
	}
	if output.Error != nil {
		return nil, output.Error
	}
	return &output, nil
}

//SkillQueueRow is an entry in a character's skill queue
type SkillQueueRow struct {
	Position  int     `xml:"queuePosition,attr"`
	TypeID    int     `xml:"typeID,attr"`
	Level     int     `xml:"level,attr"`
	StartSP   int     `xml:"startSP,attr"`
	EndSP     int     `xml:"endSP,attr"`
	StartTime eveTime `xml:"startTime,attr"`
	EndTime   eveTime `xml:"endTime,attr"`
}

func (s SkillQueueRow) String() string {
	return fmt.Sprintf("Position: %v, TypeID: %v, Level: %v, StartSP: %v, EndSP: %v, StartTime: %v, EndTime: %v", s.Position, s.TypeID, s.Level, s.StartSP, s.EndSP, s.StartTime, s.EndTime)
}

//SkillQueueResult is the result returned by the skill queue endpoint
type SkillQueueResult struct {
	APIResult
	SkillQueue []SkillQueueRow `xml:"result>rowset>row"`
}

// SkillQueue calls the API passing the parameter charID
// Returns a SkillQueueResult struct
func (api API) SkillQueue(charID string) (*SkillQueueResult, error) {
	output := SkillQueueResult{}
	arguments := url.Values{}
	arguments.Add("characterID", charID)
	err := api.Call(CharSkillQueueURL, arguments, &output)
	if err != nil {
		return nil, err
	}
	if output.Error != nil {
		return nil, output.Error
	}
	return &output, nil
}

//MarketOrdersResult is the result from calling the market orders endpoint
type MarketOrdersResult struct {
	APIResult
	Orders []MarketOrder `xml:"result>rowset>row"`
}

//MarketOrder is either a sell order or buy order
type MarketOrder struct {
	OrderID      int     `xml:"orderID,attr"`
	CharID       int64   `xml:"charID,attr"`
	StationID    int64   `xml:"stationID,attr"`
	VolEntered   int     `xml:"volEntered,attr"`
	VolRemaining int64   `xml:"volRemaining,attr"`
	MinVolume    int     `xml:"minVolume,attr"`
	TypeID       int64   `xml:"typeID,attr"`
	Range        int     `xml:"range,attr"`
	Division     int     `xml:"accountKey,attr"`
	Escrow       float64 `xml:"escrow,attr"`
	Price        float64 `xml:"price,attr"`
	IsBuyOrder   bool    `xml:"bid,attr"`
	Issued       eveTime `xml:"issued,attr"`
}

//MarketOrders returns the market orders for a given character
func (api API) MarketOrders(charID int64) (*MarketOrdersResult, error) {
	output := MarketOrdersResult{}
	args := url.Values{}
	args.Add("characterID", strconv.FormatInt(charID,10))
	err := api.Call(MarketOrdersURL, args, &output)
	if err != nil {
		return nil, err
	}
	if output.Error != nil {
		return nil, output.Error
	}
	return &output, nil
}

//MarketOrder is either a sell order or buy order
type WalletTransaction struct {
    TransactionDateTime  eveTime `xml:"transactionDateTime,attr"`  //datetime  Date and time of transaction.
    TransactionID        int64   `xml:"transactionID,attr"`        //long      Unique transaction ID.
    Quantity             int64   `xml:"quantity,attr"`             //int       Number of items bought or sold.
    TypeName             string  `xml:"typeName,attr"`             //string    Name of item bought or sold.
    TypeID               int64   `xml:"typeID,attr"`               //int       Type ID of item bought or sold.
    Price                float64 `xml:"price,attr"`                //decimal   Amount paid per unit.
    ClientID             int64   `xml:"clientID,attr"`             //long      Counterparty character or corporation ID. For NPC corporations, see the appropriate cross reference.
    ClientName           string  `xml:"clientName,attr"`           //string    Counterparty name.
    StationID            int64   `xml:"stationID,attr"`            //long      Station ID in which transaction took place.
    StationName          string  `xml:"stationName,attr"`          //string    Name of station in which transaction took place.
    TransactionType      string  `xml:"transactionType,attr"`      //string    Either "buy" or "sell" as appropriate.
    TransactionFor       string  `xml:"transactionFor,attr"`       //string    Either "personal" or "corporate" as appropriate.
    JournalTransactionID int64   `xml:"journalTransactionID,attr"` //long      Corresponding wallet journal refID.
    ClientTypeID         int64   `xml:"clientTypeID,attr"`         //long      Unknown meaning/mapping.
}
type WalletTransactionsResult struct {
    APIResult
    Transactions []WalletTransaction `xml:"result>rowset>row"`
}
//WalletTransactions returns the wallet transactions for a given character
//characterID	long	Character ID for which transactions will be requested
//accountKey	int	Account key of the wallet for which transactions will be returned. This is optional for character accounts which only have one wallet (accountKey = 1000). However, corporations have seven wallets with accountKeys numbered from 1000 through 1006. The Corp - AccountBalance call can be used to map corporation wallet to appropriate accountKey.
//fromID	use 0 to skip, long	Optional upper bound for the transaction ID of returned transactions. This argument is normally used to walk to the transaction log backwards. See Journal Walking for more information.
//rowCount	int	Optional limit on number of rows to return. Default is 1000. Maximum is 2560.
func (api API) WalletTransactions(charID int64, accountKey int64, fromID int64, rowCount int64) (*WalletTransactionsResult, error) {
	output := WalletTransactionsResult{}
	args := url.Values{}

    args.Add("characterID", strconv.FormatInt(charID,10))
    args.Add("accountKey", strconv.FormatInt(accountKey,10))
    if fromID != 0 {
        args.Add("fromID", strconv.FormatInt(fromID,10))
    }
    args.Add("rowCount", strconv.FormatInt(rowCount,10))

	err := api.Call(WalletTransactionsURL, args, &output)
	if err != nil {
		return nil, err
	}
	if output.Error != nil {
		return nil, output.Error
	}
	return &output, nil
}

func (api API) SimpleWalletTransactions(charID int64, fromID int64) (*WalletTransactionsResult, error) {
    return api.WalletTransactions(charID, 1000, fromID, 2560)
}

type Row struct {
    TypeID      int64 `xml:"typeID,attr"`
    Published   bool  `xml:"published,attr"`
    Level       int64 `xml:"level,attr"`
    SkillPoints int64 `xml:"skillpoints,attr"`
}

type Rowset struct{
    Name string `xml:"name,attr"`
    Rows []Row  `xml:"row"`
}

type CharacterSheetResult struct {
    APIResult
    Rowsets []Rowset `xml:"result>rowset"`
    Skills []Row
}

func (api API) CharacterSheet(charID int64) (*CharacterSheetResult, error) {
    output := CharacterSheetResult{}
    args := url.Values{}
    args.Add("characterID", strconv.FormatInt(charID,10))
    err := api.Call(CharacterSheetURL, args, &output)
    if err != nil {
        return nil, err
    }
    if output.Error != nil {
        return nil, output.Error
    }
    for _, v := range output.Rowsets{
        if v.Name == "skills"{
            output.Skills = v.Rows
        }
    }
    return &output, nil
}



type Job struct {
    JobID           int64  `xml:"jobID,attr"`
    InstallerID     int64  `xml:"installerID,attr"`
    InstallerName   string `xml:"installerName,attr"`
    FacilityID      int64  `xml:"facilityID,attr"`
    SolarSystemID   int64  `xml:"solarSystemID,attr"`
    SolarSystemName string `xml:"solarSystemName,attr"`
    StationID       int64  `xml:"stationID,attr"`
    ActivityID      int64  `xml:"activityID,attr"` //1 - mnf 4 - im.me
    BlueprintID     int64  `xml:"blueprintID,attr"`
    BlueprintTypeID int64  `xml:"blueprintTypeID,attr"`
    BlueprintTypeName string `xml:"blueprintTypeName,attr"`
    BlueprintLocationID int64 `xml:"blueprintLocationID,attr"`
    OutputLocationID int64 `xml:"outputLocationID,attr"`
    Runs            int64  `xml:"runs,attr"`
    Cost            float64 `xml:"cost,attr"`
    TeamID          int64  `xml:"teamID,attr"`
    LicensedRuns    int64  `xml:"licensedRuns,attr"`
    Probability     int64  `xml:"probability,attr"`
    ProductTypeID   int64  `xml:"productTypeID,attr"`
    ProductTypeName string `xml:"productTypeName,attr"`
    Status          int64  `xml:"status,attr"`
    TimeInSeconds   int64  `xml:"timeInSeconds,attr"`
    StartDate       eveTime `xml:"startDate,attr"`
    EndDate         eveTime `xml:"endDate,attr"`
    PauseDate       eveTime `xml:"pauseDate,attr"`
    CompletedDate   eveTime `xml:"completedDate,attr"`
    CompletedCharacterID int64 `xml:"completedCharacterID,attr"`
    SuccessfulRuns  int64 `xml:"successfulRuns,attr"`
}



type IndustryJobsResult struct {
    APIResult
    Jobs []Job `xml:"result>rowset>row"`
}

func (api API) IndustryJobs(charID int64) (*IndustryJobsResult, error) {
    output := IndustryJobsResult{}
    args := url.Values{}
    args.Add("characterID", strconv.FormatInt(charID,10))
    err := api.Call(IndustryJobsURL, args, &output)
    if err != nil {
        return nil, err
    }
    if output.Error != nil {
        return nil, output.Error
    }
    return &output, nil
}
