// This package has some of middleware that built-in goslim framework.
// You can use
//  g.Router.Use(middleware.LogMiddleware)
// to add middleware.
// And also write it yourself.
//
// Example
//
// You can write your middleware accord with
//  func YourMiddleware(next goslim.CHandler) goslim.CHandler {
//  	return func(c *goslim.Context) error {
//  		// do something before next handler
//
//  		err := next(c)
//
//  		// do something after next handler
//
//  		return err
//  	}
//  }
package middleware
