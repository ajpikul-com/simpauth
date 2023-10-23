# BETA IS OUT!

## Basic Idea of uWHO

uWHO will wrap your multiplexer and will not allow any request to pass unless it satisfies auth.

1. You write a state object- it contains whatever user information you want. It is alive for only the connection, but you will store it in a session (keep reading).
2. You write required methods for your state object. It must satisfy the interface for core module (`uwho`) and whatever modules you choose.
3. You choose login modules (google login, username password, email link, etc). You add it's required methods to your state object.
4. You choose session modules. Right now there is only one, it's a cookie manager, and it requires you to write two methods.
5. You have to write a function that returns a new state object.
6. You initialize the objects.

TODO: see example at ____

### Choosing Modules

```go
import (
	"github.com/ajpikul-com/uwho"
	"github.com/ajpikul-com/uwho/googlelogin"
	"github.com/ajpikul-com/uwho/usersessioncookie"
)
```

### Writing Required Methods

1. Write your state object:

```go
type myState struct {
    Email string
    IsLoggedIn bool
    // Whatever else
}
```

2. 'uwho' requires the following methods on your state object:

```go
	// LogOut will empty out user state object, that's all. It should _not_ write a response body.
	LogOut(w http.ResponseWriter, r *http.Request)
	// IsLoginAllowed tells uwho to skip login process if, for example, we're already logged in.
	// hint: you can also use it to hint that user is trying to login. It should _not_ write a response body.
	IsLoginAllowed(w http.ResponseWriter, r *http.Request) bool
	// OtherStateAction is a hook that will be called after every other source of user information has been requested. It should _not_ write a response body.
	OtherStateAction(w http.ResponseWriter, r *http.Request)
	// ChangeState will be called if the user has tried to login (with success or not), or loggedout. It must write a response body. `uwho` provides some obvious utility functions (see README.md or utilities.go) that you can use.
	ChangeState(w http.ResponseWriter, r *http.Request)
	// IsUserAuthorized will be called after session is read, the user did not login or logout, and session is updated. It's your job to check the request and see if user is authorized. If true, user will continue to the wrapped handler, Coordinator.DesiredResource. If false, you must write a response body. Maybe redirect user to a login page?
	IsUserAuthorized(w http.ResponseWriter, r *http.Request) bool
```

3. You must write a type that has a method `func (myType) New() uwho.ReqByCoord` <-- this `New()` will return a new instance of your state object, and `uwho.ReqByCoord` is the interface you are satisfying by writing the methods above.

4. `usersessioncookie` and each login module also requires you to write methods, but those methods are described in their READMEs (or see the example).

### Initialization

Here is an example:
```go
    // Set up the actual handles you want to serve
	fileserver := http.NewServeMux()
    fileserver.Handle("www.example.com", http.FileServer(http.Dir("/www/example.com/")))

    // Initialize uwho
	authMux := uwho.New("/login", "/logout", &myStateObjectFactory{}) // (loginPath, logoutPath, type according to #3 above)
	cookieSessions := usersessioncookie.New("", "/", "/some/path/to/a/private/key/)
	cookieSessions.SetID("myID") // If you want to name your cookie, otherwise it's a random UUID
	googleIdent := googlelogin.New("googleID")
    
    // Attach uwho to it's modules
	authMux.AddIdentifier(googleIdent)
	authMux.AttachSessionManager(cookieSessions)
    
    // Wrap the real handler in uwho
	authMux.DesiredResource = fileserver
    
    // Serve uwho
    serverHTTP := &http.Server{
        Addr:    ":http",
        Handler: authMux,
    }

    err := serverHTTP.ListenAndServe()
```

### Utilities

```go
var myReferrer uwho.ToReferrer = "/default/path/"
```

`uwho.ToReferrer` types have a `ServeHTTP` method which redirects user to where they wanted to go- unless it's the same as where they are, in which case it goes to the default path to avoid an infinite loop. It's good for after a login.

```go
var redirectHome uwho.RedirectHome = "/path/to/home"
```

`uwho.RedirectHome` types have a `ServeHTTP` method which redirects user to the path specified, great for after a logout.
