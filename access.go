package l27

import (
	"fmt"
)

// GET /{entityType}/{entityID}/organisations
func (c *Client) EntityGetOrganisations(entityType string, entityID int) ([]OrganisationAccess, error) {
	var response struct {
		Organisations []OrganisationAccess `json:"organisations"`
	}

	endpoint := fmt.Sprintf("%s/%d/organisations", entityType, entityID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	return response.Organisations, err
}

// POST /{entityType}/{entityID}/acls
func (c *Client) EntityAddAcl(entityType string, entityID int, add AclAdd) (Acl, error) {
	var response struct {
		Acl Acl `json:"acl"`
	}

	endpoint := fmt.Sprintf("%s/%d/acls", entityType, entityID)
	err := c.invokeAPI("POST", endpoint, add, &response)

	return response.Acl, err
}

// DELETE /{entityType}/{entityID}/acls/{organisationID}
func (c *Client) EntityRemoveAcl(entityType string, entityID int, organisationID int) error {
	endpoint := fmt.Sprintf("%s/%d/acls/%d", entityType, entityID, organisationID)
	return c.invokeAPI("DELETE", endpoint, nil, nil)
}

type AclAdd struct {
	Organisation int `json:"organisation"`
}

type Acl struct {
	ID           int             `json:"id"`
	Object       string          `json:"object"`
	ObjectID     int             `json:"objectId"`
	Permissions  interface{}     `json:"permissions"`
	Extra        interface{}     `json:"extra"`
	Type         string          `json:"type"`
	Organisation OrganisationRef `json:"organisation"`
}
