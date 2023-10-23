# usersessioncookie

should add encryption option

`usersessioncookie` takes a string, puts it in a cookie, and gives you back the string on read. It does no expiry or encryption, but it should.

`New(domain string, path string, expiry time.Duration, key string)` to create it.

If expiry is set to 0, (maybe `0 * time.Second`), there is no set (it will expire when they close the browser!). Key must be a path to a private key.

Every instance comes with a unique UUID and cookie name to not conflict other cookies or instances of usersessioncookie.

Your state object must implent:

```go

type ReqBySess interface{

	// Write a function to take your state object and convert it into a string
	StateToSession() string
	
	
	// Write a function that takes a string and a) puts it into your state object but only if 
	// b) it's valid. If it's weird, or expired, return false and pretend you never saw it.
	// usersessioncookie will treat it as dirty and delete it.
	SesstionToState(session string, expiry bool) bool
}
```


