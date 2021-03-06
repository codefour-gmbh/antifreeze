package apifakes

import "github.com/cloudfoundry/cli/cf/models"

type OldFakeAuthTokenRepo struct {
	CreatedServiceAuthTokenFields models.ServiceAuthTokenFields

	FindAllAuthTokens []models.ServiceAuthTokenFields

	FindByLabelAndProviderLabel                  string
	FindByLabelAndProviderProvider               string
	FindByLabelAndProviderServiceAuthTokenFields models.ServiceAuthTokenFields
	FindByLabelAndProviderApiResponse            error

	UpdatedServiceAuthTokenFields models.ServiceAuthTokenFields

	DeletedServiceAuthTokenFields models.ServiceAuthTokenFields
}

func (repo *OldFakeAuthTokenRepo) Create(authToken models.ServiceAuthTokenFields) (apiErr error) {
	repo.CreatedServiceAuthTokenFields = authToken
	return
}

func (repo *OldFakeAuthTokenRepo) FindAll() (authTokens []models.ServiceAuthTokenFields, apiErr error) {
	authTokens = repo.FindAllAuthTokens
	return
}
func (repo *OldFakeAuthTokenRepo) FindByLabelAndProvider(label, provider string) (authToken models.ServiceAuthTokenFields, apiErr error) {
	repo.FindByLabelAndProviderLabel = label
	repo.FindByLabelAndProviderProvider = provider

	authToken = repo.FindByLabelAndProviderServiceAuthTokenFields
	apiErr = repo.FindByLabelAndProviderApiResponse
	return
}

func (repo *OldFakeAuthTokenRepo) Delete(authToken models.ServiceAuthTokenFields) (apiErr error) {
	repo.DeletedServiceAuthTokenFields = authToken
	return
}

func (repo *OldFakeAuthTokenRepo) Update(authToken models.ServiceAuthTokenFields) (apiErr error) {
	repo.UpdatedServiceAuthTokenFields = authToken
	return
}
