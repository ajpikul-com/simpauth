# SimpAuth

SimpAuth is a modular login/session management framework with four types of modules (SimpAuth does not worry about how you register your users or store profiles):

The `AuthCoordinator`, the `LoginIdentifier`, the `ServerSessionManager`, and the `ClientSessionManager`. They all expose functions and variables, some of which are relevant for the website developer (you) and some of which are only revelant for a module developer (maybe you?). 

The `ClientSessionManager` is the API for writing and reading data to the user (like a session token). We only implement one right now: encrypted cookies. If you want to use a different one, you can see *the guide here* for writing one, the basic structure is easy, can't speak to your backend.

The `ServerSessionManager` is the API for storing sessions server side. We only implement a regular server-side in-memory dict. If you want to implement something that uses a database or other method, see *the guide here*, the basic structure is easy, can't speak to your backend.

The `LoginIdentifier` is what checks to see if a login should be successful, and returns any user data. We implement three:
* Google Log In (multiple ways to implement this, but right now only use embedded html)
* Email a Link Log In
* Expiring Link Share Log In

If you want to write your own, such as Username And Password, or GitHub Auth, please see *the guide here*.

The `AuthCoordinator` is what connects all of them and helps define user flows, so it's the most complicated for the user. You must configure it as such: 
LoginPortal // this should probably be called access denied.
LoginEndpoint
LogOutEndpoint
LogOutResult
AddServerSessionManager
AddClientSessionManager
AddLoginIdentifier
ResultHooks // so this will be a dictionary of possible login/logout/authentification errors and 

We implement one: Bouncer. Bouncer simply sits on a path and denies access to the path by showing the LoginPortal. If the login is sucessful, it redirects back to the page. It provides a LogOutEndpoint which if accessed will just logout and then move thte user to LogOutResult. Bouncer's LoginEndpoint can be accessed by ajax if you don't wat to use the LoginPortal.
 
SimpAuth is an modular login/auth framework.

It contains `bouncer`, which simply sits on a path and presents a login screen if the user tries to access something they can't.

First you must create a Bouncer:

myBouncer = Bouncer{
	// A file server or something to show the user when they want to go where they can't
	AccessDenied http.Handler
  // Like above
	LogoutResult http.Handler
  // The path where all credentials are POSTed to (or w/e)
	LoginEndpoint string
	// The path that will cause a logout by accessing
	LogoutEndpoint string
  // The http.Handler you are denying access to
  Base  http.Handler
}

We are using GoogleLogIn as a login identifier. You must set up everything with them too. Your LoginEndpoint must match their `data-login_uri`.
myGoogleLogin = GoogleLogIn{
	// The ClientID they supply to you
	ClientID
}
myBouncer.AddLoginIdentifier(myGoogleLogin)

With this setup, LoginEndpoint will simply process whatever POST data you get from having executed a google login. There are many ways to do this, according to the google's docs, embedding html, javascript. Use their guide. However, our module also supplies a default "Log In To Google Page" using their embeded HTML technique, so we do:

myBouncer.AccessDenied = myGoogleLogin.DefaultLoginPortal()

So, sessions (User and Server) must add session info the userinfo. If the session is expired, eliminate and add "expired" somehow.They will hook into loggedin, loggedout where it will record session tokens or delete them. it can hook into abouttoload if you want to supply additional information to the header, mainly this is used for informing the client about themselves, but you can also add that information at any other point that you're called to the header.

If at any point you need to communicate to receive more with the user than what you recieve in the header or through /login. you may set the user status to spoken. the process and all of its state data should eventually expire if, for example, the user initiates logout (or it is initiated on their behalf). TODO not really sure how this will work.

User has to reconcile disagreement in sessions.
