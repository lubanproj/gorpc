package auth

import (
	"context"

	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/interceptor"
	"golang.org/x/oauth2"
)

type oAuth2 struct {
	token *oauth2.Token
}

func (o *oAuth2) AuthType() string {
	return "oauth2"
}

// NewOAuth2ByToken supports the generation of an oauth2 based on a string token
func NewOAuth2ByToken(token string) *oAuth2 {
	return &oAuth2{
		token : &oauth2.Token{
			AccessToken: token,
		},
	}
}

// NewOAuth2 supports the generation of an oauth2 based on an oauth2 token
func NewOAuth2(t *oauth2.Token) *oAuth2 {
	return &oAuth2{
		token : t,
	}
}

func (o *oAuth2) GetMetadata(ctx context.Context, uri ... string) (map[string]string, error) {

	if o.token == nil {
		return nil, codes.ClientCertFailError
	}

	return map[string]string{
		"authorization": o.token.Type() + " " + o.token.AccessToken,
	}, nil
}

// AuthFunc verifies that the token is valid or not
type AuthFunc func(ctx context.Context) (context.Context, error)


// BuildAuthFilter constructs a client interceptor with an AuthFunc
func ClientInterceptor(af AuthFunc) interceptor.ClientInterceptor {
	return func(ctx context.Context, req, rsp interface{}, ivk interceptor.Invoker) error {

		newCtx, err := af(ctx)
		if err != nil {
			return codes.NewFrameworkError(codes.ClientCertFail, "token invalid ...")
		}

		return ivk(newCtx, req, rsp)
	}
}
