package l27

import (
	"encoding/json"
	"fmt"
	"net/url"
)

//gets extensions for domains
func (c *Client) Extension() ([]DomainProvider, error) {
	var extensions struct {
		Data []DomainProvider `json:"providers"`
	}

	endpoint := "domains/providers"
	err := c.invokeAPI("GET", endpoint, nil, &extensions)

	return extensions.Data, err
}

// Gets a single domain from the API
func (c *Client) Domain(id int) (Domain, error) {
	var domain struct {
		Data Domain `json:"domain"`
	}

	endpoint := fmt.Sprintf("domains/%d", id)
	err := c.invokeAPI("GET", endpoint, nil, &domain)

	return domain.Data, err
}

func (c *Client) LookupDomain(name string) ([]Domain, error) {
	results := []Domain{}
	domains, err := c.Domains(CommonGetParams{Filter: name})
	if err != nil {
		return nil, err
	}

	for _, domain := range domains {
		if domain.Fullname == name {
			results = append(results, domain)
		}
	}

	return results, err
}

//Domain gets a domain from the API
func (c *Client) Domains(getParams CommonGetParams) ([]Domain, error) {
	var domains struct {
		Data []Domain `json:"domains"`
	}

	endpoint := fmt.Sprintf("domains?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &domains)

	return domains.Data, err
}

// ------------------ /DOMAINS --------------------------

// DELETE DOMAIN
func (c *Client) DomainDelete(id int) error {
	endpoint := fmt.Sprintf("domains/%d", id)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// CREATE DOMAIN [lvl domain create <parmeters>]
func (c *Client) DomainCreate(req DomainRequest) (Domain, error) {
	if req.Action == "" {
		req.Action = "none"
	}

	var domain struct {
		Data Domain `json:"domain"`
	}

	err := c.invokeAPI("POST", "domains", req, &domain)

	return domain.Data, err
}

// TRANSFER DOMAIN [lvl domain transfer <parameters>]
func (c *Client) DomainTransfer(req DomainRequest) error {
	if req.Action == "" {
		req.Action = "transfer"
	}

	err := c.invokeAPI("POST", "domains", req, nil)

	return err
}

// INTERNAL TRANSFER
func (c *Client) DomainInternalTransfer(id int, req DomainRequest) error {
	endpoint := fmt.Sprintf("domains/%d/internaltransfer", id)
	err := c.invokeAPI("POST", endpoint, req, nil)

	return err
}

// UPDATE DOMAIN [lvl update <parameters>]
func (c *Client) DomainUpdate(id int, data map[string]interface{}) error {
	endpoint := fmt.Sprintf("domains/%d", id)
	err := c.invokeAPI("PATCH", endpoint, data, nil)

	return err
}

// PUT /domains/{domainID}
func (c *Client) DomainUpdatePut(id int, data DomainRequest) error {
	endpoint := fmt.Sprintf("domains/%d", id)
	err := c.invokeAPI("PUT", endpoint, data, nil)

	return err
}

// ------------------ /DOMAIN/RECORDS ----------------------
// GET
func (c *Client) DomainRecords(id int, recordType string, getParams CommonGetParams) ([]DomainRecord, error) {
	var records struct {
		Records []DomainRecord `json:"records"`
	}

	endpoint := fmt.Sprintf("domains/%d/records?%s", id, formatCommonGetParams(getParams))
	if recordType != "" {
		endpoint += fmt.Sprintf("&type=%s", recordType)
	}
	err := c.invokeAPI("GET", endpoint, nil, &records)

	return records.Records, err
}

func (c *Client) DomainRecord(domainId int, recordId int) (DomainRecord, error) {
	var records struct {
		Record DomainRecord `json:"record"`
	}

	endpoint := fmt.Sprintf("domains/%d/records/%d", domainId, recordId)
	err := c.invokeAPI("GET", endpoint, nil, &records)

	return records.Record, err
}

// CREATE
func (c *Client) DomainRecordCreate(id int, req DomainRecordRequest) (DomainRecord, error) {
	record := DomainRecord{}

	endpoint := fmt.Sprintf("domains/%d/records", id)
	err := c.invokeAPI("POST", endpoint, &req, &record)

	return record, err
}

// DELETE
func (c *Client) DomainRecordDelete(domainId int, recordId int) error {
	endpoint := fmt.Sprintf("domains/%d/records/%d", domainId, recordId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// UPDATE
func (c *Client) DomainRecordUpdate(domainId int, recordId int, req DomainRecordRequest) error {
	endpoint := fmt.Sprintf("domains/%d/records/%d", domainId, recordId)
	err := c.invokeAPI("PUT", endpoint, &req, nil)

	return err
}

// --------------------------------------------------- ACCESS --------------------------------------------------------
//add access to a domain

func (c *Client) DomainAccesAdd(domainId int, req DomainAccessRequest) error {
	endpoint := fmt.Sprintf("domains/%v/acls", domainId)
	err := c.invokeAPI("POST", endpoint, &req, nil)

	return err
}

//remove acces from a domain

func (c *Client) DomainAccesRemove(domainId int, organisationId int) error {
	endpoint := fmt.Sprintf("domains/%v/acls/%v", domainId, organisationId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// --------------------------------------------------- NOTIFICATIONS --------------------------------------------------------
// GET LIST OF ALL NOTIFICATIONS FOR DOMAIN
// func (c *Client) DomainNotificationGet(domainId int) []Notification {
// 	var notifications struct {
// 		Notifications []Notification `json:"notifications"`
// 	}
// 	endpoint := fmt.Sprintf("domains/%v/notifications", domainId)
// 	err := c.invokeAPI("GET", endpoint, nil, &notifications)
// 	AssertApiError(err, "notifications")
// 	return notifications.Notifications
// }

// // CREATE A NOTIFICATION
// func (c *Client) DomainNotificationAdd(domainId int, req DomainNotificationPostRequest) {
// 	endpoint := fmt.Sprintf("domains/%v/notifications", domainId)
// 	err := c.invokeAPI("POST", endpoint, req, nil)

// 	AssertApiError(err, "notifications")
// }

// --------------------------------------------------- BILLABLE ITEM --------------------------------------------------------

//--------------------------- CREATE (Turn invoicing on)
//CREATE BILLABLEITEM
func (c *Client) DomainBillableItemCreate(domainid int, req BillPostRequest) error {
	endpoint := fmt.Sprintf("domains/%v/bill", domainid)
	err := c.invokeAPI("POST", endpoint, req, nil)

	return err
}

// ---------------------------- DELETE (turn invoicing off)
//DELETE
func (c *Client) DomainBillableItemDelete(domainId int) error {
	endpoint := fmt.Sprintf("domains/%v/billableitem", domainId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

// -------------------------------------------------------CHECK AVAILABILITY---------------------------------------------------------------------
// Check domain availability
func (c *Client) DomainCheck(name string, extension string) (DomainCheckResult, error) {
	var checkResult DomainCheckResult

	endpoint := fmt.Sprintf("domains/check?name=%s&extension=%s", url.QueryEscape(name), url.QueryEscape(extension))
	err := c.invokeAPI("GET", endpoint, nil, &checkResult)

	return checkResult, err
}

// ----------- DOMAIN CONTACTS -----------

// POST /domaincontacts
func (c *Client) DomainContactCreate(request DomainContactRequest) (DomainContact, error) {
	var response struct {
		Data DomainContact `json:"domaincontact"`
	}

	endpoint := "domaincontacts"
	err := c.invokeAPI("POST", endpoint, request, &response)

	return response.Data, err
}

// GET /domaincontacts/{domainContactID}
func (c *Client) DomainContactGetSingle(id int) (DomainContact, error) {
	var response struct {
		Data DomainContact `json:"domaincontact"`
	}

	endpoint := fmt.Sprintf("domaincontacts/%d", id)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Data, err
}

// PUT /domaincontacts/{domainContactID}
func (c *Client) DomainContactUpdate(id int, request DomainContactRequest) error {
	endpoint := fmt.Sprintf("domaincontacts/%d", id)
	err := c.invokeAPI("PUT", endpoint, request, nil)

	return err
}

// DELETE /domaincontacts/{domainContactID}
func (c *Client) DomainContactDelete(id int) error {
	endpoint := fmt.Sprintf("domaincontacts/%d", id)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	return err
}

type Domain struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Fullname              string `json:"fullname"`
	TTL                   int    `json:"ttl"`
	EppCode               string `json:"eppCode"`
	Status                string `json:"status"`
	DnssecStatus          string `json:"dnssecStatus"`
	RegistrationIsHandled bool   `json:"registrationIsHandled"`
	Provider              struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		API  string `json:"api"`
	} `json:"provider"`
	DNSIsHandled    bool    `json:"dnsIsHandled"`
	DtRegister      string  `json:"dtRegister"`
	Nameserver1     *string `json:"nameserver1"`
	Nameserver2     *string `json:"nameserver2"`
	Nameserver3     *string `json:"nameserver3"`
	Nameserver4     *string `json:"nameserver4"`
	NameserverIP1   *string `json:"nameserverIp1"`
	NameserverIP2   *string `json:"nameserverIp2"`
	NameserverIP3   *string `json:"nameserverIp3"`
	NameserverIP4   *string `json:"nameserverIp4"`
	NameserverIpv61 *string `json:"nameserverIpv61"`
	NameserverIpv62 *string `json:"nameserverIpv62"`
	NameserverIpv63 *string `json:"nameserverIpv63"`
	NameserverIpv64 *string `json:"nameserverIpv64"`
	Organisation    struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Reseller string `json:"reseller"`
	} `json:"organisation"`
	Domaintype struct {
		ID                                  int    `json:"id"`
		Name                                string `json:"name"`
		Extension                           string `json:"extension"`
		RenewPeriod                         int    `json:"renewPeriod"`
		TransferAutoLicensee                bool   `json:"transferAutoLicensee"`
		RequestIncomingTransferCodePossible bool   `json:"requestIncomingTransferCodePossible"`
		RequestOutgoingTransferCodePossible bool   `json:"requestOutgoingTransferCodePossible"`
		LicenseeChangePossible              bool   `json:"licenseeChangePossible"`
		DnssecSupported                     bool   `json:"dnssecSupported"`
	} `json:"domaintype"`
	DomaincontactLicensee DomainContactRef  `json:"domaincontactLicensee"`
	DomaincontactOnsite   *DomainContactRef `json:"domaincontactOnsite"`
	Mailgroup             struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"mailgroup"`
	ExtraFields   interface{} `json:"extraFields"`
	HandleMailDNS bool        `json:"handleMailDns"`
	DtExpires     int         `json:"dtExpires"`
	BillingStatus string      `json:"billingStatus"`
	ExternalInfo  string      `json:"externalInfo"`
	Teams         []struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		AdminOnly      bool   `json:"adminOnly"`
		OrganisationId int    `json:"organisationId"`
	} `json:"teams"`
	CountTeams int   `json:"countTeams"`
	Jobs       []Job `json:"jobs"`
}

type DomainContactRef struct {
	ID               int    `json:"id,omitempty"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Fullname         string `json:"fullname"`
	OrganisationName string `json:"organisationName"`
	Street           string `json:"street"`
	HouseNumber      string `json:"houseNumber"`
	Zip              string `json:"zip"`
	City             string `json:"city"`
	State            string `json:"state"`
	Phone            string `json:"phone"`
	Fax              string `json:"fax"`
	Email            string `json:"email"`
	TaxNumber        string `json:"taxNumber"`
	Status           int    `json:"status"`
	PassportNumber   string `json:"passportNumber"`
	SocialNumber     string `json:"socialNumber"`
	BirthStreet      string `json:"birthStreet"`
	BirthZip         string `json:"birthZip"`
	BirthCity        string `json:"birthCity"`
	BirthDate        string `json:"birthDate"`
	Gender           string `json:"gender"`
	Type             string `json:"type"`
	Country          struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"country"`
}

// DomainProvider represents a single DomainProvider
type DomainProvider struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	API             string `json:"api"`
	DNSSecSupported bool   `json:"dnsSecSupported"`
	Domaintypes     []struct {
		ID        int    `json:"id"`
		Extension string `json:"extension"`
	} `json:"domaintypes"`
}

// DomainExtension represents a single DomainExtension
type DomainExtension struct {
	ID        int
	Extension string
}

// DomainRequest represents a single DomainRequest
type DomainRequest struct {
	Name                      string  `json:"name"`
	NameServer1               *string `json:"nameserver1"`
	NameServer2               *string `json:"nameserver2"`
	NameServer3               *string `json:"nameserver3"`
	NameServer4               *string `json:"nameserver4"`
	NameServer1Ip             *string `json:"nameserverIp1"`
	NameServer2Ip             *string `json:"nameserverIp2"`
	NameServer3Ip             *string `json:"nameserverIp3"`
	NameServer4Ip             *string `json:"nameserverIp4"`
	NameServer1Ipv6           *string `json:"nameserverIpv61"`
	NameServer2Ipv6           *string `json:"nameserverIpv62"`
	NameServer3Ipv6           *string `json:"nameserverIpv63"`
	NameServer4Ipv6           *string `json:"nameserverIpv64"`
	TTL                       int     `json:"ttl"`
	Action                    string  `json:"action"`
	EppCode                   string  `json:"eppCode"`
	Handledns                 bool    `json:"handleDns"`
	ExtraFields               string  `json:"extraFields"`
	Domaintype                int     `json:"domaintype"`
	Domaincontactlicensee     int     `json:"domaincontactLicensee"`
	DomainContactOnSite       *int    `json:"domaincontactOnsite"`
	Organisation              int     `json:"organisation"`
	AutoRecordTemplate        string  `json:"autorecordTemplate"`
	AutoRecordTemplateReplace bool    `json:"autorecordTemplateReplace"`
	//DomainProvider            *int    `json:"domainProvider"`
	// DtExternalCreated         string `json:"dtExternalCreated"`
	// DtExternalExpires         string `json:"dtExternalExpires"`
	// ConvertDomainRecords      string `json:"convertDomainrecords"`
	AutoTeams    string  `json:"autoTeams"`
	ExternalInfo *string `json:"externalInfo,omitempty"`
}

// request for updating a single domain
type DomainUpdateRequest struct {
	Name                      string  `json:"name"`
	NameServer1               *string `json:"nameserver1"`
	NameServer2               string  `json:"nameserver2"`
	NameServer3               string  `json:"nameserver3"`
	NameServer4               string  `json:"nameserver4"`
	NameServer1Ip             string  `json:"nameserverIp1"`
	NameServer2Ip             string  `json:"nameserverIp2"`
	NameServer3Ip             string  `json:"nameserverIp3"`
	NameServer4Ip             string  `json:"nameserverIp4"`
	NameServer1Ipv6           string  `json:"nameserverIpv61"`
	NameServer2Ipv6           string  `json:"nameserverIpv62"`
	NameServer3Ipv6           string  `json:"nameserverIpv63"`
	NameServer4Ipv6           string  `json:"nameserverIpv64"`
	TTL                       int     `json:"ttl"`
	Action                    string  `json:"action"`
	EppCode                   string  `json:"eppCode"`
	Handledns                 bool    `json:"handleDns"`
	ExtraFields               string  `json:"extraFields"`
	Domaintype                int     `json:"domaintype"`
	Domaincontactlicensee     int     `json:"domaincontactLicensee"`
	DomainContactOnSite       *int    `json:"domaincontactOnsite"`
	Organisation              int     `json:"organisation"`
	AutoRecordTemplate        string  `json:"autorecordTemplate"`
	AutoRecordTemplateReplace bool    `json:"autorecordTemplateReplace"`
	AutoTeams                 string  `json:"autoTeams"`
}

func (d DomainRequest) String() string {

	s, _ := json.Marshal(d)
	return string(s)
}

// ------------------------------------------ RECORDS ---------------------------------------

// DomainRecord represents a single Domainrecord
type DomainRecord struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Content            string `json:"content"`
	Priority           int    `json:"priority"`
	Type               string `json:"type"`
	SystemHasNetworkIP struct {
		ID int `json:"id"`
	} `json:"systemHasNetworkIp"`
	// URL            int `json:"url"`
	// SslCertificate int `json:"sslCertificate"`
	// Mailgroup      int `json:"mailgroup"`
}

// DomainRecordRequest represents a API reqest to Level27
type DomainRecordRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	Priority int    `json:"priority"`
}

// DomainContact is an object to define domain contacts at Level27
type DomainContact struct {
	ID               int     `json:"id"`
	FirstName        string  `json:"firstName"`
	LastName         string  `json:"lastName"`
	OrganisationName string  `json:"organisationName"`
	Street           string  `json:"street"`
	HouseNumber      string  `json:"houseNumber"`
	Zip              string  `json:"zip"`
	City             string  `json:"city"`
	State            *string `json:"state"`
	Phone            string  `json:"phone"`
	Fax              *string `json:"fax"`
	Email            string  `json:"email"`
	TaxNumber        string  `json:"taxNumber"`
	PassportNumber   *string `json:"passportNumber"`
	SocialNumber     *string `json:"socialNumber"`
	BirthStreet      *string `json:"birthStreet"`
	BirthZip         *string `json:"birthZip"`
	BirthCity        *string `json:"birthCity"`
	BirthDate        *string `json:"birthDate"`
	Gender           *string `json:"gender"`
	Type             string  `json:"type"`
	Country          struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"country"`
	Organisation struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
	Fullname string `json:"fullname"`
}

// DomainContactRequest is an object to define the request to create or modify a domain contact at Level27
type DomainContactRequest struct {
	Type             string  `json:"type"`
	FirstName        string  `json:"firstName"`
	LastName         string  `json:"lastName"`
	OrganisationName string  `json:"organisationName"`
	Street           string  `json:"street"`
	HouseNumber      string  `json:"houseNumber"`
	Zip              string  `json:"zip"`
	City             string  `json:"city"`
	State            *string `json:"state,omitempty"`
	Phone            string  `json:"phone"`
	Fax              *string `json:"fax,omitempty"`
	Email            string  `json:"email"`
	TaxNumber        string  `json:"taxNumber"`
	PassportNumber   *string `json:"passportNumber,omitempty"`
	SocialNumber     *string `json:"socialNumber,omitempty"`
	BirthStreet      *string `json:"birthStreet,omitempty"`
	BirthZip         *string `json:"birthZip,omitempty"`
	BirthCity        *string `json:"birthCity,omitempty"`
	BirthDate        *string `json:"birthDate,omitempty"`
	Gender           *string `json:"gender,omitempty"`
	Country          string  `json:"country"`
	Organisation     string  `json:"organisation"`
}

// ------------------------------------------ ACCESS ---------------------------------------------

// type to add acces to a domain
type DomainAccessRequest struct {
	Organisation int `json:"organisation"`
}

// ------------------------------------------ NOTIFICATIONS ---------------------------------------------
type DomainNotificationPostRequest struct {
	Type   string `json:"type"`
	Group  string `json:"group"`
	Params string `json:"params"`
}

// ------------------------------------------ CHECK/AVAILABILITY ---------------------------------------------

// Domain check
type DomainCheckResult struct {
	Success                             bool   `json:"success"`
	Status                              string `json:"status"`
	Action                              string `json:"action"`
	DomaintypeId                        int    `json:"domaintypeId"`
	DomainNameWithExtension             string `json:"domainNameWithExtension"`
	RequestIncomingTransferCodePossible bool   `json:"requestIncomingTransferCodePossible"`
	TransferAutoLicensee                bool   `json:"transferAutoLicensee"`
	TransferEppCodeRequired             bool   `json:"transferEppCodeRequired"`
	Products                            []struct {
		Id          string `json:"id"`
		Description string `json:"description"`
		Prices      []struct {
			Id       int    `json:"id"`
			Period   int    `json:"period"`
			Currency string `json:"currency"`
			Price    string `json:"price"`
			Timing   string `json:"timing"`
			Default  bool   `json:"default"`
			Status   int    `json:"sttaus"`
		} `json:"prices"`
	} `json:"products"`
}

// ------------------------------------------ JOB HISTORY ---------------------------------------------
type DomainJobHistory struct {
	Id      int           `json:"id"`
	Status  int           `json:"status"`
	Conc    int           `json:"conc"`
	Hoe     int           `json:"hoe"`
	Message string        `json:"msg"`
	Dt      string        `json:"dt"`
	Logs    []interface{} `json:"logs"`
}

type DomainJobHistoryRoot struct {
	DomainJobHistory
}

// INTEGRITY CHECKS
type DomainIntegrityCheck struct {
	IntegrityCheck
	Object   string `json:"object"`
	ObjectID int    `json:"objectId"`
	Results  struct {
		Domain struct {
			ID   int `json:"id"`
			Data struct {
				FullName          string      `json:"fullName"`
				Status            string      `json:"status"`
				StatusColor       string      `json:"statusColor"`
				Provider          string      `json:"provider"`
				HandleDNS         bool        `json:"handleDns"`
				Nameserver1       interface{} `json:"nameserver1"`
				Nameserver2       interface{} `json:"nameserver2"`
				Nameserver3       interface{} `json:"nameserver3"`
				Nameserver4       interface{} `json:"nameserver4"`
				NameserverIP1     interface{} `json:"nameserverIp1"`
				NameserverIP2     interface{} `json:"nameserverIp2"`
				NameserverIP3     interface{} `json:"nameserverIp3"`
				NameserverIP4     interface{} `json:"nameserverIp4"`
				NameserverIpv61   interface{} `json:"nameserverIpv61"`
				NameserverIpv62   interface{} `json:"nameserverIpv62"`
				NameserverIpv63   interface{} `json:"nameserverIpv63"`
				NameserverIpv64   interface{} `json:"nameserverIpv64"`
				IsRegistered      string      `json:"isRegistered"`
				IsRegisteredColor string      `json:"isRegisteredColor"`
				DcLicensee        struct {
					FullName string `json:"fullName"`
					Address  string `json:"address"`
					Country  string `json:"country"`
					Phone    string `json:"phone"`
					Email    string `json:"email"`
				} `json:"dcLicensee"`
				DcOnsite          interface{} `json:"dcOnsite"`
				DnssecStatus      string      `json:"dnssecStatus"`
				DnssecStatusColor string      `json:"dnssecStatusColor"`
				Retry             int         `json:"retry"`
				Refresh           int         `json:"refresh"`
				Expire            int         `json:"expire"`
				Minimum           int         `json:"minimum"`
				TTL               int         `json:"ttl"`
				Mailgroup         string      `json:"mailgroup"`
			} `json:"data"`
			Results struct {
				Name              string `json:"name"`
				FullName          string `json:"fullName"`
				Subdomain         string `json:"subdomain"`
				NameseversGeneral string `json:"namesevers_general"`
				Nameserver1       string `json:"nameserver1"`
				Nameserver2       string `json:"nameserver2"`
				Nameserver3       string `json:"nameserver3"`
				Nameserver4       string `json:"nameserver4"`
				NameserverIP1     string `json:"nameserverIp1"`
				NameserverIP2     string `json:"nameserverIp2"`
				NameserverIP3     string `json:"nameserverIp3"`
				NameserverIP4     string `json:"nameserverIp4"`
				NameserverIpv61   string `json:"nameserverIpv61"`
				NameserverIpv62   string `json:"nameserverIpv62"`
				NameserverIpv63   string `json:"nameserverIpv63"`
				NameserverIpv64   string `json:"nameserverIpv64"`
				Action            string `json:"action"`
				EppCode           string `json:"eppCode"`
				ExtraFields       string `json:"extraFields"`
				Domaintype        string `json:"domaintype"`
				Status            string `json:"status"`
				DnssecStatus      string `json:"dnssecStatus"`
				Retry             string `json:"retry"`
				Refresh           string `json:"refresh"`
				Expire            string `json:"expire"`
				Minimum           string `json:"minimum"`
				TTL               string `json:"ttl"`
				Organisation      string `json:"organisation"`
				Provider          string `json:"provider"`
				NsDig             string `json:"nsDig"`
			} `json:"results"`
			Records struct {
				Message string `json:"message"`
				Dig     []struct {
					ID   int `json:"id"`
					Data struct {
						Type       string      `json:"type"`
						Name       interface{} `json:"name"`
						Content    string      `json:"content"`
						Priority   interface{} `json:"priority"`
						DigContent string      `json:"digContent"`
					} `json:"data"`
					Results struct {
						Status    string `json:"status"`
						IsExpired string `json:"isExpired"`
					} `json:"results"`
					SslCertificate interface{} `json:"sslCertificate"`
				} `json:"dig"`
			} `json:"records"`
		} `json:"domain"`
		Organisations struct {
			Manual []string `json:"manual"`
			Owner  string   `json:"owner"`
		} `json:"organisations"`
		Teams []struct {
			ID             int         `json:"id"`
			Name           string      `json:"name"`
			AdminOnly      interface{} `json:"adminOnly"`
			OrganisationID int         `json:"organisationId"`
		} `json:"teams"`
		FinanceData struct {
			Billing         string `json:"billing"`
			BillingStatus   string `json:"billingStatus"`
			AutoRenew       string `json:"autoRenew"`
			AutoRenewStatus string `json:"autoRenewStatus"`
			DtExpires       string `json:"dtExpires"`
			BillingItems    []struct {
				ID          int     `json:"id"`
				Description string  `json:"description"`
				Type        string  `json:"type"`
				Period      string  `json:"period"`
				DtExpires   string  `json:"dtExpires"`
				Price       float32 `json:"price"`
			} `json:"billingItems"`
			TotalPrice float32 `json:"totalPrice"`
		} `json:"financeData"`
		Jobs []struct {
			ID           int         `json:"id"`
			Action       string      `json:"action"`
			Status       int         `json:"status"`
			Message      string      `json:"message"`
			DtEnd        interface{} `json:"dtEnd"`
			DtStamp      string      `json:"dtStamp"`
			ExceptionMsq interface{} `json:"exceptionMsq"`
		} `json:"jobs"`
		IsHealthy    bool   `json:"isHealthy"`
		ExtraMessage string `json:"extraMessage"`
	} `json:"results"`
	Dojobs               bool          `json:"dojobs"`
	Forcejobs            bool          `json:"forcejobs"`
	LocalIntegritychecks []interface{} `json:"localIntegritychecks"`
}
