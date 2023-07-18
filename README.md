# DEVELOPS ONLY 

Urgent:

* It's really onnly allowing one userdata.
* Clarify API and attachment.
* Do multipe emails work? No, everything log's in as ajpikul.
* More tests.
* Can hooks recieve user data
* Force configurations through interfaces

Concepts
We have a concept of store/retrieval interface, we call it a session.
-- the one we're using currently wants and gives back a string. Obviously that can't be required.
-- we have a login mechanism that allows us to process logins.
essentially the user should generate stagemanager that's alive with the reqest
--- it can be modified when the session is read, and the session can be updated too (validate session)
--- it can be modified when the session is updated (update session)

we need need to be able to
--- initiate session
--- renew session
--- delete session (logout)
--- expire session

Yes. That's it. How do we stop people from doing persistent storage stuff in what they make? They really can't.
