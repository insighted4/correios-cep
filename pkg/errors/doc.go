// Package errors provides a central interface to handling
// all errors in the Correios CEP domain. It closely follows Upspin's error
// handling with a few differences that reflect the system design
// of the Crawler's architecture. If you're unfamiliar with Upspin's error
// handling, we recommend that you read this article first before
// coming back here https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html.
// Crawler's errors are central around dealing with modules. So every error
// will most likely carry a Obj and a Version inside them. Furthermore,
// because Crawler is designed to run on multiple clouds, we have to design
// our errors and our logger to be friendly with each other. Therefore,
// the logger's SystemError method, although it accepts any type of error,
// it knows how to deal with errors constructed from this package in a debuggable way.
// To construct an Crawler error, call the errors.E function. The E function takes
// an Op and a variadic interface{}, but the values of the Error struct are what you can
// pass to it. Values such as the error Kind, Obj, Version, Error Message,
// and Severity (seriousness of an error) are all optional. The only truly required value is
// the errors.Op so you can construct a traceable stack that leads to where
// the error happened. However, adding more information can help catch an issue
// quicker and would help Cloud Log Monitoring services be more efficient to maintainers
// as you can run queries on Crawler Errors such as "Give me all errors of KindUnexpected"
// or "Give me all errors where caused by a particular Obj"
package errors
