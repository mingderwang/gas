// This package defined and implemented databases model.
// There have three interface:
//  interface_db.go
//  interface_builder.go
//  interface_model.go
// The DB interface defined database connection, query and transaction functions.
//
// The BuilderInterface defined how to generate sql statement and get the results.
//
// The ModelInterface defined orm functions.
//
// And now only support MySQL database, if you want to use another database,
// you can implement the three interfaces.
package model
