package client

// GrantType is the type of authorization being used by the client to obtain
// an access token.
const (
	// GrantTypeAuthorizationCode is used by confidential and public clients
	// to exchange an authorization code for an access token.
	GrantTypeAuthorizationCode = "authorization_code"

	// GrantTypeImplicit is a simplified flow that can be used by public
	// clients, where the access token is returned immediately without an
	// extra authorization code exchange step.
	GrantTypeImplicit = "implicit"

	// GrantTypePassword is used by first-party clients to exchange a user's
	// credentials for an access token.
	GrantTypePassword = "password"

	// GrantTypeClientCredentials is used by clients to obtain an access
	// token outside of the context of a user.
	GrantTypeClientCredentials = "client_credentials"

	// GrantTypeRefreshToken is used by clients to exchange a refresh token
	// for an access token when the access token has expired.
	GrantTypeRefreshToken = "refresh_token"
)

// ResponseType is the type of response being used by the authorization code
// grant type and implicit grant type flows.
const (
	// ResponseTypeCode is used for requesting an authorization code.
	ResponseTypeCode = "code"
	// ResponseTypeToken is used for requesting an access token.
	ResponseTypeToken = "token"
)

// Client represents an OAuth 2.0 client.
type Client struct {
	// ID is the identifier for this client.
	ID string `json:"id"`

	// Name is the human-readable string name of the client to be presented to the
	// end-user during authorization.
	Name string `json:"name"`

	// Secret is the client's secret.
	Secret string `json:"secret"`

	// RedirectURI is an allowed redirect url for the client.
	RedirectURI string `json:"redirectUri"`

	// GrantType is grant type the client is allowed to use.
	GrantType string `json:"grantType"`

	// ResponseType is the OAuth 2.0 response type string that the client can use at
	// the authorization endpoint.
	ResponseType string `json:"responseType"`

	// Scope is a string containing a space-separated list of scope values (as
	// described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client can use when
	// requesting access tokens.
	Scope string `json:"scope"`

	// Public is a boolean that identifies this client as public, meaning that it
	// does not have a secret. It will disable the client_credentials grant type for
	// this client if set.
	Public bool `json:"public"`
}
