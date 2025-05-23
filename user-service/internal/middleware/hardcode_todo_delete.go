package middleware

import "time"

// TODO: DELETE
const Alg string = "HS256"                                       // TODO: remove
const Secret string = "a-string-secret-at-least-256-bits-long"   // TODO: remove
const Issuer string = "user-service"                             // TODO: remove
var Audience []string = []string{"user-service", "test-service"} // TODO: remove
const ExpirationTimeDuration time.Duration = time.Minute * 15    // TODO: remove
