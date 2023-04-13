package provider

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/slices"
)

func projectWait(client *http.Client, projectId string) error {
	var operations OperationListOutput

	for {
		err := get(client, fmt.Sprintf("/projects/%s/operations?limit=1", projectId), &operations)

		if err != nil {
			return err
		}

		if operations.Operations[0].Status == "finished" {
			return nil
		}

		time.Sleep(5 * time.Second)
	}
}

func branchList(client *http.Client, projectId string) (BranchListOutput, error) {
	var branches BranchListOutput

	err := get(client, fmt.Sprintf("/projects/%s/branches", projectId), &branches)

	return branches, err
}

func branchEndpoint(client *http.Client, projectId string, branchId string) (Endpoint, error) {
	endpoints, err := endpointList(client, projectId)

	var endpoint Endpoint

	if err != nil {
		return endpoint, err
	}

	endpointIdx := slices.IndexFunc(endpoints.Endpoints, func(endpoint Endpoint) bool {
		return endpoint.BranchId == branchId
	})

	return endpoints.Endpoints[endpointIdx], nil
}

func branchGet(client *http.Client, projectId string, branchId string) (BranchOutput, error) {
	var branch BranchOutput

	err := get(client, fmt.Sprintf("/projects/%s/branches/%s", projectId, branchId), &branch)

	if err != nil {
		return branch, err
	}

	if branch.Branch.ProjectId != projectId {
		return branch, fmt.Errorf("branch %s does not belong to project %s", branchId, projectId)
	}

	return branch, nil
}

func branchUpdate(client *http.Client, projectId string, branchId string, input BranchUpdateInput) (BranchOutput, error) {
	var branch BranchOutput

	err := call(client, http.MethodPatch, fmt.Sprintf("/projects/%s/branches/%s", projectId, branchId), input, &branch)

	return branch, err
}

func endpointList(client *http.Client, projectId string) (EndpointListOutput, error) {
	var endpoints EndpointListOutput

	err := get(client, fmt.Sprintf("/projects/%s/endpoints", projectId), &endpoints)

	return endpoints, err
}

func endpointUpdate(client *http.Client, projectId string, endpointId string, input EndpointUpdateInput) (EndpointOutput, error) {
	var endpoint EndpointOutput

	err := projectWait(client, projectId)

	if err != nil {
		return endpoint, err
	}

	err = call(client, http.MethodPatch, fmt.Sprintf("/projects/%s/endpoints/%s", projectId, endpointId), input, &endpoint)

	return endpoint, err
}

func databaseCreate(client *http.Client, projectId string, branchId string, input DatabaseCreateInput) (DatabaseOutput, error) {
	var database DatabaseOutput

	err := projectWait(client, projectId)

	if err != nil {
		return database, err
	}

	err = call(client, http.MethodPost, fmt.Sprintf("/projects/%s/branches/%s/databases", projectId, branchId), input, &database)

	return database, err
}

func databaseUpdate(client *http.Client, projectId string, branchId string, name string, input DatabaseUpdateInput) (DatabaseOutput, error) {
	var database DatabaseOutput

	err := projectWait(client, projectId)

	if err != nil {
		return database, err
	}

	err = call(client, http.MethodPatch, fmt.Sprintf("/projects/%s/branches/%s/databases/%s", projectId, branchId, name), input, &database)

	return database, err
}

func databaseDelete(client *http.Client, projectId string, branchId string, name string) error {
	err := projectWait(client, projectId)

	if err != nil {
		return err
	}

	_, err = delete(client, fmt.Sprintf("/projects/%s/branches/%s/databases/%s", projectId, branchId, name))

	return err
}

func roleCreate(client *http.Client, projectId string, branchId string, input RoleCreateInput) (RoleOutput, error) {
	var role RoleOutput

	err := projectWait(client, projectId)

	if err != nil {
		return role, err
	}

	err = call(client, http.MethodPost, fmt.Sprintf("/projects/%s/branches/%s/roles", projectId, branchId), input, &role)

	return role, err
}

func roleDelete(client *http.Client, projectId string, branchId string, name string) error {
	err := projectWait(client, projectId)

	if err != nil {
		return err
	}

	_, err = delete(client, fmt.Sprintf("/projects/%s/branches/%s/roles/%s", projectId, branchId, name))

	return err
}
