# Bug hunting

In the application, the signup page stopped working.
By looking at the logs it seemed that the password field was empty.
It was discovered that the API expected the field in the request body to be named `pwd` as per the simulator specification. The bug was fixed by updating the client code to use `pwd` instead of `password` as the name.

See https://github.com/aske-w/itu-minitwit/commit/93ad407fa86d9c3f22d92eaa82adb8decc28555c
