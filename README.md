# DEVELOPERS ONLY 

But these are not developer docs, for writing modules. But still, developers only.

There are no tests.

## Basic Idea of uWHO

First, you choose what modules you want (Google Sign In, Email-Link, Encrypted-Cookies).
Second, you write a structure ("the state object") that accepts the information they give you about the user and helps make authorization decisions.
Third, you initialize your structure and module.

### Choose The Modules You Want

* github.com/ajpikul-com/uwho/googlelogin
* github.com/ajpikul-com/uwho/usersessioncookie

### You Write:

1) a state type and factory (`New()`) that produces the state objects:

```go
type myFactory struct{}

func (f *myFactory) New() uwho.ReqByCoord {
	return &myState{}
}

type myState struct { // You choose this structure's members.
	email       string
	session     uuid.NullUUID
	failedLogin bool
	failedAuth  bool
}
```

Your state objects are only alive during 1 connection/request.

2) Your `myState` struct must fulfil the `uwho.ReqByCoord` interface. `Coordinator`, the main module in *uWho* will use your `myFactory` to create a new state object every request. And it needs that the state objects fulfil (ie you write):

```go
type ReqByCoord interface {

	// We return true if the user is authorized for whatever they're asking for in the
	// http.Request. In my example, we would set myState.failedAuth to true if we deny
	// them here. This is only called if we feel the user is sufficiently "known".
	AuthorizeUser(w http.ResponseWriter, r *http.Request) bool 

    // InitState is called only right after login if it was an unknown user who called
	// the login function (supplied by a module). Already logged-in users can login again,
	// theoretically, if you provide them a way to access the endpoint. You can just ignore it.
	// Provides defaults for myState.
	InitState() error

    // Sets myState to 0, effectively, cancels everything.
	DeleteState()
}
```

The modules you use will also have interfaces that need to be fulfilled. Specified by the modules, but following good practice, it might look something like:

```go
type ReqByIdent interface{
	// Most identification modules will ask you to provide an AcceptData function,
	// It will be called by the module if a user successfully retrieves data from the modules.
	// The modules assures you that this data is legit to the user, store it how
	// you want.
	AcceptData(claims map[string]string) bool
}

type ReqBySess interface{
    // This would be called by the module so that you can move session variables to your state variable. In this case,
	// the session module is asking for your opinion about the legitimacy of the session (return true or false).
	SessionToState(sessString string) bool
	/// This would be called when the session module is trying to update whatever it has in the cookie/database.
	StateToSession() string
}
```

There are two types of modules (besides `Coordinator`): 

1) Credentialing (or Identifying) which could be Google Sign In, User/Pass, certs, or anything else a user uses to prove some information about them is legitimate (e.g. "this is my real email"). 
2) Session Management, which is how we attach credentials to a user over time/page clicks and navigating to different pages.

You generally don't call any of those required functions directly- but you can- coordinator has a bunch of hooks to hook into.

### Configuration 

*NOTE: This is for a non-dynamic login process, ie non webapp. It's possible, but I will write readme after doing it myself*

```go
import (
	"github.com/ajpikul-com/uwho"
	"github.com/ajpikul-com/uwho/googlelogin"
	"github.com/ajpikul-com/uwho/usersessioncookie"
)

func main() {
	cookieSession := usersessioncookie.New() // This has changed
	googleIdent   := googlelogin.New(clientID) // clientID is a string you get from your google cloud console
	loginScreen   := googleIdent.DefaultLoginPortal("/login") // googlelogin provides an httpHandler that displays a login form
	
	// Wrap your filesever (or any other http mux) in an `uwho.coordinator`, which forces auth before access
	fileServer := uwho.New( // new coordinator
		http.FileServer(http.Dir("/to/your/directory/to/server"), // this is what we're really serving
		&googlelogin.DefaultLoginResult{}, // send user where they originally wanted to go w/ a redirect
		loginScreen, // the default login screen we're given for Google Sign On
		&googlelogin.RedirectHome{}, // we logged out, so go home
		"/login", // endpoint that will treat user like they're trying to log in (ie look for credentials)
		"/logout", // endpoint that will log user out if they access in any way
		&myFactory) // how coordinator will create your state objects
		
	fileServer.AddIdentifier(googleIdent)
	fileServer.AttachSessionManager(cookieSession)
	
	serveMux.Handle("stage.ajpikul.com/" &fileServer) // a normal go serveMux to attach handlers to path
	
```

### Extra

Hooks!

There are four hooks available:

```
		LoggedOut
		LoggedIn
		Authorized
		AboutToLoad
```

They are part of `coordinator.Hooks` and must be of `Hook` type, see example to add them:

```go
	// We're going to create a hook. The first argument, `stateProvided`, is your `myState`, passed as the `ReqByCoord` interface.
	// We can use all its `ReqByCoord` methods, but we need to type assert it to the other interfaces which it
	// fulfills or to its bare type if we are to use it those ways.
	rightBeforeHook := func(stateProvided uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) error {
	
	// We're going to assert it as the native type we declared. This should _always work_, but theoretically
	// we could just try and handle any possible error.
	// NOTE: when you initialize, like above, all modules and coordinator test your state struct to see if it fulfils all interfaces
	// They will panic if they don't, so type-asserts after initialization should be fine
	if state, ok := stateProvided.(*stageState); ok {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		// Whatever page with these headers will not cache, good in someone situations (sensitive user data, fast-changing apps)
		return nil // no error
	}
	
	// Attach the hook
	fileServer.SetHooks(&fileServer.Hooks.AboutToLoad, []uwho.Hook{&rightBeforeHook})
}
```

// TODO: State Diagram Of Coordinator w/ hooks


