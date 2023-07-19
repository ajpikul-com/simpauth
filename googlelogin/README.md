# googelogin

This module allows you to use google login. It requires some set-up on google's end.

## Google's End

You have to get a client ID. You can google "Google Sign In" and set up your google account to allow your application to use google sign in. It will give you a client ID. It will also give you many other user-flows besides the default one implemented here.

## Using

`googlelogin.New(clientId string)` will return a new google login identifier to be attached with `uwho`'s addIdentifier function.


This interface: 

```go
type ReqByIdent interface {
	AcceptData(map[string]interface{}) bool
}
```

...must be satisfied by your state object. The `map` is whatever `pkg.go.dev/google.golang.org/api/idtoken#Payload` gives you, even if it's not explicitly documented, it definitely includes an `email` key.

Note:

`googlelogin` uses a google-suggested login interface to retrieve an `idtoken` from google, and then it uses google api to decode `idtoken` into `Claims` which has real user info. There's multiple ways to get that `idtoken`, but we only demonstrate one (see #Extras, DefaultLoginPortal). But it's not dynamic, the user fills out a form and it redirects them to a login endpoint.


TODO: check csfr token

## Extras

`RedirectHome` `DefaultLoginResult` and `GoogleLogin.DefaultLoginPortal` are three `http.Handler`s that can get us started and provide examples on how to handle login forms, login results, and logout results.
